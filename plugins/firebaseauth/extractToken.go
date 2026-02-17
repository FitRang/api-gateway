package firebaseauth

import (
	"errors"
	"net/http"
	"strings"
)

func extractToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}

	token := req.URL.Query().Get("access_token")
	if token != "" {
		return token, nil
	}

	return "", errors.New("token not found")
}
