package provider

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type AnthropicProvider struct {
	client *anthropic.Client
	model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	client := anthropic.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &AnthropicProvider{
		client: client,
		model:  model,
	}
}

func (p *AnthropicProvider) Complete(prompt, systemPrompt string, useJSON bool) (string, error) {
	var systemMessages []anthropic.TextBlockParam
	if systemPrompt != "" {
		systemMessages = []anthropic.TextBlockParam{
			{Text: systemPrompt},
		}
	}

	req := anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		System:    systemMessages,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
		MaxTokens: 4096,
	}

	resp, err := p.client.Messages.New(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("Anthropic completion error: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", fmt.Errorf("Anthropic returned no content")
	}

	return resp.Content[0].Text, nil
}

