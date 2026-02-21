package llm

import (
	"context"
	"fmt"
	"os"

	"github.com/andybarilla/skeeter/internal/config"
)

type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

type Provider interface {
	Complete(ctx context.Context, messages []Message) (string, error)
}

func NewProvider(cfg config.LLMConfig) (Provider, error) {
	apiKey := resolveAPIKey(cfg)

	switch cfg.Provider {
	case "anthropic":
		model := cfg.Model
		if model == "" {
			model = "claude-sonnet-4-20250514"
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.anthropic.com"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("API key required for Anthropic (set llm.api_key, SKEETER_LLM_API_KEY, or ANTHROPIC_API_KEY)")
		}
		return &AnthropicProvider{APIKey: apiKey, Model: model, BaseURL: baseURL}, nil

	case "openai":
		model := cfg.Model
		if model == "" {
			model = "gpt-4o"
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}
		if apiKey == "" {
			return nil, fmt.Errorf("API key required for OpenAI (set llm.api_key, SKEETER_LLM_API_KEY, or OPENAI_API_KEY)")
		}
		return &OpenAIProvider{APIKey: apiKey, Model: model, BaseURL: baseURL}, nil

	case "ollama":
		model := cfg.Model
		if model == "" {
			model = "llama3"
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:11434/v1"
		}
		return &OpenAIProvider{APIKey: apiKey, Model: model, BaseURL: baseURL}, nil

	case "lmstudio":
		model := cfg.Model
		if model == "" {
			model = "default"
		}
		baseURL := cfg.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:1234/v1"
		}
		return &OpenAIProvider{APIKey: apiKey, Model: model, BaseURL: baseURL}, nil

	case "":
		return nil, fmt.Errorf("no LLM provider configured (run: skeeter config set llm.provider anthropic)")

	default:
		return nil, fmt.Errorf("unknown LLM provider %q (supported: anthropic, openai, ollama, lmstudio)", cfg.Provider)
	}
}

func resolveAPIKey(cfg config.LLMConfig) string {
	if cfg.APIKey != "" {
		return cfg.APIKey
	}
	if key := os.Getenv("SKEETER_LLM_API_KEY"); key != "" {
		return key
	}
	switch cfg.Provider {
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}
