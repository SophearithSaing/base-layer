package auth

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RefreshTokenRepository struct {
	collection *mongo.Collection
}

func NewRefreshTokenRepository(db *mongo.Database) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		collection: db.Collection("refresh_tokens"),
	}
}

func (rtr *RefreshTokenRepository) Create(ctx context.Context, refreshToken RefreshToken) error {
	_, err := rtr.collection.InsertOne(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("error creating token: %w", err)
	}
	return nil
}

func (rtr *RefreshTokenRepository) FindOne(ctx context.Context, filter bson.D) (RefreshToken, error) {
	var token RefreshToken
	err := rtr.collection.FindOne(ctx, filter).Decode(&token)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return RefreshToken{}, ErrTokenNotFound
	}
	if err != nil {
		return RefreshToken{}, err
	}
	return token, nil
}

func (rtr *RefreshTokenRepository) UpdateOne(ctx context.Context, filter bson.M, update bson.M) (RefreshToken, error) {
	var refreshToken RefreshToken
	err := rtr.collection.FindOneAndUpdate(ctx, filter, update).Decode(&refreshToken)
	if err != nil {
		return RefreshToken{}, err
	}
	return refreshToken, nil
}
