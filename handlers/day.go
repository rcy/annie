package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/shell"
	"strings"
	"time"
)

var url = "https://www.daysoftheyear.com/today/"

func fetchDay() (string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Day$'`, url)
	return shell.Command(cmd)
}

func fetchWeek() (string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Week'`, url)
	return shell.Command(cmd)
}

func fetchMonth() (string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Month'`, url)
	return shell.Command(cmd)
}

func NationalDay(params bot.HandlerParams) error {
	r, err := fetchDay()
	if err != nil {
		return err
	}

	r = strings.TrimSpace(r)
	days := strings.Split(r, "\n")
	for _, msg := range days {
		params.Privmsgf(params.Target, "%s", msg)
		time.Sleep(30 * time.Second)
	}

	params.Privmsgf(params.Target, "according to %s", url)

	return nil
}

func NationalWeek(params bot.HandlerParams) error {
	r, err := fetchWeek()
	if err != nil {
		return err
	}

	r = strings.TrimSpace(r)
	r = strings.Join(strings.Split(r, "\n"), ", ")
	r += " (according to " + url + ")"

	params.Privmsgf(params.Target, "%s", r)

	return nil
}

func NationalMonth(params bot.HandlerParams) error {
	r, err := fetchMonth()
	if err != nil {
		return err
	}

	r = strings.TrimSpace(r)
	r = strings.Join(strings.Split(r, "\n"), ", ")
	r += " (according to " + url + ")"

	params.Privmsgf(params.Target, "%s", r)

	return nil
}
