package durfmt

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	for _, test := range []struct {
		input time.Duration
		want  string
	}{
		{
			input: 100 * time.Millisecond,
			want:  "100 ms",
		},
		{
			input: time.Second,
			want:  "1 second",
		},
		{
			input: 3 * time.Second,
			want:  "3 seconds",
		},
		{
			input: 60 * time.Second,
			want:  "1 minute",
		},
		{
			input: 61 * time.Second,
			want:  "1 minute",
		},
		{
			input: 25 * time.Minute,
			want:  "25 minutes",
		},
		{
			input: 60 * time.Minute,
			want:  "1 hour",
		},
		{
			input: 120 * time.Minute,
			want:  "2 hours",
		},
		{
			input: 12 * time.Hour,
			want:  "12 hours",
		},
		{
			input: 24 * time.Hour,
			want:  "1 day",
		},
		{
			input: 2 * day,
			want:  "2 days",
		},
		{
			input: 7 * day,
			want:  "1 week",
		},
		{
			input: 8 * day,
			want:  "1 week",
		},
		{
			input: 13 * day,
			want:  "1 week",
		},
		{
			input: 14 * day,
			want:  "2 weeks",
		},
		{
			input: 29 * day,
			want:  "1 month",
		},
		{
			input: 30 * day,
			want:  "1 month",
		},
		{
			input: 40 * day,
			want:  "1 month",
		},
		{
			input: 2 * month,
			want:  "2 months",
		},
		{
			input: 12 * month,
			want:  "1 year",
		},
		{
			input: 364 * day,
			want:  "1 year",
		},
		{
			input: 365 * day,
			want:  "1 year",
		},
		{
			input: 1*year + 6*month,
			want:  "1 year",
		},
		{
			input: 2 * year,
			want:  "2 years",
		},
	} {
		t.Run(test.want, func(t *testing.T) {
			got := Format(test.input)

			if got != test.want {
				t.Errorf("want %s got %s", test.want, got)
			}
		})
	}
}
