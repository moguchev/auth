package main

import "context"

type user struct {
	Name   string
	Email  string
	Scopes []string
}

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
