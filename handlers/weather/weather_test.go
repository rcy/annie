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
			want:     "Chicago, United States -6.8°C (feels like -11.2°C), snow, mist, 1.1mm pcpn over last hour, visibility 2.0km, wind 9 kph wsw",
		},
		{
			filename: "testdata/creston,ca.json",
			want:     "Creston, Canada -10.9°C, overcast clouds, visibility 6.0km, wind 4 kph nne (gust 4 kph)",
		},
		{
			filename: "testdata/gibsons,ca.json",
			want:     "Gibsons, Canada -0.6°C (feels like -4.9°C), overcast clouds, wind 14 kph nne (gust 16 kph)",
		},
		{
			filename: "testdata/shanghai,cn.json",
			want:     "Shanghai, China 9.2°C (feels like 7.5°C), moderate rain, 1.9mm pcpn over last hour, visibility 7.0km, wind 11 kph n",
		},
		{
			filename: "testdata/toronto,ca.json",
			want:     "Toronto, Canada -8.1°C (feels like -15.1°C), overcast clouds, wind 30 kph wsw (gust 39 kph)",
		},
		{
			filename: "testdata/victoria,ca.json",
			want:     "Victoria, Canada 0.6°C (feels like -5.9°C), overcast clouds, wind 32 kph n",
		},
		{
			filename: "testdata/vancouver,ca.json",
			want:     "Vancouver, Canada -1.7°C (feels like -5.2°C), overcast clouds, wind 9 kph e",
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

func TestMakeWeatherAPIURL(t *testing.T) {
	got, _ := makeWeatherAPIURL("APIKEY", "san francisco")
	want := "http://api.openweathermap.org/data/2.5/weather?appid=APIKEY&q=san+francisco&units=metric"
	if want != got {
		t.Errorf("want: %s, got: %s", want, got)
	}
}
