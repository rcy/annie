package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func CreateNote(params bot.HandlerParams) error {
	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "not your personal secretary")
		return nil
	}

	text := params.Matches[1]

	err := notes.Create(params.Target, params.Nick, "note", text)
	if err != nil {
		return err
	} else {
		params.Privmsgf(params.Target, "recorded note")
	}
	return nil
}
