package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/durfmt"
	"goirc/model"
	"goirc/model/notes"
	"goirc/util"
	"time"
)

func candidateLinks(age time.Duration) ([]notes.Note, error) {
	notes := []notes.Note{}

	query := `select id, created_at, nick, text, kind from notes where created_at <= datetime(?) and nick = target order by random() limit 420`

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
	COOLOFF   = time.Hour * 5
)

var lastSentAt = time.Unix(0, 0)

func FeedMe(params bot.HandlerParams) error {
	if time.Now().Sub(lastSentAt) < COOLOFF {
		if params.Nick != "" {
			params.Privmsgf(params.Target, "throttled until %s", durfmt.Format(time.Now().Sub(lastSentAt.Add(COOLOFF))))
		}
		return nil
	}

	notes, err := candidateLinks(MINAGE)
	if err != nil {
		return err
	}

	if len(notes) < THRESHOLD {
		if params.Nick != "" {
			params.Privmsgf(params.Target, "not enough links to feed the channel")
		}
		return nil
	}

	note := notes[0]

	_, err = model.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.Id)
	if err != nil {
		return err
	}

	ready, fermenting, err := health()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (%s ago) [pipe=%d+%d]", note.Text, util.Since(note.CreatedAt), ready, fermenting)

	lastSentAt = time.Now()

	return nil
}

func PipeHealth(params bot.HandlerParams) error {
	ready, fermenting, err := health()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%d links ready to serve (%d fermenting)", ready, fermenting)

	return nil
}

// Return ready, fermenting, error
func health() (int, int, error) {
	readyNotes, err := candidateLinks(MINAGE)
	if err != nil {
		return 0, 0, err
	}
	totalNotes, err := candidateLinks(0)
	if err != nil {
		return 0, 0, err
	}

	ready := len(readyNotes)
	fermenting := len(totalNotes) - ready

	return ready, fermenting, nil
}
