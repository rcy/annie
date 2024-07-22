package weather

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"goirc/bot"
	db "goirc/db/model"
	"goirc/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pariz/gountries"
)

type forecast struct {
	Cod     string `json:"cod"`
	Message int    `json:"message"`
	Cnt     int    `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			FeelsLike float64 `json:"feels_like"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  int     `json:"pressure"`
			SeaLevel  int     `json:"sea_level"`
			GrndLevel int     `json:"grnd_level"`
			Humidity  int     `json:"humidity"`
			TempKf    float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Visibility int     `json:"visibility"`
		Pop        float64 `json:"pop"`
		Sys        struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Country    string `json:"country"`
		Population int    `json:"population"`
		Timezone   int    `json:"timezone"`
		Sunrise    int    `json:"sunrise"`
		Sunset     int    `json:"sunset"`
	} `json:"city"`
}

func (f forecast) String() string {
	//components := []string{}

	var countryStr string
	country, err := gountries.New().FindCountryByAlpha(f.City.Country)
	if err != nil {
		countryStr = "??"
	} else {
		countryStr = country.Name.Common
	}

	city := fmt.Sprintf("%s, %s", f.City.Name, countryStr)

	// show temperatures for the next 3 days, every 5 hours, or whatever is in the result
	temps := []string{}
	for _, fc := range f.List {
		temps = append(temps, fmt.Sprintf("%0.1f", fc.Main.Temp))
	}

	return city + ": " + strings.Join(temps, ", ")
}

func fetchForecast(q string) (*forecast, error) {
	weather, err := fetchWeather(q)
	if err != nil {
		return nil, err
	}
	fmt.Println(weather)
	return fetchForecastByCoords(weather.Coord.Lat, weather.Coord.Lon)
}

func fetchForecastByCoords(lat, lon float64) (*forecast, error) {
	key := os.Getenv("OPENWEATHERMAP_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("bad api key")
	}

	resp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%f&lon=%f&appid=%s", lat, lon, key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("location not found")
	}

	var f forecast
	err = json.NewDecoder(resp.Body).Decode(&f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func HandleForecast(params bot.HandlerParams) error {
	ctx := context.TODO()
	queries := db.New(model.DB)

	var q string
	if len(params.Matches) > 1 {
		q = params.Matches[1]
	}

	if q == "" {
		last, err := queries.LastNickWeatherRequest(ctx, params.Nick)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("no previous weather station to report on")
			}
			return err
		}
		if params.Nick == last.Nick {
			if strings.HasPrefix(last.City, q) {
				q = last.City + "," + last.Country
			}
		}
	} else {
		last, err := queries.LastWeatherRequestByPrefix(ctx, sql.NullString{String: q, Valid: true})
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if last.ID != 0 {
			q = last.City + "," + last.Country
		}
	}

	forecast, err := fetchForecast(q)
	if err != nil {
		return err
	}
	str, err := forecast.Format()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s, %s forecast: %s", forecast.City.Name, forecast.City.Country, str)

	err = queries.InsertNickWeatherRequest(ctx, db.InsertNickWeatherRequestParams{
		Nick:    params.Nick,
		Query:   q,
		City:    forecast.City.Name,
		Country: forecast.City.Country,
	})
	if err != nil {
		return err
	}

	return nil
}

func (f *forecast) Format() (string, error) {
	type report struct {
		Temp float64
		Time time.Time
	}

	reps := make([]report, len(f.List))

	location := time.FixedZone("", f.City.Timezone)

	for i, r := range f.List {
		tp, err := time.Parse("2006-01-02 15:04:05", r.DtTxt)
		if err != nil {
			return "", err
		}
		reps[i] = report{Temp: r.Main.Temp, Time: tp}
	}

	type dayHighLow struct {
		Day  time.Time
		High float64
		Low  float64
	}

	dhlsa := []dayHighLow{}
	curr := ""
	i := -1
	for _, r := range reps {
		date := r.Time.In(location).Format(time.DateOnly)
		if curr != date {
			curr = date
			i++
			dhl := dayHighLow{Day: r.Time, High: r.Temp, Low: r.Temp}
			dhlsa = append(dhlsa, dhl)
		} else {
			if r.Temp > dhlsa[i].High {
				dhlsa[i].High = r.Temp
			}
			if r.Temp < dhlsa[i].Low {
				dhlsa[i].Low = r.Temp
			}
		}
	}

	arr := []string{}
	for _, v := range dhlsa {
		arr = append(arr, fmt.Sprintf("%s:%0.0f°/%0.0f°", v.Day.In(location).Format("Mon"), v.Low, v.High))
	}
	return strings.Join(arr, " "), nil
}
