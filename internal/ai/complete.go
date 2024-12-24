package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func Complete(ctx context.Context, model string, system string, message string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: system,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		})
	if err != nil {
		if strings.Contains(err.Error(), "billing") {
			return "", fmt.Errorf("I need money: https://ko-fi.com/rcyemb")
		}

		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
