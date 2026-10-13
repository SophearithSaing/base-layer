package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootRoute(t *testing.T) {
	mux := http.NewServeMux()

	HandleRoutes(mux)

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected application/json, got %q", contentType)
	}

	var body RootResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json body: %v", err)
	}

	if body.Name != "BaseLayer" {
		t.Errorf("expected BaseLayer, got %q", body.Name)
	}

	if body.Status != "ok" {
		t.Errorf("expected ok, got %q", body.Status)
	}
}

func TestHealthRoute(t *testing.T) {
	mux := http.NewServeMux()

	HandleRoutes(mux)

	res := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, res)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Fatalf("expected application/json, got %q", contentType)
	}

	var body HealthResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode json: %v", err)
	}

	if body.Ok != true {
		t.Errorf("expected true, got %v", body.Ok)
	}

	name, _ := body.Timestamp.Zone()
	if name != "UTC" {
		t.Errorf("expected UTC time zone, got %v", name)
	}
}
