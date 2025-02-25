package handlers

import (
	"context"
	"goirc/bot"
	"goirc/util"
	"time"
)

func Worldcup(ctx context.Context, params bot.HandlerParams) error {
	then, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
	if err != nil {
		return err
	}
	until := util.Ago(time.Until(then))
	params.Privmsgf(params.Target, "the world cup will start in %s", until)
	return nil
}
