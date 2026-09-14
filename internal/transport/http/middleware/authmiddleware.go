package middleware

import (
	"context"
	"net/http"

	"github.com/PashakArt/file-server/internal/service"
	"github.com/PashakArt/file-server/internal/transport/http/types"
)

const TOKEN_HEADER_NAME = "x-token"

type contextKey string

const userIDContextKey contextKey = "userID"

func AuthMiddleware(authService *service.AuthService) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get(TOKEN_HEADER_NAME)
			if token == "" {
				types.SendError(w, r, http.StatusUnauthorized, "header X-Token is required")
				return
			}

			userID, err := authService.ValidateToken(r.Context(), token)
			if err != nil {
				types.SendError(w, r, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDContextKey).(string)
	return userID, ok
}
