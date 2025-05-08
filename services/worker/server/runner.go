package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gbh007/buttoners/core/rabbitmq"
	handlerdto "github.com/gbh007/buttoners/services/handler/dto"
	notificationServerClient "github.com/gbh007/buttoners/services/notification/client"
	"github.com/gbh007/buttoners/services/worker/internal/storage"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type runner struct {
	tracer trace.Tracer

	notification *notificationServerClient.Client
	db           *storage.Database
	queue        chan rabbitmq.Read[handlerdto.RabbitMQData]
}

func (r *runner) run(ctx context.Context) {
	for {
		select {
		case dataReader := <-r.queue:
			r.handle(ctx, dataReader)

		case <-ctx.Done():
			if len(r.queue) == 0 {
				return
			}
		}
	}
}

func (r *runner) handle(ctx context.Context, dataReader rabbitmq.Read[handlerdto.RabbitMQData]) {
	activeTaskTotal.Inc()
	defer activeTaskTotal.Dec()

	ctx, data, err := dataReader(ctx)
	if err != nil {
		log.Println(err)

		return
	}

	ctx, span := r.tracer.Start(ctx, "handle msg")
	defer span.End()

	log.Printf("accept %#+v\n", data)

	startTime := time.Now()

	n := &notificationServerClient.Notification{
		Kind: notificationServerClient.ButtonKind,
	}

	errText := ""

	result, resultText, err := r.someBusinessLogic(ctx, data.Duration, data.Chance)
	if err != nil {
		n.Level = notificationServerClient.ErrorLevel
		n.Title = "Ошибка"
		n.Body = fmt.Sprintf("Ошибка во время выполнения:\n%s", err.Error())

		errText = err.Error()

		span.RecordError(err)
		span.SetStatus(codes.Error, "business")
	} else {
		n.Level = notificationServerClient.SuccessLevel
		n.Title = "Завершено"
		n.Body = resultText
	}

	businessEndTime := time.Now()

	log.Printf("finished %s = %#+v\n", data.RequestID, n)

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
		log.Println(err)

		span.RecordError(err)
		span.SetStatus(codes.Error, "insert result")
	}

	notificationCtx, notificationCnl := context.WithTimeout(ctx, time.Second*10)
	defer notificationCnl()

	err = r.notification.New(notificationCtx, data.UserID, n)
	if err != nil {
		log.Println(err)

		span.RecordError(err)
		span.SetStatus(codes.Error, "new notification")
	}

	// Общее время выполнения
	registerHandleTime(time.Since(startTime))
	// Бизнесовое время выполнения
	registerBusinessHandleTime(errText == "", businessEndTime.Sub(startTime))
}
