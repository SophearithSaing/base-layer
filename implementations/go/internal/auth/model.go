package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type RegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Id       bson.ObjectID `json:"id"`
	Username string        `json:"username"`
}

type RefreshToken struct {
	UserId          bson.ObjectID `json:"userId"`
	TokenLookupHash string        `json:"tokenLookupHash"`
	TokenHash       string        `json:"tokenHash"`
	ExpiresAt       time.Time     `json:"expiresAt"`
	RevokedAt       time.Time     `json:"revokedAt"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}
