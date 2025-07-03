package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
)

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type Config struct {
	SelfAddress       string
	SelfToken         string
	PrometheusAddress string
	DB                DBConfig
}

func Run(ctx context.Context, l *slog.Logger, cfg Config) error {
	go metrics.Run(l, metrics.Config{Addr: cfg.PrometheusAddress})

	tracer := otel.GetTracerProvider().Tracer("notification-server")

	httpServerMetrics, err := metrics.NewHTTPServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	s := &server{
		db:      db,
		token:   cfg.SelfToken,
		tracer:  tracer,
		logger:  l,
		metrics: httpServerMetrics,
	}
	// FIXME: добавить авторизацию
	server := &fasthttp.Server{
		Handler: s.handle,
	}

	go func() {
		<-ctx.Done()
		sCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := server.ShutdownWithContext(sCtx)
		l.Error("shutdown web server", "error", err.Error())
	}()

	err = server.ListenAndServe(cfg.SelfAddress)
	if err != nil {
		return err
	}

	return nil
}
