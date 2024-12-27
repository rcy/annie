package handlers

import (
	"context"
	"goirc/bot"
	"goirc/db/model"
	"goirc/handlers/linkpool"
	"goirc/image"
	db "goirc/model"
)

const (
	// how old a message must be before it is delivered
	minAge = 0 //7 * time.Hour * 24

	// how long to wait after an anon message is posted to send one from the queue
	FutureMessageInterval = "+1 hour"
)

func AnonLink(params bot.HandlerParams) error {
	q := model.New(db.DB)
	pool := linkpool.New(q, minAge)
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
	pool := linkpool.New(q, minAge)
	note, err := pool.PopRandomNote(context.Background(), params.Target, "quote")
	if err != nil {
		return err
	}

	img, err := image.GenerateDALLE(context.TODO(), note.Text.String)
	if err != nil {
		params.Privmsgf(params.Target, "%s", note.Text.String)
		return err
	}

	params.Privmsgf(params.Target, "%s %s", note.Text.String, img.URL())

	return nil
}

func AnonStatus(params bot.HandlerParams) error {
	ctx := context.TODO()
	q := model.New(db.DB)
	allPool := linkpool.New(q, 0)
	allLinks, err := allPool.Notes(ctx, "link")
	if err != nil {
		return err
	}
	allQuotes, err := allPool.Notes(ctx, "quote")
	if err != nil {
		return err
	}

	dayPool := linkpool.New(q, minAge)
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
