package ddate

import (
	"testing"
	"time"
)

func TestOn(t *testing.T) {
	for _, tc := range []struct {
		when string
		want string
	}{
		{
			when: "1970-01-01",
			want: "Sweetmorn, Chaos 1, 3136 YOLD",
		},
		{
			when: "1999-12-31",
			want: "Setting Orange, The Aftermath 73, 3165 YOLD",
		},
		{
			when: "2000-01-01",
			want: "Sweetmorn, Chaos 1, 3166 YOLD",
		},
		{
			when: "2025-06-08",
			want: "Prickle-Prickle, Confusion 13, 3191 YOLD",
		},
		{
			when: "2025-06-09",
			want: "Setting Orange, Confusion 14, 3191 YOLD",
		},
		{
			when: "2025-06-10",
			want: "Sweetmorn, Confusion 15, 3191 YOLD",
		},
		{
			when: "2025-06-11",
			want: "Boomtime, Confusion 16, 3191 YOLD",
		},
		{
			when: "2025-06-12",
			want: "Pungenday, Confusion 17, 3191 YOLD",
		},
		{
			when: "2020-02-28",
			want: "Prickle-Prickle, Chaos 59, 3186 YOLD",
		},
		{
			when: "2020-02-29",
			want: "St. Tib's Day, 3186 YOLD",
		},
		{
			when: "2020-03-01",
			want: "Setting Orange, Chaos 60, 3186 YOLD",
		},
		{
			when: "2020-02-19",
			//want: "Setting Orange, Chaos 50, 3186 YOLD (Chaoflux)",
			want: "Setting Orange, Chaos 50, 3186 YOLD",
		},
		{
			when: "2020-06-11",
			want: "Boomtime, Confusion 16, 3186 YOLD",
		},
		{
			when: "2020-12-31",
			want: "Setting Orange, The Aftermath 73, 3186 YOLD",
		},
	} {
		t.Run(tc.when, func(t *testing.T) {
			date, err := time.Parse(time.DateOnly, tc.when)
			if err != nil {
				t.Errorf("time.Parse: %s", err)
			}

			// verify with classic ddate command
			got, err := ddateCmd(date)
			if err != nil {
				t.Errorf("On: %s", err)
			}
			if got != tc.want {
				t.Errorf("(classic) want: %s, got: %s", tc.want, got)
			}

			// verify go version
			got = FromTime(date).Format(false)
			if got != tc.want {
				t.Errorf("(go) want: %s, got: %s", tc.want, got)
			}
		})
	}
}
