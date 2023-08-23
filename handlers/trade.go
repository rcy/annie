package handlers

import (
	"goirc/bot"
	"goirc/trader"
)

func Trade(params bot.HandlerParams) bool {
	reply, err := trader.Trade(params.Nick, params.Matches[1])
	if err != nil {
		params.Privmsgf(params.Target, "error: %s", err)
		return true
	}

	if reply != "" {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
		return true
	}

	return false
}
