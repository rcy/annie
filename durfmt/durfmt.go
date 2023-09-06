package durfmt

import (
	"fmt"
	"time"
)

const (
	day   = 24 * time.Hour
	week  = 7 * day
	month = 4 * week
	year  = 12 * month
)

func Format(dur time.Duration) string {
	if dur >= 0 && dur < time.Second {
		return fmt.Sprintf("%d ms", dur/time.Millisecond)
	}
	if dur >= time.Second && dur < time.Minute {
		if dur >= 2*time.Second {
			return fmt.Sprintf("%d seconds", dur/time.Second)
		}
		return fmt.Sprintf("%d second", dur/time.Second)
	}
	if dur >= time.Minute && dur < time.Hour {
		if dur >= 2*time.Minute {
			return fmt.Sprintf("%d minutes", dur/time.Minute)
		}
		return fmt.Sprintf("%d minute", dur/time.Minute)
	}
	if dur >= time.Hour && dur < day {
		if dur >= 2*time.Hour {
			return fmt.Sprintf("%d hours", dur/time.Hour)
		}
		return fmt.Sprintf("%d hour", dur/time.Hour)
	}
	if dur >= day && dur < week {
		if dur >= 2*day {
			return fmt.Sprintf("%d days", dur/day)
		}
		return fmt.Sprintf("%d day", dur/day)
	}
	if dur >= week && dur < month {
		if dur >= 2*week {
			return fmt.Sprintf("%d weeks", dur/week)
		}
		return fmt.Sprintf("%d week", dur/week)
	}
	if dur >= month && dur < year {
		if dur >= 2*month {
			return fmt.Sprintf("%d months", dur/month)
		}
		return fmt.Sprintf("%d month", dur/month)
	}
	if dur >= 2*year {
		return fmt.Sprintf("%d years", dur/year)
	}
	return fmt.Sprintf("%d year", dur/year)
}
