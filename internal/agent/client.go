package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alexandrmotologa/consensai/internal/config"
)

// CompletionRequest contains the prompt payload for an LLM turn.
type CompletionRequest struct {
	SystemPrompt string  `json:"system_prompt"`
	UserPrompt   string  `json:"user_prompt"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
}

// ProviderClient abstracts interaction with diverse LLM inference endpoints.
type ProviderClient interface {
	GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error)
	Name() string
}

// NewProviderClient creates an LLM client based on the provider configuration.
func NewProviderClient(cfg *config.ProviderConfig) (ProviderClient, error) {
	provider := strings.ToLower(cfg.Default)
	if provider == "" {
		provider = "mock"
	}

	switch provider {
	case "mock", "offline", "heuristic":
		return NewMockProvider(), nil

	case "openai":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("openai provider requires an API key (set %s or provider.api_key)", cfg.APIKeyEnv)
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		model := cfg.Model
		if model == "" {
			model = "gpt-4o"
		}
		return &OpenAIProvider{
			apiKey:  cfg.APIKey,
			baseURL: baseURL,
			model:   model,
			client:  &http.Client{Timeout: 60 * time.Second},
		}, nil

	case "anthropic":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("anthropic provider requires an API key (set %s or provider.api_key)", cfg.APIKeyEnv)
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.anthropic.com/v1"
		}
		model := cfg.Model
		if model == "" {
			model = "claude-3-5-sonnet-20241022"
		}
		return &AnthropicProvider{
			apiKey:  cfg.APIKey,
			baseURL: baseURL,
			model:   model,
			client:  &http.Client{Timeout: 60 * time.Second},
		}, nil

	case "gemini":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("gemini provider requires an API key (set %s or provider.api_key)", cfg.APIKeyEnv)
		}
		model := cfg.Model
		if model == "" {
			model = "gemini-1.5-pro"
		}
		return &GeminiProvider{
			apiKey: cfg.APIKey,
			model:  model,
			client: &http.Client{Timeout: 60 * time.Second},
		}, nil

	case "ollama":
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434"
		}
		model := cfg.Model
		if model == "" {
			model = "llama3.1"
		}
		return &OllamaProvider{
			baseURL: baseURL,
			model:   model,
			client:  &http.Client{Timeout: 90 * time.Second},
		}, nil

	default:
		return nil, fmt.Errorf("unsupported provider %q (options: mock, openai, anthropic, gemini, ollama)", cfg.Default)
	}
}

// -------------------------------------------------------------------------
// OpenAI Provider Implementation
// -------------------------------------------------------------------------

type OpenAIProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func (p *OpenAIProvider) Name() string { return "openai:" + p.model }

func (p *OpenAIProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", strings.TrimRight(p.baseURL, "/"))

	type msg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	messages := []msg{
		{Role: "system", Content: req.SystemPrompt},
		{Role: "user", Content: req.UserPrompt},
	}

	payload := map[string]interface{}{
		"model":       p.model,
		"messages":    messages,
		"temperature": req.Temperature,
	}
	if req.MaxTokens > 0 {
		payload["max_tokens"] = req.MaxTokens
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("openai request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading openai response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("openai error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", fmt.Errorf("parsing openai response: %w", err)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("openai returned empty choices")
	}

	return parsed.Choices[0].Message.Content, nil
}

// -------------------------------------------------------------------------
// Anthropic Provider Implementation
// -------------------------------------------------------------------------

type AnthropicProvider struct {
	apiKey  string
	baseURL string
	model   string
	client  *http.Client
}

func (p *AnthropicProvider) Name() string { return "anthropic:" + p.model }

func (p *AnthropicProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	url := fmt.Sprintf("%s/messages", strings.TrimRight(p.baseURL, "/"))

	type msg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	payload := map[string]interface{}{
		"model":       p.model,
		"system":      req.SystemPrompt,
		"messages":    []msg{{Role: "user", Content: req.UserPrompt}},
		"max_tokens":  2048,
		"temperature": req.Temperature,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("anthropic request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading anthropic response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("anthropic error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}

	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", fmt.Errorf("parsing anthropic response: %w", err)
	}

	var sb strings.Builder
	for _, block := range parsed.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}

	return sb.String(), nil
}

// -------------------------------------------------------------------------
// Gemini Provider Implementation
// -------------------------------------------------------------------------

type GeminiProvider struct {
	apiKey string
	model  string
	client *http.Client
}

func (p *GeminiProvider) Name() string { return "gemini:" + p.model }

func (p *GeminiProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)

	payload := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": req.SystemPrompt},
			},
		},
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]interface{}{
					{"text": req.UserPrompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": req.Temperature,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading gemini response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("gemini error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", fmt.Errorf("parsing gemini response: %w", err)
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty candidates")
	}

	return parsed.Candidates[0].Content.Parts[0].Text, nil
}

// -------------------------------------------------------------------------
// Ollama Provider Implementation
// -------------------------------------------------------------------------

type OllamaProvider struct {
	baseURL string
	model   string
	client  *http.Client
}

func (p *OllamaProvider) Name() string { return "ollama:" + p.model }

func (p *OllamaProvider) GenerateCompletion(ctx context.Context, req CompletionRequest) (string, error) {
	url := fmt.Sprintf("%s/api/chat", strings.TrimRight(p.baseURL, "/"))

	type msg struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	messages := []msg{
		{Role: "system", Content: req.SystemPrompt},
		{Role: "user", Content: req.UserPrompt},
	}

	payload := map[string]interface{}{
		"model":    p.model,
		"messages": messages,
		"stream":   false,
		"options": map[string]interface{}{
			"temperature": req.Temperature,
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading ollama response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}

	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return "", fmt.Errorf("parsing ollama response: %w", err)
	}

	return parsed.Message.Content, nil
}
