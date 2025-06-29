package main

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	// state - это произвольная строка, которую мы генерируем при редиректе пользователя на OIDC-провайдер.
	// Она потом вернётся нам в /callback вместе с code.

	// Зачем? - Защита от CSRF-атак:
	// если злоумышленник попытается подделать запрос, у него не будет правильного state.
	const state = "some random state"
	redirectURL := oauth2Config.AuthCodeURL(state)

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
