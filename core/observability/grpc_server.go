package observability

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCServerInterceptor struct {
	logger     *slog.Logger
	metrics    *metrics.GRPCServerMetrics
	serverName string
}

func NewGRPCServerInterceptor(
	logger *slog.Logger,
	metrics *metrics.GRPCServerMetrics,
	serverName string,
) *GRPCServerInterceptor {
	return &GRPCServerInterceptor{
		logger:     logger,
		metrics:    metrics,
		serverName: serverName,
	}
}

func (h *GRPCServerInterceptor) Unary(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	tStart := time.Now()

	h.metrics.IncActive(h.serverName, info.FullMethod)
	defer h.metrics.DecActive(h.serverName, info.FullMethod)

	var (
		requestLog, responseLog []any
	)

	requestLog = append(
		requestLog,
		slog.String("host", h.serverName),
		slog.String("path", info.FullMethod),
	)

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		headers := make(map[string]string)

		for k, v := range md {
			headers[k] = strings.Join(v, ";")
		}

		requestLog = append(
			requestLog,
			slog.Any("headers", headers),
		)
	}

	if data, err := json.Marshal(req); err == nil {
		requestLog = append(
			requestLog,
			slog.String("body", string(data)),
		)
	}

	var statusCode codes.Code

	defer func() {
		h.metrics.AddHandle(h.serverName, info.FullMethod, statusCode.String(), time.Since(tStart))
	}()

	resp, err = handler(ctx, req)

	statusCode = status.Code(err)

	responseLog = append(
		responseLog,
		slog.String("status", statusCode.String()),
	)

	// FIXME: поддержать заголовки ответа (через grpc.StatsHandler)

	if data, err := json.Marshal(resp); err == nil {
		responseLog = append(
			responseLog,
			slog.String("body", string(data)),
		)
	}

	h.logger.InfoContext(
		ctx, h.serverName+" server request",
		slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
		slog.Group("request", requestLog...),
		slog.Group("response", responseLog...),
	)

	return
}
