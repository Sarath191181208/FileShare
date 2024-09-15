package api

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"sarath/backend_project/internal/jwt"
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
			ctx = context.WithValue(ctx, "id", claims.Id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getLogginMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Log the request details before passing to the next handler
			logger.Printf("Incoming request: Method=%s, URL=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)

			// Call the next handler
			next.ServeHTTP(w, r)

			// Log after the request is completed
			duration := time.Since(start)
			logger.Printf("Completed request: Method=%s, URL=%s, Duration=%v", r.Method, r.URL.Path, duration)
		})
	}
}
