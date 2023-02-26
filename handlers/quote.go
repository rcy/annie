package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/twitter"
	"log"
	"regexp"
)

func Quote(params bot.HandlerParams) bool {
	// match anything that starts with a quote and has no subsequent quotes
	re := regexp.MustCompile(`^("[^"]+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "not your personal secretary")
			return false
		}

		text := string(matches[1])

		err := notes.Create(params.Target, params.Nick, "quote", text)
		if err != nil {
			log.Print(err)
		}

		twitter.Post(text)

		return true
	}
	return false
}
