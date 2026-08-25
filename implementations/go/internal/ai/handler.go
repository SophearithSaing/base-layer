package ai

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

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var payload MessagePayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
	}
	response, err := h.service.SendMessage(r.Context(), payload.Message)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
	}
	api.JSONResponseWriter(w, http.StatusOK, api.GenericResponse{
		Message: response,
	})
}
