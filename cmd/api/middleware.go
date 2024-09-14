package api

import (
	"context"
	"net/http"
	"sarath/backend_project/internal/jwt"
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
			claims, err := jwt.ValidateJWT(jwtToken, tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "email", claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
