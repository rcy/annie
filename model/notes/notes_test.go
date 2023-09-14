package notes

import (
	"testing"
)

func TestCreate(t *testing.T) {
	for _, tc := range []struct {
		name      string
		target    string
		nick      string
		kind      string
		text      string
		wantError string
	}{
		{
			name:   "good link to channel",
			target: "#chan",
			nick:   "nick",
			kind:   "link",
			text:   "https://www.gnu.org",
		},
		{
			name:   "good link to nick",
			target: "nick",
			nick:   "nick",
			kind:   "link",
			text:   "https://www.gnu.org",
		},
		{
			name:      "missing target",
			nick:      "nick",
			kind:      "link",
			text:      "https://www.gnu.org",
			wantError: "target cannot be empty",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			note, err := Create(CreateParams{
				Target: tc.target,
				Nick:   tc.nick,
				Kind:   tc.kind,
				Text:   tc.text,
			})
			if err != nil {
				if err.Error() != tc.wantError {
					t.Errorf("wanted error '%s', got '%s'", tc.wantError, err.Error())
					return
				}
				return
			}
			if note.Kind != tc.kind {
				t.Errorf("wanted kind '%s', got '%s'", tc.kind, note.Kind)
			}
			if note.Nick != tc.nick {
				t.Errorf("wanted nick '%s', got '%s'", tc.nick, note.Nick)
			}
			if note.Text != tc.text {
				t.Errorf("wanted text '%s', got '%s'", tc.text, note.Text)
			}
			if note.Target != tc.target {
				t.Errorf("wanted target '%s', got '%s'", tc.target, note.Target)
			}
		})
	}
}
