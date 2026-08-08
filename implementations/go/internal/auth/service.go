package auth

import (
	"baselayer/internal/user"
	"context"
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
	// Add check existing validation
	// filter := bson.D{{Key: "username", Value: payload.Username}}
	// queryResult, err := s.userService.FindOne(ctx, filter)
	// if err == nil {
	// 	return user.User{}, err
	// }

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
