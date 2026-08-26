package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
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

	return chatResponse.Choices[0].Message.Content.(string), nil
}

func (s *Service) ExtractText(ctx context.Context, file io.Reader, contentType string) (string, error) {
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	dataURL := fmt.Sprintf(
		"data:%s;base64,%s",
		contentType,
		base64.StdEncoding.EncodeToString(data),
	)
	payload := ChatRequest{
		Model: MiniMaxM3,
		Messages: []Message{
			{
				Role:    MessageRoleSystem,
				Content: "You are an OCR text extraction assistant. Your task is to extract all visible text from the provided image as accurately as possible.",
			},
			// {
			// 	Role: MessageRoleSystem,
			// 	Content: `
			//  				You are a receipt extraction assistant. Extract only purchased line items from the receipt image.
			//        			Return only a valid JSON array. Each object must have exactly:

			//           		* name: item name as shown
			//             		* amount: final line item price as a number

			//               	Rules:
			//                	* Exclude subtotal, tax, discounts, tips, fees, totals, payment info, dates, store info, and receipt metadata.
			//                 	* If quantity is shown, return one object per purchased unit.
			//                  * If the receipt shows 2 Chicken Sandwich 11.98, return two objects, each with "amount": 5.99\.
			//                  * If only the total line amount is shown for multiple units, divide it by the quantity.
			//                  * Omit items with unreadable names or prices.
			//                  * Do not guess, explain, or add Markdown.
			//                  * If no items are found, return [].

			//                  Example:
			//                  [
			//                  	{
			//                   		"name": "Chicken Sandwich",
			//                     		"amount": 5.99
			//                       }
			//                  ]
			// 	`,
			// },
			{
				Role: MessageRoleUser,
				Content: []map[string]any{
					{
						"type": "image_url",
						"image_url": map[string]string{
							"url": dataURL,
						},
					},
				},
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
		bytes.NewReader(body),
	)
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

	return chatResponse.Choices[0].Message.Content.(string), nil
}
