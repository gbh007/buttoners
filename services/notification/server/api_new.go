package server

import (
	"database/sql"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
)

func (s *server) New(ctx *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.NewRequest](ctx.Request.Body())
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		marshal(ctx, notificationclient.ErrorResponse{
			Code:    "unmarshal",
			Details: err.Error(),
		})

		return
	}

	if req.UserID < 1 {
		ctx.SetStatusCode(http.StatusBadRequest)
		marshal(ctx, notificationclient.ErrorResponse{
			Code:    "validation",
			Details: "user id",
		})

		return
	}

	err = s.db.CreateNotification(ctx, &storage.Notification{
		UserID: req.UserID,
		Kind:   req.Kind,
		Level:  req.Level,
		Title:  req.Title,
		Body: sql.NullString{
			String: req.Body,
			Valid:  req.Body != "",
		},
		Created: req.Created,
	})
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		marshal(ctx, notificationclient.ErrorResponse{
			Code:    "logic",
			Details: err.Error(),
		})

		return
	}

	ctx.SetStatusCode(http.StatusNoContent)
}
