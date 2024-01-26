package handlers

import (
	"context"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"goirc/util"
	"strings"
	"time"
)

func Catchup(params bot.HandlerParams) (string, error) {
	ctx := context.TODO()
	q := model.New(db.DB)

	notes, err := q.LastDaysNotes(ctx)
	if err != nil {
		return "", err
	}
	lines := []string{}
	if len(notes) >= 1 {
		for _, note := range notes {
			lines = append(lines, fmt.Sprintf("%s (from %s %s ago)", note.Text.String, note.Nick.String, util.Ago(time.Since(note.CreatedAt).Round(time.Second))))
		}
	}
	lines = append(lines, fmt.Sprintf("--- %d total from last 24 hours", len(notes)))

	return strings.Join(lines, "\n"), nil
}
