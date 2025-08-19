package provider

import (
	"fmt"
	"strings"

	"github.com/autodevops/verifier-go/internal/config"
)

// LLMProvider is the interface for AI providers.
type LLMProvider interface {
	Complete(prompt string, systemPrompt string, useJSON bool) (string, error)
}

// ProviderFactory creates an instance of an LLMProvider.
func ProviderFactory(model string, cfg *config.Config) (LLMProvider, error) {
	if strings.HasPrefix(model, "gpt") {
		if cfg.Providers.OpenAI.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key is not configured")
		}
		return NewOpenAIProvider(cfg.Providers.OpenAI.APIKey, model), nil
	}
	if strings.HasPrefix(model, "claude") {
		if cfg.Providers.Anthropic.APIKey == "" {
			return nil, fmt.Errorf("Anthropic API key is not configured")
		}
		return NewAnthropicProvider(cfg.Providers.Anthropic.APIKey, model), nil
	}
	return nil, fmt.Errorf("unsupported model provider for model: %s", model)
}
