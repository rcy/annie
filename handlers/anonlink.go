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
	threshold = 1
)

func AnonLink(params bot.HandlerParams) error {
	q := model.New(db.DB)
	pool := linkpool.New(q, threshold, minAge)
	note, err := pool.PopRandomNote(context.Background(), params.Target, "link")
	if err != nil {
		return err
	}

	text, err := note.Link()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, text)
	return nil
}

func AnonQuote(params bot.HandlerParams) error {
	q := model.New(db.DB)
	pool := linkpool.New(q, threshold, minAge)
	note, err := pool.PopRandomNote(context.Background(), params.Target, "quote")
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, note.Text.String)
	return nil
}

func AnonStatus(params bot.HandlerParams) error {
	ctx := context.TODO()
	q := model.New(db.DB)
	allPool := linkpool.New(q, 0, 0)
	allLinks, err := allPool.Notes(ctx, "link")
	if err != nil {
		return err
	}
	allQuotes, err := allPool.Notes(ctx, "quote")
	if err != nil {
		return err
	}

	dayPool := linkpool.New(q, 0, minAge)
	dayLinks, err := dayPool.Notes(ctx, "link")
	if err != nil {
		return err
	}
	dayQuotes, err := dayPool.Notes(ctx, "quote")
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "links=%d+%d quotes=%d+%d",
		len(dayLinks), len(allLinks)-len(dayLinks),
		len(dayQuotes), len(allQuotes)-len(dayQuotes),
	)

	return nil
}
