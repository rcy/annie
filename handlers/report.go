package handlers

import (
	"goirc/bot"
	"goirc/trader"
)

func Report(params bot.HandlerParams) error {
	reply, err := trader.Report(params.Nick)
	if err != nil {
		return err
	}

	if reply != "" {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
	}

	return nil
}
