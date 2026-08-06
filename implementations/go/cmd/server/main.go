package main

import (
	"baselayer/internal/api"
	"baselayer/internal/config"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port, err := config.GetPort()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	allowedOrigins, err := config.GetClientOrigins()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      api.Cors(mux, allowedOrigins),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	api.HandleRoutes(mux)

	fmt.Println("Listening on port " + port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
