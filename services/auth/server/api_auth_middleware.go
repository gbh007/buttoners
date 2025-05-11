package server

import (
	"net/http"

	"github.com/gbh007/buttoners/core/clients/authclient"
)

func (s *server) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.Header().Set("Content-Type", authclient.ContentType)
			w.WriteHeader(http.StatusUnauthorized)
			marshal(w, authclient.ErrorResponse{
				Code:    "unauthorized",
				Details: "empty token",
			})

			return
		}

		if token != s.token {
			w.Header().Set("Content-Type", authclient.ContentType)
			w.WriteHeader(http.StatusForbidden)
			marshal(w, authclient.ErrorResponse{
				Code:    "forbidden",
				Details: "invalid token",
			})

			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
