package db

import (
	"baselayer/internal/config"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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

func (m *Mongo) EnsureIndexes(ctx context.Context) error {
	users := m.DB.Collection("users")
	refreshTokens := m.DB.Collection("refresh_tokens")

	userIndexName, err := users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("username_uq"),
	})
	if err != nil {
		return err
	}
	log.Printf("index created: %v", userIndexName)

	rtIndexNames, err := refreshTokens.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "userId", Value: 1}},
			Options: options.Index().SetName("userId_idx"),
		},
		{
			Keys:    bson.D{{Key: "hashedToken", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("hashedToken_uq"),
		},
		{
			Keys:    bson.D{{Key: "expiresAt", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0).SetName("expiresAt_ttl"),
		},
	})
	if err != nil {
		return err
	}
	log.Printf("indexes created: %v", strings.Join(rtIndexNames, ", "))

	return nil
}
