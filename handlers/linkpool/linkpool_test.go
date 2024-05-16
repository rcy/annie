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

	_, err := db.DB.Exec("delete from notes")
	if err != nil {
		t.Fatal(err)
	}

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
				Kind:   "link",
				Text:   fmt.Sprintf("hello note %d", i),
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		notes, err := pool.Notes(ctx, "link")
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

	_, err := db.DB.Exec("delete from notes")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		seed      int64
		numLinks  int
		numQuotes int
		wantNote  string
		wantError bool
	}{
		{
			numLinks:  0,
			wantNote:  "",
			wantError: true,
		},
		{
			numLinks: 1,
			wantNote: "hello note 0",
		},
		{
			numLinks: 2,
			wantNote: "hello note 0",
		},
		{
			seed:     29,
			numLinks: 100,
			wantNote: "hello note 82",
		},
	} {
		t.Run(tc.wantNote, func(t *testing.T) {
			pool.Seed(tc.seed)
			for i := 0; i < tc.numLinks; i++ {
				_, err := pool.PushNote(ctx, PushNoteParams{
					Target: "nick",
					Nick:   "nick",
					Kind:   "link",
					Text:   fmt.Sprintf("hello note %d", i),
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			note, err := pool.PeekRandomNote(ctx, "link")
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

	for _, tc := range []struct {
		name      string
		seed      int64
		numLinks  int
		numQuotes int
		noteKind  string
		wantNote  string
		wantError bool
	}{
		{
			name:      "no notes",
			wantNote:  "",
			wantError: true,
		},
		{
			name:     "one link",
			numLinks: 1,
			noteKind: "link",
			wantNote: "link 0",
		},
		{
			name:      "one quote",
			numQuotes: 1,
			noteKind:  "quote",
			wantNote:  "quote 0",
		},
		{
			name:     "two links",
			numLinks: 2,
			noteKind: "link",
			wantNote: "link 0",
		},
		{
			name:     "100 links",
			numLinks: 100,
			noteKind: "link",
			wantNote: "link 74",
		},
		{
			name:      "100 quotes",
			numQuotes: 100,
			noteKind:  "quote",
			wantNote:  "quote 74",
		},
		{
			name:      "1 quote, 100 links",
			numQuotes: 1,
			numLinks:  100,
			noteKind:  "quote",
			wantNote:  "quote 0",
		},
		{
			name:      "1 link, 100 quotes",
			numLinks:  1,
			numQuotes: 100,
			noteKind:  "link",
			wantNote:  "link 0",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := db.DB.Exec("delete from notes")
			if err != nil {
				t.Fatal(err)
			}

			pool := New(q, 0, 0*time.Hour)

			pool.Seed(tc.seed)
			for i := 0; i < tc.numLinks; i++ {
				_, err := pool.PushNote(ctx, PushNoteParams{
					Target: "nick",
					Nick:   "nick",
					Kind:   "link",
					Text:   fmt.Sprintf("link %d", i),
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			for i := 0; i < tc.numQuotes; i++ {
				_, err := pool.PushNote(ctx, PushNoteParams{
					Target: "nick",
					Nick:   "nick",
					Kind:   "quote",
					Text:   fmt.Sprintf("quote %d", i),
				})
				if err != nil {
					t.Fatal(err)
				}
			}

			note, err := pool.PopRandomNote(ctx, "#chann", tc.noteKind)
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
