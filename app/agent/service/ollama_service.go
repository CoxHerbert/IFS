package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type OllamaService struct {
	BaseURL string
	Client  *http.Client
}

type OllamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Stream   bool            `json:"stream"`
	Messages []OllamaMessage `json:"messages"`
}

type ollamaChatResponse struct {
	Message OllamaMessage `json:"message"`
	Error   string        `json:"error,omitempty"`
}

func NewOllamaService() *OllamaService {
	return &OllamaService{BaseURL: "http://localhost:11434", Client: &http.Client{Timeout: 90 * time.Second}}
}

func (s *OllamaService) Chat(modelName string, messages []OllamaMessage) (string, error) {
	payload, err := json.Marshal(ollamaChatRequest{Model: modelName, Stream: false, Messages: messages})
	if err != nil {
		return "", err
	}
	resp, err := s.Client.Post(s.BaseURL+"/api/chat", "application/json", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if result.Error != "" {
			return "", errors.New(result.Error)
		}
		return "", fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}
	if result.Error != "" {
		return "", errors.New(result.Error)
	}
	return result.Message.Content, nil
}
