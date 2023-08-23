package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/twitter"
	"log"
)

func Quote(params bot.HandlerParams) bool {
	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "not your personal secretary")
		return false
	}

	text := string(params.Matches[1])

	err := notes.Create(params.Target, params.Nick, "quote", text)
	if err != nil {
		log.Print(err)
	}

	twitter.Post(text)

	return true
}
