package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/fetch"
	db "goirc/model"
	"log"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
)

type queries interface {
	AddOpenGraphData(context.Context, model.AddOpenGraphDataParams) (model.Opengraph, error)
}

func storeOpenGraphData(q queries, url string) error {
	var og = opengraph.NewOpenGraph()

	code, body, err := fetch.Get(url, time.Minute)
	if err != nil {
		return err
	}
	if code >= 300 {
		return fmt.Errorf("got status code: %d", code)
	}
	err = og.ProcessHTML(bytes.NewReader(body))
	json, err := og.ToJSON()
	if err != nil {
		return err
	}
	_, err = q.AddOpenGraphData(context.TODO(), model.AddOpenGraphDataParams{Url: url, Data: json})
	if err != nil {
		return err
	}
	return nil
}

func Link(params bot.HandlerParams) error {
	q := model.New(db.DB)

	url := params.Matches[1]

	go func() {
		err := storeOpenGraphData(q, url)
		if err != nil {
			log.Printf("error: storeOpenGraphData: %s: %s", url, err)
		}
	}()

	note, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   "link",
		Text:   sql.NullString{String: url, Valid: true},
	})
	if err != nil {
		return err
	}

	if params.Target == params.Nick {
		// posted in a private message
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
