package firebaseauth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type FirebaseClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	jwt.RegisteredClaims
}

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
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if !claims.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return claims, nil
}
