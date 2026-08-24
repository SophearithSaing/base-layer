package group

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Group struct {
	ID        bson.ObjectID `bson:"_id" json:"_id"`
	CreatorID bson.ObjectID `bson:"creatorId" json:"creatorId"`
	Name      string        `bson:"name" json:"name"`
	People    []string      `bson:"people" json:"people"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type CreateGroupPayload struct {
	Name   string   `bson:"name" json:"name"`
	People []string `bson:"people" json:"people"`
}

type UpdatePayload[T any] struct {
	Payload   T         `bson:",inline"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type UpdateGroupPayload struct {
	Name   string   `bson:"name,omitempty" json:"name,omitempty"`
	People []string `bson:"people,omitempty" json:"people,omitempty"`
}
