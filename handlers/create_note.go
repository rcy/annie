package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func CreateNote(params bot.HandlerParams) error {
	text := params.Matches[1]

	err := notes.Create(params.Target, params.Nick, "note", text)
	if err != nil {
		return err
	} else {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "recorded note to share later, maybe")
		} else {
			params.Privmsgf(params.Target, "recorded note")
		}
	}
	return nil
}
