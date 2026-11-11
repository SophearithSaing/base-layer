package movie

import (
	"baselayer/internal/auth"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *handler, authMiddleware auth.AuthMiddleware) {
	mux.Handle("POST /movies/create", authMiddleware(http.HandlerFunc(handler.createMovie)))
	mux.Handle("GET /movies", authMiddleware(http.HandlerFunc(handler.listMovies)))
	mux.Handle("GET /movies/{id}", authMiddleware(http.HandlerFunc(handler.getMovieByID)))
	mux.Handle("PATCH /movies/{id}", authMiddleware(http.HandlerFunc(handler.updateMovie)))
}
