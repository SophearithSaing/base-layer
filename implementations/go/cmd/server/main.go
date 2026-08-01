package main

import (
	"baselayer/internal/api"
	"baselayer/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	portValue := config.GetEnv("PORT", "8000")
	port, err := strconv.ParseInt(portValue, 10, 32)
	if port < 1 || port > 65535 {
		fmt.Println("invalid port " + portValue)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + portValue,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	api.HandleRoutes(mux)

	fmt.Println("Listening on port " + portValue)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
