package api

import (
	"net/http"
	"slices"
)

func Cors(next http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCorsHeaders(w, r, allowedOrigins)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		setCorsHeaders(w, r, allowedOrigins)
		next.ServeHTTP(w, r)
	})
}

func setCorsHeaders(w http.ResponseWriter, r *http.Request, allowedOrigins []string) {
	origin := r.Header.Get("Origin")

	if slices.Contains(allowedOrigins, origin) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
