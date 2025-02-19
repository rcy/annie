package handlers

import (
	"context"
	"goirc/bot"
	"time"
)

func Nice(ctx context.Context, params bot.HandlerParams) error {
	go func() {
		time.Sleep(10 * time.Second)
		params.Privmsgf(params.Target, "%s: nice", params.Nick)
	}()

	return nil
}
