package ai

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	// mux.Handle("POST /ai/chat")
	// mux.Handle("POST /ai/extract-text")
	// mux.Handle("POST /ai/extract-receipt")
	// mux.Handle("GET /ai/extractions")
}
