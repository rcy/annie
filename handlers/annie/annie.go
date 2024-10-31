package annie

import (
	"context"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

func Handle(params bot.HandlerParams) error {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.TODO()

	q := model.New(db.DB.DB)

	notes, err := q.Notes(ctx)
	if err != nil {
		return err
	}

	lines := make([]string, len(notes))
	for i, n := range notes {
		lines[i] = fmt.Sprintf("%s <%s> %s", n.CreatedAt, n.Nick.String, n.Text.String)
	}

	systemPrompt := "You are annie, a friend hanging out in an irc channel. Respond with single sentences, in lower case, with no punctuation. You are generally knowledgable but give special importance to the facts you learned from this chat history:"
	systemPrompt += strings.Join(lines, "\n")

	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
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
