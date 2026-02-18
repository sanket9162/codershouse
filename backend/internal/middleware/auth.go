package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/sanket9162/codershouse/internal/utils"
)

type userKey string

const UserKey userKey = "user_id"

func (m *Middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := m.App.Auth.GetTokenFromHeader(r)
		if err != nil {
			utils.ErrorJSON(w, errors.New("unauthorized - no token"), http.StatusUnauthorized)
			return
		}

		claims, err := m.App.Auth.ValidateToken(token)
		if err != nil {
			utils.ErrorJSON(w, errors.New("unauthorized - invalid token"), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, claims.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
