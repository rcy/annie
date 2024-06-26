package handlers

import (
	"context"
	"database/sql"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
)

func Link(params bot.HandlerParams) error {
	q := model.New(db.DB)

	url := params.Matches[1]

	// posted in a private message?
	isAnonymous := params.Target == params.Nick

	note, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   "link",
		Text:   sql.NullString{String: url, Valid: true},
		Anon:   isAnonymous,
	})
	if err != nil {
		return err
	}

	if isAnonymous {
		link, err := note.Link()
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s will be shared later, maybe", link)
		params.Publish("anonnoteposted", note)
		return nil
	}

	return nil
}
