package main

import (
	"baselayer/internal/api"
	"baselayer/internal/config"
	"baselayer/internal/db"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	port, err := config.GetPort()
	if err != nil {
		return fmt.Errorf("error getting port: %w", err)
	}

	allowedOrigins, err := config.GetClientOrigins()
	if err != nil {
		return fmt.Errorf("error getting allowed origins: %w", err)
	}

	mongodbConfig, err := config.GetMongoDBConfig()
	if err != nil {
		return fmt.Errorf("error getting mongodb config: %w", err)
	}

	mongo, err := db.NewMongo(mongodbConfig)
	if err != nil {
		return fmt.Errorf("error initializing mongodb: %w", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongo.Close(ctx); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()
	log.Printf("db connected: %v", mongo.DB.Name())

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      api.Cors(mux, allowedOrigins),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	api.HandleRoutes(mux)

	log.Printf("listening on port: %v", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}
