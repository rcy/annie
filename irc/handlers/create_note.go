package handlers

import (
	"log"
	"regexp"
)

func createNote(params Params) bool {
	re := regexp.MustCompile(`^,(.+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Irccon.Privmsg(params.Target, "not your personal secretary")
			return false
		}

		text := string(matches[1])

		err := insertNote(params.Db, params.Target, params.Nick, "note", text)
		if err != nil {
			log.Print(err)
			params.Irccon.Privmsg(params.Target, err.Error())
		} else {
			params.Irccon.Privmsg(params.Target, "recorded note")
		}
		return true
	}
	return false
}
