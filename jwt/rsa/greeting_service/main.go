package main

import (
	"errors"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// защищённый маршрут
	mux.Handle("GET /hello", authMiddleware(http.HandlerFunc(hello)))

	srv := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("Server running at http://localhost:8081")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
