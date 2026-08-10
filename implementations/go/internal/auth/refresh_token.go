package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
)

func generateRefreshToken() (string, string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", "", err
	}
	rawToken := hex.EncodeToString(b)
	hashedToken := hashRefreshToken(rawToken)

	return rawToken, hashedToken, nil
}

func hashRefreshToken(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))
	return hashedToken
}

func verifyRefreshToken(rawToken string, refreshToken RefreshToken) error {
	if refreshToken.IsRevoked {
		return errors.New("token is revoked")
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		return errors.New("token is expired")
	}

	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	if hashedToken != refreshToken.HashedToken {
		return errors.New("token is invalid")
	}

	return nil
}
