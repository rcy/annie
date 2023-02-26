package handlers

import (
	"goirc/model"
	"goirc/util"
	"regexp"
	"time"
)

func catchup(params Params) bool {
	re := regexp.MustCompile(`^!catchup`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	notes := []model.Note{}

	// TODO: markAsSeen
	err := params.Db.Select(&notes, `select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc`)
	if err != nil {
		params.Irccon.Privmsgf(params.Target, "%v", err)
		return false
	}
	if len(notes) >= 1 {
		for _, note := range notes {
			params.Irccon.Privmsgf(params.Nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
	params.Irccon.Privmsgf(params.Nick, "--- %d total from last 24 hours", len(notes))
	return true
}
