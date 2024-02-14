package handlers

import (
	"context"
	"goirc/bot"
	"goirc/db/model"
	"goirc/handlers/linkpool"
	db "goirc/model"
	"time"
)

const (
	minAge    = time.Hour * 24
	threshold = 25
)

func AnonLink(params bot.HandlerParams) error {
	q := model.New(db.DB)
	pool := linkpool.New(q, threshold, minAge)
	note, err := pool.PopRandomNote(context.Background(), params.Target)
	if err != nil {
		return err
	}

	var text string
	if note.Kind == "link" {
		text, err = note.Link()
		if err != nil {
			return err
		}
	} else {
		text = note.Text.String
	}

	params.Privmsgf(params.Target, text)
	return nil
}

func AnonStatus(params bot.HandlerParams) error {
	ctx := context.TODO()
	q := model.New(db.DB)
	allPool := linkpool.New(q, 0, 0)
	allNotes, err := allPool.Notes(ctx)
	if err != nil {
		return err
	}
	dayPool := linkpool.New(q, 0, minAge)
	dayNotes, err := dayPool.Notes(ctx)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Target, "ready=%d fermenting=%d", len(dayNotes), len(allNotes)-len(dayNotes))
	return nil
}

func rot13(s string) string {
	rot := ""
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			r = 'a' + (r-'a'+13)%26
		} else if r >= 'A' && r <= 'Z' {
			r = 'A' + (r-'A'+13)%26
		}
		rot += string(r)
	}
	return rot
}
