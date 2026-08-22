package project

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Legend struct {
	Difficulty map[string]string `json:"difficulty"`
	Type       map[string]string `json:"type"`
}

type Item struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"string"`
}

type Phase struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Type          string   `json:"type"`
	Difficulty    int      `json:"difficulty"`
	Summary       string   `json:"summary"`
	Concepts      []Item   `json:"concepts"`
	Tools         []Item   `json:"tools"`
	Practice      []Item   `json:"practice"`
	MasteryChecks []string `json:"masterChecks"`
	Prerequisites []string `json:"prerequisites"`
}

type Capstone struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Difficulty    int      `json:"difficulty"`
	Summary       string   `json:"summary"`
	Build         []string `json:"build"`
	Concepts      []Item   `json:"concepts"`
	Tools         []Item   `json:"tools"`
	Prerequisites []string `json:"prerequisites"`
}

type NoteContent struct {
	Notes []struct {
		Text      string    `json:"text"`
		Timestamp time.Time `json:"timestamp"`
	} `json:"notes"`
	Links []struct {
		Text string `json:"text"`
		URL  string `json:"url"`
	} `json:"links"`
}

type Project struct {
	ID                 bson.ObjectID `bson:"_id" json:"id"`
	Title              string        `bson:"title" json:"title"`
	Description        string        `bson:"description" json:"description"`
	Legend             Legend        `bson:"legend" json:"legend"`
	Phases             []Phase       `bson:"phases" json:"phases"`
	Capstones          []Capstone    `bson:"capstones" json:"capstones"`
	RecommendedOrder   []string      `bson:"recommendedOrder" json:"recommendedOrder"`
	MasteryDefinitions []string      `bson:"masteryDefintions" json:"masteryDefinitions"`
	CreatedAt          time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt          time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type ProjectProgress struct {
	ID             bson.ObjectID          `bson:"_id" json:"id"`
	UserID         bson.ObjectID          `bson:"userId" json:"userId"`
	ProjectID      bson.ObjectID          `bson:"projectId" json:"projectId"`
	Title          string                 `bson:"title" json:"title"`
	Description    string                 `bson:"description" json:"description"`
	CompletedItems map[string]bool        `bson:"completedItems" json:"completedItems"`
	Progress       int                    `bson:"progress" json:"progress"`
	Notes          map[string]NoteContent `bson:"notes" json:"notes"`
	CreatedAt      time.Time              `bson:"createdAt" json:"createdAt"`
	UpdatedAt      time.Time              `bson:"updatedAt" json:"updatedAt"`
}
