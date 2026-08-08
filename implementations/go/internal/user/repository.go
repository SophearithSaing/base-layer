package user

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("users")}
}

func (r *Repository) Create(ctx context.Context, user User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (r *Repository) GetById(ctx context.Context, id string) (User, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return User{}, err
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	var user User
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *Repository) FindOne(ctx context.Context, filter bson.D) (User, error) {
	var user User
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}
