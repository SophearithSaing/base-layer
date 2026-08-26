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
	Item   string `json:"item"`
	Amount int    `json:"amount"`
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

type ChatRequest struct {
	Model    AIModel   `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
