package server

import (
	"context"
	"log/slog"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	db                *storage.Database
	l                 *slog.Logger
	tracer            trace.Tracer
	cfg               Config
	kafka             *kafka.Consumer[dto.KafkaLogData]
	httpServerMetrics *metrics.HTTPServerMetrics
}

func New(l *slog.Logger) *Server {
	return &Server{
		l: l,
	}
}

func (s *Server) Init(ctx context.Context, cfg Config) error {
	httpServerMetrics, err := metrics.NewHTTPServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
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

	kafkaClient := kafka.NewConsumer(s.l, cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, queueReaderMetrics, s.handle)

	s.db = db
	s.httpServerMetrics = httpServerMetrics
	s.tracer = otel.GetTracerProvider().Tracer(cfg.ServiceName)
	s.kafka = kafkaClient
	s.cfg = cfg

	return nil
}
