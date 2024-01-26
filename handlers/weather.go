package handlers

import (
	"context"
	"goirc/bot"
	"goirc/internal/weather"
)

func WeatherHandler(params bot.HandlerParams) (string, error) {
	var q string
	if len(params.Matches) > 1 {
		q = params.Matches[1]
	}

	w, err := weather.Fetch(context.TODO(), q, params.Nick)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}
