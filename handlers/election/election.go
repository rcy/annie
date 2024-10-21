package election

import (
	"encoding/json"
	"fmt"
	"goirc/bot"
	"io"
	"net/http"

	"golang.org/x/text/message"
)

type PartyE struct {
	PartyCode              string `json:"partyCode"`
	PartyName              string `json:"partyName"`
	PartyNameShort         string `json:"partyNameShort"`
	WinnerParty            string `json:"winnerParty"`
	ColourDarkElectedText  string `json:"colourDarkElectedText"`
	ColourDarkLeadingText  string `json:"colourDarkLeadingText"`
	ColourLightElectedText string `json:"colourLightElectedText"`
	ColourLightLeadingText string `json:"colourLightLeadingText"`
	DisplayParty           string `json:"displayParty"`
	LeaderID               string `json:"leaderId"`
	ColourLightLeading     string `json:"colourLightLeading"`
	ColourLightElected     string `json:"colourLightElected"`
	ColourDarkLeading      string `json:"colourDarkLeading"`
	ColourDarkElected      string `json:"colourDarkElected"`
}

type Data struct {
	Timestamp int  `json:"ts"`
	Live      bool `json:"live"`
	Data      struct {
		Config struct {
			EnableTracking bool   `json:"enableTracking"`
			WidgetCTA      string `json:"widgetCTA"`
			RsURL          string `json:"rsURL"`
			VideoID        string `json:"video_id"`
			VideoEnabled   bool   `json:"video_enabled"`
			Updates        []any  `json:"updates"`
			Language       struct {
			} `json:"language"`
			ClientFetchInterval string `json:"clientFetchInterval"`
			WidgetTitle         string `json:"widgetTitle"`
			WidgetDescription   string `json:"widgetDescription"`
			EnableSharing       bool   `json:"enableSharing"`
			Visibility          struct {
				Hero         bool `json:"hero"`
				Developments bool `json:"developments"`
				Headline     bool `json:"headline"`
				Sticky       bool `json:"sticky"`
				Standings    bool `json:"standings"`
				Ridings      bool `json:"ridings"`
				Datatable    bool `json:"datatable"`
				Video        bool `json:"video"`
				Close        bool `json:"close"`
				Gains        bool `json:"gains"`
				GainsDetail  bool `json:"gains_detail"`
				Leaders      bool `json:"leaders"`
				Cabinet      bool `json:"cabinet"`
				Polls        bool `json:"polls"`
				Faq          bool `json:"faq"`
			} `json:"visibility"`
			ElectionState string `json:"electionState"`
			HeadlineBody  string `json:"headlineBody"`
			HeadlineID    string `json:"headlineId"`
			Headline      string `json:"headline"`
		} `json:"config"`
		Ridings []struct {
			ID                       int    `json:"id"`
			RidingNumber             int    `json:"ridingNumber"`
			EnglishName              string `json:"englishName"`
			TotalVoters              int    `json:"totalVoters"`
			TotalPolls               int    `json:"totalPolls"`
			PreviousElectedPartyCode string `json:"previousElectedPartyCode"`
			CallerCode               string `json:"callerCode"`
			IsCandidateElected       int    `json:"isCandidateElected"`
			PollsReported            int    `json:"pollsReported"`
			TotalVotesReported       int    `json:"totalVotesReported"`
			CandidateVotesLead       int    `json:"candidateVotesLead"`
			DecisionDesk             struct {
				PollsRequired    float64 `json:"pollsRequired"`
				MajorityRequired int     `json:"majorityRequired"`
				Caller           string  `json:"caller"`
				Party            string  `json:"party"`
				Suggestion       string  `json:"suggestion"`
				SuretyLevel      string  `json:"suretyLevel"`
			} `json:"decisionDesk"`
			Parties []struct {
				ID                       int    `json:"id"`
				FirstName                string `json:"firstName"`
				LastName                 string `json:"lastName"`
				PartyID                  int    `json:"partyId"`
				Incumbent                string `json:"incumbent"`
				TotalVotes               int    `json:"totalVotes"`
				TotalVotesPosition       int    `json:"totalVotesPosition"`
				TotalVotesPercentage     string `json:"totalVotesPercentage"`
				TotalVotesLead           int    `json:"totalVotesLead"`
				TotalVotesPercentageLead string `json:"totalVotesPercentageLead"`
				IsAcclaimed              int    `json:"isAcclaimed"`
				IsCandidateElected       int    `json:"isCandidateElected"`
				VotesPercentage          string `json:"votesPercentage"`
				VotesLead                int    `json:"votesLead"`
				Votes                    int    `json:"votes"`
				PartyCode                string `json:"partyCode"`
				E                        struct {
					MinisterTitle string `json:"ministerTitle"`
				} `json:"e,omitempty"`
			} `json:"parties"`
		} `json:"ridings"`
		Parties []struct {
			ID                           int     `json:"id"`
			EnglishName                  string  `json:"englishName"`
			EnglishCode                  string  `json:"englishCode"`
			Priority                     int     `json:"priority"`
			ElectedSeats                 int     `json:"electedSeats"`
			LeadingSeats                 int     `json:"leadingSeats"`
			TotalElectedLeadingSeats     int     `json:"totalElectedLeadingSeats"`
			TotalVotes                   int     `json:"totalVotes"`
			TotalVotesPercentage         float64 `json:"totalVotesPercentage"`
			PreviousElected              int     `json:"previousElected"`
			PreviousTotalVotes           int     `json:"previousTotalVotes"`
			PreviousTotalVotesPercentage float64 `json:"previousTotalVotesPercentage"`
			PercentageDifference         float64 `json:"percentageDifference"`
			Seats                        int     `json:"seats"`
			CurrentSeats                 int     `json:"currentSeats"`
			DisplayOrder                 int     `json:"displayOrder"`
			Gain                         int     `json:"gain"`
			Loss                         int     `json:"loss"`
			Hold                         int     `json:"hold"`
			Net                          int     `json:"net"`
			Undecided                    int     `json:"undecided"`
			UpstreamSlug                 string  `json:"upstream_slug"`
			E                            any     `json:"e"`
		} `json:"parties"`
		Election struct {
			PollsReported      int `json:"pollsReported"`
			TotalPolls         int `json:"totalPolls"`
			TotalVotesReported int `json:"totalVotesReported"`
		} `json:"election"`
		Extras struct {
			Cabinet []struct {
				FirstName                  string `json:"firstName"`
				LastName                   string `json:"lastName"`
				RidingID                   int    `json:"ridingId"`
				Title                      string `json:"title"`
				RidingName                 string `json:"ridingName"`
				CandidateID                int    `json:"candidateId"`
				IsCandidateElectedInRiding int    `json:"isCandidateElectedInRiding"`
				TotalVotesPosition         int    `json:"totalVotesPosition"`
			} `json:"cabinet"`
			CloseRaces []struct {
				ID       int  `json:"id"`
				IsPinned bool `json:"isPinned"`
			} `json:"closeRaces"`
		} `json:"extras"`
	} `json:"data"`
}

func electionData() (*Data, error) {
	resp, err := http.Get("https://canopy.cbc.ca/live/elections/prov/BC2024/all/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := Data{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func partyResults() ([]string, error) {
	data, err := electionData()
	if err != nil {
		return nil, err
	}

	rows := []string{}

	rows = append(rows, fmt.Sprintf("%s | %s", data.Data.Config.Headline, data.Data.Config.HeadlineBody))

	p := message.NewPrinter(message.MatchLanguage("en"))

	for _, party := range data.Data.Parties {
		if party.ElectedSeats+party.LeadingSeats == 0 {
			continue
		}
		rows = append(rows, p.Sprintf("%s: elected %02d, leading %02d, votes %d, share %04.1f%%",
			party.EnglishCode,
			party.ElectedSeats,
			party.LeadingSeats,
			party.TotalVotes,
			party.TotalVotesPercentage*100,
		))
	}
	//rows = append(rows, fmt.Sprintf("timestamp %s", time.Unix(int64(data.Timestamp), 0)))
	return rows, nil
}

func Handle(params bot.HandlerParams) error {
	rows, err := partyResults()
	if err != nil {
		return err
	}

	for _, row := range rows {
		params.Privmsgf(params.Target, "%s", row)
	}

	return nil
}
