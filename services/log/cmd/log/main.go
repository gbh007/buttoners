package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/config"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/tracer"
	"github.com/gbh007/buttoners/services/log/server"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self           config.Service
	Kafka          config.Kafka
	DB             config.Database
	PrometheusAddr string `envconfig:"default=pushgateway:9091"`
	Jaeger         config.Jaeger
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

	metrics.InstanceName = "log"

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	err = server.Run(
		ctx,
		server.Config{
			ServiceName:       metrics.InstanceName,
			SelfAddress:       cfg.Self.Addr,
			SelfToken:         cfg.Self.Token,
			PrometheusAddress: cfg.PrometheusAddr,
			Kafka: server.KafkaConfig{
				Addr:    cfg.Kafka.Addr,
				Topic:   cfg.Kafka.LogTopic,
				GroupID: cfg.Kafka.GroupID,
			},
			DB: server.DBConfig{
				Username:     cfg.DB.User,
				Password:     cfg.DB.Pass,
				Addr:         cfg.DB.Addr,
				DatabaseName: cfg.DB.Name,
			},
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
