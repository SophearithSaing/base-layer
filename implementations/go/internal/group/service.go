package group

import (
	"baselayer/internal/auth"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateGroup(ctx context.Context, payload CreateGroupPayload) (string, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	group := Group{
		ID:        bson.NewObjectID(),
		CreatorID: userID,
		Name:      payload.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, group)
}

func (s *Service) ListGroup(ctx context.Context) (*[]Group, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"creatorId": userID}
	sort := bson.D{{Key: "createdAt", Value: -1}}
	return s.repo.Search(ctx, filter, sort)
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
