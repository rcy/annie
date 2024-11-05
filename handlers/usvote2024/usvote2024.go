package usvote2024

import (
	"fmt"
	"goirc/bot"
	"strings"

	"github.com/gocolly/colly"
)

type party struct {
	Name  string
	Tally string
	Votes string
}

func (p party) String() string {
	return fmt.Sprintf("%s %s %s", p.Name, p.Tally, p.Votes)
}

type results struct {
	Rep party
	Dem party
}

func count() (*results, error) {
	base := "https://www.reuters.com/graphics/USA-ELECTION/RESULTS/zjpqnemxwvx/"
	c := colly.NewCollector()

	results := results{}

	c.OnHTML("div.div-rep", func(e *colly.HTMLElement) {
		results.Rep.Tally = e.ChildText("div.tally-numbers")
		results.Rep.Name = e.ChildText("div.tally-party")
	})
	c.OnHTML("p.div-rep", func(e *colly.HTMLElement) {
		results.Rep.Votes = strings.TrimSuffix(strings.TrimSpace(e.Text), " votes")
	})

	c.OnHTML("div.div-dem", func(e *colly.HTMLElement) {
		results.Dem.Tally = e.ChildText("div.tally-numbers")
		results.Dem.Name = e.ChildText("div.tally-party")
	})
	c.OnHTML("p.div-dem", func(e *colly.HTMLElement) {
		results.Dem.Votes = strings.TrimSuffix(strings.TrimSpace(e.Text), " votes")
	})

	err := c.Visit(base)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

func HandleTrump(params bot.HandlerParams) error {
	results, err := count()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (%s) %s | %s (%s) %s",
		results.Rep.Name, results.Rep.Votes, results.Rep.Tally, results.Dem.Tally, results.Dem.Votes, results.Dem.Name,
	)

	return nil
}

func HandleHarris(params bot.HandlerParams) error {
	results, err := count()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s (%s) %s | %s (%s) %s",
		results.Dem.Name, results.Dem.Votes, results.Dem.Tally, results.Rep.Tally, results.Rep.Votes, results.Rep.Name,
	)

	return nil
}
