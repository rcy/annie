package handlers

import (
	"context"
	"goirc/bot"
	"goirc/internal/weather"
)

func WeatherHandler(params bot.HandlerParams) error {
	var q string
	if len(params.Matches) > 1 {
		q = params.Matches[1]
	}

	w, err := weather.Fetch(context.TODO(), q, params.Nick)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", w.String())

	return nil
}
