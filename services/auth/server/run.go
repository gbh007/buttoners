package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/redis"
	"github.com/gbh007/buttoners/services/auth/internal/storage"
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

func Run(ctx context.Context, comCfg CommunicationConfig, cfg DBConfig) error {
	go metrics.Run(metrics.Config{Addr: comCfg.PrometheusAddress})

	redisClient := redis.New[dto.UserInfo](comCfg.RedisAddress)

	err := redisClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer redisClient.Close()

	db, err := storage.New(ctx, cfg.Username, cfg.Password, cfg.Addr, cfg.DatabaseName)
	if err != nil {
		return err
	}

	s := &server{
		db:    db,
		redis: redisClient,
		token: comCfg.SelfToken,
	}

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
		Addr:    comCfg.SelfAddress,
		Handler: otelhttp.NewHandler(router, "Auth server"),
	}

	go func() {
		<-ctx.Done()
		sCtx, _ := context.WithTimeout(context.Background(), time.Second*10)
		server.Shutdown(sCtx)
	}()

	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
