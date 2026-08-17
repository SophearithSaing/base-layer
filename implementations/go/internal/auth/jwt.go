package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
}

func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{
		secret: []byte(secret),
	}
}

func (j *JWTProvider) signJWT(userId string) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userId,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTProvider) verifyJWT(raw string) (jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(
		raw, &claims, func(token *jwt.Token) (any, error) {
			return j.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}
	if !token.Valid {
		return jwt.RegisteredClaims{}, ErrTokenIsInvalid
	}
	return claims, nil
}
