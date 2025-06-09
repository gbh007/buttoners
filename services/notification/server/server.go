package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/core/metrics"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type server struct {
	db      *storage.Database
	token   string
	tracer  trace.Tracer
	logger  *slog.Logger
	metrics *metrics.HTTPServerMetrics
}

func (s *server) handle(rc *fasthttp.RequestCtx) {
	tStart := time.Now()

	s.metrics.IncActive(string(rc.Request.Host()), string(rc.Request.URI().Path()), string(rc.Request.Header.Method()))
	defer s.metrics.DecActive(string(rc.Request.Host()), string(rc.Request.URI().Path()), string(rc.Request.Header.Method()))

	defer func() {
		s.metrics.AddHandle(string(rc.Request.Host()), string(rc.Request.URI().Path()), string(rc.Request.Header.Method()), rc.Response.StatusCode(), time.Since(tStart))
	}()

	var ctx context.Context = rc

	// Распространение трассировки
	ctx = otel.GetTextMapPropagator().Extract(ctx, &headerWraper{
		raw: &rc.Request.Header,
	})

	ctx, span := s.tracer.Start(ctx, "notification-server:"+string(rc.Path()))
	defer span.End()

	defer logData(ctx, s.logger, &rc.Request, &rc.Response)

	rc.SetContentType(notificationclient.ContentType)

	if !rc.IsPost() {
		rc.SetStatusCode(http.StatusNotFound)
		marshal(rc, notificationclient.ErrorResponse{
			Code:    "not found",
			Details: string(rc.Method()),
		})

		return
	}

	// FIXME: использовать более быстрые библиотеки для json

	p := string(rc.Path())

	switch p {
	case notificationclient.NewPath:
		s.New(ctx, rc)
	case notificationclient.ListPath:
		s.List(ctx, rc)
	case notificationclient.ReadPath:
		s.Read(ctx, rc)
	default:
		rc.SetStatusCode(http.StatusNotFound)
		marshal(rc, notificationclient.ErrorResponse{
			Code:    "not found",
			Details: p,
		})
	}
}

func marshal[T any](w io.Writer, v T) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return err
	}

	return nil
}

func unmarshal[T any](data []byte) (T, error) {
	var v T

	err := json.Unmarshal(data, &v)
	if err != nil {
		return v, err
	}

	return v, nil
}
