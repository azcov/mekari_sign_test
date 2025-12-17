package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"
)

// AI Service
type PredictCategoryRequest struct {
	File string `json:"file"`
}

type PredictCategoryResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Data    PredictCategoryData `json:"data"`
}

type PredictCategoryData struct {
	Category string `json:"category"`
}

type AiService struct {
	BaseURL    string
	HttpClient *http.Client
}

func NewAiServiceClient() *AiService {
	return &AiService{
		BaseURL: "https://ai.example.com",
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (p *AiService) Call(req []byte) string {
	minDelay := 10 // milliseconds
	maxDelay := 500
	delay := rand.IntN(maxDelay - minDelay)
	time.Sleep(time.Duration(delay+minDelay) * time.Millisecond)
	category := categories[rand.IntN(len(categories))]
	return fmt.Sprintf(`
{
	"data": {
		"category": "%s"
	},
	"status": "success",
	"message": "Category predicted successfully"
}
	`, category.Name)
}

func (p *AiService) PredictCategory(todo Todo) (Category, error) {
	// Mock implementation: randomly assign a category
	data, err := GeneratePdfTodo(todo)
	if err != nil {
		return Category{}, err
	}
	b64Pdf, err := EncodePDFToBase64(data)
	if err != nil {
		return Category{}, err
	}

	req := PredictCategoryRequest{
		File: b64Pdf,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return Category{}, err
	}

	resp := p.Call(reqBytes)
	var response PredictCategoryResponse
	if err := json.Unmarshal([]byte(resp), &response); err != nil {
		return Category{}, err
	}

	return Category{Name: response.Data.Category}, nil
}

func GeneratePdfTodo(todo Todo) (string, error) {
	// Mock implementation: return a base64 encoded string as PDF data
	return fmt.Sprintf("%s_%s.pdf", todo.Text, todo.Description), nil
}

func EncodePDFToBase64(path string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(path)), nil
}
