package auth

import (
	"baselayer/internal/user"
	"context"
	"errors"
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

func (s *Service) Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, string, string, error) {
	err := validateInput(payload.Username, payload.Password)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}

	filter := bson.D{{Key: "username", Value: payload.Username}}
	_, err = s.userService.FindOne(ctx, filter)
	if err == nil {
		return RegisterResponse{}, "", "", ErrUsernameAlreadyExists
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	createUserPayload := user.CreateUserPayload{
		Username:     payload.Username,
		PasswordHash: passwordHash,
	}
	result, err := s.userService.Create(ctx, createUserPayload)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	token, err := s.signJWT(result.Id.String())
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	refreshToken, err := s.issueRefreshToken(ctx, result.Id)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	return RegisterResponse{Id: result.Id, Username: result.Username}, token, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, payload LoginPayload) (string, string, error) {
	err := validateInput(payload.Username, payload.Password)
	if err != nil {
		return "", "", err
	}

	filter := bson.D{{Key: "username", Value: payload.Username}}
	existing, err := s.userService.FindOne(ctx, filter)
	if err != nil {
		return "", "", err
	}

	result := verifyPassword(payload.Password, existing.PasswordHash)
	if !result {
		return "", "", err
	}

	accessToken, err := s.signJWT(existing.Id.String())
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.issueRefreshToken(ctx, existing.Id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func validateInput(username, password string) error {
	if len(username) == 0 {
		return ErrInvalidUsername
	}
	if len(password) < 8 {
		return ErrInvalidPassword
	}
	return nil
}

func (s *Service) Refresh(ctx context.Context, token string) (string, string, error) {
	existing, err := s.validateRefreshToken(ctx, token)
	if err != nil {
		return "", "", err
	}
	err = s.revokeRefreshToken(ctx, token)
	if err != nil {
		return "", "", err
	}
	accessToken, err := s.signJWT(existing.UserId.String())
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.issueRefreshToken(ctx, existing.UserId)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
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

func (s *Service) issueRefreshToken(ctx context.Context, userId bson.ObjectID) (string, error) {
	token, hashedToken, err := generateRefreshToken()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	err = s.refreshTokenRepo.Create(ctx, RefreshToken{
		UserId:      userId,
		HashedToken: hashedToken,
		IsRevoked:   false,
		RevokedAt:   nil,
		ExpiresAt:   now.AddDate(0, 0, 30),
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) validateRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	hashedToken := hashRefreshToken(token)
	filter := bson.D{{Key: "hashedToken", Value: hashedToken}}
	refreshToken, err := s.refreshTokenRepo.FindOne(ctx, filter)
	if err != nil {
		return RefreshToken{}, err
	}
	if refreshToken.IsRevoked {
		return RefreshToken{}, errors.New("token is revoked")
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		return RefreshToken{}, errors.New("token is expired")
	}
	return refreshToken, nil
}

func (s *Service) revokeRefreshToken(ctx context.Context, token string) error {
	hashedToken := hashRefreshToken(token)
	filter := bson.M{"hashedToken": hashedToken, "isRevoked": false}
	update := bson.M{"$set": bson.M{"isRevoked": true, "revokedAt": time.Now().UTC()}}
	_, err := s.refreshTokenRepo.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
