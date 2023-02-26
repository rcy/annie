package handlers

import (
	"goirc/bot"
	"goirc/util"
	"regexp"
	"time"
)

func OpeningDay(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`(opening.?day)|(baseball)`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	then, err := time.Parse(time.RFC3339, "2023-03-30T13:05:00-05:00")
	if err != nil {
		params.Privmsgf(params.Target, "error: %v", err)
		return true
	}
	until := util.Ago(time.Until(then))
	params.Privmsgf(params.Target, "opening day is in %s", until)
	return true
}
