package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"github.com/gbh007/buttoners/services/worker/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	tracer trace.Tracer
	logger *slog.Logger

	notification *notificationclient.Client
	rabbitClient *rabbitmq.Reader[dto.RabbitMQData]

	db  *storage.Database
	cfg Config
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Init(ctx context.Context, cfg Config) error {
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
		s.logger, otel.GetTracerProvider().Tracer("notification-client"), httpClientMetrics,
		cfg.NotificationService.Addr, cfg.NotificationService.Token, "worker-service",
	)
	if err != nil {
		return err
	}

	s.notification = notificationClient
	s.db = db
	s.tracer = otel.GetTracerProvider().Tracer(cfg.ServiceName)
	s.cfg = cfg

	rabbitClient, err := rabbitmq.NewReader[dto.RabbitMQData](
		s.logger,
		cfg.RabbitMQ.Username,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Addr,
		cfg.RabbitMQ.QueueName,
		queueReaderMetrics,
		s.handle,
	)
	if err != nil {
		return err
	}

	s.rabbitClient = rabbitClient

	return nil
}

func (s *Server) Close() error {
	var errs []error
	if s.notification != nil {
		err := s.notification.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.rabbitClient != nil {
		err := s.rabbitClient.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}
