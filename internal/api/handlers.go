package api

import (
	"net/http"
	"time"
)

type RootResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type HealthResponse struct {
	Ok        bool      `json:"ok"`
	Timestamp time.Time `json:"timestamp"`
}

type BasicResponse struct {
	Message string `json:"message"`
}

func HandleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		data := RootResponse{
			Name:   "BaseLayer",
			Status: "ok",
		}
		JSONResponseWriter(w, http.StatusOK, data)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		data := HealthResponse{
			Ok:        true,
			Timestamp: time.Now().UTC(),
		}
		JSONResponseWriter(w, http.StatusOK, data)
	})

	mux.HandleFunc("/{path...}", func(w http.ResponseWriter, r *http.Request) {
		data := BasicResponse{
			Message: "Route not found",
		}
		JSONResponseWriter(w, http.StatusNotFound, data)
	})
}
