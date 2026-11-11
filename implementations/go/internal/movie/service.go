package movie

import (
	"baselayer/internal/auth"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type movieRepository interface {
	create(context.Context, movie) (string, error)
	bulkCreate(context.Context, []movie) ([]string, error)
	getByID(context.Context, string) (*movie, error)
	search(context.Context, bson.M, bson.D) (*[]movie, error)
	update(context.Context, bson.M, bson.M) (*movie, error)
}

type service struct {
	repo movieRepository
}

func NewService(repo movieRepository) *service {
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

func (s *service) bulkCreateMovies(ctx context.Context, payloads []createMoviePayload) ([]string, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if len(payloads) == 0 {
		return nil, ErrEmptyMovieBulk
	}
	now := time.Now()
	movies := make([]movie, len(payloads))
	for i, payload := range payloads {
		if payload.Rating == 0 {
			return nil, fmt.Errorf("movie %d: %w", i+1, ErrInvalidRating)
		}
		movies[i] = movie{
			ID:        bson.NewObjectID(),
			CreatorID: userID,
			Date:      payload.Date,
			Metadata:  payload.Metadata,
			Rating:    payload.Rating,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return s.repo.bulkCreate(ctx, movies)
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
