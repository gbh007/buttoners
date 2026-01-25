package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/tracer"
	"github.com/gbh007/buttoners/services/legacy/internal/controller"
	"github.com/kelseyhightower/envconfig"
)

type Config struct{}

func main() {
	cfg := controller.Config{}

	envconfig.MustProcess("", &cfg)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	const serviceName = "legacy"

	l := logger.New(serviceName, "debug") // FIXME: level
	metrics.InstanceName = serviceName

	_, _, err := tracer.InitTracer(cfg.JaegerURL, metrics.InstanceName)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail init tracer", "error", err.Error())
		os.Exit(1)
	}

	go metrics.Run(l, metrics.Config{Addr: cfg.PrometheusAddr})

	c, err := controller.New(l, cfg)
	if err != nil {
		l.Error("create controller", "error", err)
		os.Exit(1)
	}

	l.Info("start server")

	err = c.Serve(ctx)
	if err != nil {
		l.Error("serve http", "error", err)
		os.Exit(1)
	}

	l.Info("have a nice day")
}
