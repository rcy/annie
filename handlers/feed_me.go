package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
)

func FeedMe(params bot.HandlerParams) bool {
	rows := []notes.Note{}

	err := model.DB.Select(&rows, `select id, created_at, nick, text, kind from notes order by random() limit 1`)
	if err != nil {
		params.Privmsgf(params.Target, "%v", err)
	} else if len(rows) >= 1 {
		note := rows[0]
		params.Privmsgf(params.Target, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
		err = notes.MarkAsSeen(note.Id, params.Target)
		if err != nil {
			params.Privmsgf(params.Target, err.Error())
		}
	}
	return true
}
