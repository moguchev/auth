package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /refresh", refresh)
	// защищённый маршрут
	mux.Handle("GET /hello", authMiddleware(http.HandlerFunc(hello)))

	srv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Server running at http://localhost:8080")
	srv.ListenAndServe()
}
