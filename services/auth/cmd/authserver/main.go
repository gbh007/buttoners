package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/gbh007/buttoners/core/config"
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

	metrics.InstanceName = "auth"

	_, _, err = tracer.InitTracer(cfg.Jaeger.URL, metrics.InstanceName)
	if err != nil {
		log.Fatalln(err) //nolint:gocritic
	}

	err = server.Run(ctx,
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
		log.Println(err)
	}

	log.Println("server stop")
}
