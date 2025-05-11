package server

import (
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/valyala/fasthttp"
)

func (s *server) List(ctx *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.ListRequest](ctx.Request.Body())
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

	rawNotifications, err := s.db.GetNotificationsByUserID(ctx, req.UserID)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		marshal(ctx, notificationclient.ErrorResponse{
			Code:    "logic",
			Details: err.Error(),
		})

		return
	}

	notifications := make([]notificationclient.NotificationData, len(rawNotifications))

	for index, raw := range rawNotifications {
		notifications[index] = notificationclient.NotificationData{
			ID:      raw.ID,
			Kind:    raw.Kind,
			Level:   raw.Level,
			Title:   raw.Title,
			Body:    raw.Body.String,
			Created: raw.Created,
		}
	}

	ctx.SetStatusCode(http.StatusOK)
	marshal(ctx, notificationclient.ListResponse{
		Notifications: notifications,
	})
}
