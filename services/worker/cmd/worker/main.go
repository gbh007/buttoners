package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/config"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/tracer"
	"github.com/gbh007/buttoners/services/worker/server"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	RabbitMQ         config.RabbitMQ
	DB               config.Database
	NotificationAddr string `envconfig:"default=notification:50051"`
	PrometheusAddr   string `envconfig:"default=pushgateway:9091"`
	RunnerCount      int    `envconfig:"default=20"`
	Jaeger           config.Jaeger
}

func main() {
	cfg := new(Config)

	err := envconfig.Init(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancelNotify := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancelNotify()

	log.Println("server start")

	metrics.InstanceName = "worker"

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	err = server.Run(
		ctx,
		server.Config{
			ServiceName:         metrics.InstanceName,
			NotificationAddress: cfg.NotificationAddr,
			PrometheusAddress:   cfg.PrometheusAddr,
			DB: server.DBConfig{
				Username:     cfg.DB.User,
				Password:     cfg.DB.Pass,
				Addr:         cfg.DB.Addr,
				DatabaseName: cfg.DB.Name,
			},
			RabbitMQ: server.RabbitMQConfig{
				Username:  cfg.RabbitMQ.User,
				Password:  cfg.RabbitMQ.Pass,
				Addr:      cfg.RabbitMQ.Addr,
				QueueName: cfg.RabbitMQ.Queue,
			},
			RunnerCount: cfg.RunnerCount,
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
