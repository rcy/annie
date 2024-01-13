package weather

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"net/http"
	"os"
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
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
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

	var data response
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func Handle(params bot.HandlerParams) error {
	q := params.Matches[1]

	resp, err := Weather(q)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, fmt.Sprintf("%f", resp.Main.Temp))

	return nil
}
