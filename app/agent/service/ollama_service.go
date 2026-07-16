package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	baseURL, timeoutSeconds := ollamaRuntimeConfig()
	return &OllamaService{BaseURL: baseURL, Client: &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}}
}

func ollamaRuntimeConfig() (string, int) {
	baseURL := "http://localhost:11434"
	timeoutSeconds := 90
	config := GetConfigService().GetOllamaConfig()
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	if config.Timeout > 0 {
		timeoutSeconds = config.Timeout
	}
	return strings.TrimRight(baseURL, "/"), timeoutSeconds
}

func (s *OllamaService) Chat(modelName string, messages []OllamaMessage) (string, error) {
	payload, err := json.Marshal(ollamaChatRequest{Model: modelName, Stream: false, Messages: messages})
	if err != nil {
		return "", err
	}
	baseURL, timeoutSeconds := ollamaRuntimeConfig()
	client := s.Client
	if timeoutSeconds > 0 {
		client = &http.Client{Timeout: time.Duration(timeoutSeconds) * time.Second}
	}
	resp, err := client.Post(baseURL+"/api/chat", "application/json", bytes.NewReader(payload))
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
