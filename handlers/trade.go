package handlers

import (
	"goirc/bot"
	"goirc/trader"
)

func Trade(params bot.HandlerParams) error {
	reply, err := trader.Trade(params.Nick, params.Matches[1])
	if err != nil {
		return err
	}

	if reply != "" {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
	}

	return nil
}
