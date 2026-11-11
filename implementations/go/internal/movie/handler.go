package movie

import (
	"baselayer/internal/api"
	"encoding/json"
	"net/http"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{service: service}
}

func (h *handler) createMovie(w http.ResponseWriter, r *http.Request) {
	var payload createMoviePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	id, err := h.service.createMovie(r.Context(), payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusCreated, api.GenericCreatedResponse{
		ID:      id,
		Message: "movie created",
	})
}

func (h *handler) listMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.listMovies(r.Context())
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, movies)
}

func (h *handler) getMovieByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	movie, err := h.service.getMovieByID(r.Context(), id)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, movie)
}

func (h *handler) updateMovie(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var payload updateMoviePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	movie, err := h.service.updateMovie(r.Context(), id, payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, movie)
}
