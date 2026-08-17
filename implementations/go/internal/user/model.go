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
	ID           bson.ObjectID `bson:"_id"`
	Username     string        `bson:"username"`
	PasswordHash string        `bson:"passwordHash"`
	CreatedAt    time.Time     `bson:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt"`
}
