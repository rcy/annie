package handlers

import (
	"goirc/bot"
	"goirc/util"
	"regexp"
	"time"
)

func Worldcup(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`world.?cup`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	then, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
	if err != nil {
		params.Privmsgf(params.Target, "error: %v", err)
		return true
	}
	until := util.Ago(time.Until(then))
	params.Privmsgf(params.Target, "the world cup will start in %s", until)
	return true
}
