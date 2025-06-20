package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/config"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/tracer"
	"github.com/gbh007/buttoners/services/gate/server"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	Self                config.Addr
	Kafka               config.Kafka
	AuthService         config.Service
	NotificationService config.Service
	LogService          config.Service
	RedisAddr           string `envconfig:"default=redis:6379"`
	PrometheusAddr      string `envconfig:"default=pushgateway:9091"`
	Jaeger              config.Jaeger
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

	metrics.InstanceName = "gate"

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	err = server.Run(
		ctx,
		server.Config{
			SelfAddress:         cfg.Self.Full(),
			AuthService:         cfg.AuthService,
			LogService:          cfg.LogService,
			NotificationService: cfg.NotificationService,
			RedisAddress:        cfg.RedisAddr,
			PrometheusAddress:   cfg.PrometheusAddr,
			Kafka: server.KafkaConfig{
				Addr:          cfg.Kafka.Addr,
				TaskTopic:     cfg.Kafka.TaskTopic,
				LogTopic:      cfg.Kafka.LogTopic,
				NumPartitions: cfg.Kafka.NumPartitions,
			},
		},
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("server stop")
}
