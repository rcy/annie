package handlers

import (
	"goirc/bot"
	"goirc/internal/epigrams"
)

func EpigramHandler(params bot.HandlerParams) error {
	params.Privmsgf(params.Target, "%s", epigrams.Random())
	return nil
}
