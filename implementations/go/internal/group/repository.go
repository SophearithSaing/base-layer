package group

import "go.mongodb.org/mongo-driver/v2/mongo"

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("group")}
}
