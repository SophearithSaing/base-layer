package movie

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{collection: db.Collection("movie")}
}

func (r *repository) create(ctx context.Context, movie movie) (string, error) {
	result, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return "", err
	}
	id, err := getStringID(result.InsertedID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *repository) getByID(ctx context.Context, id string) (*movie, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID}
	var movie movie
	err = r.collection.FindOne(ctx, filter).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *repository) search(ctx context.Context, filter bson.M, sort bson.D) (*[]movie, error) {
	opts := options.Find().SetSort(sort)
	var movies []movie
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &movies)
	if err != nil {
		return nil, err
	}
	return &movies, nil
}

func (r *repository) update(ctx context.Context, filter bson.M, update bson.M) (*movie, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var movie movie
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&movie)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func getStringID(raw any) (string, error) {
	objectID, ok := raw.(bson.ObjectID)
	if !ok {
		return "", ErrInvalidID
	}
	return objectID.Hex(), nil
}
