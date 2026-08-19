package project

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("GET /projects", authMiddleware(http.HandlerFunc(handler.ListProjects)))
	mux.Handle("POST /projects/create", authMiddleware(http.HandlerFunc(handler.CreateProject)))
}
