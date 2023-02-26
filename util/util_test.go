package util

import (
	"testing"
	"time"
)

func mustParseTime(str string) time.Time {
	result, err := ParseTime(str)
	if err != nil {
		panic(err)
	}
	return result
}

func TestAgo(t *testing.T) {
	type cases struct {
		t1   time.Time
		t2   time.Time
		want string
	}

	for _, scenario := range []cases{
		{
			t1:   mustParseTime("2010-01-01 00:00:00"),
			t2:   mustParseTime("2020-01-01 00:00:00"),
			want: "3652d",
		},
		{
			t1:   mustParseTime("2010-01-01 00:00:00"),
			t2:   mustParseTime("2010-01-02 00:00:00"),
			want: "24h0m0s",
		},
		{
			t1:   mustParseTime("2010-01-01 00:00:00"),
			t2:   mustParseTime("2010-01-08 00:00:00"),
			want: "7d",
		},
		{
			t1:   mustParseTime("2010-01-01 00:00:00"),
			t2:   mustParseTime("2010-02-01 00:00:00"),
			want: "31d",
		},
	} {
		dur := scenario.t2.Sub(scenario.t1)
		result := Ago(dur)
		if result != scenario.want {
			t.Errorf("want %s got %s", scenario.want, result)
		}
	}
}
