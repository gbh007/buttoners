package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type DBConfig struct {
	Username, Password, Addr, DatabaseName string
}

type CommunicationConfig struct {
	SelfAddress       string
	SelfToken         string
	RedisAddress      string
	PrometheusAddress string
}

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.logger, metrics.Config{Addr: s.comCfg.PrometheusAddress})

	defer s.Close(ctx)

	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		s.authMiddleWare,
	)

	router.Post(authclient.LoginPath, s.Login)
	router.Post(authclient.LogoutPath, s.Logout)
	router.Post(authclient.RegisterPath, s.Register)
	router.Post(authclient.InfoPath, s.Info)

	server := &http.Server{
		Addr:    s.comCfg.SelfAddress,
		Handler: otelhttp.NewHandler(observability.NewHTTPMiddleware(s.logger, s.httpServerMetrics, "auth", router), s.serviceName),
	}

	go func() {
		<-ctx.Done()
		sCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		server.Shutdown(sCtx)
	}()

	err := server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
