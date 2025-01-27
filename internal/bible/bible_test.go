package bible

import (
	"reflect"
	"testing"
)

func TestParseRef(t *testing.T) {
	for _, tc := range []struct {
		text      string
		want      ref
		wantError bool
	}{
		{
			text: "Genesis 1:1",
			want: ref{"Genesis", 1, 1, 1},
		},
		{
			text: "Judges 3:12-30",
			want: ref{"Judges", 3, 12, 30},
		},
		{
			text: "1 Samuel 15:3",
			want: ref{"1 Samuel", 15, 3, 3},
		},
		{
			text:      "Acts 28:2-1",
			want:      ref{},
			wantError: true,
		},
	} {
		t.Run(tc.text, func(t *testing.T) {
			got, err := parseRef(tc.text)
			if err != nil {
				if tc.wantError {
					return
				}
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, &tc.want) {
				t.Errorf("expected %v, got %v", tc.want, got)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	b := New()
	err := b.SetActiveTranslation("KJV")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		text      string
		want      string
		wantError bool
	}{
		{
			text: "Genesis 1:1",
			want: "In the beginning God created the heaven and the earth.",
		},
		{
			text: "2 Corinthians 5:14-15",
			want: "For the love of Christ constraineth us; because we thus judge, that if one died for all, then were all dead: And that he died for all, that they which live should not henceforth live unto themselves, but unto him which died for them, and rose again.",
		},
		{
			text: "Matthew 5:8",
			want: "Blessed are the pure in heart: for they shall see God.",
		},
	} {
		t.Run(tc.text, func(t *testing.T) {
			got, err := b.Lookup(tc.text)
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("\nexpected: %v\ngot: %v", tc.want, got)
			}
		})
	}
}
