package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserIDKey contextKey = "userID"
	ContextUsernameKey contextKey = "username"
)

func JWTMiddleware(jwtManager *JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			claims, err := jwtManager.Verify(token)
			if err != nil {
				http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Inject claims into context
			ctx := context.WithValue(r.Context(), ContextUserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ContextUsernameKey, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
