package server

import (
	"net/http"

	"github.com/gbh007/buttoners/core/clients/notificationclient"
	"github.com/valyala/fasthttp"
)

func (s *server) Read(ctx *fasthttp.RequestCtx) {
	req, err := unmarshal[notificationclient.ReadRequest](ctx.Request.Body())
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

	if req.ID > 0 {
		err = s.db.MarkReadByID(ctx, req.UserID, req.ID)
	} else {
		err = s.db.MarkReadByUserID(ctx, req.UserID)
	}

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
