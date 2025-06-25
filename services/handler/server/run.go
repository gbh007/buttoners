package server

import (
	"context"
	"log/slog"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"go.opentelemetry.io/otel"
)

func Run(ctx context.Context, l *slog.Logger, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	queueReaderMetrics, err := metrics.NewQueueReaderMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	queueWriterMetrics, err := metrics.NewQueueWriterMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	kafkaClient := kafka.New(l, cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions, queueReaderMetrics, queueWriterMetrics)

	err = kafkaClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	rabbitClient := rabbitmq.New[dto.RabbitMQData](
		l,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Addr,
		cfg.RabbitMQ.QueueName,
		queueReaderMetrics,
		queueWriterMetrics,
	)

	err = rabbitClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer rabbitClient.Close()

	h := handler{
		tracer: otel.GetTracerProvider().Tracer(cfg.ServiceName),
	}

label1:
	for {
		data := new(dto.KafkaTaskData)
		ctx, key, err := kafkaClient.Read(ctx, data)
		if err != nil {
			logger.LogWithMeta(l, ctx, slog.LevelWarn, "kafka read", "error", err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		h.handle(ctx, key, data, rabbitClient)
	}

	return nil
}
