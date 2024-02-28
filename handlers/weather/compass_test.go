package weather

import (
	"fmt"
	"testing"
)

func TestCompass16(t *testing.T) {
	for _, tc := range []struct {
		deg  uint
		want string
	}{
		{deg: 0, want: "n"}, {deg: 11, want: "n"},
		{deg: 12, want: "nne"}, {deg: 33, want: "nne"},
		{deg: 34, want: "ne"}, {deg: 56, want: "ne"},
		{deg: 57, want: "ene"}, {deg: 78, want: "ene"},
		{deg: 79, want: "e"}, {deg: 101, want: "e"},
		{deg: 102, want: "ese"}, {deg: 123, want: "ese"},
		{deg: 124, want: "se"}, {deg: 146, want: "se"},
		{deg: 147, want: "sse"}, {deg: 168, want: "sse"},
		{deg: 169, want: "s"}, {deg: 191, want: "s"},
		{deg: 192, want: "ssw"}, {deg: 213, want: "ssw"},
		{deg: 214, want: "sw"}, {deg: 236, want: "sw"},
		{deg: 237, want: "wsw"}, {deg: 258, want: "wsw"},
		{deg: 259, want: "w"}, {deg: 281, want: "w"},
		{deg: 282, want: "wnw"}, {deg: 303, want: "wnw"},
		{deg: 304, want: "nw"}, {deg: 326, want: "nw"},
		{deg: 327, want: "nnw"}, {deg: 348, want: "nnw"},
		{deg: 349, want: "n"}, {deg: 359, want: "n"},

		{deg: 360, want: "n"},
		{deg: 540, want: "s"},
	} {
		t.Run(fmt.Sprint(tc.deg), func(t *testing.T) {
			got := compass16(tc.deg)

			if tc.want != got {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}

}
