package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

const issuer = "best.hotel.com"

var publicKey *rsa.PublicKey

func init() {
	key, err := loadPublicKey("public.pem")
	if err != nil {
		log.Fatal(err)
	}
	publicKey = key
}

var ErrInvalidToken = errors.New("invalid token")

// функция провайдер ключа для верификации подписи
//
// в RSA мы проверяем подпись публичным
func keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (interface{}, error) { return publicKey, nil }
}

func verifyAccessToken(accessToken string) (user, error) {
	token, err := jwt.Parse(accessToken, keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}), // проверяем соответсвие алгоритма подписи
		jwt.WithIssuer(issuer),       // проверяем соответствие автора токена
		jwt.WithExpirationRequired(), // проверяем время жизни токена
	)
	if err != nil {
		return user{}, fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return user{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return user{}, ErrInvalidToken
	}

	userEmail, _ := claims["user_email"].(string)
	userName, _ := claims["user_name"].(string)
	return user{
		Name:  userName,
		Email: userEmail,
	}, nil
}
