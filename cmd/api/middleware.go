package api

import (
	"context"
	"net/http"
	"strings"
)

func getAuthMiddlewarewithJWT(jwtToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.Split(authHeader, "Bearer ")[1]
			claims, err := ValidateJWT(jwtToken, tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "username", claims.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
