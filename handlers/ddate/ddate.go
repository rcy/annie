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
)

func Handle(params bot.HandlerParams) error {
	tz, err := getNickTimezone(context.TODO(), params.Nick)
	if err != nil {
		return fmt.Errorf("getNickTimezone: %w", err)
	}

	str, err := ddate.NowInZone(tz)
	if err != nil {
		return fmt.Errorf("ddate.InZone: %w", err)
	}

	params.Privmsgf(params.Target, "%s", str)

	return nil
}

func getNickTimezone(ctx context.Context, nick string) (string, error) {
	q := model.New(db.DB.DB)
	nickTimezone, err := q.GetNickTimezone(ctx, nick)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "America/Los_Angeles", nil
		}
		return "", fmt.Errorf("GetNickTimezone: %w", err)
	}
	return nickTimezone.Tz, nil
}
