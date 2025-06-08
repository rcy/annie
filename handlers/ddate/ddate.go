package ddate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/internal/ddate"
	db "goirc/model"
	"time"
)

func Handle(params bot.HandlerParams) error {
	location, err := getNickLocation(context.TODO(), params.Nick)
	if err != nil {
		return fmt.Errorf("getNickLocation: %w", err)
	}
	now := time.Now().In(location)
	str := ddate.FromTime(now).Format(true)

	params.Privmsgf(params.Target, "%s", str)

	return nil
}

func getNickLocation(ctx context.Context, nick string) (*time.Location, error) {
	q := model.New(db.DB.DB)
	nickTimezone, err := q.GetNickTimezone(ctx, nick)
	var tz string
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tz = "America/Los_Angeles"
		} else {
			return nil, fmt.Errorf("GetNickTimezone: %w", err)
		}
	} else {
		tz = nickTimezone.Tz
	}
	location, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("LoadLocation: %w", err)
	}
	return location, nil
}
