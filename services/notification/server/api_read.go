package server

import (
	"context"
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/valyala/fasthttp"
)

func (s *Server) Read(ctx context.Context, rc *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.ReadRequest](rc.Request.Body())
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

	if req.ID > 0 {
		err = s.db.MarkReadByID(ctx, req.UserID, req.ID)
	} else {
		err = s.db.MarkReadByUserID(ctx, req.UserID)
	}

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
