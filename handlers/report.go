package handlers

import (
	"goirc/bot"
	"goirc/trader"
)

func Report(params bot.HandlerParams) bool {
	reply, err := trader.Report(params.Nick)
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
