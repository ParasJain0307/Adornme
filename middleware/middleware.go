package middleware

import (
	"Adornme/Auth"
	"context"
	"net/http"
	"strings"
)

type key string

const UserIDKey key = "userID"

// JWTMiddleware validates the access token from the Authorization header
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
			return
		}

		userID, err := auth.ValidateAccessToken(parts[1])
		if err != nil {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Store userID in context for downstream handlers
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
