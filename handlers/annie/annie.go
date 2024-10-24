package annie

import (
	"context"
	"goirc/bot"
	"os"

	"github.com/sashabaranov/go-openai"
)

func Handle(params bot.HandlerParams) error {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.TODO()

	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a personal assistant named Annie.  You are very terse, but friendly, with dry humour.  Answer in once sentence always.  Sometimes your name will be mentioned in the third person, as you exist in a group chat with multiple humans.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: params.Msg,
				},
			},
		})
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s: %s", params.Nick, resp.Choices[0].Message.Content)

	return nil
}
