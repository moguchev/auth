package main

import (
	"encoding/json"
	"net/http"

	oauth2 "golang.org/x/oauth2"
)

func refresh(w http.ResponseWriter, r *http.Request) {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	ts := oauth2Config.TokenSource(ctx, &oauth2.Token{RefreshToken: req.RefreshToken})
	newToken, err := ts.Token()
	if err != nil {
		http.Error(w, "Failed to refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewEncoder(w).Encode(response{
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
	})
}
