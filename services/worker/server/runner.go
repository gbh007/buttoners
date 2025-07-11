package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/core/rabbitmq"
	"github.com/gbh007/buttoners/services/worker/internal/storage"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type runner struct {
	tracer trace.Tracer
	logger *slog.Logger

	notification *notificationclient.Client
	rabbitClient *rabbitmq.Reader[dto.RabbitMQData]

	db *storage.Database
}

func (r *runner) run(ctx context.Context) error {
	return r.rabbitClient.Start(ctx)
}

func (r *runner) handle(ctx context.Context, key string, data dto.RabbitMQData) error {
	activeTaskTotal.Inc()
	defer activeTaskTotal.Dec()

	ctx, span := r.tracer.Start(ctx, "handle msg")
	defer span.End()

	startTime := time.Now()

	n := notificationclient.NewRequest{
		UserID: data.UserID,
		Kind:   "button",
	}

	errText := ""

	result, resultText, err := r.someBusinessLogic(ctx, data.Duration, data.Chance)
	if err != nil {
		n.Level = "error"
		n.Title = "Ошибка"
		n.Body = fmt.Sprintf("Ошибка во время выполнения:\n%s", err.Error())

		errText = err.Error()

		span.RecordError(err)
		span.SetStatus(codes.Error, "business")
	} else {
		n.Level = "success"
		n.Title = "Завершено"
		n.Body = resultText
	}

	businessEndTime := time.Now()

	defer func() {
		// Общее время выполнения
		registerHandleTime(time.Since(startTime))
		// Бизнесовое время выполнения
		registerBusinessHandleTime(errText == "", businessEndTime.Sub(startTime))
	}()

	logger.LogWithMeta(r.logger, ctx, slog.LevelInfo, "finished", "data_request_id", data.RequestID, "notification", n)

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err = r.db.InsertTaskResult(dbCtx, &storage.TaskResult{
		UserID:     data.UserID,
		Chance:     data.Chance,
		Duration:   data.Duration,
		Result:     result,
		ResultText: resultText,
		ErrorText:  errText,
		StartTime:  startTime,
		EndTime:    businessEndTime,
	})
	if err != nil {
		logger.LogWithMeta(r.logger, ctx, slog.LevelError, "write to task result", "error", err.Error(), "data_request_id", data.RequestID)

		span.RecordError(err)
		span.SetStatus(codes.Error, "insert result")

		return err
	}

	notificationCtx, notificationCnl := context.WithTimeout(ctx, time.Second*10)
	defer notificationCnl()

	err = r.notification.New(notificationCtx, n)
	if err != nil {
		logger.LogWithMeta(r.logger, ctx, slog.LevelError, "send notification", "error", err.Error(), "data_request_id", data.RequestID)

		span.RecordError(err)
		span.SetStatus(codes.Error, "new notification")

		return err
	}

	return nil
}
