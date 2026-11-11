package movie

import (
	"baselayer/internal/auth"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type service struct {
	repo *repository
}

func NewService(repo *repository) *service {
	return &service{repo: repo}
}

func (s *service) createMovie(ctx context.Context, payload createMoviePayload) (string, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	movie := movie{
		ID:        bson.NewObjectID(),
		CreatorID: userID,
		Date:      payload.Date,
		Metadata:  payload.Metadata,
		Rating:    payload.Rating,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.create(ctx, movie)
}

func (s *service) listMovies(ctx context.Context) (*[]movie, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"creatorId": userID}
	sort := bson.D{{Key: "createdAt", Value: -1}}
	return s.repo.search(ctx, filter, sort)
}

func (s *service) getMovieByID(ctx context.Context, id string) (*movie, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	movie, err := s.repo.getByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if movie.CreatorID != userID {
		return nil, ErrUserDontHavePermissionToView
	}
	return movie, nil
}

func (s *service) updateMovie(ctx context.Context, id string, payload updateMoviePayload) (*movie, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"_id":       objectID,
		"creatorId": userID,
	}
	update := updatePayload[updateMoviePayload]{
		Payload:   payload,
		UpdatedAt: time.Now(),
	}
	return s.repo.update(ctx, filter, bson.M{"$set": update})
}

func currentUserID(ctx context.Context) (bson.ObjectID, error) {
	rawID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return bson.NilObjectID, err
	}
	objectID, err := bson.ObjectIDFromHex(rawID)
	if err != nil {
		return bson.NilObjectID, err
	}
	return objectID, nil
}
