package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const issuer = "best.hotel.com"

// Не храните в исходниках кода!
var secret = []byte("your-secret-string")

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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

var ErrInvalidToken = errors.New("invalid token")

// функция провайдер ключа для верификации подписи
//
// в HMAC мы проверяем подпись тем же ключом, что и подписывали
func keyFunc() jwt.Keyfunc {
	return func(_ *jwt.Token) (interface{}, error) { return secret, nil }
}

func verifyAccessToken(accessToken string) (user, error) {
	token, err := jwt.Parse(accessToken, keyFunc(),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}), // проверяем соответсвие алгоритма подписи
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
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
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
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
