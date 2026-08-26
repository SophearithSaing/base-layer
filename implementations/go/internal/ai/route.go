package ai

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("POST /ai/chat", authMiddleware(http.HandlerFunc(handler.SendMessage)))
	mux.Handle("POST /ai/extract-text", authMiddleware(http.HandlerFunc(handler.ExtractText)))
	mux.Handle("POST /ai/extract-receipt", authMiddleware(http.HandlerFunc(handler.ExtractReceipt)))
	// mux.Handle("GET /ai/extractions")
}
