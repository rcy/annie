package handlers

import (
	"goirc/bot"
	"goirc/util"
	"time"
)

func Worldcup(params bot.HandlerParams) bool {
	then, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
	if err != nil {
		params.Privmsgf(params.Target, "error: %v", err)
		return true
	}
	until := util.Ago(time.Until(then))
	params.Privmsgf(params.Target, "the world cup will start in %s", until)
	return true
}
