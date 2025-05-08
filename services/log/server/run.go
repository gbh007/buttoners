package server

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/log/internal/pb"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
)

func Run(ctx context.Context, cfg Config) error {
	go metrics.Run(metrics.Config{Addr: cfg.PrometheusAddress})

	db, err := storage.Init(ctx, cfg.DB.Username, cfg.DB.Password, cfg.DB.Addr, cfg.DB.DatabaseName)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", cfg.SelfAddress)
	if err != nil {
		return err
	}

	kafkaClient := kafka.New(cfg.Kafka.Addr, cfg.Kafka.Topic, cfg.Kafka.GroupID, cfg.Kafka.NumPartitions)

	err = kafkaClient.Connect(cfg.Kafka.NumPartitions > 0)
	if err != nil {
		return err
	}

	defer kafkaClient.Close()

	handler := &handler{
		kafka:  kafkaClient,
		db:     db,
		tracer: otel.GetTracerProvider().Tracer(cfg.ServiceName),
	}

	server := &pbServer{
		db: db,
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(logInterceptor),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterLogServer(grpcServer, server)

	go func() {
		<-ctx.Done()
		grpcServer.GracefulStop()
	}()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := handler.Run(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()

		err := grpcServer.Serve(lis)
		if err != nil {
			log.Println(err)
		}
	}()

	wg.Wait()

	return nil
}
