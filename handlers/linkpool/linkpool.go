package linkpool

import (
	"context"
	"database/sql"
	"errors"
	"goirc/db/model"
	"math/rand"
	"time"
)

type pool struct {
	rnd       *rand.Rand
	threshold int
	minAge    time.Duration
	queries   queries
}

type queries interface {
	UnsentAnonymousNotes(context.Context, time.Time) ([]model.Note, error)
	MarkAnonymousNoteDelivered(context.Context, model.MarkAnonymousNoteDeliveredParams) (model.Note, error)
	InsertNote(context.Context, model.InsertNoteParams) (model.Note, error)
}

func New(queries *model.Queries, threshold int, minAge time.Duration) pool {
	return pool{
		queries:   queries,
		rnd:       rand.New(rand.NewSource(time.Now().UnixNano())),
		threshold: threshold,
		minAge:    minAge,
	}
}

func (p *pool) Seed(seed int64) {
	p.rnd = rand.New(rand.NewSource(seed))
}

func (p *pool) Notes(ctx context.Context) ([]model.Note, error) {
	olderThan := time.Now().UTC().Add(-p.minAge)
	notes, err := p.queries.UnsentAnonymousNotes(ctx, olderThan)
	if err != nil {
		return []model.Note{}, err
	}
	return notes, nil
}

func (p *pool) PeekRandomNote(ctx context.Context) (model.Note, error) {
	notes, err := p.Notes(ctx)
	if err != nil {
		return model.Note{}, err
	}
	if len(notes) <= p.threshold {
		return model.Note{}, errors.New("no note found")
	}
	r := p.rnd.Intn(len(notes))
	return notes[r], nil
}

func (p *pool) PopRandomNote(ctx context.Context, target string) (model.Note, error) {
	note, err := p.PeekRandomNote(ctx)
	if err != nil {
		return model.Note{}, err
	}
	note, err = p.queries.MarkAnonymousNoteDelivered(ctx, model.MarkAnonymousNoteDeliveredParams{
		ID:     note.ID,
		Target: target,
	})
	return note, nil
}

type PushNoteParams struct {
	Target string
	Nick   string
	Kind   string
	Text   string
}

func (p *pool) PushNote(ctx context.Context, params PushNoteParams) (model.Note, error) {
	return p.queries.InsertNote(ctx, model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   params.Kind,
		Text:   sql.NullString{String: params.Text, Valid: true},
	})
}
