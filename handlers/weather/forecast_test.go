package weather

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestForecast(t *testing.T) {
	for _, tc := range []struct {
		filename string
		want     string
	}{
		{
			filename: "testdata/chicago,us,forecast.json",
			want:     "Chicago, United States: 3.5, 3.2, 3.0, 2.6, 2.6, 3.6, 3.8, 3.0, 3.0, 2.1, 1.3, 0.9, 2.0, 4.6, 5.9, 4.6, 3.9, 2.7, 1.9, 1.5, 2.8, 4.7, 5.3, 4.0, 3.7, 3.1, 2.4, 2.0, 3.3, 5.3, 5.4, 3.9, 3.3, 3.0, 2.7, 2.3, 3.2, 6.0, 7.7, 6.0",
		},
	} {
		t.Run(tc.filename, func(t *testing.T) {
			payload, err := os.ReadFile(tc.filename)
			if err != nil {
				t.Fatal(err)
			}
			w := forecast{}

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

func TestFetchForecast(t *testing.T) {
	if os.Getenv("NETWORK") == "" {
		t.Skip("NETWORK env var not set")
	}
	forecast, err := fetchForecast("creston,ca")
	if err != nil {
		t.Fatal(err)
	}
	str, err := forecast.Format()
	fmt.Println(str)
}
