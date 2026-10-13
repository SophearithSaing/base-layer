package api

import (
	"baselayer/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	t.Setenv("CLIENT_ORIGINS", "http://localhost:4200")

	allowedOrigins, err := config.GetClientOrigins()
	if err != nil {
		t.Fatalf("expected allowed origins, got nil")
	}

	HandleRoutes(mux)
	handler := Cors(mux, allowedOrigins)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("Origin", "https://example.com")
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	allowOriginHeader := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOriginHeader != "" {
		t.Fatalf("expected to be empty, got %v", allowOriginHeader)
	}
}

func TestCorsAllowed(t *testing.T) {
	mux := http.NewServeMux()
	t.Setenv("CLIENT_ORIGINS", "http://localhost:5173")

	allowedOrigins, err := config.GetClientOrigins()
	if err != nil {
		t.Fatalf("expected allowed origins, got nil")
	}

	HandleRoutes(mux)
	handler := Cors(mux, allowedOrigins)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	req.Header.Set("Origin", "http://localhost:5173")
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	allowOriginHeader := rec.Header().Get("Access-Control-Allow-Origin")
	if allowOriginHeader != "http://localhost:5173" {
		t.Fatalf("expected to be http://localhost:5173, got %v", allowOriginHeader)
	}
}
