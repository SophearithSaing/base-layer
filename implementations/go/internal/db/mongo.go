package db

import (
	"baselayer/internal/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

type Mongo struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongo(config config.MongoDBConfig) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(config.URI))
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("error pinging mongodb: %w", err)
	}

	db := client.Database(config.DBName)
	return &Mongo{
		Client: client,
		DB:     db,
	}, nil
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
