package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"math"
	"testing"
	"time"
)

func TestFeedMe(t *testing.T) {
	type message struct {
		text      string
		createdAt time.Time
	}

	threshold = 5
	cooloff = time.Hour * 5

	for _, tc := range []struct {
		name     string
		messages []message
		want     int
	}{
		{
			name: "6 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: 5,
		},
		{
			name: "5 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: 4,
		},
		{
			name: "4 old messages and 1 new message",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now()},
			},
			want: 5,
		},
		{
			name: "3 old messages and 2 new message",
			messages: []message{
				{"", time.Now()},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now()},
			},
			want: 5,
		},
		{
			name: "2 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: 2,
		},
		{
			name: "1 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: 1,
		},
		{
			name:     "no messages",
			messages: []message{},
			want:     0,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reset(t)

			for _, x := range tc.messages {
				query := `insert into notes(target, nick, text, kind, created_at) values(?, ?, ?, ?, datetime(?))`
				createdAt := x.createdAt.UTC().Format("2006-01-02T15:04:05Z")
				_, err := model.DB.Exec(query, "nick", "nick", "link", x.text, createdAt)
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
			}
			var notes []notes.Note
			err := model.DB.Select(&notes, "select * from notes")
			if err != nil {
				panic(err)
			}
			err = FeedMe(bot.HandlerParams{
				Privmsgf: dummyPrivmsgf,
				Msg:      "",
				Target:   "",
				Matches:  []string{},
			})
			if err != nil {
				t.Fatalf("error running FeedMe %s", err)
			}

			var count int
			err = model.DB.Get(&count, `select count(*) from notes where nick = target`)
			if err != nil {
				t.Fatalf("error getting note count %s", err)
			}
			if count != tc.want {
				t.Fatalf("want %d got %d", tc.want, count)
			}
		})
	}
}

func TestPipeHealth(t *testing.T) {
	type message struct {
		text      string
		createdAt time.Time
	}

	for _, tc := range []struct {
		name           string
		messages       []message
		wantReady      int
		wantFermenting int
	}{
		{
			name: "6 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			wantReady: 6,
		},
		{
			name: "5 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			wantReady: 5,
		},
		{
			name: "3 old messages and 2 new message",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 48)},
				{"", time.Now().Add(-time.Hour * 1)},
				{"", time.Now()},
			},
			wantReady:      3,
			wantFermenting: 2,
		},
		{
			name: "3 old messages and 2 new message",
			messages: []message{
				{"", time.Now()},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now()},
			},
			wantReady:      3,
			wantFermenting: 2,
		},
		{
			name: "2 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			wantReady: 2,
		},
		{
			name: "1 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
			},
			wantReady: 1,
		},
		{
			name:     "no messages",
			messages: []message{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reset(t)
			for _, x := range tc.messages {
				query := `insert into notes(target, nick, text, kind, created_at) values(?, ?, ?, ?, datetime(?))`
				createdAt := x.createdAt.UTC().Format("2006-01-02T15:04:05Z")
				_, err := model.DB.Exec(query, "nick", "nick", "link", x.text, createdAt)
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
			}
			err := PipeHealth(bot.HandlerParams{
				Privmsgf: func(x string, y string, z ...interface{}) {
					ready := z[0]
					fermenting := z[1]
					if ready != tc.wantReady {
						t.Errorf("ready want %d got %d", tc.wantReady, ready)
					}
					if fermenting != tc.wantFermenting {
						t.Errorf("fermenting want %d got %d", tc.wantFermenting, fermenting)
					}
				},
				Msg:     "",
				Target:  "",
				Matches: []string{},
			})
			if err != nil {
				t.Fatalf("error running PipeHealth %s", err)
			}
		})
	}
}

func TestCandidateLinks(t *testing.T) {
	reset(t)

	err := Link(bot.HandlerParams{
		Privmsgf: dummyPrivmsgf,
		Matches:  []string{"", "http://www.example.com"},
		Target:   "theguy",
		Nick:     "theguy",
	})

	if err != nil {
		t.Fatal(err)
	}

	var notes []notes.Note
	notes, err = candidateLinks(0)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 1 {
		t.Fatalf("candidateLinks 0: want 1 got %d", len(notes))
	}

	notes, err = candidateLinks(time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if len(notes) != 0 {
		t.Fatalf("candidate links 1 hour: want 0 got %d", len(notes))
	}
}

func dummyPrivmsgf(x string, y string, z ...interface{}) {
}

func reset(t *testing.T) {
	lastSentAt = time.Unix(0, 0)

	_, err := model.DB.Exec(`delete from notes`)
	if err != nil {
		t.Fatalf("error deleting notes %s", err)
	}
}

func TestCanSendIn(t *testing.T) {
	for _, tc := range []struct {
		start time.Time
		want  time.Duration
	}{
		{
			start: time.Now(),
			want:  5 * time.Hour,
		},
		{
			start: time.Now().Add(-2 * time.Hour),
			want:  3 * time.Hour,
		},
		{
			start: time.Now().Add(-6 * time.Hour),
			want:  -1 * time.Hour,
		},
	} {
		got := canSendIn(tc.start)
		if math.Abs(float64(tc.want-got)) > float64(time.Second) {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}
