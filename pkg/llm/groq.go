package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
	GenerateJSON(ctx context.Context, system, user string) (string, error)
}

type GroqLLM struct {
	BaseURL string
	Model   string
	APIKey  string
	Client  *http.Client
}

func NewGroqLLM(baseURL, model, apiKey string) *GroqLLM {
	if baseURL == "" {
		baseURL = "https://api.groq.com"
	}
	return &GroqLLM{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Model:   model,
		APIKey:  apiKey,
		Client:  &http.Client{Timeout: 60 * time.Second},
	}
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqChatRequest struct {
	Model          string           `json:"model"`
	Messages       []groqMessage    `json:"messages"`
	Temperature    float64          `json:"temperature,omitempty"`
	ResponseFormat *groqResponseFmt `json:"response_format,omitempty"`
}

type groqResponseFmt struct {
	Type string `json:"type"`
}

type groqChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (l *GroqLLM) Generate(ctx context.Context, prompt string) (string, error) {
	return l.chat(ctx, groqChatRequest{
		Model: l.Model,
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
	})
}

func (l *GroqLLM) GenerateJSON(ctx context.Context, system, user string) (string, error) {
	messages := make([]groqMessage, 0, 2)
	if system != "" {
		messages = append(messages, groqMessage{Role: "system", Content: system})
	}
	messages = append(messages, groqMessage{Role: "user", Content: user})

	return l.chat(ctx, groqChatRequest{
		Model:       l.Model,
		Messages:    messages,
		Temperature: 0.4,
		ResponseFormat: &groqResponseFmt{
			Type: "json_object",
		},
	})
}

func (l *GroqLLM) chat(ctx context.Context, payload groqChatRequest) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, l.BaseURL+"/openai/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+l.APIKey)

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to generate response: %s: %s", resp.Status, string(respBody))
	}

	var result groqChatResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty llm response")
	}

	return result.Choices[0].Message.Content, nil
}
