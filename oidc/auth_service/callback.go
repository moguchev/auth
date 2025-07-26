package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
)

func callback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// защита от CSRF атаки
	if r.URL.Query().Get("state") != "some random state" {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	// Получаем временный authorization code, который отправил Keycloak
	code := r.URL.Query().Get("code")
	// Обмениваем authorization code на настоящие токены

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// id_token - это JWT, который содержит данные о пользователе (имя, email и т.д.)
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field", http.StatusInternalServerError)
		return
	}

	// Этот метод из библиотеки github.com/coreos/go-oidc выполняет проверку ID Token,
	// которую мы получили от Identity Provider (в нашем случае — Keycloak).
	// Что включает эта проверка:
	// - Подпись токена — проверяет, что токен действительно подписан Keycloak-ом (через JWK).
	// - Поле aud (аудитория) — проверяет, что токен предназначен именно для нашего клиента (clientID).
	// - Поле exp (время жизни) — не истёк ли токен.
	// - Поле iss (issuer) — что токен выпущен
	verifier := oidcProvider.Verifier(&oidc.Config{ClientID: os.Getenv("OIDC_CLIENT_ID")})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token", http.StatusInternalServerError)
		return
	}

	var claims map[string]any
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
		return
	}

	log.Printf("Пользователь %s авторизован", claims["name"])

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
	json.NewEncoder(w).Encode(resp)
}
