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
			want:     "Chicago, US -6.8°C (feels like -11.2°C), snow, mist, 1.1mm pcpn over last hour, visibility 2.0km",
		},
		{
			filename: "testdata/creston,ca.json",
			want:     "Creston, CA -10.9°C, overcast clouds, visibility 6.0km",
		},
		{
			filename: "testdata/gibsons,ca.json",
			want:     "Gibsons, CA -0.6°C (feels like -4.9°C), overcast clouds",
		},
		{
			filename: "testdata/shanghai,cn.json",
			want:     "Shanghai, CN 9.2°C (feels like 7.5°C), moderate rain, 1.9mm pcpn over last hour, visibility 7.0km",
		},
		{
			filename: "testdata/toronto,ca.json",
			want:     "Toronto, CA -8.1°C (feels like -15.1°C), overcast clouds",
		},
		{
			filename: "testdata/victoria,ca.json",
			want:     "Victoria, CA 0.6°C (feels like -5.9°C), overcast clouds",
		},
		{
			filename: "testdata/vancouver,ca.json",
			want:     "Vancouver, CA -1.7°C (feels like -5.2°C), overcast clouds",
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
