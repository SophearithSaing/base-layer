package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GenericResponse struct {
	Message string `json:"message"`
}

type GenericCreatedResponse struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func JSONResponseWriter[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println(err)
	}
}
