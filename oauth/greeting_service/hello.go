package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
)

func hello(w http.ResponseWriter, r *http.Request) {
	user, ok := getUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	helloAllowedScopes := func(s string) bool { return s == "read" || s == "read:hello" }
	hasScope := slices.ContainsFunc(user.Scopes, helloAllowedScopes)
	if !hasScope {
		http.Error(w, "invalid scope", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	type response struct {
		Message string `json:"message"`
	}
	json.NewEncoder(w).Encode(response{Message: fmt.Sprintf("Hello %s!", user.Name)})
}
