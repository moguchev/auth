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

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running at http://localhost:8080")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
