package handlers

import (
	"fmt"
	"goirc/bot"
	"time"
)

func Nice(params bot.HandlerParams) (string, error) {
	time.Sleep(10 * time.Second)
	return fmt.Sprintf("%s: nice", params.Nick), nil
}
