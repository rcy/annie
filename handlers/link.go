package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/twitter"
	"log"
)

func Link(params bot.HandlerParams) bool {
	url := string(params.Matches[1])

	if params.Target == params.Nick {
		// posted in a private message
		params.Privmsgf(params.Target, "this will be shared later")
		return false
	}

	// posted to channel

	err := notes.Create(params.Target, params.Nick, "link", url)
	if err != nil {
		log.Print(err)
		params.Privmsgf(params.Target, err.Error())
	} else {
		log.Printf("recorded url %s", url)
	}

	twitter.Post(url)

	return true
}
