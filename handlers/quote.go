package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/twitter"
)

func Quote(params bot.HandlerParams) error {
	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "not your personal secretary")
		return nil
	}

	text := params.Matches[1]

	err := notes.Create(params.Target, params.Nick, "quote", text)
	if err != nil {
		return err
	}

	err = twitter.Post(text)
	if err != nil {
		return err
	}

	return nil
}
