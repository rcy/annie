package weather

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"io"
	"net/http"
	"os"
	"strings"
)

type response struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		OneH   float64 `json:"1h"`
		ThreeH float64 `json:"3h"`
	} `json:"snow"`
	Snow struct {
		OneH   float64 `json:"1h"`
		ThreeH float64 `json:"3h"`
	} `json:"snow"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Message int    `json:"message"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func (r *response) String() string {
	str := fmt.Sprintf("%s, %s %.1f°C ", r.Name, r.Sys.Country, r.Main.Temp)
	if r.Main.FeelsLike != r.Main.Temp {
		str += fmt.Sprintf("(feels like %.1f°C) ", r.Main.FeelsLike)
	}
	descs := []string{}
	for _, w := range r.Weather {
		descs = append(descs, w.Description)
	}
	str += strings.Join(descs, ", ")

	return str
}

const iconURLFmt = "https://openweathermap.org/img/wn/%s@2x.png"

func Weather(q string) (*response, error) {
	key := os.Getenv("OPENWEATHERMAP_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("bad api key")
	}

	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?units=metric&q=%s&appid=%s", q, key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("city not found")
	}

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func XWeather(q string) ([]byte, error) {
	key := os.Getenv("OPENWEATHERMAP_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("bad api key")
	}

	resp, err := http.Get(fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?units=metric&q=%s&appid=%s", q, key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("city not found")
	}

	return io.ReadAll(resp.Body)
}

func Handle(params bot.HandlerParams) error {
	q := params.Matches[1]

	resp, err := Weather(q)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, resp.String())

	return nil
}

func XHandle(params bot.HandlerParams) error {
	q := params.Matches[1]

	resp, err := XWeather(q)
	if err != nil {
		return err
	}

	chunks := splitBytes(resp, 420)

	for _, chunk := range chunks {
		params.Privmsgf(params.Target, string(chunk))
	}

	return nil
}

func splitBytes(data []byte, chunkSize int) [][]byte {
	var chunks [][]byte

	for len(data) > 0 {
		if len(data) < chunkSize {
			chunkSize = len(data)
		}
		chunks = append(chunks, data[:chunkSize])
		data = data[chunkSize:]
	}

	return chunks
}
