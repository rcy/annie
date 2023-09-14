package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func Quote(params bot.HandlerParams) error {
	text := params.Matches[1]

	_, err := notes.Create(notes.CreateParams{Target: params.Target, Nick: params.Nick, Kind: "quote", Text: text})
	if err != nil {
		return err
	}

	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "stored quote to share later, maybe")
	}

	return nil
}
