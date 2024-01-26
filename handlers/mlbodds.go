package handlers

import (
	"goirc/bot"
	"goirc/internal/mlb"
)

func MLBOddsHandler(params bot.HandlerParams) (string, error) {
	return mlb.PlayoffOdds()
}
