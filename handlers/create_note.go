package handlers

import (
	"context"
	"database/sql"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
)

func CreateNote(params bot.HandlerParams) (string, error) {
	q := model.New(db.DB)

	text := params.Matches[1]

	_, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   "note",
		Text:   sql.NullString{String: text, Valid: true},
	})
	if err != nil {
		return "", err
	}
	if params.Target == params.Nick {
		return "recorded note to share later, maybe", nil
	}
	return "recorded note", nil
}
