package mlb

import (
	"encoding/json"
	"goirc/bot"
	"net/http"
	"time"
)

type Game struct {
	Schedule struct {
		GameDateTimeUTC time.Time `json:"GameDateTimeUTC"`
		GameDateParam   string    `json:"GameDateParam"`
		Season          int       `json:"Season"`
		GameID          int       `json:"GameId"`
		TempDate        string    `json:"TempDate"`
		TempTime        string    `json:"TempTime"`
		Dh              int       `json:"DH"`
		AwayTeamID      int       `json:"AwayTeamId"`
		HomeTeamID      int       `json:"HomeTeamId"`
		GameDate        string    `json:"GameDate"`
		Hour            int       `json:"Hour"`
		Minute          int       `json:"Minute"`
		ProbHome        int       `json:"ProbHome"`
		ProbAway        int       `json:"ProbAway"`
		MLBGameID       int       `json:"MLBGameId"`
		Gamestate       string    `json:"gamestate"`
		HomeStarter     string    `json:"HomeStarter"`
		AwayStarter     string    `json:"AwayStarter"`
		HomeTeamName    string    `json:"HomeTeamName"`
		AwayTeamName    string    `json:"AwayTeamName"`
		HomeTeamAbbName string    `json:"HomeTeamAbbName"`
		AwayTeamAbbName string    `json:"AwayTeamAbbName"`
		League          string    `json:"League"`
		HomeGameOdds    float64   `json:"HomeGameOdds"`
		AwayGameOdds    float64   `json:"AwayGameOdds"`
	} `json:"schedule"`
	Scores struct {
		GameDate   string  `json:"GameDate"`
		HomeTeamID int     `json:"HomeTeamId"`
		AwayTeamID int     `json:"AwayTeamId"`
		Dh         int     `json:"DH"`
		WinTeamID  int     `json:"WinTeamId"`
		IsFinal    int     `json:"isFinal"`
		Inning     int     `json:"Inning"`
		InningHalf int     `json:"InningHalf"`
		HomeScore  int     `json:"HomeScore"`
		AwayScore  int     `json:"AwayScore"`
		Outs       int     `json:"Outs"`
		Has1B      int     `json:"has1B"`
		Has2B      int     `json:"has2B"`
		Has3B      int     `json:"has3B"`
		LiveWEHome float64 `json:"LiveWEHome"`
		LiveWEAway float64 `json:"LiveWEAway"`
		Ih         string  `json:"IH"`
		HomeAbb    string  `json:"HomeAbb"`
		AwayAbb    string  `json:"AwayAbb"`
		HomeName   string  `json:"HomeName"`
		AwayName   string  `json:"AwayName"`
		HomeLeague string  `json:"HomeLeague"`
	} `json:"scores"`
	IsLiveData bool `json:"isLiveData"`
}

func GameOdds(params bot.HandlerParams) error {
	g, err := fetchGameOdds()
	if err != nil {
		return err
	}

	half := "top"
	if g[0].Scores.InningHalf == 1 {
		half = "bot"
	}

	bases := ""
	if g[0].Scores.Has3B == 1 {
		bases += "<"
	} else {
		bases += "."
	}
	if g[0].Scores.Has2B == 1 {
		bases += "^"
	} else {
		bases += "."
	}
	if g[0].Scores.Has1B == 1 {
		bases += ">"
	} else {
		bases += "."
	}

	params.Privmsgf(params.Target, "%s %d [%.0f%%], %s %d [%.0f%%], %d out %s %d %s",
		g[0].Schedule.AwayTeamName,
		g[0].Scores.AwayScore,
		g[0].Scores.LiveWEAway*100,
		g[0].Schedule.HomeTeamName,
		g[0].Scores.HomeScore,
		g[0].Scores.LiveWEHome*100,
		g[0].Scores.Outs,
		half,
		g[0].Scores.Inning,
		bases)

	return nil
}

func fetchGameOdds() ([]Game, error) {
	resp, err := http.Get("https://www.fangraphs.com/api/scores/most-recent")
	if err != nil {
		return nil, err
	}

	games := []Game{}

	err = json.NewDecoder(resp.Body).Decode(&games)
	if err != nil {
		return nil, err
	}

	return games, nil
}
