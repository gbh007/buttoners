package server

import (
	"context"
	"log/slog"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Server struct {
	db      *storage.Database
	tracer  trace.Tracer
	logger  *slog.Logger
	metrics *metrics.HTTPServerMetrics
	cfg     Config
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Init(ctx context.Context, cfg Config) error {
	tracer := otel.GetTracerProvider().Tracer("notification-server")

	httpServerMetrics, err := metrics.NewHTTPServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	s.db = db
	s.tracer = tracer
	s.metrics = httpServerMetrics
	s.cfg = cfg

	return nil
}
