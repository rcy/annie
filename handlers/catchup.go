package handlers

import (
	"context"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"goirc/util"
	"time"
)

func Catchup(params bot.HandlerParams) error {
	ctx := context.TODO()
	q := model.New(db.DB)

	notes, err := q.LastDaysNotes(ctx)
	if err != nil {
		return err
	}
	if len(notes) >= 1 {
		for _, note := range notes {
			params.Privmsgf(params.Nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Ago(time.Since(note.CreatedAt).Round(time.Second)))
			time.Sleep(1 * time.Second)
		}
	}
	params.Privmsgf(params.Nick, "--- %d total from last 24 hours", len(notes))
	return nil
}
