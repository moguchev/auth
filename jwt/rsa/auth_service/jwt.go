package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "best.hotel.com"

var privateKey *rsa.PrivateKey

func init() {
	key, err := loadPrivateKey("private.pem")
	if err != nil {
		log.Fatal(err)
	}
	privateKey = key
}

func createAccessToken(u user) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		// стандартные JWT claims
		"iss": issuer,                           // кто выдал токен
		"sub": u.Email,                          // кому выдан токен
		"iat": now.Unix(),                       // время создания токена
		"exp": now.Add(15 * time.Minute).Unix(), // время жизни токена
		// наши произвольные claims
		"user_email": u.Email,
		"user_name":  u.Name,
	}
	/* ... */

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

var ErrInvalidToken = errors.New("invalid token")

// функция провайдер ключа для верификации подписи
//
// в RSA мы проверяем подпись публичном ключом
func keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (interface{}, error) { return privateKey.PublicKey, nil }
}

var (
	refreshTokens = make(map[string]struct{})
	mx            sync.RWMutex
)

func createRefreshToken(u user) (string, error) {
	tokenID := uuid.New().String() // уникальный ID для refresh токена

	now := time.Now()
	claims := jwt.MapClaims{
		// стандартные JWT claims
		"iss": issuer,                             // кто выдал токен
		"sub": u.Email,                            // кому выдан токен
		"iat": now.Unix(),                         // время создания токена
		"exp": now.Add(7 * 24 * time.Hour).Unix(), // время жизни токена
		"jti": tokenID,                            // "JWT ID" — идентификатор токена
		// наши произвольные claims
		"type": "refresh",
	}

	/* ... */

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	// храним токен
	mx.Lock()
	refreshTokens[tokenID] = struct{}{}
	mx.Unlock()

	return signed, nil
}

func verifyRefreshToken(refreshToken string) (user, error) {
	token, err := jwt.Parse(refreshToken, keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name}),
		jwt.WithIssuer(issuer),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return user{}, fmt.Errorf("parse token failed: %w", err)
	}

	if !token.Valid {
		return user{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		return user{}, ErrInvalidToken
	}

	tokenID, ok := claims["jti"].(string)
	if !ok {
		return user{}, ErrInvalidToken
	}

	email, ok := claims["sub"].(string)
	if !ok {
		return user{}, ErrInvalidToken
	}

	mx.RLock()
	_, exists := refreshTokens[tokenID]
	mx.RUnlock()

	if !exists {
		return user{}, ErrInvalidToken
	}

	mx.Lock()
	delete(refreshTokens, tokenID) // больше нельзя использовать этот refresh
	mx.Unlock()

	idx := slices.IndexFunc(usersDB, func(u user) bool { return strings.EqualFold(email, u.Email) })
	if idx == -1 {
		return user{}, ErrInvalidToken
	}

	return usersDB[idx], nil
}
