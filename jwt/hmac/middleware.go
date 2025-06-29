package main

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey = contextKey("user")

// Добавляем пользователя в контекст
func putUserToContext(ctx context.Context, u user) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

// Достаем пользователя из контекста
func getUserFromContext(ctx context.Context) (user, bool) {
	u, ok := ctx.Value(userContextKey).(user)
	return u, ok
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		const bearerPrefix = "Bearer "

		if authHeader == "" || !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Missing access token", http.StatusUnauthorized)
			return
		}

		rawToken := strings.TrimPrefix(authHeader, bearerPrefix)
		u, err := verifyAccessToken(rawToken)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := putUserToContext(r.Context(), u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
