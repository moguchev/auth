package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type authCode struct {
	Code        string
	ClientID    string
	User        user
	Scope       string
	RedirectURI string
	ExpiresAt   time.Time
}

var (
	codeStore = map[string]authCode{}
	codeMx    sync.Mutex
)

func generateCode() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func authorizeCode(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Проверяем client_id
	clientID := r.Form.Get("client_id")
	client, ok := clients[clientID]
	if !ok {
		http.Error(w, "unknown_client_id", http.StatusBadRequest)
		return
	}

	// Проверяем что redirect_uri в whitelist для этого client
	redirectURI := r.Form.Get("redirect_uri")
	if !client.IsAllowedRedirectURI(redirectURI) {
		http.Error(w, "invalid_redirect_uri", http.StatusBadRequest)
		return
	}

	// Проверяем что scope в whitelist для этого client
	scope := r.Form.Get("scope")
	if !client.IsValidScope(scope) {
		http.Error(w, "invalid_scope", http.StatusBadRequest)
		return
	}

	var (
		email    = r.Form.Get("email")
		password = r.Form.Get("password")
	)
	usr, err := authUser(email, password)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// генерируем access code для обмена на токен
	code := generateCode()

	// сохраняем code с параметрами аутентификации
	codeMx.Lock()
	codeStore[code] = authCode{
		Code:        code,
		ClientID:    clientID,
		User:        usr,
		Scope:       scope,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().Add(2 * time.Minute),
	}
	codeMx.Unlock()

	// Редиректим пользователя обратно на указанный URI
	redirect := fmt.Sprintf("%s?code=%s", redirectURI, code)

	// state — это непрозрачная строка, которую клиент генерирует сам,
	// отправляет в /oauth2/authorize, а Authorization Server обязан вернуть
	// без изменений при редиректе назад.
	//
	// Защита от CSRF
	state := r.Form.Get("state")
	if state != "" {
		redirect += "&state=" + state
	}

	http.Redirect(w, r, redirect, http.StatusFound)
}
