package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"
)

func New(serviceName, level string) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: getLogLevel(level),
	}))
	logger = logger.With("service_name", serviceName)

	return logger
}

func getLogLevel(logLevel string) slog.Level {
	logLevel = strings.ToLower(logLevel)

	switch logLevel {
	case "trace", "debug", "dbg":
		return slog.LevelDebug
	case "info", "inf":
		return slog.LevelInfo
	case "err", "error":
		return slog.LevelError
	case "wrn", "warn", "warning":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}

func LogWithMeta(l *slog.Logger, ctx context.Context, level slog.Level, msg string, args ...any) {
	args = append([]any{}, args...)
	args = append(args, slog.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String()))
	l.Log(ctx, level, msg, args...)
}
