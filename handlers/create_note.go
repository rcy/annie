package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func CreateNote(params bot.HandlerParams) error {
	text := params.Matches[1]

	_, err := notes.Create(notes.CreateParams{Target: params.Target, Nick: params.Nick, Kind: "note", Text: text})
	if err != nil {
		return err
	}
	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "recorded note to share later, maybe")
	} else {
		params.Privmsgf(params.Target, "recorded note")
	}
	return nil
}
