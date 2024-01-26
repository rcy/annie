package handlers

import (
	"goirc/bot"
	db "goirc/model"
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
				_, err := db.DB.Exec(query, "nick", "nick", "link", x.text, createdAt)
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
			}
			_, err := FeedMe(bot.HandlerParams{
				Msg:     "",
				Target:  "",
				Matches: []string{},
			})
			if err != nil {
				t.Fatalf("error running FeedMe %s", err)
			}

			var count int
			err = db.DB.Get(&count, `select count(*) from notes where nick = target`)
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
		name     string
		messages []message
		want     string
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
			want: "6 links ready to serve (0 fermenting)",
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
			want: "5 links ready to serve (0 fermenting)",
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
			want: "3 links ready to serve (2 fermenting)",
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
			want: "3 links ready to serve (2 fermenting)",
		},
		{
			name: "2 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: "2 links ready to serve (0 fermenting)",
		},
		{
			name: "1 old messages",
			messages: []message{
				{"", time.Now().Add(-time.Hour * 24)},
			},
			want: "1 links ready to serve (0 fermenting)",
		},
		{
			name:     "no messages",
			messages: []message{},
			want:     "0 links ready to serve (0 fermenting)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			reset(t)
			for _, x := range tc.messages {
				query := `insert into notes(target, nick, text, kind, created_at) values(?, ?, ?, ?, datetime(?))`
				createdAt := x.createdAt.UTC().Format("2006-01-02T15:04:05Z")
				_, err := db.DB.Exec(query, "nick", "nick", "link", x.text, createdAt)
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
			}
			got, err := PipeHealth(bot.HandlerParams{
				Msg:     "",
				Target:  "",
				Matches: []string{},
			})
			if err != nil {
				t.Fatalf("error running PipeHealth %s", err)
			}
			if got != tc.want {
				t.Errorf("want %s got %s", got, tc.want)
			}

		})
	}
}

func TestCandidateLinks(t *testing.T) {
	reset(t)

	_, err := Link(bot.HandlerParams{
		Matches: []string{"", "http://www.example.com"},
		Target:  "theguy",
		Nick:    "theguy",
	})

	if err != nil {
		t.Fatal(err)
	}

	notes, err := candidateLinks(0)
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

func reset(t *testing.T) {
	lastSentAt = time.Unix(0, 0)

	_, err := db.DB.Exec(`delete from notes`)
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
