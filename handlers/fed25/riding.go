package fed25

import (
	"bytes"
	"cmp"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/fetch"
	"html/template"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Riding struct {
	RidingID   int `json:"riding_id"`
	Candidates []struct {
		CandidateID   int    `json:"candidate_id"`
		Firstname     string `json:"firstname"`
		Lastname      string `json:"lastname"`
		PartyShortEng string `json:"party_short_eng"`
		PartyShortFr  string `json:"party_short_fr"`
		Votes         int    `json:"votes"`
		Margin        int    `json:"margin"`
		Previous      int    `json:"previous"`
		Prominent     string `json:"prominent"`
	} `json:"candidates"`
	ElectedCandidateID   int    `json:"elected_candidate_id"`
	ElectedPartyShortEng string `json:"elected_party_short_eng"`
	LastUpdated          string `json:"last_updated"`
	LeadingParties       []any  `json:"leading_parties"`
	PollsReporting       int    `json:"polls_reporting"`
	RidingNameEng        string `json:"riding_name_eng"`
	TotalPolls           int    `json:"total_polls"`
	TotalVotes           int    `json:"total_votes"`
	RegionID             int    `json:"region_id"`
	RegionNameEng        string `json:"region_name_eng"`
	RidingCode           int    `json:"riding_code"`
	TgamKeyFlag          bool   `json:"tgam_key_flag"`
}

func fetchRidings() ([]Riding, error) {
	_, data, err := fetch.Get("https://federal-election-2025-prod-dc5q4g5x5w7l.s3.ca-central-1.amazonaws.com/federal-election-ridings.json", time.Minute)
	if err != nil {
		return nil, err
	}

	summary := make(map[string]Riding)

	err = json.Unmarshal(data, &summary)
	if err != nil {
		return nil, err
	}

	ridings := make([]Riding, 0, len(summary))
	for _, riding := range summary {
		ridings = append(ridings, riding)
	}

	slices.SortStableFunc(ridings, func(a Riding, b Riding) int {
		return cmp.Compare(a.RidingNameEng, b.RidingNameEng)
	})

	return ridings, nil
}

// Look up ridings by text that matches the riding name
func findRidingsByName(text string) ([]Riding, error) {
	ridings, err := fetchRidings()
	if err != nil {
		return nil, err
	}

	var matches []Riding

	text = strings.ToLower(text)

	for _, riding := range ridings {
		if strings.Contains(strings.ToLower(riding.RidingNameEng), text) {
			matches = append(matches, riding)
			// if len(matches) >= 5 {
			// 	break
			// }
		}
	}

	return matches, nil
}

//go:embed riding_short.tmpl
var ridingShortTemplateContent string
var ridingShortTemplate = template.Must(template.New("").Parse(ridingShortTemplateContent))

func findRidingsByNameSummaries(text string) ([]string, error) {
	ridings, err := findRidingsByName(text)
	if err != nil {
		return nil, err
	}

	summaries := make([]string, 0, len(ridings))

	for _, riding := range ridings {
		fmt.Println("riding", riding.RidingID)
		var buf bytes.Buffer
		err := ridingShortTemplate.Execute(&buf, riding)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, buf.String())
	}

	return summaries, nil
}

func FindRidingsByNameHandler(params bot.HandlerParams) error {
	text := params.Matches[1]
	ridings, err := findRidingsByName(text)
	if err != nil {
		return err
	}

	if len(ridings) == 0 {
		params.Privmsgf(params.Target, "%s", "no matching ridings")
		return nil
	}

	if len(ridings) == 1 {
		summary, err := findRidingByIDSummary(ridings[0].RidingID)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s", summary)
		return nil
	}

	const max = 9
	for i, riding := range ridings {
		var buf bytes.Buffer
		err := ridingShortTemplate.Execute(&buf, riding)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s", buf.String())
		if i > max {
			params.Privmsgf(params.Target, "%d more '%s' not shown", len(ridings)-max, text)
			return nil
		}
	}

	return nil
}

// Look up ridings by text that matches the riding name or candidate name
func findRidingByID(id int) (Riding, error) {
	ridings, err := fetchRidings()
	if err != nil {
		return Riding{}, err
	}

	for _, riding := range ridings {
		if riding.RidingID == id {
			return riding, nil
		}
	}

	return Riding{}, errors.New("riding not found")
}

//go:embed riding.tmpl
var ridingTemplateContent string
var ridingTemplate = template.Must(template.New("").Parse(ridingTemplateContent))

func findRidingByIDSummary(id int) (string, error) {
	riding, err := findRidingByID(id)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = ridingTemplate.Execute(&buf, riding)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func RidingHandler(params bot.HandlerParams) error {
	id, _ := strconv.Atoi(params.Matches[1])
	summary, err := findRidingByIDSummary(id)
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", summary)

	return nil
}
