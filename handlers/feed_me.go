package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
	"time"
)

func FeedMe(params bot.HandlerParams) error {
	notes := []notes.Note{}

	query := `select id, created_at, nick, text, kind from notes where created_at <= ? and nick = target order by random() limit 5`

	t := time.Now().Add(-time.Hour * 24).UTC().Format("2006-01-02T15:04:05Z")
	err := model.DB.Select(&notes, query, t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	if len(notes) < 5 {
		params.Privmsgf(params.Target, "not enough links to feed the channel")
		return nil
	}

	note := notes[0]

	_, err = model.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.Id)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (from ??? %s ago)", note.Text, util.Since(note.CreatedAt))

	return nil
}
