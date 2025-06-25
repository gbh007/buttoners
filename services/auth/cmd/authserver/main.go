package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/config"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/tracer"
	"github.com/gbh007/buttoners/services/auth/server"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self           config.Service
	DB             config.Database
	RedisAddr      string `envconfig:"default=redis:6379"`
	PrometheusAddr string `envconfig:"default=pushgateway:9091"`
	Jaeger         config.Jaeger
}

func main() {
	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	const serviceName = "auth"

	l := logger.New(serviceName, "debug") // FIXME: level
	metrics.InstanceName = serviceName

	cfg := new(Config)

	err := envconfig.Init(cfg)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail parse config", "error", err.Error())
		os.Exit(1)
	}

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, serviceName)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail init tracer", "error", err.Error())
		os.Exit(1)
	}

	logger.LogWithMeta(l, ctx, slog.LevelInfo, "server start")
	defer logger.LogWithMeta(l, ctx, slog.LevelInfo, "server stop")

	err = server.Run(
		ctx,
		l,
		server.CommunicationConfig{
			SelfAddress:       cfg.Self.Addr,
			SelfToken:         cfg.Self.Token,
			RedisAddress:      cfg.RedisAddr,
			PrometheusAddress: cfg.PrometheusAddr,
		},
		server.DBConfig{
			Username:     cfg.DB.User,
			Password:     cfg.DB.Pass,
			Addr:         cfg.DB.Addr,
			DatabaseName: cfg.DB.Name,
		},
	)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "unsuccess server run result", "error", err.Error())
		os.Exit(1)
	}
}
