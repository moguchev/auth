package main

import (
	"html/template"
	"net/http"
)

var loginTmpl = template.Must(template.New("login").Parse(`
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>OAuth Login</title>
</head>
<body>
  <h3>OAuth Login</h3>
  <form method="post" action="/oauth2/authorize">
    <input type="hidden" name="client_id" value="{{.ClientID}}">
    <input type="hidden" name="redirect_uri" value="{{.RedirectURI}}">
    <input type="hidden" name="scope" value="{{.Scope}}">
    <input type="hidden" name="state" value="{{.State}}">
    <label>Email:    <input name="email" autocomplete="username"></label><br/>
    <label>Password: <input type="password" name="password" autocomplete="current-password"></label><br/>
    <button type="submit">Login</button>
  </form>
</body>
</html>
`))

type loginViewModel struct {
	ClientID    string
	RedirectURI string
	Scope       string
	State       string
}

func authorizeLogin(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Проверяем client_id
	clientID := q.Get("client_id")
	client, ok := clients[clientID]
	if !ok {
		http.Error(w, "unknown_client_id", http.StatusBadRequest)
		return
	}

	// Проверяем что redirect_uri в whitelist для этого client
	redirectURI := q.Get("redirect_uri")
	if !client.IsAllowedRedirectURI(redirectURI) {
		http.Error(w, "invalid_redirect_uri", http.StatusBadRequest)
		return
	}

	// Проверяем что scope в whitelist для этого client
	scope := q.Get("scope")
	if !client.IsValidScope(scope) {
		http.Error(w, "invalid_scope", http.StatusBadRequest)
		return
	}

	// state — это строка, которую клиент генерирует сам, отправляет в /oauth2/authorize,
	// а Authorization Server обязан вернуть без изменений при редиректе назад.
	// Защита от CSRF
	state := q.Get("state")

	vm := loginViewModel{
		ClientID:    clientID,    // добавляем в форму client_id
		RedirectURI: redirectURI, // добавляем в форму redirect_uri
		Scope:       scope,       // добавляем в форму scope
		State:       state,       // добавляем в форму state
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := loginTmpl.Execute(w, vm); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
}
