package server

import (
	"context"
	"log/slog"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/gbh007/buttoners/core/redis"
	"github.com/gbh007/buttoners/services/auth/internal/storage"
)

type Server struct {
	logger            *slog.Logger
	httpServerMetrics *metrics.HTTPServerMetrics
	comCfg            CommunicationConfig
	serviceName       string

	db    *storage.Database
	redis *redis.Client[dto.UserInfo]
	token string
}

func New(logger *slog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Init(ctx context.Context, comCfg CommunicationConfig, cfg DBConfig, serviceName string) error {
	httpServerMetrics, err := metrics.NewHTTPServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	redisMetrics, err := metrics.NewRedisMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	redisClient := redis.New[dto.UserInfo](comCfg.RedisAddress)

	err = redisClient.Connect(ctx, observability.NewRedisHook(s.logger, redisMetrics, comCfg.RedisAddress, serviceName))
	if err != nil {
		return err
	}

	defer redisClient.Close()

	db, err := storage.New(ctx, cfg.Username, cfg.Password, cfg.Addr, cfg.DatabaseName)
	if err != nil {
		return err
	}

	s.db = db
	s.redis = redisClient
	s.token = comCfg.SelfToken
	s.httpServerMetrics = httpServerMetrics
	s.comCfg = comCfg
	s.serviceName = serviceName

	return nil
}

func (s *Server) Close(ctx context.Context) error {
	if s.redis != nil {
		return s.redis.Close()
	}

	return nil
}
