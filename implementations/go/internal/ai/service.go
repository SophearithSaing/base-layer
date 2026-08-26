package ai

import (
	"baselayer/internal/auth"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AIModel string

const (
	Gemma_4   AIModel = "google/gemma-4-31B-it"
	GLM_5_2   AIModel = "zai-org/GLM-5.2"
	MiniMaxM3 AIModel = "MiniMaxAI/MiniMax-M3"
)

const ChatURL = "https://api.together.ai/v1/chat/completions"
const ChatURLv2 = "https://api-inference.together.ai/v2/chat/completions"

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
			Timeout: 60 * time.Second,
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

func (s *Service) ExtractText(ctx context.Context, file io.Reader, fileName string, mimeType string) (*AIExtraction, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	dataURL := fmt.Sprintf(
		"data:%s;base64,%s",
		mimeType,
		base64.StdEncoding.EncodeToString(data),
	)
	payload := ChatRequest{
		Model: MiniMaxM3,
		Messages: []Message{
			{
				Role:    MessageRoleSystem,
				Content: "You are an OCR text extraction assistant. Your task is to extract all visible text from the provided image as accurately as possible.",
			},
			{
				Role: MessageRoleUser,
				Content: []ImageURLContent{
					{
						Type: "image_url",
						ImageURL: ImageURL{
							URL: dataURL,
						},
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ChatURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("together ai error: %s", responseBody)
	}
	defer res.Body.Close()

	var chatResponse ChatResponse
	err = json.NewDecoder(res.Body).Decode(&chatResponse)
	if err != nil {
		return nil, err
	}
	aiExtraction := AIExtraction{
		ID:            bson.NewObjectID(),
		UserID:        userID,
		ExtractedText: chatResponse.Choices[0].Message.Content,
		FileName:      fileName,
		MimeType:      mimeType,
		CreatedAt:     time.Now(),
	}
	_, err = s.repo.CreateExtraction(ctx, aiExtraction)
	if err != nil {
		return nil, err
	}

	return &aiExtraction, nil
}

func (s *Service) ExtractReceipt(ctx context.Context, file io.Reader, fileName string, mimeType string) (*AIExtraction, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	dataURL := fmt.Sprintf(
		"data:%s;base64,%s",
		mimeType,
		base64.StdEncoding.EncodeToString(data),
	)
	additionalProperties := false
	schema := JSONSchema{
		Type: "array",
		Items: &JSONSchema{
			Type:                 "object",
			AdditionalProperties: &additionalProperties,
			Required:             []string{"name", "amount"},
			Properties: map[string]*JSONSchema{
				"name": {
					Type: "string",
				},
				"amount": {
					Type: "number",
				},
			},
		},
	}
	payload := ChatRequest{
		ResponseFormat: &ResponseFormat{
			Type: "json_schema",
			JSONSchema: ResponseSchema{
				Name:   "receipt",
				Schema: schema,
			},
		},
		Model: MiniMaxM3,
		Messages: []Message{
			{
				Role: MessageRoleSystem,
				Content: `
			 				You are a receipt extraction assistant. Extract only purchased line items from the receipt image.
			       			Return only a valid JSON array. Each object must have exactly:

			          		* name: item name as shown
			            	* amount: final line item price as a number

			              	Rules:
			               	* Exclude subtotal, tax, discounts, tips, fees, totals, payment info, dates, store info, and receipt metadata.
			                * If quantity is shown, return one object per purchased unit.
			                * If the receipt shows 2 Chicken Sandwich 11.98, return two objects, each with "amount": 5.99\.
			                * If only the total line amount is shown for multiple units, divide it by the quantity.
			                * Omit items with unreadable names or prices.
			                * Do not guess, explain, or add Markdown.
			                * If no items are found, return [].

			                Example:
			                [
			                 	{
			                  		"name": "Chicken Sandwich",
			                    		"amount": 5.99
			                      }
			                ]
				`,
			},
			{
				Role: MessageRoleUser,
				Content: []ImageURLContent{
					{
						Type: "image_url",
						ImageURL: ImageURL{
							URL: dataURL,
						},
					},
				},
			},
		},
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ChatURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("together ai error: %s", responseBody)
	}
	defer res.Body.Close()

	var chatResponse ChatResponse
	err = json.NewDecoder(res.Body).Decode(&chatResponse)
	if err != nil {
		return nil, err
	}
	var extractedItems []Item
	err = json.Unmarshal([]byte(chatResponse.Choices[0].Message.Content), &extractedItems)
	if err != nil {
		return nil, err
	}
	aiExtraction := AIExtraction{
		ID:             bson.NewObjectID(),
		UserID:         userID,
		FileName:       fileName,
		MimeType:       mimeType,
		ExtractedItems: extractedItems,
		CreatedAt:      time.Now(),
	}
	_, err = s.repo.CreateExtraction(ctx, aiExtraction)

	return &aiExtraction, nil
}

func (s *Service) GetExtractions(ctx context.Context) (*[]AIExtraction, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"userId": userID}
	sort := bson.D{{Key: "createdAt", Value: -1}}

	return s.repo.GetExtractions(ctx, filter, sort)
}

func currentUserID(ctx context.Context) (bson.ObjectID, error) {
	rawID, err := auth.CurrentUserID(ctx)
	if err != nil {
		return bson.NilObjectID, err
	}
	objectID, err := bson.ObjectIDFromHex(rawID)
	if err != nil {
		return bson.NilObjectID, err
	}
	return objectID, nil
}
