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
		jwks.Keyfunc,
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*FirebaseClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	if !token.Valid {
		return nil, errors.New("invalid token signature")
	}

	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}

	if !claims.VerifyAudience(ProjectID, true) {
		return nil, errors.New("invalid audience")
	}

	if !claims.VerifyExpiresAt(time.Now(), true) {
		return nil, errors.New("token expired")
	}

	if !claims.VerifyIssuedAt(time.Now(), true) {
		return nil, errors.New("invalid issued-at")
	}

	if claims.Subject == "" {
		return nil, errors.New("invalid subject")
	}

	if !claims.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return claims, nil
}
