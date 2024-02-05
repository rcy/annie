package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/shell"
	"strings"
)

var url = "https://www.daysoftheyear.com/today/"

func fetchDay() (string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Day$'`, url)
	return shell.Command(cmd)
}

func NationalDay(params bot.HandlerParams) error {
	r, err := fetchDay()
	if err != nil {
		return err
	}

	r = strings.TrimSpace(r)
	r = strings.Join(strings.Split(r, "\n"), ", ")
	r += " (according to " + url + ")"

	params.Privmsgf(params.Target, "%s", r)

	return nil
}
