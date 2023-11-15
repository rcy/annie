package handlers

import (
	"goirc/bot"
	"goirc/durfmt"
	"time"
)

func TimeLeft(params bot.HandlerParams) error {
	params.Privmsgf(params.Target, durfmt.Format(time.Unix(2<<30, 0).Sub(time.Now())))

	return nil
}
