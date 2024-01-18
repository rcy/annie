package weather

import (
	"encoding/json"
	"os"
	"testing"
)

func TestWeather(t *testing.T) {
	for _, tc := range []struct {
		filename string
		want     string
	}{
		{
			filename: "testdata/chicago,us.json",
			want:     "Chicago, US -6.6°C (feels like -10.9°C) light snow, mist",
		},
		{
			filename: "testdata/creston,ca.json",
			want:     "Creston, CA -10.9°C overcast clouds",
		},
		{
			filename: "testdata/gibsons,ca.json",
			want:     "Gibsons, CA -0.6°C (feels like -4.9°C) overcast clouds",
		},
		{
			filename: "testdata/toronto,ca.json",
			want:     "Toronto, CA -8.1°C (feels like -15.1°C) overcast clouds",
		},
		{
			filename: "testdata/victoria,ca.json",
			want:     "Victoria, CA 0.6°C (feels like -5.9°C) overcast clouds",
		},
		{
			filename: "testdata/vancouver,ca.json",
			want:     "Vancouver, CA -1.7°C (feels like -5.2°C) overcast clouds",
		},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			payload, err := os.ReadFile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			w := weather{}

			err = json.Unmarshal(payload, &w)
			if err != nil {
				t.Fatal(err)
			}
			got := w.String()

			if tc.want != got {
				t.Errorf("expected:\n%s\ngot:\n%s", tc.want, got)
			}
		})
	}

}
