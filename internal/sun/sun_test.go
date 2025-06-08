package sun

import (
	"testing"
	"time"
)

func TestSunriseSunset(t *testing.T) {
	for _, tc := range []struct {
		zone     string
		day      string
		coords   []float64
		wantRise string
		wantSet  string
	}{
		{
			day:      "2025-06-08",
			zone:     "America/Creston",
			coords:   []float64{49.09987, -116.50211},
			wantRise: "2025-06-08T04:42:14-07:00",
			wantSet:  "2025-06-08T20:48:18-07:00",
		},
		{
			day:      "2025-06-09",
			zone:     "America/Creston",
			coords:   []float64{49.09987, -116.50211},
			wantRise: "2025-06-09T04:41:53-07:00",
			wantSet:  "2025-06-09T20:49:01-07:00",
		},
		{
			day:      "2025-06-09",
			zone:     "America/Toronto",
			coords:   []float64{43.64487, -79.38429},
			wantRise: "2025-06-09T05:36:08-04:00",
			wantSet:  "2025-06-09T20:57:48-04:00",
		},
		{
			// sunrise and sunset in toronto but reported in vancouver time
			day:      "2025-06-09",
			zone:     "America/Vancouver",
			coords:   []float64{43.64487, -79.38429},
			wantRise: "2025-06-09T02:36:08-07:00",
			wantSet:  "2025-06-09T17:57:48-07:00",
		},
		{
			// sunrise and sunset in toronto but reported in UTC
			day:      "2025-06-09",
			zone:     "",
			coords:   []float64{43.64487, -79.38429},
			wantRise: "2025-06-09T09:36:08Z",
			wantSet:  "2025-06-10T00:57:48Z",
		},
	} {
		t.Run(tc.day+" "+tc.zone, func(t *testing.T) {
			day, err := time.Parse(time.DateOnly, tc.day)
			if err != nil {
				t.Errorf("time.Parse: %s", err)
			}
			rise, set, err := SunriseSunset(day, tc.zone, tc.coords[0], tc.coords[1])
			if err != nil {
				t.Errorf("SunriseSunset: %s", err)
			}

			gotRise := rise.Format(time.RFC3339)
			if gotRise != tc.wantRise {
				t.Errorf("rise want: %s got: %s", tc.wantRise, gotRise)
			}

			gotSet := set.Format(time.RFC3339)
			if gotSet != tc.wantSet {
				t.Errorf("set want: %s got: %s", tc.wantSet, gotSet)
			}
		})
	}
}
