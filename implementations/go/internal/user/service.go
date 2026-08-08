package user

import (
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

func (s *Service) Create(ctx context.Context, payload CreateUserPayload) (User, error) {
	user := User{
		Id:           bson.NewObjectID(),
		Username:     payload.Username,
		PasswordHash: payload.PasswordHash,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *Service) GetById(ctx context.Context, id string) (User, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service) FindOne(ctx context.Context, filter bson.D) (User, error) {
	return s.repo.FindOne(ctx, filter)
}
