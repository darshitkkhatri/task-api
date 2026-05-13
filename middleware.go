package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

// contextKey is a custom type for context keys — avoids collisions
type contextKey string

const userContextKey contextKey = "user"

// loggerMiddleware logs every request
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf("%s %s → %d (%s)", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

// authMiddleware validates JWT token and puts claims in context
func authMiddleware(next http.Handler, cfg *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			unauthorizedError("missing authorization header").respond(w)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorizedError("invalid authorization format").respond(w)
			return
		}

		claims, err := ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			unauthorizedError("invalid or expired token").respond(w)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts claims from request context
func GetUserFromContext(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(userContextKey).(*Claims)
	return claims, ok
}
