package server

import (
	"context"
	"net"

	"github.com/gbh007/buttoners/core/clients/gateclient/gen/pb"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/core/observability"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func (s *Server) Run(ctx context.Context) error {
	go metrics.Run(s.logger, metrics.Config{Addr: s.cfg.PrometheusAddress})

	lis, err := net.Listen("tcp", s.cfg.SelfAddress)
	if err != nil {
		return err
	}

	obs := observability.NewGRPCServerInterceptor(s.logger, s.grpcServerMetrics, metrics.InstanceName)

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
