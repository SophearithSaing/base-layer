package group

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("POST /groups/create", authMiddleware(http.HandlerFunc(handler.CreateGroup)))
	mux.Handle("GET /groups", authMiddleware(http.HandlerFunc(handler.ListGroup)))
	// mux.Handle("GET /groups/{id}")
	// mux.Handle("PATCH /groups/{id}")
}
