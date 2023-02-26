package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
	"regexp"
	"time"
)

func Catchup(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^!catchup`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	notes := []notes.Note{}

	// TODO: markAsSeen
	err := model.DB.Select(&notes, `select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc`)
	if err != nil {
		params.Privmsgf(params.Target, "%v", err)
		return false
	}
	if len(notes) >= 1 {
		for _, note := range notes {
			params.Privmsgf(params.Nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
	params.Privmsgf(params.Nick, "--- %d total from last 24 hours", len(notes))
	return true
}
