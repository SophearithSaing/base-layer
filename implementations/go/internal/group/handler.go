package group

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

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var payload CreateGroupPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	id, err := h.service.CreateGroup(r.Context(), payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusCreated, api.GenericCreatedResponse{
		ID:      id,
		Message: "group created",
	})
}

func (h *Handler) ListGroup(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.ListGroup(r.Context())
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, groups)
}

func (h *Handler) GetGroupByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	group, err := h.service.GetGroupByID(r.Context(), id)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, group)
}
