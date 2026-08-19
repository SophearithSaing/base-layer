package project

import (
	"baselayer/internal/api"
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.ListProjects(r.Context())
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
	}
	api.JSONResponseWriter(w, http.StatusOK, projects)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var payload Project
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	err = h.service.CreateProject(r.Context(), payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusCreated, api.GenericResponse{
		Message: "project created",
	})
}

func (h *Handler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := h.service.GetProjectByID(r.Context(), id)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, project)
}
