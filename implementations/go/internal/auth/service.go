package auth

import (
	"baselayer/internal/user"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	refreshTokenRepo *RefreshTokenRepository
	userService      *user.Service
	jwtSecret        []byte
}

func NewService(refreshTokenRepo *RefreshTokenRepository, userService *user.Service, jwtSecret string) *Service {
	return &Service{
		refreshTokenRepo: refreshTokenRepo,
		userService:      userService,
		jwtSecret:        []byte(jwtSecret),
	}
}

func (s *Service) Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, string, error) {
	filter := bson.D{{Key: "username", Value: payload.Username}}
	_, err := s.userService.FindOne(ctx, filter)
	if err == nil {
		return RegisterResponse{}, "", fmt.Errorf("username already exists")
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		return RegisterResponse{}, "", err
	}
	createUserPayload := user.CreateUserPayload{
		Username:     payload.Username,
		PasswordHash: passwordHash,
	}
	result, err := s.userService.Create(ctx, createUserPayload)
	if err != nil {
		return RegisterResponse{}, "", err
	}
	token, err := s.signJWT(result.Id.String())
	if err != nil {
		return RegisterResponse{}, "", err
	}
	return RegisterResponse{Id: result.Id, Username: result.Username}, token, nil
}

func (s *Service) Login(ctx context.Context, payload LoginPayload) (string, error) {
	filter := bson.D{{Key: "username", Value: payload.Username}}
	existing, err := s.userService.FindOne(ctx, filter)
	if err != nil {
		return "", err
	}

	result := verifyPassword(payload.Password, existing.PasswordHash)
	if result {
		return s.signJWT(existing.Id.String())
	} else {
		return "", err
	}
}

func (s *Service) signJWT(userId string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
