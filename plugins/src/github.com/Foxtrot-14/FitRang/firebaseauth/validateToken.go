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
	token, err := jwt.Parse(tokenString, keyFunc)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	mapClaims := token.Claims.(jwt.MapClaims)

	now := time.Now()

	claims := &FirebaseClaims{
		Email:         getString(mapClaims, "email"),
		EmailVerified: getBool(mapClaims, "email_verified"),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   getString(mapClaims, "iss"),
			Subject:  getString(mapClaims, "sub"),
			Audience: jwt.ClaimStrings{getString(mapClaims, "aud")},
			ExpiresAt: jwt.NewNumericDate(
				time.Unix(int64(getFloat(mapClaims, "exp")), 0),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Unix(int64(getFloat(mapClaims, "iat")), 0),
			),
		},
	}

	if claims.Issuer != Issuer {
		return nil, errors.New("invalid issuer")
	}

	if !claims.VerifyAudience(ProjectID, true) {
		return nil, errors.New("invalid audience")
	}

	if claims.ExpiresAt.Before(now) {
		return nil, errors.New("token expired")
	}

	if claims.Subject == "" {
		return nil, errors.New("invalid subject")
	}

	if !claims.EmailVerified {
		return nil, errors.New("email not verified")
	}

	return claims, nil
}

func getString(m jwt.MapClaims, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getBool(m jwt.MapClaims, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}

func getFloat(m jwt.MapClaims, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}
