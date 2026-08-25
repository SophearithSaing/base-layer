package ai

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("POST /ai/chat", authMiddleware(http.HandlerFunc(handler.SendMessage)))
	// mux.Handle("POST /ai/extract-text")
	// mux.Handle("POST /ai/extract-receipt")
	// mux.Handle("GET /ai/extractions")
}
