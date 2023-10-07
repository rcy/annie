package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
	"time"
)

func candidateLinks(age time.Duration) ([]notes.Note, error) {
	notes := []notes.Note{}

	query := `select id, created_at, nick, text, kind from notes where created_at <= datetime(?) and nick = target order by random() limit 69`

	t := time.Now().UTC().Add(-age).Format(time.RFC3339)
	err := model.DB.Select(&notes, query, t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return notes, nil
}

const (
	MINAGE    = time.Hour * 24
	THRESHOLD = 5
)

func FeedMe(params bot.HandlerParams) error {
	notes, err := candidateLinks(MINAGE)
	if err != nil {
		return err
	}

	if len(notes) < THRESHOLD {
		params.Privmsgf(params.Target, "not enough links to feed the channel")
		return nil
	}

	note := notes[0]

	_, err = model.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.Id)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (from ??? %s ago)", note.Text, util.Since(note.CreatedAt))

	return nil
}

func PipeHealth(params bot.HandlerParams) error {
	readyNotes, err := candidateLinks(MINAGE)
	if err != nil {
		return err
	}
	totalNotes, err := candidateLinks(0)
	if err != nil {
		return err
	}

	ready := len(readyNotes)
	fermenting := len(totalNotes) - ready

	params.Privmsgf(params.Target, "%d links ready to serve (%d fermenting)", ready, fermenting)

	return nil
}
