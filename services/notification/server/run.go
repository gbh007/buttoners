package server

import (
	"context"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
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

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.logger, metrics.Config{Addr: s.cfg.PrometheusAddress})

	// FIXME: добавить авторизацию
	server := &fasthttp.Server{
		Handler: s.handle,
	}

	go func() {
		<-ctx.Done()
		sCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		err := server.ShutdownWithContext(sCtx)
		s.logger.Error("shutdown web server", "error", err.Error())
	}()

	err := server.ListenAndServe(s.cfg.SelfAddress)
	if err != nil {
		return err
	}

	return nil
}
