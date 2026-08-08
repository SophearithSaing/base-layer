package auth

import "go.mongodb.org/mongo-driver/v2/mongo"

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}
