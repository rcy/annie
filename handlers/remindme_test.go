package handlers

import (
	"goirc/model/reminders"
	"math"
	"testing"
	"time"
)

func TestRemind(t *testing.T) {
	tests := []struct {
		Nick      string
		Duration  string
		What      string
		Want      reminders.Reminder
		WantError string
	}{
		{
			Nick:     "preethi",
			Duration: "69s",
			What:     "start breathing again",
			Want: reminders.Reminder{
				Nick:     "preethi",
				RemindAt: time.Now().Add(69 * time.Second),
				What:     "start breathing again",
			},
		},
		{
			Nick:     "gerulf",
			Duration: "12m",
			What:     "check egg",
			Want: reminders.Reminder{
				Nick:     "gerulf",
				RemindAt: time.Now().Add(12 * time.Minute),
				What:     "check egg",
			},
		},
		{
			Nick:     "mary",
			Duration: "8h",
			What:     "wake up",
			Want: reminders.Reminder{
				Nick:     "mary",
				RemindAt: time.Now().Add(8 * time.Hour),
				What:     "wake up",
			},
		},
		{
			Nick:     "bob",
			Duration: "2w",
			What:     "return from vacation",
			Want: reminders.Reminder{
				Nick:     "bob",
				RemindAt: time.Now().Add(14 * 24 * time.Hour),
				What:     "return from vacation",
			},
		},
		{
			Nick:     "slamet",
			Duration: "1y",
			What:     "be a year older",
			Want: reminders.Reminder{
				Nick:     "slamet",
				RemindAt: time.Now().Add(14 * 24 * time.Hour),
				What:     "be a year older",
			},
			WantError: `time: unknown unit "y" in duration "1y"`,
		},
	}

	for _, test := range tests {
		t.Run(test.Nick+" "+test.What, func(t *testing.T) {
			initRows, err := reminders.All()
			if err != nil {
				t.Fatal(err)
			}

			result, err := remind(test.Nick, test.Duration, test.What)
			if err != nil {
				if test.WantError != err.Error() {
					t.Fatalf("wanted error %s got %s", test.WantError, err)
				}
				return
			}
			if test.WantError != "" {
				t.Fatalf("wanted error %s, didn't get one", test.WantError)
			}

			rows, err := reminders.All()
			if err != nil {
				t.Fatal(err)
			}
			if len(rows) != len(initRows)+1 {
				t.Fatalf("want %d row, got %d", len(initRows)+1, len(rows))
			}

			got := rows[len(rows)-1]

			diff := secondsApart(*result, test.Want.RemindAt)
			if diff > 1 {
				t.Fatalf("result=%v and want.RemindAt=%v are not close: %f", result, test.Want.RemindAt, diff)
			}

			diff = secondsApart(test.Want.RemindAt, got.RemindAt)
			if diff > 1 {
				t.Fatalf("want.RemindAt=%v and got.RemindAt=%v are not close: %f", test.Want.RemindAt, got.RemindAt, diff)
			}

			if got.Nick != test.Want.Nick {
				t.Fatalf("bad nick: want %s got %s", test.Want.Nick, got.Nick)
			}

			if got.What != test.Want.What {
				t.Fatalf("bad what: want %s got %s", test.Want.What, got.What)
			}
		})
	}
}

func secondsApart(t1, t2 time.Time) float64 {
	return math.Abs(t1.Sub(t2).Seconds())
}
