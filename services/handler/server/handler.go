package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	tracer trace.Tracer
	logger *slog.Logger
	rabbitClient *rabbitmq.Client[dto.RabbitMQData]
}

func (h *handler) handle(
	ctx context.Context, key string, data *dto.KafkaTaskData,
) {
	ctx, span := h.tracer.Start(ctx, "handle msg")
	defer span.End()

	startTime := time.Now()

	rabbitCtx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
	defer rabbitCnl()

	err := h.rabbitClient.Write(rabbitCtx, dto.RabbitMQData{
		RequestID: key,
		UserID:    data.UserID,
		Chance:    data.Chance,
		Duration:  data.Duration,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "handle error")

		logger.LogWithMeta(h.logger, ctx, slog.LevelError, "write to rabbitmq", "error", err.Error(), "msg_key", key)
	}

	registerHandleTime(time.Since(startTime))
}
