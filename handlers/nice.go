package handlers

import (
	"fmt"
	"goirc/bot"
	"regexp"
	"time"
)

func Nice(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`\b69\b`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	go func() {
		time.Sleep(10 * time.Second)
		params.Privmsgf(params.Target, fmt.Sprintf("%s: nice", params.Nick))
	}()

	return true
}
