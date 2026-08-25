package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AIModel string

const (
	Gemma_4   AIModel = "google/gemma-4-31B-it"
	KimiK_2_6 AIModel = "moonshotai/Kimi-K2.6"
	GLM_5_2   AIModel = "zai-org/GLM-5.2"
)

const ChatURL = "https://api.together.ai/v1/chat/completions"

type Service struct {
	repo   *Repository
	apiKey string
	client *http.Client
}

func NewService(repo *Repository, apiKey string) *Service {
	return &Service{
		repo:   repo,
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *Service) SendMessage(ctx context.Context, message string) (string, error) {
	payload := ChatRequest{
		Model: GLM_5_2,
		Messages: []Message{
			{
				Role:    "user",
				Content: message,
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ChatURL,
		bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var chatResponse ChatResponse
	err = json.NewDecoder(res.Body).Decode(&chatResponse)
	if err != nil {
		return "", err
	}

	return chatResponse.Choices[0].Message.Content, nil
}
