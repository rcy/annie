package day

import (
	"context"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

func mustGetenv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

var openaiAPIKey = mustGetenv("OPENAI_API_KEY")

func generateImage(ctx context.Context, prompt string) (string, error) {
	client := openai.NewClient(openaiAPIKey)

	req := openai.ImageRequest{
		Prompt:         prompt,
		Model:          openai.CreateImageModelDallE3,
		N:              1,
		Size:           "1024x1024",
		ResponseFormat: "url",
	}

	resp, err := client.CreateImage(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Data[0].URL, nil
}
