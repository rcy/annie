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
		t0   string
		t1   string
		want string
	}

	for _, scenario := range []cases{
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2020-01-01 00:00:00",
			want: "10 years",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2019-01-01 00:00:00",
			want: "9 years",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-02 00:00:00",
			want: "1 day",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-08 00:00:00",
			want: "1 week",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-02-01 00:00:00",
			want: "4 weeks",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-03-01 00:00:00",
			want: "8 weeks", // 2 months would be better here
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-01 01:01:01",
			want: "1 hour",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-01 00:45:01",
			want: "45 minutes",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-01 00:01:01",
			want: "1 minute",
		},
		{
			t0:   "2010-01-01 00:00:00",
			t1:   "2010-01-01 00:00:10",
			want: "10 seconds",
		},
		{
			t0:   "2023-02-26 11:59:13",
			t1:   "2026-06-01T15:00:00Z",
			want: "3 years",
		},
	} {
		then := mustParseTime(scenario.t0)
		now := mustParseTime(scenario.t1)
		dur := now.Sub(then)
		result := Ago(dur)
		if result != scenario.want {
			t.Errorf("want %s got %s", scenario.want, result)
		}
	}
}
