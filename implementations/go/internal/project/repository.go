package project

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	ProjectCollection  *mongo.Collection
	ProgressCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		ProjectCollection:  db.Collection("projects"),
		ProgressCollection: db.Collection("project_progresses"),
	}
}

func (r *Repository) CreateProject(ctx context.Context, project Project) error {
	_, err := r.ProjectCollection.InsertOne(ctx, project)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) UpdateProject(ctx context.Context, filter bson.M, update bson.M) (Project, error) {
	var project Project
	err := r.ProgressCollection.FindOneAndUpdate(ctx, filter, update).Decode(&project)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (r *Repository) GetProjectByID(ctx context.Context, id string) (*Project, error) {
	var project Project
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = r.ProjectCollection.FindOne(ctx, filter).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *Repository) SearchProjects(ctx context.Context, filter bson.D, sort bson.D) ([]Project, error) {
	opts := options.Find().SetSort(sort)
	cursor, err := r.ProjectCollection.Find(ctx, filter, opts)
	if err != nil {
		return []Project{}, err
	}
	var projects []Project
	err = cursor.All(ctx, &projects)
	if err != nil {
		return []Project{}, err
	}
	return projects, nil
}

func (r *Repository) CreateProgress(ctx context.Context, progress ProjectProgress) (string, error) {
	result, err := r.ProgressCollection.InsertOne(ctx, progress)
	if err != nil {
		return "", err
	}
	id, err := getStringID(result.InsertedID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) UpdateProgress(ctx context.Context, filter bson.M, update bson.M) error {
	var progress ProjectProgress
	err := r.ProgressCollection.FindOneAndUpdate(ctx, filter, update).Decode(&progress)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetProgressByID(ctx context.Context, id string) (ProjectProgress, error) {
	var progress ProjectProgress
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return ProjectProgress{}, err
	}
	filter := bson.D{{Key: "_id", Value: objectID}}
	err = r.ProgressCollection.FindOne(ctx, filter).Decode(&progress)
	if err != nil {
		return ProjectProgress{}, err
	}
	return progress, nil
}

func (r *Repository) SearchProgresses(ctx context.Context, filter bson.D, sort bson.D) (*[]ProjectProgress, error) {
	opts := options.Find().SetSort(sort)
	cursor, err := r.ProgressCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var progresses []ProjectProgress
	err = cursor.All(ctx, &progresses)
	if err != nil {
		return nil, err
	}
	return &progresses, nil
}

func getStringID(raw any) (string, error) {
	objectID, ok := raw.(bson.ObjectID)
	if !ok {
		return "", ErrInvalidID
	}
	return objectID.Hex(), nil
}
