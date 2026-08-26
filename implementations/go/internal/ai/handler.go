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

func (h *Handler) ExtractText(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	result, err := h.service.ExtractText(
		r.Context(),
		file,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusOK, api.GenericResponse{
		Message: result,
	})
}

func (h *Handler) ExtractReceipt(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("image")
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	defer file.Close()

	extraction, err := h.service.ExtractReceipt(
		r.Context(),
		file,
		header.Filename,
		header.Header.Get("Content-Type"),
	)
	if err != nil {
		api.JSONResponseWriter(w, http.StatusInternalServerError, api.GenericResponse{
			Message: err.Error(),
		})
		return
	}
	api.JSONResponseWriter(w, http.StatusCreated, extraction)
}
