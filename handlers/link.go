package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
)

func Link(params bot.HandlerParams) error {
	url := params.Matches[1]

	note, err := notes.Create(notes.CreateParams{Target: params.Target, Nick: params.Nick, Kind: "link", Text: url})
	if err != nil {
		return err
	}

	if params.Target == params.Nick {
		// posted in a private message
		link, err := note.Link()
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s will be shared later, maybe", link)
		return nil
	}

	return nil
}
