package handlers

import (
	"fmt"
	"goirc/bot"
	"sort"
	"strconv"
	"strings"
	"time"
)

func MLBOdds(params bot.HandlerParams) error {
	team := strings.ToUpper(params.Matches[1])

	date := time.Now().Format(time.DateOnly)

	r, err := shell(fmt.Sprintf(`
curl -s 'https://www.fangraphs.com/api/playoff-odds/odds?dateEnd=%s&dateDelta=&projectionMode=2&standingsType=div' |\
jq .[] |\
jq 'select(.abbName == "%s")'.endData.poffTitle
`, date, team))
	if err != nil {
		return err
	}

	f, err := strconv.ParseFloat(strings.TrimSpace(r), 32)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "FanGraphs projects %s to have a %.3f%% chance of making the playoffs", team, f*100)

	return nil
}

func MLBTeams(params bot.HandlerParams) error {
	date := time.Now().Format(time.DateOnly)

	r, err := shell(fmt.Sprintf(`
curl -s 'https://www.fangraphs.com/api/playoff-odds/odds?dateEnd=%s&dateDelta=&projectionMode=2&standingsType=div' | jq -r .[].abbName
`, date))
	if err != nil {
		return err
	}

	teams := strings.Split(r, "\n")
	sort.Slice(teams, func(i, j int) bool {
		return teams[i] < teams[j]
	})

	params.Privmsgf(params.Target, "Teams: %s", strings.Join(teams, " "))

	return nil
}
