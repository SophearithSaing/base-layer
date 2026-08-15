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

type LoginResponse struct {
	Success bool `json:"success"`
}

type RegisterResponse struct {
	Id       bson.ObjectID `json:"id"`
	Username string        `json:"username"`
}

type RefreshResponse struct {
	Success bool `json:"success"`
}

type RefreshToken struct {
	UserId      bson.ObjectID `bson:"userId"`
	HashedToken string        `bson:"hashedToken"`
	IsRevoked   bool          `bson:"isRevoked"`
	ExpiresAt   time.Time     `bson:"expiresAt"`
	RevokedAt   *time.Time    `bson:"revokedAt"`
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
}
