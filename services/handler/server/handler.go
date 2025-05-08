package server

import (
	"context"
	"log"
	"time"

	gatedto "github.com/gbh007/buttoners/services/gate/dto"
	handlerdto "github.com/gbh007/buttoners/services/handler/dto"

	"github.com/gbh007/buttoners/core/rabbitmq"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	tracer trace.Tracer
}

func (h *handler) handle(
	ctx context.Context, key string, data *gatedto.KafkaTaskData,
	rabbitClient *rabbitmq.Client[handlerdto.RabbitMQData],
) {
	ctx, span := h.tracer.Start(ctx, "handle msg")
	defer span.End()

	startTime := time.Now()

	log.Printf("accept %s %#+v\n", key, data)

	rabbitCtx, rabbitCnl := context.WithTimeout(ctx, time.Second*10)
	defer rabbitCnl()

	err := rabbitClient.Write(rabbitCtx, handlerdto.RabbitMQData{
		RequestID: key,
		UserID:    data.UserID,
		Chance:    data.Chance,
		Duration:  data.Duration,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "handle error")

		log.Println(key, err)
	}

	log.Printf("send to RabbitMQ %s\n", key)
	registerHandleTime(time.Since(startTime))
}
