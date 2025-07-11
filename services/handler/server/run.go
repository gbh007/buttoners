package server

import (
	"context"
	"log/slog"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context, l *slog.Logger, cfg Config) error {
	go metrics.Run(l, metrics.Config{Addr: cfg.PrometheusAddress})

	queueReaderMetrics, err := metrics.NewQueueReaderMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	queueWriterMetrics, err := metrics.NewQueueWriterMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	rabbitClient, err := rabbitmq.NewWriter[dto.RabbitMQData](
		l,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Addr,
		cfg.RabbitMQ.QueueName,
		queueWriterMetrics,
	)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	h := handler{
		tracer:       otel.GetTracerProvider().Tracer(cfg.ServiceName),
		logger:       l,
		rabbitClient: rabbitClient,
	}

	kafkaClient := kafka.NewConsumer(l, cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, queueReaderMetrics, func(ctx context.Context, key string, v dto.KafkaTaskData) error {
		h.handle(ctx, key, &v)

		return nil
	})

	defer kafkaClient.Close()

	return kafkaClient.Start(ctx)
}
