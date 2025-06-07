package server

import (
	"context"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/valyala/fasthttp"
)

func (s *server) List(ctx context.Context, rc *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.ListRequest](rc.Request.Body())
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

	rawNotifications, err := s.db.GetNotificationsByUserID(ctx, req.UserID)
	if err != nil {
		rc.SetStatusCode(http.StatusInternalServerError)
		marshal(rc, notificationclient.ErrorResponse{
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

	rc.SetStatusCode(http.StatusOK)
	marshal(rc, notificationclient.ListResponse{
		Notifications: notifications,
	})
}
