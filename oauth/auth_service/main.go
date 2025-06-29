package main

import (
	"errors"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /refresh", refresh)
	mux.HandleFunc("GET /.well-known/jwks.json", jwksEndpoint)

	mux.HandleFunc("GET /oauth2/authorize", authorizeLogin)
	mux.HandleFunc("POST /oauth2/authorize", authorizeCode)
	mux.HandleFunc("POST /oauth2/token", token)

	srv := http.Server{
		Addr:    ":8080",
		Handler: withCORS(mux),
	}

	log.Println("Server running at http://localhost:8080")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
