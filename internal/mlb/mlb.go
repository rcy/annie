package mlb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goirc/fetch"
	"goirc/util"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type teamEndData struct {
	PoffTitle float64
	WsWin     float64
	CsWin     float64
}

type team struct {
	TeamID   int `json:"teamId"`
	AbbName  string
	League   string
	Division string
	EndData  teamEndData
}

type teamList []team

func (tl teamList) String() string {
	arr := []string{}
	for _, team := range tl {
		if team.EndData.WsWin != 0.0 {
			arr = append(arr, fmt.Sprintf("%s:%.0f%%", team.AbbName, 100*team.EndData.WsWin))
		}
	}
	return strings.Join(arr, " ")
}

func fetchTeams() (teamList, error) {
	date := time.Now().Format(time.DateOnly)

	url := fmt.Sprintf("https://www.fangraphs.com/api/playoff-odds/odds?dateEnd=%s&dateDelta=&projectionMode=2&standingsType=lg", date)

	_, bytes, err := fetch.Get(url, time.Minute)
	if err != nil {
		return nil, err
	}

	var teams teamList
	err = json.Unmarshal(bytes, &teams)
	if err != nil {
		return nil, err
	}

	sort.Slice(teams, func(i, j int) bool {
		return teams[i].EndData.WsWin > teams[j].EndData.WsWin
	})

	return teams, nil
}

func fetchLeagueTeams(league string) (teamList, error) {
	teams, err := fetchTeams()
	if err != nil {
		return nil, err
	}

	var lt teamList

	for _, t := range teams {
		if t.League == league {
			lt = append(lt, t)
		}
	}

	return lt, nil
}

func PlayoffOdds() (string, error) {
	teams, err := fetchLeagueTeams("AL")
	if err != nil {
		return "", err
	}
	al := fmt.Sprintf("AL %s", teams.String())

	teams, err = fetchLeagueTeams("NL")
	if err != nil {
		return "", err
	}
	nl := fmt.Sprintf("NL %s", teams.String())

	at, err := lastUpdatedAt()
	if err != nil {
		return "", err
	}
	lastUpdatedStr := fmt.Sprintf("%s ago", util.Ago(time.Since(*at)))

	return fmt.Sprintf("%s - %s - %s - %s", al, nl, lastUpdatedStr, "https://www.mlb.com/postseason"), nil
}

func lastUpdatedAt() (*time.Time, error) {
	code, data, err := fetch.Get("https://www.fangraphs.com/standings/playoff-odds", time.Minute)
	if err != nil {
		return nil, err
	}
	if code > 299 {
		return nil, fmt.Errorf("Bad status: %d", code)
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	text := doc.Find(".footer-bar-item").Find(".mobile-hide").Text()
	// Updated: Friday, October 20, 2023 11:12 PM ET

	timeStr, ok := strings.CutPrefix(text, "Updated: ")
	if !ok {
		return nil, fmt.Errorf("couldn't cut prefix from %s", timeStr)
	}
	timeStr, ok = strings.CutSuffix(timeStr, " ET")
	if !ok {
		return nil, fmt.Errorf("couldn't cut prefix from %s", timeStr)
	}

	location, err := time.LoadLocation("America/New_York")
	if err != nil {
		return nil, err
	}

	result, err := time.ParseInLocation("Monday, January 2, 2006 15:04 PM", timeStr, location)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
