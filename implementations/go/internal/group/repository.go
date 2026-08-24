package group

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("group")}
}

func (r *Repository) Create(ctx context.Context, group Group) (string, error) {
	result, err := r.collection.InsertOne(ctx, group)
	if err != nil {
		return "", err
	}
	id, err := getStringID(result.InsertedID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Group, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	var group Group
	err = r.collection.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *Repository) Search(ctx context.Context, filter bson.M, sort bson.D) (*[]Group, error) {
	opts := options.Find().SetSort(sort)
	var groups []Group
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &groups)
	if err != nil {
		return nil, err
	}
	return &groups, nil
}

func (r *Repository) Update(ctx context.Context, id string, update bson.M) (*Group, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var group Group
	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&group)
	return &group, nil
}

func getStringID(raw any) (string, error) {
	objectID, ok := raw.(bson.ObjectID)
	if !ok {
		return "", ErrInvalidID
	}
	return objectID.Hex(), nil
}
