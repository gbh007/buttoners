package server

import (
	"context"
	"log/slog"
	"sync"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"github.com/gbh007/buttoners/services/worker/internal/storage"

	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context, l *slog.Logger, cfg Config) error {
	go metrics.Run(l, metrics.Config{Addr: cfg.PrometheusAddress})

	httpClientMetrics, err := metrics.NewHTTPClientMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	queueReaderMetrics, err := metrics.NewQueueReaderMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	notificationClient, err := notificationclient.New(
		l, otel.GetTracerProvider().Tracer("notification-client"), httpClientMetrics,
		cfg.NotificationService.Addr, cfg.NotificationService.Token, "worker-service",
	)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	rn := &runner{
		notification: notificationClient,
		db:           db,
		tracer:       otel.GetTracerProvider().Tracer(cfg.ServiceName),
		logger:       l,
	}

	rabbitClient, err := rabbitmq.NewReader[dto.RabbitMQData](
		l,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Addr,
		cfg.RabbitMQ.QueueName,
		queueReaderMetrics,
		rn.handle,
	)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	rn.rabbitClient = rabbitClient

	runnerCtx, runnerCnl := context.WithCancel(context.TODO())
	runnerWg := new(sync.WaitGroup)

	for i := 0; i < cfg.RunnerCount; i++ {
		runnerWg.Add(1)

		go func() {
			defer runnerWg.Done()

			rnErr := rn.run(runnerCtx)
			if rnErr != nil {
				l.Error("runner end unsuccessful", "error", rnErr.Error())
			}
		}()
	}

	<-ctx.Done()

	runnerCnl()

	runnerWg.Wait()

	return nil
}
