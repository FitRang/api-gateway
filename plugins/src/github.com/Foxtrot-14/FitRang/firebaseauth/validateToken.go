package firebaseauth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type FirebaseClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}

const (
	ProjectID = "fitrang-6c0aa"
	Issuer    = "https://securetoken.google.com/" + ProjectID
)

func validateToken(tokenString string) (*FirebaseClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&FirebaseClaims{},
		keyFunc,
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*FirebaseClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	rc := claims.RegisteredClaims

	// if rc.Issuer != Issuer {
	// 	return nil, errors.New("invalid issuer")
	// }
	//
	// if !rc.VerifyAudience(ProjectID, true) {
	// 	return nil, errors.New("invalid audience")
	// }

	if rc.ExpiresAt == nil || time.Now().After(rc.ExpiresAt.Time) {
		return nil, errors.New("token expired")
	}

	if rc.Subject == "" {
		return nil, errors.New("invalid subject")
	}

	return claims, nil
}
