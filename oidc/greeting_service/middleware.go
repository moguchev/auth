package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

var oidcProvider *oidc.Provider

func init() {
	var err error
	oidcProvider, err = oidc.NewProvider(context.Background(), "http://localhost:8080/realms/example")
	if err != nil {
		log.Fatal(err)
	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, "Missing access token", http.StatusUnauthorized)
			return
		}
		accessToken := strings.TrimPrefix(authHeader, bearerPrefix)

		verifier := oidcProvider.Verifier(&oidc.Config{
			ClientID:          os.Getenv("OIDC_CLIENT_ID"),
			SkipClientIDCheck: true, // Access tokens don't require client ID
		})

		token, err := verifier.Verify(ctx, accessToken)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		if err := token.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
			return
		}

		user := user{
			Name:  claims.Name,
			Email: claims.Email,
		}
		next.ServeHTTP(w, r.WithContext(putUserToContext(ctx, user)))
	})
}
