package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
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

	// Swagger UI
	mux.Handle("/swagger/doc.json", http.HandlerFunc(swagger))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	srv := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	log.Println("Server running at http://localhost:8081")
	log.Println("Swagger running at http://localhost:8081/swagger/index.html")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
