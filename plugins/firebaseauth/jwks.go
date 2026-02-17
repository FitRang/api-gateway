package firebaseauth

import (
	"time"

	"github.com/MicahParks/keyfunc"
)

var jwks *keyfunc.JWKS

func initJWKS() error {
	options := keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshTimeout:    10 * time.Second,
		RefreshUnknownKID: true,
	}

	var err error
	jwks, err = keyfunc.Get(
		"https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com",
		options,
	)

	return err
}
