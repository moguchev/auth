package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

const kid = "k2" // вы можете менять его по мере ротации ключей

var jwks = jwk.NewSet()

func init() {
	key1, err := loadPrivateKey("private_1.pem") // наш первый ключ
	if err != nil {
		log.Fatal(err)
	}
	jwkKey1, err := jwk.FromRaw(key1)
	if err != nil {
		log.Fatal(err)
	}

	// Указываем дополнительные поля
	_ = jwkKey1.Set(jwk.KeyIDKey, "k1")
	_ = jwkKey1.Set(jwk.KeyUsageKey, "sig")
	_ = jwkKey1.Set(jwk.AlgorithmKey, "RS256")
	jwks.AddKey(jwkKey1)

	key2, err := loadPrivateKey("private_2.pem") // наш второй ключ
	if err != nil {
		log.Fatal(err)
	}
	jwkKey2, err := jwk.FromRaw(key2)
	if err != nil {
		log.Fatal(err)
	}

	// Указываем дополнительные поля
	_ = jwkKey2.Set(jwk.KeyIDKey, "k2")
	_ = jwkKey2.Set(jwk.KeyUsageKey, "sig")
	_ = jwkKey2.Set(jwk.AlgorithmKey, "RS256")
	jwks.AddKey(jwkKey2)
}

// getPrivateKey - отдает приватный ключ по kid
func getPrivateKey(kid string) (crypto.PrivateKey, error) {
	k, ok := jwks.LookupKeyID(kid)
	if !ok {
		return nil, errors.New("pk not found")
	}
	var rawKey any
	if err := k.Raw(&rawKey); err != nil {
		return nil, errors.New("не удалось достать raw ключ")
	}

	return rawKey, nil
}

// getPublicKey - отдает публичный ключ по kid
func getPublicKey(kid string) (crypto.PublicKey, error) {
	k, ok := jwks.LookupKeyID(kid)
	if !ok {
		return nil, errors.New("pk not found")
	}

	pk, err := k.PublicKey()
	if err != nil {
		return nil, err
	}

	var rawKey any
	if err := pk.Raw(&rawKey); err != nil {
		return nil, errors.New("не удалось достать raw ключ")
	}

	switch key := rawKey.(type) {
	case *rsa.PublicKey:
		return key, nil
	default:
		return nil, fmt.Errorf("unsupported key type %T", rawKey)
	}
}

func jwksEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
