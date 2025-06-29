package main

import (
	"crypto"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// вы можете менять его по мере ротации ключей
const kid = "k2" // "k1"

type Key struct {
	id    string
	path  string
	usage string
	alg   string
}

var jwks = jwk.NewSet()

var keys = []Key{
	{
		id:    "k1",
		path:  "private_1.pem",
		usage: "sig", // для подписи
		alg:   "RS256",
	},
	{
		id:    "k2",
		path:  "private_2.pem",
		usage: "sig", // для подписи
		alg:   "RS256",
	},
}

func loadKeys(keys []Key) error {
	for _, k := range keys {
		key, err := loadPrivateKey(k.path)
		if err != nil {
			return err
		}

		jwkKey, err := jwk.Import(key)
		if err != nil {
			return err
		}

		// Указываем дополнительные поля
		_ = jwkKey.Set(jwk.KeyIDKey, k.id)       // идентификатор
		_ = jwkKey.Set(jwk.KeyUsageKey, k.usage) // назначение
		_ = jwkKey.Set(jwk.AlgorithmKey, k.alg)  // алгоритм
		jwks.AddKey(jwkKey)
	}
	return nil
}

func init() {
	if err := loadKeys(keys); err != nil {
		log.Fatal(err)
	}
}

// getPrivateKey - отдает приватный ключ по kid
func getPrivateKey(kid string) (crypto.PrivateKey, error) {
	// Ищем ключ с нужным нам идентификатором
	k, ok := jwks.LookupKeyID(kid)
	if !ok {
		return nil, fmt.Errorf("kid %s: not found", kid)
	}

	// Получаем тип ключа
	alg, ok := k.Algorithm()
	if !ok {
		return nil, fmt.Errorf("kid %s: unknown alg", kid)
	}

	sa, ok := jwa.LookupSignatureAlgorithm(alg.String())
	if !ok {
		return nil, fmt.Errorf("kid %s: unknown alg", kid)
	}

	// Отдаем ключ для подписи
	switch sa {
	case jwa.RS256():
		var rawKey rsa.PrivateKey
		if err := jwk.Export(k, &rawKey); err != nil {
			return nil, fmt.Errorf("kid %s: не удалось достать RS256 ключ: %s", kid, err)
		}
		return &rawKey, nil
	default:
		return nil, fmt.Errorf("kid %s: unsupported alg", kid)
	}
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

	alg, ok := pk.Algorithm()
	if !ok {
		return nil, fmt.Errorf("kid %s: unknown alg", kid)
	}

	sa, ok := jwa.LookupSignatureAlgorithm(alg.String())
	if !ok {
		return nil, fmt.Errorf("kid %s: unknown alg", kid)
	}

	switch sa {
	case jwa.RS256():
		var rawKey rsa.PublicKey
		if err := jwk.Export(k, &rawKey); err != nil {
			return nil, fmt.Errorf("kid %s: не удалось достать RS256 ключ: %s", kid, err)
		}
		return &rawKey, nil
	default:
		return nil, fmt.Errorf("kid %s: unsupported alg", kid)
	}
}

func jwksEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jwks)
}
