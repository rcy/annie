package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"log"
	"regexp"
)

func CreateNote(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^,(.+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "not your personal secretary")
			return false
		}

		text := string(matches[1])

		err := notes.Create(params.Target, params.Nick, "note", text)
		if err != nil {
			log.Print(err)
			params.Privmsgf(params.Target, err.Error())
		} else {
			params.Privmsgf(params.Target, "recorded note")
		}
		return true
	}
	return false
}
