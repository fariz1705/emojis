package main

import (
	"context"
	"fmt"

	"github.com/ServiceWeaver/weaver"
	openai "github.com/sashabaranov/go-openai"
)
type ChatGPT interface {
	Complete(ctx context.Context, prompt string) (string, error)
}
type chatgpt struct {
	weaver.Implements[ChatGPT]
	weaver.WithConfig[config]
}
type config struct {
	APIKey string `toml:"api_key"`
}

func (gpt *chatgpt) Complete(ctx context.Context, prompt string) (string, error) {

	if gpt.Config().APIKey == "" {
		return "", fmt.Errorf("ChatGPT api_key not provided")
	}

	client := openai.NewClient(gpt.Config().APIKey)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
	}
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ChatGPT completion error: %w", err)
	}

	return resp.Choices[0].Message.Content, nil
}