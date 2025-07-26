package main

import (
	"encoding/json"
	"net/http"
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

	user, err := verifyRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Bad refresh token", http.StatusBadRequest)
		return
	}

	newAccess, err := createAccessToken(user)
	if err != nil {
		http.Error(w, "Error creating access token", http.StatusInternalServerError)
		return
	}

	newRefresh, err := createRefreshToken(user)
	if err != nil {
		http.Error(w, "Error creating refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	json.NewEncoder(w).Encode(response{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}
