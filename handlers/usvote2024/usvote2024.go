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
	return fmt.Sprintf("%s %s (%s)", p.Tally, p.Name, p.Votes)
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
		results.Rep.Votes = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.div-dem", func(e *colly.HTMLElement) {
		results.Dem.Tally = e.ChildText("div.tally-numbers")
		results.Dem.Name = e.ChildText("div.tally-party")
	})
	c.OnHTML("p.div-dem", func(e *colly.HTMLElement) {
		results.Dem.Votes = strings.TrimSpace(e.Text)
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

	params.Privmsgf(params.Target, "%s", results.Rep)
	params.Privmsgf(params.Target, "%s", results.Dem)

	return nil
}

func HandleHarris(params bot.HandlerParams) error {
	results, err := count()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", results.Dem)
	params.Privmsgf(params.Target, "%s", results.Rep)

	return nil
}
