package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
)

func FeedMe(params bot.HandlerParams) error {
	rows := []notes.Note{}

	err := model.DB.Select(&rows, `select id, created_at, nick, text, kind from notes order by random() limit 1`)
	if err != nil {
		return err
	}
	return nil
}
