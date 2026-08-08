package user

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CreateUserPayload struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type User struct {
	Id           bson.ObjectID `json:"id"`
	Username     string        `json:"username"`
	PasswordHash string        `json:"passwordHash"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}
