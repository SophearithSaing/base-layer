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
	rawID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return "", err
	}
	objectID, err := bson.ObjectIDFromHex(rawID)
	if err != nil {
		return "", err
	}
	now := time.Now()
	group := Group{
		ID:        bson.NewObjectID(),
		CreatorID: objectID,
		Name:      payload.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return s.repo.Create(ctx, group)
}
