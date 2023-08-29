package handlers

import (
	"goirc/bot"
	"time"
)

func Nice(params bot.HandlerParams) error {
	go func() {
		time.Sleep(10 * time.Second)
		params.Privmsgf(params.Target, "%s: nice", params.Nick)
	}()

	return nil
}
