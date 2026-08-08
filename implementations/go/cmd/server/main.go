package main

import (
	"baselayer/internal/api"
	"baselayer/internal/auth"
	"baselayer/internal/config"
	"baselayer/internal/db"
	"baselayer/internal/user"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Printf("application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

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

	jwtSecret, err := config.GetJWTSecret()
	if err != nil {
		return fmt.Errorf("error getting jwt secret: %w", err)
	}

	// User
	userRepo := user.NewRepository(mongo.DB)
	userService := user.NewService(userRepo)

	// Auth
	refreshTokenRepo := auth.NewRefreshTokenRepository(mongo.DB)
	authService := auth.NewService(refreshTokenRepo, userService, jwtSecret)
	authHandler := auth.NewHandler(authService)

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      api.Cors(mux, allowedOrigins),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	api.HandleRoutes(mux)
	auth.RegisterRoutes(mux, authHandler)

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("listening on port: %v", port)
		serverErr <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server failed: %w", err)
		}
	case <-ctx.Done():
		log.Printf("shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	log.Printf("server shutdown complete")

	return nil
}
