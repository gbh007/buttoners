package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"github.com/gbh007/buttoners/core/redis"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, l *slog.Logger, cfg Config) error {
	go metrics.Run(l, metrics.Config{Addr: cfg.PrometheusAddress})

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

	authClient, err := authclient.New(l, httpClientMetrics, cfg.AuthService.Addr, cfg.AuthService.Token, metrics.InstanceName)
	if err != nil {
		return err
	}

	defer authClient.Close()

	redisClient := redis.New[dto.UserInfo](cfg.RedisAddress)

	err = redisClient.Connect(ctx, observability.NewRedisHook(l, redisMetrics, cfg.RedisAddress, metrics.InstanceName))
	if err != nil {
		return err
	}

	defer redisClient.Close()

	notificationClient, err := notificationclient.New(
		l, otel.GetTracerProvider().Tracer("notification-client"), httpClientMetrics,
		cfg.NotificationService.Addr, cfg.NotificationService.Token, metrics.InstanceName,
	)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	logClient, err := logclient.New(l, otel.GetTracerProvider().Tracer("log-client"), httpClientMetrics, cfg.LogService.Addr, cfg.LogService.Token, metrics.InstanceName)
	if err != nil {
		return err
	}

	defer logClient.Close()

	kafkaTaskClient := kafka.NewProducer[dto.KafkaTaskData](l, cfg.Kafka.Addr, cfg.Kafka.TaskTopic, queueWriterMetrics)

	defer kafkaTaskClient.Close()

	kafkaLogClient := kafka.NewProducer[dto.KafkaLogData](l, cfg.Kafka.Addr, cfg.Kafka.LogTopic, queueWriterMetrics)

	defer kafkaLogClient.Close()

	lis, err := net.Listen("tcp", cfg.SelfAddress)
	if err != nil {
		return err
	}

	s := &pbServer{
		auth:         authClient,
		kafkaTask:    kafkaTaskClient,
		kafkaLog:     kafkaLogClient,
		notification: notificationClient,
		log:          logClient,
		redis:        redisClient,
		logger:       l,
	}

	obs := observability.NewGRPCServerInterceptor(l, grpcServerMetrics, metrics.InstanceName)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			obs.Unary,
			s.logInterceptor,
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterGateServer(grpcServer, s)
	pb.RegisterNotificationServer(grpcServer, s)
	pb.RegisterLogServer(grpcServer, s)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}
