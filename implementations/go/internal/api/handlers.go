package api

import (
	"encoding/json"
	"fmt"
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

func HandleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := RootResponse{
			Name:   "BaseLayer",
			Status: "ok",
		}
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := HealthResponse{
			Ok:        true,
			Timestamp: time.Now().UTC(),
		}
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("/{path...}", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "route not found", http.StatusNotFound)
	})
}
