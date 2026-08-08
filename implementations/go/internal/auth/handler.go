package auth

import (
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(authService *Service) *Handler {
	return &Handler{
		service: authService,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload RegisterPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("invalid payload: %v", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	result, err := h.service.Register(r.Context(), payload)
	if err != nil {
		log.Printf("failed to register: %v", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var payload LoginPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Printf("invalid payload: %v", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), payload)
	if err != nil {
		log.Printf("failed to login: %v", err)
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token)

	w.WriteHeader(http.StatusOK)
	res := LoginResponse{
		Success: true,
	}
	json.NewEncoder(w).Encode(res)
}
