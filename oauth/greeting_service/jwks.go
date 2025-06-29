package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/lestrrat-go/httprc/v3"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

const certs = `http://localhost:8080/.well-known/jwks.json`

var jwks jwk.Set

func RunJWKSRefresher(ctx context.Context) error {
	cache, err := jwk.NewCache(ctx, httprc.NewClient())
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	if err := cache.Register(ctx, certs); err != nil {
		return fmt.Errorf("failed to register google JWKS: %w", err)
	}

	init := make(chan struct{})
	once := sync.Once{}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			keyset, err := cache.Lookup(ctx, certs)
			if err != nil {
				log.Printf("failed to fetch google JWKS: %s\n", err)
				return
			}

			jwks = keyset

			once.Do(func() { close(init) })
			time.Sleep(time.Minute)
		}
	}()

	<-init
	return nil
}
