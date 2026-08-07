package main

import (
	"baselayer/internal/api"
	"baselayer/internal/config"
	"baselayer/internal/db"
	"context"
	"fmt"
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

	mongodbConfig, err := config.GetMongoDBConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	allowedOrigins, err := config.GetClientOrigins()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mongo, err := db.NewMongo(mongodbConfig)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongo.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println("connected:", mongo.DB.Name())

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
		fmt.Println(err)
	}
}
