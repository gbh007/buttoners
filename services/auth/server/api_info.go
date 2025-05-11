package server

import (
	"net/http"

	"github.com/gbh007/buttoners/core/clients/authclient"
)

func (s *server) Info(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshal[authclient.InfoRequest](r.Body)
	if err != nil {
		w.Header().Set("Content-Type", authclient.ContentType)
		w.WriteHeader(http.StatusBadRequest)
		marshal(w, authclient.ErrorResponse{
			Code:    "unmarshal",
			Details: err.Error(),
		})

		return
	}

	if req.Token == "" {
		w.Header().Set("Content-Type", authclient.ContentType)
		w.WriteHeader(http.StatusBadRequest)
		marshal(w, authclient.ErrorResponse{
			Code:    "validate",
			Details: "empty token",
		})

		return
	}

	user, err := s.getUser(ctx, req.Token)
	if err != nil {
		w.Header().Set("Content-Type", authclient.ContentType)
		w.WriteHeader(http.StatusInternalServerError)
		marshal(w, authclient.ErrorResponse{
			Code:    "logic",
			Details: err.Error(),
		})

		return
	}

	w.Header().Set("Content-Type", authclient.ContentType)
	w.WriteHeader(http.StatusOK)
	marshal(w, authclient.InfoResponse{
		Login:  user.Login,
		UserID: user.ID,
	})
}
