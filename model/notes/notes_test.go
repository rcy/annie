package notes

import (
	"testing"
)

func TestCreate(t *testing.T) {
	note, err := Create(CreateParams{
		Nick: "nick",
		Kind: "link",
		Text: "https://www.gnu.org",
	})
	if err != nil {
		t.Error(err)
		return
	}
	want := "https://www.gnu.org"
	if note.Text != want {
		t.Errorf("want %s got %s", want, note.Text)
	}
}
