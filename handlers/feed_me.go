package handlers

import (
	"goirc/bot"
	"goirc/model/notes"
	"goirc/util"
	"regexp"
)

func FeedMe(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^!feedme`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	notes := []notes.Note{}

	err := params.Db.Select(&notes, `select id, created_at, nick, text, kind from notes order by random() limit 1`)
	if err != nil {
		params.Privmsgf(params.Target, "%v", err)
	} else if len(notes) >= 1 {
		note := notes[0]
		params.Privmsgf(params.Target, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
		err = markAsSeen(params.Db, note.Id, params.Target)
		if err != nil {
			params.Privmsgf(params.Target, err.Error())
		}
	}
	return true
}
