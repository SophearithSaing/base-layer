package ai

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

type AIMessage struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	UserID    bson.ObjectID `bson:"userId" json:"userId"`
	Role      MessageRole   `bson:"role" json:"role"`
	Content   string        `bson:"content" json:"content"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
}

type Item struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

type AIExtraction struct {
	ID             bson.ObjectID `bson:"_id" json:"id"`
	UserID         bson.ObjectID `bson:"userId" json:"userId"`
	FileName       string        `bson:"fileName" json:"fileName"`
	MimeType       string        `bson:"mimeType" json:"mimeType"`
	ExtractedText  string        `bson:"extractedText" json:"extractedText"`
	ExtractedItems []Item        `bson:"extractedItems" json:"extractedItems"`
	CreatedAt      time.Time     `bson:"createdAt" json:"createdAt"`
}

type MessagePayload struct {
	Message string `json:"message"`
}

type Message struct {
	Role    MessageRole `json:"role"`
	Content any         `json:"content"`
}

type ResponseSchema struct {
	Name   string     `json:"name"`
	Schema JSONSchema `json:"schema"`
}

type ResponseFormat struct {
	Type       string         `json:"type"`
	JSONSchema ResponseSchema `json:"json_schema"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ImageURLContent struct {
	Type     string   `json:"type"`
	ImageURL ImageURL `json:"image_url"`
}

type ChatRequest struct {
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	Model          AIModel         `json:"model"`
	Messages       []Message       `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type JSONSchema struct {
	Type                 string                 `json:"type,omitempty"`
	Items                *JSONSchema            `json:"items,omitempty"`
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties *bool                  `json:"additionalProperties,omitempty"`
}
