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

	redisClient := redis.New[dto.UserInfo](cfg.RedisAddress)

	err = redisClient.Connect(ctx, observability.NewRedisHook(s.logger, redisMetrics, cfg.RedisAddress, metrics.InstanceName))
	if err != nil {
		return err
	}

	notificationClient, err := notificationclient.New(
		s.logger, otel.GetTracerProvider().Tracer("notification-client"), httpClientMetrics,
		cfg.NotificationService.Addr, cfg.NotificationService.Token, metrics.InstanceName,
	)
	if err != nil {
		return err
	}

	logClient, err := logclient.New(s.logger, otel.GetTracerProvider().Tracer("log-client"), httpClientMetrics, cfg.LogService.Addr, cfg.LogService.Token, metrics.InstanceName)
	if err != nil {
		return err
	}

	kafkaTaskClient := kafka.NewProducer[dto.KafkaTaskData](s.logger, cfg.Kafka.Addr, cfg.Kafka.TaskTopic, queueWriterMetrics)

	kafkaLogClient := kafka.NewProducer[dto.KafkaLogData](s.logger, cfg.Kafka.Addr, cfg.Kafka.LogTopic, queueWriterMetrics)

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

func (s *Server) Close(ctx context.Context) error {
	var errs []error

	if s.auth != nil {
		err := s.auth.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.redis != nil {
		err := s.redis.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.notification != nil {
		err := s.notification.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.log != nil {
		err := s.log.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.kafkaTask != nil {
		err := s.kafkaTask.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if s.kafkaLog != nil {
		err := s.kafkaLog.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}

	return nil
}
