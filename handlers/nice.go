package handlers

import (
	"fmt"
	"goirc/bot"
	"time"
)

func Nice(params bot.HandlerParams) bool {
	go func() {
		time.Sleep(10 * time.Second)
		params.Privmsgf(params.Target, fmt.Sprintf("%s: nice", params.Nick))
	}()

	return true
}
