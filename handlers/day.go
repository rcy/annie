package handlers

import (
	"goirc/bot"
	"goirc/shell"
)

func NationalDay(params bot.HandlerParams) error {
	r, err := shell.Command(`curl -s https://nationaltoday.com/ | pup '.holiday-title-text json{}' | jq -r .[0].text`)
	if err != nil {
		params.Privmsgf(params.Target, "error: %v", err)
	}

	params.Privmsgf(params.Target, "%s", r)

	return nil
}
