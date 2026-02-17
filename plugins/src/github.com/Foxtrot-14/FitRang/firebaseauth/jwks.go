package firebaseauth

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const jwksURL = "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

var (
	publicKeys  map[string]*rsa.PublicKey
	keysMutex   sync.RWMutex
	lastRefresh time.Time
	refreshTTL  = time.Hour
)

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
	Kty string `json:"kty"`
}

func initJWKS() error {
	keys, err := fetchJWKS()
	if err != nil {
		return err
	}

	keysMutex.Lock()
	publicKeys = keys
	lastRefresh = time.Now()
	keysMutex.Unlock()

	go backgroundRefresh()

	return nil
}

func backgroundRefresh() {
	ticker := time.NewTicker(refreshTTL)
	defer ticker.Stop()

	for range ticker.C {
		keys, err := fetchJWKS()
		if err != nil {
			continue
		}

		keysMutex.Lock()
		publicKeys = keys
		lastRefresh = time.Now()
		keysMutex.Unlock()
	}
}

func fetchJWKS() (map[string]*rsa.PublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", jwksURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks jwksResponse

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	keys := make(map[string]*rsa.PublicKey)

	for _, key := range jwks.Keys {
		pubKey, err := buildPublicKey(key)
		if err != nil {
			continue
		}

		keys[key.Kid] = pubKey
	}

	return keys, nil
}

func buildPublicKey(jwk jwk) (*rsa.PublicKey, error) {
	nBytes, err := jwt.DecodeSegment(jwk.N)
	if err != nil {
		return nil, err
	}

	eBytes, err := jwt.DecodeSegment(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)

	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

func keyFunc(token *jwt.Token) (any, error) {
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid")
	}

	keysMutex.RLock()
	key, exists := publicKeys[kid]
	keysMutex.RUnlock()

	if !exists {
		return nil, errors.New("unknown kid")
	}

	return key, nil
}
