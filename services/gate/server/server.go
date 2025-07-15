package server

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/gbh007/buttoners/core/redis"
	"go.opentelemetry.io/otel"
)

var errInvalidInputData = errors.New("invalid")

type Server struct {
	pb.UnimplementedGateServer
	pb.UnimplementedNotificationServer
	pb.UnimplementedLogServer

	cfg Config

	logger            *slog.Logger
	grpcServerMetrics *metrics.GRPCServerMetrics

	auth         *authclient.Client
	notification *notificationclient.Client
	log          *logclient.Client
	kafkaTask    *kafka.Producer[dto.KafkaTaskData]
	kafkaLog     *kafka.Producer[dto.KafkaLogData]
	redis        *redis.Client[dto.UserInfo]
}

func New(l *slog.Logger) *Server {
	return &Server{
		logger: l,
	}
}

func (s *Server) Init(ctx context.Context, cfg Config) error {
	httpClientMetrics, err := metrics.NewHTTPClientMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	grpcServerMetrics, err := metrics.NewGRPCServerMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	queueWriterMetrics, err := metrics.NewQueueWriterMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	redisMetrics, err := metrics.NewRedisMetrics(metrics.DefaultRegistry, metrics.DefaultTimeBuckets)
	if err != nil {
		return err
	}

	authClient, err := authclient.New(s.logger, httpClientMetrics, cfg.AuthService.Addr, cfg.AuthService.Token, metrics.InstanceName)
	if err != nil {
		return err
	}

	defer authClient.Close()

	redisClient := redis.New[dto.UserInfo](cfg.RedisAddress)

	err = redisClient.Connect(ctx, observability.NewRedisHook(s.logger, redisMetrics, cfg.RedisAddress, metrics.InstanceName))
	if err != nil {
		return err
	}

	defer redisClient.Close()

	notificationClient, err := notificationclient.New(
		s.logger, otel.GetTracerProvider().Tracer("notification-client"), httpClientMetrics,
		cfg.NotificationService.Addr, cfg.NotificationService.Token, metrics.InstanceName,
	)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	logClient, err := logclient.New(s.logger, otel.GetTracerProvider().Tracer("log-client"), httpClientMetrics, cfg.LogService.Addr, cfg.LogService.Token, metrics.InstanceName)
	if err != nil {
		return err
	}

	defer logClient.Close()

	kafkaTaskClient := kafka.NewProducer[dto.KafkaTaskData](s.logger, cfg.Kafka.Addr, cfg.Kafka.TaskTopic, queueWriterMetrics)

	defer kafkaTaskClient.Close()

	kafkaLogClient := kafka.NewProducer[dto.KafkaLogData](s.logger, cfg.Kafka.Addr, cfg.Kafka.LogTopic, queueWriterMetrics)

	defer kafkaLogClient.Close()

	s.auth = authClient
	s.kafkaTask = kafkaTaskClient
	s.kafkaLog = kafkaLogClient
	s.notification = notificationClient
	s.log = logClient
	s.redis = redisClient
	s.cfg = cfg
	s.grpcServerMetrics = grpcServerMetrics

	return nil
}
