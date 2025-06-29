package main

import "strings"

type OAuthClient struct {
	id            string
	redirectURIs  map[string]bool
	allowedScopes map[string]bool
}

func (c *OAuthClient) IsAllowedRedirectURI(uri string) bool {
	return c.redirectURIs[uri]
}

func (c *OAuthClient) IsValidScope(scope string) bool {
	scopes := strings.Fields(scope)
	for _, s := range scopes {
		if !c.allowedScopes[s] {
			return false
		}
	}

	return true
}

var clients = map[string]OAuthClient{
	"greeting_service": {
		id: "greeting_service",
		redirectURIs: map[string]bool{
			"http://localhost:8081/swagger/oauth2-redirect.html": true,
		},
		allowedScopes: map[string]bool{
			"read":       true,
			"read:hello": true,
		},
	},
}
