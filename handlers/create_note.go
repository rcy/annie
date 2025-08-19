package handlers

import (
	"context"
	"database/sql"
	"goirc/db/model"
	"goirc/internal/responder"
	db "goirc/model"
)

func CreateNote(params responder.Responder) error {
	q := model.New(db.DB)

	//text := strings.TrimSpace(params.Matches[1])
	text := params.Match(1)

	_, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target(),
		Nick:   sql.NullString{String: params.Nick(), Valid: true},
		Kind:   "note",
		Text:   sql.NullString{String: text, Valid: true},
	})
	if err != nil {
		return err
	}
	if params.Target() == params.Nick() {
		params.Privmsgf(params.Target(), "recorded note to share later, maybe")
	} else {
		params.Privmsgf(params.Target(), "recorded note")
	}
	return nil
}
