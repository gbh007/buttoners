package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
)

type server struct {
	db    *storage.Database
	token string
}

func (s *server) handle(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType(notificationclient.ContentType)

	if !ctx.IsPost() {
		ctx.SetStatusCode(http.StatusNotFound)
		marshal(ctx, notificationclient.ErrorResponse{
			Code:    "not found",
			Details: string(ctx.Method()),
		})

		return
	}

	// FIXME: поддержать телеметрию, использовать более быстрые библиотеки для json

	p := string(ctx.Path())

	switch p {
	case notificationclient.NewPath:
		s.New(ctx)
	case notificationclient.ListPath:
		s.List(ctx)
	case notificationclient.ReadPath:
		s.Read(ctx)
	default:
		ctx.SetStatusCode(http.StatusNotFound)
		marshal(ctx, notificationclient.ErrorResponse{
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
