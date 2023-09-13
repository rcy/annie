package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func Link(params bot.HandlerParams) error {
	url := params.Matches[1]

	err := notes.Create(params.Target, params.Nick, "link", url)
	if err != nil {
		return err
	}

	if params.Target == params.Nick {
		// posted in a private message
		params.Privmsgf(params.Target, "stored link to share later, maybe")
		return nil
	}

	return nil
}
