package bedtime

import (
	"goirc/bot"
	"goirc/model"
)

func Handle(params bot.HandlerParams) error {
	_, err := model.DB.Exec(`insert into bedtimes(nick, message) values(?, ?)`, params.Nick, params.Msg)
	return err
}
