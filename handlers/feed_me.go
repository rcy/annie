package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"goirc/util"
	"time"

	"github.com/rcy/durfmt"
)

func candidateLinks(age time.Duration) ([]model.Note, error) {
	q := model.New(db.DB)

	t := time.Now().UTC().Add(-age)
	notes, err := q.UnsentAnonymousNotes(context.TODO(), t)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return notes, nil
}

const (
	MINAGE = time.Hour * 24
)

// these are vars so they can be changed in the tests
var (
	threshold = 40
	cooloff   = time.Hour
)

var lastSentAt = time.Unix(0, 0)

func canSendIn(startTime time.Time) time.Duration {
	return time.Until(startTime.Add(cooloff))
}

func FeedMe(params bot.HandlerParams) (string, error) {
	waitFor := canSendIn(lastSentAt)

	if waitFor > 0 {
		if params.Nick != "" {
			return fmt.Sprintf("throttled for another %s", durfmt.Format(waitFor)), nil
		}
		return "", nil
	}

	notes, err := candidateLinks(MINAGE)
	if err != nil {
		return "", err
	}

	if len(notes) < threshold {
		if params.Nick != "" {
			return "not enough links to feed the channel", nil
		}
		return "", nil
	}

	candidates := make([]string, len(notes))
	for i, n := range notes {
		candidates[i] = n.Text.String
	}

	note := notes[0]

	_, err = db.DB.Exec(`update notes set target = ? where id = ?`, params.Target, note.ID)
	if err != nil {
		return "", err
	}

	ready, fermenting, err := health()
	if err != nil {
		return "", err
	}

	var text string
	if note.Kind == "link" {
		text, err = note.Link()
		if err != nil {
			return "", err
		}
	} else {
		text = note.Text.String
	}

	lastSentAt = time.Now()

	return fmt.Sprintf("%s (%s ago) [pipe=%d+%d]", text, util.Ago(time.Since(note.CreatedAt)), ready, fermenting), nil
}

func PipeHealth(params bot.HandlerParams) (string, error) {
	ready, fermenting, err := health()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d links ready to serve (%d fermenting)", ready, fermenting), nil
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
