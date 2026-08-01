package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func HandleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]string{
			"name":   "BaseLayer",
			"status": "ok",
		}
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data := map[string]any{
			"ok":        true,
			"timestamp": time.Now().UTC(),
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
