package server

import (
	"context"
	"net"

	"github.com/gbh007/buttoners/core/clients/authclient"
	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/clients/logclient"
	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/redis"
	"github.com/gbh007/buttoners/services/gate/dto"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})
	const serviceName = "gate-service"

	authClient, err := authclient.New(cfg.AuthService.Addr, cfg.AuthService.Token, serviceName)
	if err != nil {
		return err
	}

	defer authClient.Close()

	redisClient := redis.New[dto.UserInfo](cfg.RedisAddress)

	err = redisClient.Connect(ctx)
	if err != nil {
		return err
	}

	defer redisClient.Close()

	notificationClient, err := notificationclient.New(cfg.NotificationService.Addr, cfg.NotificationService.Token, serviceName)
	if err != nil {
		return err
	}

	defer notificationClient.Close()

	logClient, err := logclient.New(cfg.LogService.Addr, cfg.LogService.Token, serviceName)
	if err != nil {
		return err
	}

	defer logClient.Close()

	kafkaTaskClient := kafka.New(cfg.Kafka.Addr, cfg.Kafka.TaskTopic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions)

	err = kafkaTaskClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaTaskClient.Close()

	kafkaLogClient := kafka.New(cfg.Kafka.Addr, cfg.Kafka.LogTopic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions)

	err = kafkaLogClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

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
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.logInterceptor),
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
