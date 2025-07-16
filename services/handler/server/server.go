package server

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	tracer       trace.Tracer
	logger       *slog.Logger
	rabbitClient *rabbitmq.Writer[dto.RabbitMQData]
	kafkaClient  *kafka.Consumer[dto.KafkaTaskData]
	cfg          Config
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Init(ctx context.Context, cfg Config) error {
	queueReaderMetrics, err := metrics.NewQueueReaderMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	queueWriterMetrics, err := metrics.NewQueueWriterMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	rabbitClient, err := rabbitmq.NewWriter[dto.RabbitMQData](
		s.logger,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Addr,
		cfg.RabbitMQ.QueueName,
		queueWriterMetrics,
	)
	if err != nil {
		return err
	}

	s.tracer = otel.GetTracerProvider().Tracer(cfg.ServiceName)
	s.rabbitClient = rabbitClient
	s.kafkaClient = kafka.NewConsumer(s.logger, cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, queueReaderMetrics, s.handle)
	s.cfg = cfg

	return nil
}
func (s *Server) Close(ctx context.Context) error {
	var errs []error

	if s.rabbitClient != nil {
		err := s.rabbitClient.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.kafkaClient != nil {
		err := s.kafkaClient.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (s *Server) handle(
	ctx context.Context, key string, data dto.KafkaTaskData,
) error {
	ctx, span := s.tracer.Start(ctx, "handle msg")
	defer span.End()

	rabbitCtx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
	defer rabbitCnl()

	err := s.rabbitClient.Write(rabbitCtx, key, dto.RabbitMQData{
		RequestID: key,
		UserID:    data.UserID,
		Chance:    data.Chance,
		Duration:  data.Duration,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "handle error")

		logger.LogWithMeta(s.logger, ctx, slog.LevelError, "write to rabbitmq", "error", err.Error(), "msg_key", key)

		return err
	}

	return nil
}
