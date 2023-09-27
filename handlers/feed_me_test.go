package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"testing"
	"time"
)

func TestFeedMe(t *testing.T) {
	type message struct {
		text      string
		createdAt time.Time
	}

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
			_, err := model.DB.Exec(`delete from notes`)
			if err != nil {
				t.Fatalf("error deleting notes %s", err)
			}

			for _, x := range tc.messages {
				query := `insert into notes(target, nick, text, kind, created_at) values(?, ?, ?, ?, ?)`
				createdAt := x.createdAt.UTC().Format("2006-01-02T15:04:05Z")
				_, err := model.DB.Exec(query, "nick", "nick", "link", x.text, createdAt)
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
			}
			var notes []notes.Note
			err = model.DB.Select(&notes, "select * from notes")
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

func dummyPrivmsgf(x string, y string, z ...interface{}) {
	return
}
