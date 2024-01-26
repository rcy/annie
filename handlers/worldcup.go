package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/util"
	"time"
)

func Worldcup(params bot.HandlerParams) (string, error) {
	then, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
	if err != nil {
		return "", err
	}
	until := util.Ago(time.Until(then))

	return fmt.Sprintf("the world cup will start in %s", until), nil
}
