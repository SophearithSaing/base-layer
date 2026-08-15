package auth

import (
	"baselayer/internal/api"
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

	result, token, refreshToken, err := h.service.Register(r.Context(), payload)
	if err != nil {
		log.Printf("failed to register: %v", err)
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token, refreshToken)
	api.JSONResponseWriter(w, http.StatusCreated, result)
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

	token, refreshToken, err := h.service.Login(r.Context(), payload)
	if err != nil {
		log.Printf("failed to login: %v", err)
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, token, refreshToken)
	api.JSONResponseWriter(w, http.StatusOK, LoginResponse{
		Success: true,
	})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("token not found: %v", err)
		http.Error(w, "token not found", http.StatusInternalServerError)
		return
	}
	accessToken, refreshToken, err := h.service.Refresh(r.Context(), cookie.Value)
	if err != nil {
		log.Printf("failed to refresh: %v", err)
		http.Error(w, "failed to refresh", http.StatusInternalServerError)
		return
	}

	setAuthCookie(w, accessToken, refreshToken)
	api.JSONResponseWriter(w, http.StatusOK, RefreshResponse{
		Success: true,
	})
}
