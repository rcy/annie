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
		{deg: 0, want: "N"}, {deg: 11, want: "N"},
		{deg: 12, want: "NNE"}, {deg: 33, want: "NNE"},
		{deg: 34, want: "NE"}, {deg: 56, want: "NE"},
		{deg: 57, want: "ENE"}, {deg: 78, want: "ENE"},
		{deg: 79, want: "E"}, {deg: 101, want: "E"},
		{deg: 102, want: "ESE"}, {deg: 123, want: "ESE"},
		{deg: 124, want: "SE"}, {deg: 146, want: "SE"},
		{deg: 147, want: "SSE"}, {deg: 168, want: "SSE"},
		{deg: 169, want: "S"}, {deg: 191, want: "S"},
		{deg: 192, want: "SSW"}, {deg: 213, want: "SSW"},
		{deg: 214, want: "SW"}, {deg: 236, want: "SW"},
		{deg: 237, want: "WSW"}, {deg: 258, want: "WSW"},
		{deg: 259, want: "W"}, {deg: 281, want: "W"},
		{deg: 282, want: "WNW"}, {deg: 303, want: "WNW"},
		{deg: 304, want: "NW"}, {deg: 326, want: "NW"},
		{deg: 327, want: "NNW"}, {deg: 348, want: "NNW"},
		{deg: 349, want: "N"}, {deg: 359, want: "N"},

		{deg: 360, want: "N"},
		{deg: 540, want: "S"},
	} {
		t.Run(fmt.Sprint(tc.deg), func(t *testing.T) {
			got := compass16(tc.deg)

			if tc.want != got {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}

}
