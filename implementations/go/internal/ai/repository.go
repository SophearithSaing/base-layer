package ai

import "go.mongodb.org/mongo-driver/v2/mongo"

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
