package auth

import (
	"baselayer/internal/user"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	refreshTokenRepo *RefreshTokenRepository
	userService      *user.Service
}

func NewService(refreshTokenRepo *RefreshTokenRepository, userService *user.Service) *Service {
	return &Service{
		refreshTokenRepo: refreshTokenRepo,
		userService:      userService,
	}
}

func (s *Service) Register(ctx context.Context, payload RegisterPayload) (user.User, error) {
	filter := bson.D{{Key: "username", Value: payload.Username}}
	_, err := s.userService.FindOne(ctx, filter)
	if err == nil {
		return user.User{}, fmt.Errorf("username already exists")
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		return user.User{}, err
	}
	createUserPayload := user.CreateUserPayload{
		Username:     payload.Username,
		PasswordHash: passwordHash,
	}
	result, err := s.userService.Create(ctx, createUserPayload)
	if err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (s *Service) Login(ctx context.Context, payload LoginPayload) (user.User, error) {
	filter := bson.D{{Key: "username", Value: payload.Username}}
	existing, err := s.userService.FindOne(ctx, filter)
	if err != nil {
		return user.User{}, err
	}

	result := verifyPassword(payload.Password, existing.PasswordHash)
	if result {
		return existing, nil
	} else {
		return user.User{}, err
	}
}
