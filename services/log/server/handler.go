package server

import (
	"context"
	"database/sql"
	"log"
	"time"

	gatedto "github.com/gbh007/buttoners/services/gate/dto"

	"github.com/gbh007/buttoners/core/kafka"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	kafka *kafka.Client

	db *storage.Database

	tracer trace.Tracer
}

func (h *handler) Run(ctx context.Context) error {
label1:
	for {
		data := new(gatedto.KafkaLogData)
		ctx, key, err := h.kafka.Read(ctx, data)
		if err != nil {
			log.Println(err.Error())

			select {
			case <-ctx.Done():
				break label1
			default:
				continue
			}
		}

		h.handle(ctx, key, data)
	}

	return nil
}

func (h *handler) handle(ctx context.Context, key string, data *gatedto.KafkaLogData) {
	ctx, span := h.tracer.Start(ctx, "handle msg")
	defer span.End()

	startTime := time.Now()

	log.Printf("accept %s %#+v\n", key, data)

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err := h.db.InsertUserLog(dbCtx, &storage.UserLog{
		RequestID: key,
		Addr:      data.Addr,
		UserID: sql.NullInt64{
			Int64: data.UserID,
			Valid: data.UserID != 0,
		},
		SessionToken: sql.NullString{
			String: data.SessionToken,
			Valid:  data.SessionToken != "",
		},
		Action: data.Action,
		Chance: sql.NullInt64{
			Int64: data.Chance,
			Valid: data.Chance != 0,
		},
		Duration: sql.NullInt64{
			Int64: data.Duration,
			Valid: data.Duration != 0,
		},
		RequestTime: data.RequestTime,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "handle error")

		log.Println(key, err)
	}

	registerHandleTime(time.Since(startTime))
}
