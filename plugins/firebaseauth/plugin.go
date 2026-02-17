package firebaseauth

import (
	"context"
	"fmt"
	"net/http"
)

// Config holds the plugin configuration.
type Config struct {
	HeaderName  string `json:"headerName,omitempty"`
	HeaderValue string `json:"headerValue,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		HeaderName:  "X-Custom-Header",
		HeaderValue: "default-value",
	}
}

// Plugin is the plugin struct.
type Plugin struct {
	next        http.Handler
	name        string
	headerName  string
	headerValue string
}

// New creates a new plugin instance.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	err := initJWKS()
	if err != nil {
		return nil, err
	}

	if len(config.HeaderName) == 0 {
		return nil, fmt.Errorf("headerName cannot be empty")
	}

	return &Plugin{
		next:        next,
		name:        name,
		headerName:  config.HeaderName,
		headerValue: config.HeaderValue,
	}, nil
}

func (p *Plugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	tokenString, err := extractToken(req)
	if err != nil {
		http.Error(rw, "Unauthorized: token missing", http.StatusUnauthorized)
		return
	}

	claims, err := validateToken(tokenString)
	if err != nil {
		http.Error(rw, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	req.Header.Set("X-User-Email", claims.Email)
	req.Header.Set("X-Email-Verified", "true")

	p.next.ServeHTTP(rw, req)
}
