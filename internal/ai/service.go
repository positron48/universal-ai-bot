package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Service handles AI provider interactions
type Service struct {
	client *http.Client
	url    string
	model  string
	apiKey string
	prompt string
	logger *zap.Logger
}

// NewService creates a new AI service
func NewService(url, model, apiKey, prompt string, logger *zap.Logger) *Service {
	// Process prompt to handle escaped newlines
	processedPrompt := strings.ReplaceAll(prompt, "\\n", "\n")

	return &Service{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		url:    url,
		model:  model,
		apiKey: apiKey,
		prompt: processedPrompt,
		logger: logger,
	}
}

// ChatRequest represents the OpenAI-compatible chat request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents the OpenAI-compatible chat response
type ChatResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Choice represents a response choice
type Choice struct {
	Message Message `json:"message"`
}

// Error represents an API error
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GenerateResponse sends a message to the AI provider and returns the response
func (s *Service) GenerateResponse(ctx context.Context, userMessage string) (string, error) {
	// Prepare messages with system prompt
	messages := []Message{
		{
			Role:    "system",
			Content: s.prompt,
		},
		{
			Role:    "user",
			Content: userMessage,
		},
	}

	// Create request
	req := ChatRequest{
		Model:       s.model,
		Messages:    messages,
		MaxTokens:   1000,
		Temperature: 0.7,
	}

	// Marshal request
	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.url+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.apiKey)

	s.logger.Debug("sending request to AI provider",
		zap.String("url", s.url),
		zap.String("model", s.model),
		zap.String("user_message", userMessage),
	)

	// Send request
	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		s.logger.Error("AI provider returned error",
			zap.Int("status_code", resp.StatusCode),
			zap.String("response", string(respBody)),
		)
		return "", fmt.Errorf("AI provider returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for API error
	if chatResp.Error != nil {
		return "", fmt.Errorf("AI provider error: %s", chatResp.Error.Message)
	}

	// Check if we have choices
	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no response choices received")
	}

	response := chatResp.Choices[0].Message.Content
	s.logger.Debug("received response from AI provider",
		zap.String("response", response),
	)

	return response, nil
}
