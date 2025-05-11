package server

import (
	"net/http"
	"strings"

	"github.com/gbh007/buttoners/core/clients/authclient"
)

func (s *server) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshal[authclient.RegisterRequest](r.Body)
	if err != nil {
		w.Header().Set("Content-Type", authclient.ContentType)
		w.WriteHeader(http.StatusBadRequest)
		marshal(w, authclient.ErrorResponse{
			Code:    "unmarshal",
			Details: err.Error(),
		})

		return
	}

	if req.Login == "" || req.Password == "" {
		w.Header().Set("Content-Type", authclient.ContentType)
		w.WriteHeader(http.StatusBadRequest)
		marshal(w, authclient.ErrorResponse{
			Code:    "validate",
			Details: "empty login or password",
		})

		return
	}

	login := strings.ToLower(req.Login)
	pass := req.Password

	_, err = s.createUser(ctx, login, pass)
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
	w.WriteHeader(http.StatusNoContent)
}
