package server

import (
	"context"
	"log"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func logInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	addr := "unknown"
	routeName := "unknown"

	p, ok := peer.FromContext(ctx)
	if ok {
		addr = p.Addr.String()
	}

	if info != nil {
		routeName = info.FullMethod
	}

	log.Printf("handle %s %s\n", routeName, addr)

	requestStart := time.Now()

	resp, err = handler(ctx, req)

	metrics.LogRequest(routeName, err == nil, time.Since(requestStart))

	return
}
