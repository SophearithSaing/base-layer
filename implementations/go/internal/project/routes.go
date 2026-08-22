package project

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("POST /projects/create", authMiddleware(http.HandlerFunc(handler.CreateProject)))
	mux.Handle("GET /projects", authMiddleware(http.HandlerFunc(handler.ListProjects)))
	mux.Handle("GET /projects/{id}", authMiddleware(http.HandlerFunc(handler.GetProjectByID)))
	mux.Handle("PATCH /projects/{id}", authMiddleware(http.HandlerFunc(handler.UpdateProject)))
	mux.Handle("POST /projects/{id}/start", authMiddleware(http.HandlerFunc(handler.StartProject)))

	mux.Handle("GET /projects/progresses", authMiddleware(http.HandlerFunc(handler.ListProgresses)))
	mux.Handle("GET /projects/progresses/{id}", authMiddleware(http.HandlerFunc(handler.GetProgressByID)))
}
