package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
)

func FeedMe(params bot.HandlerParams) error {
	note := notes.Note{}

	err := model.DB.Get(&note, `select id, created_at, nick, text, kind from notes where nick = target order by random() limit 1`)
	if err != nil {
		if err == sql.ErrNoRows {
			params.Privmsgf(params.Target, "no")
			return nil
		}
		return err
	}
	_, err = model.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.Id)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (from ??? %s ago)", note.Text, util.Since(note.CreatedAt))

	return nil
}
