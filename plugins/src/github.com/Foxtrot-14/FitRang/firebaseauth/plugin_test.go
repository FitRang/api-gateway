package firebaseauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const TEST_TOKEN = ""

func mockNextHandler(t *testing.T) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		email := r.Header.Get("X-User-Email")

		if email == "" {
			t.Error("X-User-Email header not set")
		}

		w.WriteHeader(http.StatusOK)
	})
}

func TestPlugin_HTTP(t *testing.T) {
	config := CreateConfig()

	handler, err := New(context.Background(), mockNextHandler(t), config, "firebaseauth")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	req := httptest.NewRequest("GET", "http://localhost/graphql", nil)

	req.Header.Set("Authorization", "Bearer "+TEST_TOKEN)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}
}

func TestPlugin_WebSocket(t *testing.T) {
	config := CreateConfig()

	handler, err := New(context.Background(), mockNextHandler(t), config, "firebaseauth")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	req := httptest.NewRequest(
		"GET",
		"http://localhost/ws?access_token="+TEST_TOKEN,
		nil,
	)

	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Connection", "Upgrade")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}
}

func TestPlugin_InvalidToken(t *testing.T) {
	config := CreateConfig()

	handler, err := New(context.Background(), mockNextHandler(t), config, "firebaseauth")
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	req := httptest.NewRequest("GET", "http://localhost/graphql", nil)

	req.Header.Set("Authorization", "Bearer invalid.token.here")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code == http.StatusOK {
		t.Fatal("Expected failure for invalid token")
	}
}
