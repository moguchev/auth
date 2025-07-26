package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /login/oauth/authorize", login)
	mux.HandleFunc("GET /oidc/callback", callback)
	mux.HandleFunc("POST /refresh", refresh)

	srv := http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	fmt.Println("Server running at http://localhost:3000")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
