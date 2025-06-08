package sun

import (
	"fmt"
	"time"

	"github.com/nathan-osman/go-sunrise"
)

// Return the time of sunrise and sunset in the requested zone at lat and long
func SunriseSunset(day time.Time, zone string, lat float64, long float64) (time.Time, time.Time, error) {
	rise, set := sunrise.SunriseSunset(
		lat, long,
		day.Year(), day.Month(), day.Day(),
	)
	location, err := time.LoadLocation(zone)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("LoadLocation: %w", err)
	}

	return rise.In(location), set.In(location), nil
}
