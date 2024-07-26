package weather

import (
	"context"
	"database/sql"
	"errors"
	db "goirc/db/model"
	"goirc/model"
	"strings"
)

func weatherQueryByNick(ctx context.Context, q string, nick string) (string, error) {
	queries := db.New(model.DB)

	if q == "" {
		last, err := queries.LastNickWeatherRequest(ctx, nick)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return "", errors.New("no previous weather station to report on")
			}
			return "", err
		}
		if nick == last.Nick {
			if strings.HasPrefix(last.City, q) {
				q = last.City + "," + last.Country
			}
		}
	} else {
		last, err := queries.LastWeatherRequestByPrefix(ctx, sql.NullString{String: q, Valid: true})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return "", err
			}
		}
		if last.ID != 0 {
			q = last.City + "," + last.Country
		}
	}
	return q, nil
}
