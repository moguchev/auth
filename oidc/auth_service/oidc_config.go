package main

import (
	"context"
	"log"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	oauth2 "golang.org/x/oauth2"
)

var (
	// Конфиг OAuth2.0 для клиента, который отправляет запросы на авторизацию и получение токенов.
	oauth2Config *oauth2.Config
	// OpenID-провайдера (в нашем случае — Keycloak)
	oidcProvider *oidc.Provider
)

func init() {
	var err error
	const realm = "http://localhost:8080/realms/example"
	oidcProvider, err = oidc.NewProvider(context.Background(), realm)
	if err != nil {
		log.Fatal(err)
	}

	oauth2Config = &oauth2.Config{
		ClientID:     os.Getenv("OIDC_CLIENT_ID"),
		ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  "http://localhost:3000/oidc/callback", // Адрес, куда Keycloak вернёт пользователя после логина.
		Scopes: []string{
			oidc.ScopeOpenID,        // обязательно для OIDC
			oidc.ScopeOfflineAccess, // чтобы получить Refresh Token
			"profile", "email",      // чтобы получить имя и email в ID Token.
		},
	}
}
