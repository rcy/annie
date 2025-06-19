package timeoff

import (
	"goirc/internal/sun"
	"time"

	"github.com/rcy/ddate"
)

func IsTimeoff(when time.Time, zone string, lat float64, long float64) (bool, error) {
	location, err := time.LoadLocation(zone)
	if err != nil {
		return false, err
	}

	when = when.In(location)
	sunrise, sunset, err := sun.SunriseSunset(when, zone, lat, long)
	if err != nil {
		return false, err
	}

	weekday := ddate.FromTime(when).WeekDay

	return (weekday == ddate.SettingOrange && when.After(sunset)) ||
		(weekday == ddate.Sweetmorn && when.Before(sunrise)), nil
}
