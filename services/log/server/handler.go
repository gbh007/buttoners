package server

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/gbh007/buttoners/core/dto"
	"github.com/gbh007/buttoners/core/logger"
	"github.com/gbh007/buttoners/services/log/internal/storage"
	"go.opentelemetry.io/otel/codes"
)

func (s *Server) handle(ctx context.Context, key string, data dto.KafkaLogData) error {
	ctx, span := s.tracer.Start(ctx, "handle msg")
	defer span.End()

	startTime := time.Now()

	dbCtx, dbCnl := context.WithTimeout(ctx, time.Second*5)
	defer dbCnl()

	err := s.db.InsertUserLog(dbCtx, &storage.UserLog{
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

		logger.LogWithMeta(s.l, ctx, slog.LevelError, "add user log", "error", err.Error(), "msg_key", key)
	}

	registerHandleTime(time.Since(startTime))

	return err
}
