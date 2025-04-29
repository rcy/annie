package fed25

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"fmt"
	"goirc/bot"
	"goirc/fetch"
	"html/template"
	"reflect"
	"slices"
	"strings"
	"time"
)

type OverviewPartyDetails struct {
	PartyShortEng         string  `json:"party_short_eng"`
	PartyNameEng          string  `json:"party_name_eng"`
	Votes                 int     `json:"votes"`
	PopularVotePercentage float32 `json:"popular_vote_percentage"`
	Elected               int     `json:"elected"`
	Leading               int     `json:"leading"`
}

type FederalElectionSummary struct {
	ElectionID      string `json:"election_id"`
	SummaryID       int    `json:"summary_id"`
	ElectionNameEng string `json:"election_name_eng"`
	LastUpdated     string `json:"last_updated"`
	TotalRidings    int    `json:"total_ridings"`
	RidingsCalled   int    `json:"ridings_called"`
	TotalPolls      int    `json:"total_polls"`
	ReportedPolls   int    `json:"reported_polls"`
	ElectedParty    struct {
		PartyShortEng        string `json:"party_short_eng"`
		CallType             string `json:"call_type"`
		EnglishCallStatement string `json:"english_call_statement"`
	} `json:"elected_party"`
	PercentagePolls      int                    `json:"percentage_polls"`
	CurrentLeaderParties []any                  `json:"current_leader_parties"`
	OverviewPartyDetails []OverviewPartyDetails `json:"overview_party_details"`
}

func fetchSummary() (*FederalElectionSummary, error) {
	_, data, err := fetch.Get("https://federal-election-2025-prod-dc5q4g5x5w7l.s3.ca-central-1.amazonaws.com/federal-election-summary.json", 10*time.Second)
	if err != nil {
		return nil, err
	}

	summary := &FederalElectionSummary{}

	err = json.Unmarshal(data, summary)
	if err != nil {
		return nil, err
	}

	return summary, nil
}

type partySeats struct {
	party string
	seats int
}

var leaderboard = []partySeats{}

func LeaderboardHandler(params bot.HandlerParams) error {
	summary, err := fetchSummary()
	if err != nil {
		return err
	}

	var newLeaderboard []partySeats

	for _, p := range summary.OverviewPartyDetails {
		newLeaderboard = append(newLeaderboard, partySeats{party: p.PartyShortEng, seats: p.Elected + p.Leading})
	}

	slices.SortStableFunc(newLeaderboard, func(a partySeats, b partySeats) int {
		return cmp.Compare(b.seats, a.seats)
	})

	if !reflect.DeepEqual(leaderboard, newLeaderboard) {
		leaderboard = newLeaderboard

		display := []string{}
		seats := 0
		for _, i := range leaderboard {
			if i.seats > 0 {
				display = append(display, fmt.Sprintf("%s %d", i.party, i.seats))
			}
		}

		params.Privmsgf(params.Target, "%s (%d seats to come)", strings.Join(display, ", "), 343-seats)
	}

	return nil
}

//go:embed fed25.tmpl
var summaryTemplateContent string
var summaryTemplate = template.Must(template.New("").Parse(summaryTemplateContent))

func formatSummary(results *FederalElectionSummary) (string, error) {
	var buf bytes.Buffer

	slices.SortStableFunc(results.OverviewPartyDetails, func(a OverviewPartyDetails, b OverviewPartyDetails) int {
		res := cmp.Compare(b.Elected+b.Leading, a.Elected+a.Leading)
		if res == 0 {
			return cmp.Compare(b.PopularVotePercentage, a.PopularVotePercentage)
		}
		return res
	})

	err := summaryTemplate.Execute(&buf, results)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Handler(params bot.HandlerParams) error {
	summary, err := fetchSummary()
	if err != nil {
		return err
	}
	str, err := formatSummary(summary)

	params.Privmsgf(params.Target, "%s", str)

	return nil
}
