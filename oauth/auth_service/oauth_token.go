package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func token(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid_request", http.StatusBadRequest)
		return
	}

	// Разрешаем только authorization_code
	if r.Form.Get("grant_type") != "authorization_code" {
		http.Error(w, "unsupported grant_type", http.StatusBadRequest)
		return
	}

	// 1) опять проверяем client_id
	clientID := r.Form.Get("client_id")
	client, ok := clients[clientID]
	if !ok {
		http.Error(w, "invalid_client", http.StatusUnauthorized)
		return
	}

	// 2) валидируем redirect_uri
	redirectURI := r.Form.Get("redirect_uri")
	if !client.IsAllowedRedirectURI(redirectURI) {
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	// 3) проверяем наличие кода и сразу его удаляем (одноразовый)
	code := r.Form.Get("code")
	codeMx.Lock()
	ac, ok := codeStore[code]
	if ok {
		delete(codeStore, code)
	}
	codeMx.Unlock()

	// 4) валидируем code
	if !ok || time.Now().After(ac.ExpiresAt) || ac.ClientID != clientID || ac.RedirectURI != redirectURI {
		http.Error(w, "invalid_grant", http.StatusBadRequest)
		return
	}

	// 5) выписываем токены
	access, expiresIn, err := createAccessToken(createAccessTokenParams{
		User:     ac.User,
		ClientID: clientID,
		Scope:    ac.Scope,
		Audience: clientID,
	})
	if err != nil {
		http.Error(w, "token error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"access_token":  access,
		"token_type":    "Bearer",
		"expires_in":    expiresIn,
		"scope":         ac.Scope,
		"refresh_token": "TODO",
	})
}
