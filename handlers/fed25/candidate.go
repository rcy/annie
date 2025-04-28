package fed25

import (
	"bytes"
	_ "embed"
	"goirc/bot"
	"html/template"
	"strings"
)

type XCandidateWithRiding struct {
	Candidate Candidate
	Riding    Riding
}

// Find candidates and their riding
func findCandidatesWithRidingByName(text string) ([]XCandidateWithRiding, error) {
	ridings, err := fetchRidings()
	if err != nil {
		return nil, err
	}

	var matches []XCandidateWithRiding

	text = strings.ToLower(text)

	for _, riding := range ridings {
		for _, candidate := range riding.Candidates {
			if strings.Contains(strings.ToLower(candidate.Firstname), text) ||
				strings.Contains(strings.ToLower(candidate.Lastname), text) {
				matches = append(matches, XCandidateWithRiding{
					Candidate: candidate,
					Riding:    riding,
				})
			}
		}
	}

	return matches, nil
}

//go:embed candidate.tmpl
var candidateTemplateContent string
var candidateTemplate = template.Must(template.New("").Parse(candidateTemplateContent))

func FindCandidatesHandler(params bot.HandlerParams) error {
	text := params.Matches[1]
	candidates, err := findCandidatesWithRidingByName(text)
	if err != nil {
		return err
	}

	if len(candidates) == 0 {
		params.Privmsgf(params.Target, "%s", "no matching candidates")
		return nil
	}

	if len(candidates) == 1 {
		summary, err := findRidingByIDSummary(candidates[0].Riding.RidingID)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s", summary)
		return nil
	}

	const max = 9
	for i, candidate := range candidates {
		var buf bytes.Buffer
		err := candidateTemplate.Execute(&buf, candidate)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "%s", buf.String())
		if i > max {
			params.Privmsgf(params.Target, "%d more '%s' not shown", len(candidates)-max, text)
			return nil
		}
	}

	return nil
}
