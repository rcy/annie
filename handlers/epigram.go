package handlers

import (
	"goirc/bot"
	"goirc/internal/epigrams"
)

func EpigramHandler(bot.HandlerParams) (string, error) {
	return epigrams.Random(), nil
}
