package handlers

import (
	"goirc/bot"
	"regexp"
)

func Nice(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`\b69\b`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	params.Privmsgf(params.Target, "nice")
	return true
}
