package handlers

import (
	"goirc/bot"
	"goirc/internal/mlb"
)

func MLBOddsHandler(params bot.HandlerParams) error {
	s, err := mlb.PlayoffOdds()
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "%s", s)
	return nil
}
