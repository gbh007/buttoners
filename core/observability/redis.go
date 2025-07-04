package observability

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/metrics"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
)

type RedisHook struct {
	logger     *slog.Logger
	metrics    *metrics.RedisMetrics
	addr       string
	clientName string
}

func NewRedisHook(
	logger *slog.Logger,
	metrics *metrics.RedisMetrics,
	addr string,
	clientName string,
) *RedisHook {
	return &RedisHook{
		logger:     logger,
		metrics:    metrics,
		addr:       addr,
		clientName: clientName,
	}
}

func (rh *RedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (rh *RedisHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		tStart := time.Now()
		method := cmd.Name()

		rh.metrics.IncActive(rh.addr, method)
		defer rh.metrics.DecActive(rh.addr, method)

		err := next(ctx, cmd)

		var (
			notFound, hasError      bool
			requestLog, responseLog []any
		)

		requestLog = append(
			requestLog,
			slog.String("host", rh.addr),
			slog.String("method", method),
		)

		cmdArgs := cmd.Args()

		if method == "set" && len(cmdArgs) >= 3 {
			requestLog = append(
				requestLog,
				slog.Any("path", cmdArgs[1]),
				slog.Any("body", cmdArgs[2]),
			)
		}

		if method == "get" && len(cmdArgs) >= 2 {
			requestLog = append(
				requestLog,
				slog.Any("path", cmdArgs[1]),
			)
		}

		if err != nil || cmd.Err() != nil {
			hasError = true

			if errors.Is(cmd.Err(), redis.Nil) ||
				errors.Is(err, redis.Nil) {
				notFound = true
			}
		}

		status := "ok"
		switch {
		case notFound:
			status = "not_found"
		case hasError:
			status = "err"
		}

		defer func() {
			rh.metrics.AddHandle(rh.addr, method, status, time.Since(tStart))
		}()

		responseLog = append(
			responseLog,
			slog.String("status", status),
		)

		switch tVal := cmd.(type) {
		case interface{ Val() string }:
			responseLog = append(
				responseLog,
				slog.String("body", tVal.Val()),
			)
		}

		rh.logger.InfoContext(
			ctx, rh.clientName+" redis request",
			slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()),
			slog.Group("request", requestLog...),
			slog.Group("response", responseLog...),
		)

		return err
	}
}

func (rh *RedisHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
