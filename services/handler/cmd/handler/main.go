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
	"github.com/gbh007/buttoners/services/handler/server"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	RabbitMQ       config.RabbitMQ
	Kafka          config.Kafka
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

	const serviceName = "handler"

	l := logger.New(serviceName, "debug") // FIXME: level
	metrics.InstanceName = serviceName

	cfg := new(Config)

	err := envconfig.Init(cfg)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail parse config", "error", err.Error())
		os.Exit(1)
	}

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail init tracer", "error", err.Error())
		os.Exit(1)
	}

	logger.LogWithMeta(l, ctx, slog.LevelInfo, "server start")
	defer logger.LogWithMeta(l, ctx, slog.LevelInfo, "server stop")

	srvConf := server.Config{
		ServiceName:       metrics.InstanceName,
		PrometheusAddress: cfg.PrometheusAddr,
		Kafka: server.KafkaConfig{
			Addr:    cfg.Kafka.Addr,
			Topic:   cfg.Kafka.TaskTopic,
			GroupID: cfg.Kafka.GroupID,
		},
		RabbitMQ: server.RabbitMQConfig{
			Username:  cfg.RabbitMQ.User,
			Password:  cfg.RabbitMQ.Pass,
			Addr:      cfg.RabbitMQ.Addr,
			QueueName: cfg.RabbitMQ.Queue,
		},
	}

	server := server.New(l)

	err = server.Init(ctx, srvConf)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "fail server init", "error", err.Error())
		os.Exit(1)
	}

	err = server.Run(ctx)
	if err != nil {
		logger.LogWithMeta(l, ctx, slog.LevelWarn, "unsuccess server run result", "error", err.Error())
		os.Exit(1)
	}
}
