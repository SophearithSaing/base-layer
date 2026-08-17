package auth

import (
	"baselayer/internal/user"
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Service struct {
	refreshTokenRepo *RefreshTokenRepository
	jwtProvider      *JWTProvider
	userService      *user.Service
}

func NewService(refreshTokenRepo *RefreshTokenRepository, jwtProvider *JWTProvider, userService *user.Service) *Service {
	return &Service{
		refreshTokenRepo: refreshTokenRepo,
		userService:      userService,
		jwtProvider:      jwtProvider,
	}
}

func (s *Service) Register(ctx context.Context, payload RegisterPayload) (RegisterResponse, string, string, error) {
	username := strings.ToLower(strings.TrimSpace(payload.Username))
	err := validateInput(username, payload.Password)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}

	filter := bson.D{{Key: "username", Value: username}}
	_, err = s.userService.FindOne(ctx, filter)
	if err == nil {
		return RegisterResponse{}, "", "", ErrUsernameAlreadyExists
	}
	if !errors.Is(err, user.ErrUserNotFound) {
		return RegisterResponse{}, "", "", err
	}

	passwordHash, err := hashPassword(payload.Password)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	createUserPayload := user.CreateUserPayload{
		Username:     username,
		PasswordHash: passwordHash,
	}
	result, err := s.userService.Create(ctx, createUserPayload)
	if errors.Is(err, user.ErrUserAlreadyExists) {
		return RegisterResponse{}, "", "", ErrUsernameAlreadyExists
	}
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	token, err := s.jwtProvider.signJWT(result.Id.Hex())
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	refreshToken, err := s.issueRefreshToken(ctx, result.Id)
	if err != nil {
		return RegisterResponse{}, "", "", err
	}
	return RegisterResponse{ID: result.Id, Username: result.Username}, token, refreshToken, nil
}

func (s *Service) Login(ctx context.Context, payload LoginPayload) (string, string, error) {
	username := strings.ToLower(strings.TrimSpace(payload.Username))
	err := validateInput(username, payload.Password)
	if err != nil {
		return "", "", err
	}

	filter := bson.D{{Key: "username", Value: username}}
	existing, err := s.userService.FindOne(ctx, filter)
	if err != nil {
		return "", "", err
	}

	result := verifyPassword(payload.Password, existing.PasswordHash)
	if !result {
		return "", "", ErrIncorrectPassword
	}

	accessToken, err := s.jwtProvider.signJWT(existing.Id.Hex())
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.issueRefreshToken(ctx, existing.Id)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	_, err := s.validateRefreshToken(ctx, token)
	if errors.Is(err, ErrTokenNotFound) || errors.Is(err, ErrTokenIsRevoked) || errors.Is(err, ErrTokenIsExpired) {
		return nil
	}
	if err != nil {
		return err
	}
	return s.revokeRefreshToken(ctx, token)
}

func validateInput(username, password string) error {
	if len(username) == 0 {
		return ErrUsernameNotProvided
	}
	if len(password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (s *Service) Me(ctx context.Context) (UserResponse, error) {
	userId, err := CurrentUserID(ctx)
	if err != nil {
		return UserResponse{}, err
	}
	return s.getUser(ctx, userId)
}

func (s *Service) getUser(ctx context.Context, userId string) (UserResponse, error) {
	user, err := s.userService.GetById(ctx, userId)
	if err != nil {
		return UserResponse{}, err
	}
	return UserResponse{ID: user.Id, Username: user.Username}, nil
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
	accessToken, err := s.jwtProvider.signJWT(existing.UserID.Hex())
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.issueRefreshToken(ctx, existing.UserID)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Service) issueRefreshToken(ctx context.Context, userId bson.ObjectID) (string, error) {
	token, hashedToken, err := generateRefreshToken()
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	err = s.refreshTokenRepo.Create(ctx, RefreshToken{
		UserID:      userId,
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
		return RefreshToken{}, ErrTokenIsRevoked
	}
	if time.Now().After(refreshToken.ExpiresAt) {
		return RefreshToken{}, ErrTokenIsExpired
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
