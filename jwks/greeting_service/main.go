package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := RunJWKSRefresher(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// защищённый маршрут
	mux.Handle("GET /hello", authMiddleware(http.HandlerFunc(hello)))

	srv := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	fmt.Println("Server running at http://localhost:8081")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
