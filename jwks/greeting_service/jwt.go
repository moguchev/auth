package main

import (
	"errors"
	"fmt"

	"github.com/lestrrat-go/jwx/v3/jwt"
)

const issuer = "best.hotel.com"

var ErrInvalidToken = errors.New("invalid token")

func verifyAccessToken(accessToken string) (user, error) {
	token, err := jwt.Parse([]byte(accessToken),
		jwt.WithKeySet(jwks),
		jwt.WithIssuer(issuer),
		jwt.WithRequiredClaim("user_email"),
		jwt.WithRequiredClaim("user_name"),
	)
	if err != nil {
		return user{}, fmt.Errorf("failed to verify JWS: %s", err)
	}

	if err := jwt.Validate(token); err != nil {
		return user{}, ErrInvalidToken
	}

	var userEmail string
	token.Get("user_email", &userEmail)
	var userName string
	token.Get("user_name", &userName)

	return user{
		Name:  userName,
		Email: userEmail,
	}, nil
}
