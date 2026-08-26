package ai

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	messagesCollection    *mongo.Collection
	extractionsCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		messagesCollection:    db.Collection("ai_messages"),
		extractionsCollection: db.Collection("ai_extractions"),
	}
}

func (r *Repository) CreateMessage(ctx context.Context, message AIMessage) (string, error) {
	result, err := r.messagesCollection.InsertOne(ctx, message)
	if err != nil {
		return "", err
	}
	id, err := getStringID(result.InsertedID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) CreateExtraction(ctx context.Context, extraction AIExtraction) (string, error) {
	result, err := r.extractionsCollection.InsertOne(ctx, extraction)
	if err != nil {
		return "", err
	}
	id, err := getStringID(result.InsertedID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetExtractions(ctx context.Context, filter bson.M, sort bson.D) (*[]AIExtraction, error) {
	opts := options.Find().SetSort(sort)
	cursor, err := r.extractionsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var aiExtractions []AIExtraction
	err = cursor.All(ctx, &aiExtractions)
	if err != nil {
		return nil, err
	}
	return &aiExtractions, nil
}

func getStringID(raw any) (string, error) {
	objectID, ok := raw.(bson.ObjectID)
	if !ok {
		return "", ErrInvalidID
	}
	return objectID.Hex(), nil
}
