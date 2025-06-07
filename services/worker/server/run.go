package server

import (
	"context"
	"sync"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"github.com/gbh007/buttoners/services/worker/internal/storage"

	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	rabbitClient := rabbitmq.New[dto.RabbitMQData](
		cfg.RabbitMQ.Username, cfg.RabbitMQ.Password, cfg.RabbitMQ.Addr, cfg.RabbitMQ.QueueName,
	)

	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	notificationClient, err := notificationclient.New(cfg.NotificationService.Addr, cfg.NotificationService.Token, "worker-service")
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	messages, err := rabbitClient.StartRead(ctx)
	if err != nil {
		return err
	}

	runnerCtx, runnerCnl := context.WithCancel(context.TODO())
	runnerWg := new(sync.WaitGroup)

	for i := 0; i < cfg.RunnerCount; i++ {
		runnerWg.Add(1)

		r := &runner{
			notification: notificationClient,
			db:           db,
			queue:        messages,
			tracer:       otel.GetTracerProvider().Tracer(cfg.ServiceName),
		}

		go func() {
			defer runnerWg.Done()
			r.run(runnerCtx)
		}()
	}

	<-ctx.Done()

	runnerCnl()

	runnerWg.Wait()

	return nil
}
