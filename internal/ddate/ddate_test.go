package ddate

import (
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestFromTime(t *testing.T) {
	for _, tc := range []struct {
		when string
		want Date
	}{
		{
			when: "2025-09-21",
			want: Date{
				Year:      3191,
				Season:    3,
				SeasonDay: 45,
				WeekDay:   3,
				HolyDay:   "",
			},
		},
		{
			when: "2030-10-24",
			want: Date{
				Year:      3196,
				Season:    4,
				SeasonDay: 5,
				WeekDay:   1,
				HolyDay:   "Maladay",
			},
		},
		{
			when: "2020-02-29",
			want: Date{
				Year:      3186,
				Season:    0,
				SeasonDay: 0,
				WeekDay:   0,
				HolyDay:   "St. Tib's Day",
			},
		},
		{
			when: time.Unix(2<<30, 0).Format(time.DateOnly), // epocholypse day
			want: Date{
				Year:      3204,
				Season:    0,
				SeasonDay: 18,
				WeekDay:   2,
				HolyDay:   "",
			},
		},
	} {
		t.Run(tc.when, func(t *testing.T) {
			date, err := time.Parse(time.DateOnly, tc.when)
			if err != nil {
				t.Fatalf("time.Parse: %s", err)
			}

			got := FromTime(date)
			if got != tc.want {
				t.Errorf("got: %v, want: %v", got, tc.want)
			}
		})
	}
}

// Return a random time.Time between the year 1 and year 9999
func randomTime() time.Time {
	min := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
	max := time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
	span := max.Sub(min)
	randDur := time.Duration(rand.Int63n(span.Nanoseconds()))
	return min.Add(randDur)
}

// Return the string from the classic ddate command line tool found in the path
func ddateCmd(day time.Time) (string, error) {
	cmd := exec.Command("ddate", strconv.Itoa(day.Day()), strconv.Itoa(int(day.Month())), strconv.Itoa(day.Year()))

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func TestCompareDDateCmd(t *testing.T) {
	// compare every day for 2 years, include 2020 which is a leap year
	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC)
	for when := start; when.Before(end); when = when.AddDate(0, 0, 1) {
		t.Run(when.String(), func(t *testing.T) {
			got := FromTime(when).Format(false)
			gotLegacy, err := ddateCmd(when)
			if err != nil {
				t.Errorf("ddateCmd: %s", err)
			}

			if got != gotLegacy {
				t.Errorf("original command output mismatch: got: %s, gotLegacy: %s", got, gotLegacy)
			}
		})
	}
}

func TestFormatNoHoly(t *testing.T) {
	for _, tc := range []struct {
		when         string
		want         string
		wantWithHoly string
	}{
		{when: "1970-01-01", want: "Sweetmorn, Chaos 1, 3136 YOLD"},
		{when: "1999-12-31", want: "Setting Orange, The Aftermath 73, 3165 YOLD"},
		{when: "2000-01-01", want: "Sweetmorn, Chaos 1, 3166 YOLD"},
		{when: "2020-02-19", want: "Setting Orange, Chaos 50, 3186 YOLD"},
		{when: "2020-02-28", want: "Prickle-Prickle, Chaos 59, 3186 YOLD"},
		{when: "2020-02-29", want: "St. Tib's Day, 3186 YOLD"}, // even with holy days off, St. Tib's Day is shown
		{when: "2020-03-01", want: "Setting Orange, Chaos 60, 3186 YOLD"},
		{when: "2020-06-11", want: "Boomtime, Confusion 16, 3186 YOLD"},
		{when: "2020-12-31", want: "Setting Orange, The Aftermath 73, 3186 YOLD"},
		{when: "2025-01-05", want: "Setting Orange, Chaos 5, 3191 YOLD"},
		{when: "2025-02-19", want: "Setting Orange, Chaos 50, 3191 YOLD"},
		{when: "2025-03-19", want: "Pungenday, Discord 5, 3191 YOLD"},
		{when: "2025-05-03", want: "Pungenday, Discord 50, 3191 YOLD"},
		{when: "2025-05-31", want: "Sweetmorn, Confusion 5, 3191 YOLD"},
		{when: "2025-06-08", want: "Prickle-Prickle, Confusion 13, 3191 YOLD"},
		{when: "2025-06-09", want: "Setting Orange, Confusion 14, 3191 YOLD"},
		{when: "2025-06-10", want: "Sweetmorn, Confusion 15, 3191 YOLD"},
		{when: "2025-06-11", want: "Boomtime, Confusion 16, 3191 YOLD"},
		{when: "2025-06-12", want: "Pungenday, Confusion 17, 3191 YOLD"},
		{when: "2025-07-15", want: "Sweetmorn, Confusion 50, 3191 YOLD"},
		{when: "2025-08-12", want: "Prickle-Prickle, Bureaucracy 5, 3191 YOLD"},
		{when: "2025-09-26", want: "Prickle-Prickle, Bureaucracy 50, 3191 YOLD"},
		{when: "2025-10-24", want: "Boomtime, The Aftermath 5, 3191 YOLD"},
		{when: "2025-12-08", want: "Boomtime, The Aftermath 50, 3191 YOLD"},
	} {
		t.Run(tc.when, func(t *testing.T) {
			date, err := time.Parse(time.DateOnly, tc.when)
			if err != nil {
				t.Fatalf("time.Parse: %s", err)
			}

			got := FromTime(date).Format(false)
			if got != tc.want {
				t.Errorf("got: %s, want: %s", got, tc.want)
			}
		})
	}
}

func TestFormatHoly(t *testing.T) {
	for _, tc := range []struct {
		when string
		want string
	}{
		{when: "2020-01-05", want: "Setting Orange, Chaos 5, 3186 YOLD (Mungday)"},
		{when: "2020-02-19", want: "Setting Orange, Chaos 50, 3186 YOLD (Chaoflux)"},
		{when: "2020-02-29", want: "St. Tib's Day, 3186 YOLD"},
		{when: "2020-03-19", want: "Pungenday, Discord 5, 3186 YOLD (Mojoday)"},
		{when: "2020-05-03", want: "Pungenday, Discord 50, 3186 YOLD (Discoflux)"},
		{when: "2020-05-31", want: "Sweetmorn, Confusion 5, 3186 YOLD (Syaday)"},
		{when: "2020-07-15", want: "Sweetmorn, Confusion 50, 3186 YOLD (Confuflux)"},
		{when: "2020-08-12", want: "Prickle-Prickle, Bureaucracy 5, 3186 YOLD (Zaraday)"},
		{when: "2020-09-26", want: "Prickle-Prickle, Bureaucracy 50, 3186 YOLD (Bureflux)"},
		{when: "2020-10-24", want: "Boomtime, The Aftermath 5, 3186 YOLD (Maladay)"},
		{when: "2020-12-08", want: "Boomtime, The Aftermath 50, 3186 YOLD (Afflux)"},
	} {
		t.Run(tc.when, func(t *testing.T) {
			date, err := time.Parse(time.DateOnly, tc.when)
			if err != nil {
				t.Fatalf("time.Parse: %s", err)
			}

			got := FromTime(date).Format(true)
			if got != tc.want {
				t.Errorf("got: %s, want: %s", got, tc.want)
			}
		})
	}
}
