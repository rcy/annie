package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/twitter"
	"log"
)

func Link(params bot.HandlerParams) error {
	url := params.Matches[1]

	if params.Target == params.Nick {
		// posted in a private message
		params.Privmsgf(params.Target, "this will be shared later")
		return nil
	}

	// posted to channel

	err := notes.Create(params.Target, params.Nick, "link", url)
	if err != nil {
		return err
	}

	log.Printf("recorded url %s", url)

	err = twitter.Post(url)
	if err != nil {
		return err
	}

	return nil
}
