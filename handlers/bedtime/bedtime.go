package bedtime

import (
	"context"
	"goirc/bot"
	"goirc/model"
)

func Handle(ctx context.Context, params bot.HandlerParams) error {
	_, err := model.DB.Exec(`insert into bedtimes(nick, message) values(?, ?)`, params.Nick, params.Msg)
	return err
}
