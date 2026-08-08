package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	} else {
		return fallback
	}
}

func GetPort() (string, error) {
	portValue := GetEnv("PORT", "8000")
	port, err := strconv.ParseInt(portValue, 10, 32)
	if err != nil {
		return "", fmt.Errorf("failed to parse port: %v", err)
	}

	if port < 1 || port > 65535 {
		return "", fmt.Errorf("invalid port: %d", port)
	}

	return portValue, nil
}

func GetClientOrigins() ([]string, error) {
	clientOrigins := GetEnv("CLIENT_ORIGINS", "http://localhost:5173")

	if clientOrigins == "" {
		return nil, fmt.Errorf("no client origin found")
	}

	var origins []string
	for origin := range strings.SplitSeq(clientOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins, nil
}

type MongoDBConfig struct {
	URI    string
	DBName string
}

func GetMongoDBConfig() (MongoDBConfig, error) {
	var mongoDBConfig MongoDBConfig
	uri := GetEnv("MONGODB_URI", "")
	if uri == "" {
		return mongoDBConfig, fmt.Errorf("MONGODB_URI not found")
	}

	dbName := GetEnv("MONGODB_DB_NAME", "baselayer")

	mongoDBConfig = MongoDBConfig{
		URI:    uri,
		DBName: dbName,
	}
	return mongoDBConfig, nil
}
