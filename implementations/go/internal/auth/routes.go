package auth

import (
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /auth/register", handler.Register)
	mux.HandleFunc("POST /auth/login", handler.Login)
	mux.HandleFunc("POST /auth/refresh", handler.Refresh)
	mux.HandleFunc("POST /auth/logout", handler.Logout)
	mux.HandleFunc("GET /auth/me", handler.Me)
}
