package server

import (
	"context"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
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

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	s := &server{
		db:    db,
		token: cfg.SelfToken,
	}
	// FIXME: добавить логирование и авторизацию
	server := &fasthttp.Server{
		Handler: s.handle,
	}

	go func() {
		<-ctx.Done()
		sCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		server.ShutdownWithContext(sCtx)
	}()

	err = server.ListenAndServe(cfg.SelfAddress)
	if err != nil {
		return err
	}

	return nil
}
