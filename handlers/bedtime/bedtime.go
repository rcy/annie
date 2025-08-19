package bedtime

import (
	"goirc/internal/responder"
	"goirc/model"
)

func Handle(params responder.Responder) error {
	_, err := model.DB.Exec(`insert into bedtimes(nick, message) values(?, ?)`, params.Nick, params.Msg())
	return err
}
