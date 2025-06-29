package main

import (
	_ "embed"
	"net/http"
)

//go:embed swagger/doc.json
var swaggerJSON []byte

func swagger(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(swaggerJSON)
}
