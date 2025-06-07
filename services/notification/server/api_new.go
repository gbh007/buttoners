package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/gbh007/buttoners/services/notification/internal/storage"
	"github.com/valyala/fasthttp"
)

func (s *server) New(ctx context.Context, rc *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.NewRequest](rc.Request.Body())
	if err != nil {
		rc.SetStatusCode(http.StatusBadRequest)
		marshal(rc, notificationclient.ErrorResponse{
			Code:    "unmarshal",
			Details: err.Error(),
		})

		return
	}

	if req.UserID < 1 {
		rc.SetStatusCode(http.StatusBadRequest)
		marshal(rc, notificationclient.ErrorResponse{
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
		rc.SetStatusCode(http.StatusInternalServerError)
		marshal(rc, notificationclient.ErrorResponse{
			Code:    "logic",
			Details: err.Error(),
		})

		return
	}

	rc.SetStatusCode(http.StatusNoContent)
}
