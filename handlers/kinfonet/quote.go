package kinfonet

import (
	"goirc/bot"
	"time"

	"github.com/gocolly/colly"
)

func TodaysQuoteHandler(params bot.HandlerParams) error {
	quote, err := todaysQuote()
	if err != nil {
		return err
	}

	params.Privmsgf(params.Target, "%s", quote.url)

	return nil
}

type quote struct{ date, text, url string }

type date string

func todaysQuote() (quote, error) {
	quotes, err := getQuotes()
	if err != nil {
		return quote{}, err
	}
	today := time.Now().Format("Jan 02")
	return quotes[date(today)], nil
}

func getQuotes() (map[date]quote, error) {
	base := "https://www.kinfonet.org"
	c := colly.NewCollector()

	quotes := map[date]quote{}

	c.OnHTML("a.daily-quote-slug", func(e *colly.HTMLElement) {
		quotes[date(e.ChildText(".daily-quote-date"))] = quote{
			url:  base + e.Attr("href"),
			text: e.ChildText(".daily-quote-slug"),
		}
	})

	err := c.Visit(base + "/daily_quotes")
	if err != nil {
		return nil, err
	}

	return quotes, nil
}
