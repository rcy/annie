package handlers

import (
	"context"
	"database/sql"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
)

func Quote(params bot.HandlerParams) error {
	q := model.New(db.DB)
	text := params.Matches[1]

	// posted to private channel
	isAnonymous := params.Target == params.Nick

	_, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   "quote",
		Text:   sql.NullString{String: text, Valid: true},
		Anon:   isAnonymous,
	})
	if err != nil {
		return err
	}

	if isAnonymous {
		_, err = q.ScheduleFutureMessage(context.TODO(), "quote")
		if err != nil {
			return err
		}

		params.Privmsgf(params.Target, "thanks for the quote")

		return nil
	}

	return nil
}
