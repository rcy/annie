package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/model"
	"goirc/model/notes"
	"testing"
)

func TestFeedMe(t *testing.T) {
	for i, tc := range []struct {
		messages []string
		want     int
	}{
		{
			messages: []string{"abc", "def", "ghi", "jkl", "mno"},
			want:     4,
		},
		{
			messages: []string{"abc", "def", "ghi", "jkl"},
			want:     4,
		},
		{
			messages: []string{"abc", "def"},
			want:     2,
		},
		{
			messages: []string{"abc"},
			want:     1,
		},
		{
			messages: []string{},
			want:     0,
		},
	} {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			_, err := model.DB.Exec(`delete from notes`)
			if err != nil {
				t.Fatalf("error deleting notes %s", err)
			}

			for _, x := range tc.messages {
				_, err := notes.Create(notes.CreateParams{Target: "nick", Nick: "nick", Kind: "link", Text: x})
				if err != nil {
					t.Fatalf("error creating note %s", err)
				}
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
