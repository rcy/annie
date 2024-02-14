package linkpool

import (
	"context"
	"fmt"
	"goirc/db/model"
	db "goirc/model"
	"testing"
	"time"
)

func TestNotes(t *testing.T) {
	ctx := context.Background()
	q := model.New(db.DB)

	db.DB.Exec("delete from notes")

	pool := New(q, 0, 0*time.Hour)

	for _, tc := range []struct {
		numNotes       int
		wantFirstNote  string
		wantRandomNote string
	}{
		{
			numNotes:      0,
			wantFirstNote: "",
		},
		{
			numNotes:      1,
			wantFirstNote: "hello note 0",
		},
		{
			numNotes:      100,
			wantFirstNote: "hello note 0",
		},
	} {
		for i := 0; i < tc.numNotes; i++ {
			_, err := pool.PushNote(ctx, PushNoteParams{
				Target: "nick",
				Nick:   "nick",
				Kind:   "note",
				Text:   fmt.Sprintf("hello note %d", i),
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		notes, err := pool.Notes(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if len(notes) < tc.numNotes {
			t.Errorf("fail less than 1 %d", len(notes))
		}

		if len(notes) == tc.numNotes && tc.numNotes != 0 {
			got := notes[0].Text.String
			if got != tc.wantFirstNote {
				t.Errorf("expected %s, got %s", tc.wantFirstNote, got)
			}
		}
	}
}

func TestPeekRandomNote(t *testing.T) {
	ctx := context.Background()

	q := model.New(db.DB)
	pool := New(q, 0, 0*time.Hour)

	db.DB.Exec("delete from notes")

	for _, tc := range []struct {
		seed      int64
		numNotes  int
		wantNote  string
		wantError bool
	}{
		{
			numNotes:  0,
			wantNote:  "",
			wantError: true,
		},
		{
			numNotes: 1,
			wantNote: "hello note 0",
		},
		{
			numNotes: 2,
			wantNote: "hello note 0",
		},
		{
			seed:     29,
			numNotes: 100,
			wantNote: "hello note 82",
		},
	} {
		t.Run(tc.wantNote, func(t *testing.T) {
			pool.Seed(tc.seed)
			for i := 0; i < tc.numNotes; i++ {
				_, err := pool.PushNote(ctx, PushNoteParams{
					Target: "nick",
					Nick:   "nick",
					Kind:   "note",
					Text:   fmt.Sprintf("hello note %d", i),
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			note, err := pool.PeekRandomNote(ctx)
			if err != nil {
				if !tc.wantError {
					t.Errorf("want error=%v, got error=%s", tc.wantError, err)
				}
			}
			got := note.Text.String
			if got != tc.wantNote {
				t.Errorf("expected %s, got %s", tc.wantNote, got)
			}
		})
	}
}

func TestPopRandomNote(t *testing.T) {
	ctx := context.Background()
	q := model.New(db.DB)

	db.DB.Exec("delete from notes")

	pool := New(q, 0, 0*time.Hour)

	for _, tc := range []struct {
		seed      int64
		numNotes  int
		wantNote  string
		wantError bool
	}{
		{
			numNotes:  0,
			wantNote:  "",
			wantError: true,
		},
		{
			numNotes: 1,
			wantNote: "hello note 0",
		},
		{
			numNotes: 2,
			wantNote: "hello note 0",
		},
		{
			seed:     29,
			numNotes: 100,
			wantNote: "hello note 64",
		},
	} {
		t.Run(tc.wantNote, func(t *testing.T) {
			pool.Seed(tc.seed)
			for i := 0; i < tc.numNotes; i++ {
				_, err := pool.PushNote(ctx, PushNoteParams{
					Target: "nick",
					Nick:   "nick",
					Kind:   "note",
					Text:   fmt.Sprintf("hello note %d", i),
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			note, err := pool.PopRandomNote(ctx, "#chann")
			if (err != nil) != tc.wantError {
				t.Errorf("want error=%v, got error=%s", tc.wantError, err)
			}
			got := note.Text.String
			if got != tc.wantNote {
				t.Errorf("expected %s, got %s", tc.wantNote, got)
			}
		})
	}
}
