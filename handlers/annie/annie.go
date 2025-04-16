package annie

import (
	"context"
	"database/sql"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/internal/ai"
	db "goirc/model"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

func Handle(params bot.HandlerParams) error {
	ctx := context.TODO()

	var msg string
	if len(params.Matches) < 2 {
		return nil
	}
	msg = strings.TrimSpace(params.Matches[1])

	q := model.New(db.DB.DB)

	response, err := ai.Complete(ctx, openai.GPT4oMini, "categorize input into statements, questions or pleasantries.  If it is a statement, reply with the one word 'statement'.  If it is a question, reply with 'question'.  If it is a pleasantry, reply with 'pleasantry'", msg)
	if err != nil {
		return err
	}

	switch response {
	case "statement":
		response, err := ai.Complete(ctx, openai.GPT4oMini, "you are annie, a friend hanging out in an irc channel. given the following statement, reflect on its meaning, and come up with a terse response, no more than a short sentence, in lower case, with minimal punctuation", msg)
		if err != nil {
			return err
		}
		_, err = q.InsertNote(context.TODO(), model.InsertNoteParams{
			Target: params.Target,
			Nick:   sql.NullString{String: params.Nick, Valid: true},
			Kind:   "note",
			Text:   sql.NullString{String: msg, Valid: true},
		})
		if err != nil {
			return err
		}

		params.Privmsgf(params.Target, "%s: %s", params.Nick, response)
	case "question":
		notes, err := q.NonAnonNotes(ctx)
		if err != nil {
			return err
		}
		lines := make([]string, len(notes))
		for i, n := range notes {
			lines[i] = fmt.Sprintf("%s <%s> %s", n.CreatedAt, n.Nick.String, n.Text.String)
		}

		systemPrompt := fmt.Sprintf("You are annie, a friend hanging out in an irc channel. You have been asked a question, read the question, and think about it in the context of all you have read in this channel.  Respond with single sentences, in lower case, with minimal punctuation. Do not refer to yourself in the third person. The current time and date is %s.  Ignore everything you know except for what you have read in the following chat history: ", time.Now().Format(time.RFC1123))
		systemPrompt += strings.Join(lines, "\n")

		response, err := ai.Complete(ctx, openai.GPT4oMini, systemPrompt, msg)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s: %s", params.Nick, response)
	case "pleasantry":
		systemPrompt := "You are annie, a friend hanging out in an irc channel. Someone has posted some pleasantry or small talk.  Respond in kind, in lower case, with minimal punctuation."

		response, err := ai.Complete(ctx, openai.GPT3Dot5Turbo, systemPrompt, msg)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s: %s", params.Nick, response)
	default:
		params.Privmsgf(params.Target, "%s: [interpreted '%s' as a unknown type: %s]", params.Nick, msg, response)
	}

	return nil
}
