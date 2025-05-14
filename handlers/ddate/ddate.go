package ddate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"os"
	"os/exec"
	"strings"
)

func Handle(params bot.HandlerParams) error {
	tz, err := getNickTimezone(context.TODO(), params.Nick)
	if err != nil {
		return fmt.Errorf("getNickTimezone: %w", err)
	}

	cmd := exec.Command("ddate")
	cmd.Env = append(os.Environ(), "TZ="+tz)

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	str := out.String()

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
