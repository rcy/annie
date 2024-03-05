package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/shell"
	"strings"
	"time"
)

var url = "https://www.daysoftheyear.com/today/"

func fetchDays() ([]string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Day$'`, url)
	r, err := shell.Command(cmd)
	if err != nil {
		return nil, err
	}
	r = strings.TrimSpace(r)
	return strings.Split(r, "\n"), nil
}

func fetchWeeks() ([]string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Week'`, url)
	r, err := shell.Command(cmd)
	if err != nil {
		return nil, err
	}
	r = strings.TrimSpace(r)
	return strings.Split(r, "\n"), nil
}

func fetchMonths() ([]string, error) {
	cmd := fmt.Sprintf(`curl -s %s | pup 'body picture img json{}' | jq -r .[].alt | grep -E ' Month'`, url)
	r, err := shell.Command(cmd)
	if err != nil {
		return nil, err
	}
	r = strings.TrimSpace(r)
	return strings.Split(r, "\n"), nil
}

func NationalDay(params bot.HandlerParams) error {
	days, err := fetchDays()
	if err != nil {
		return err
	}

	for _, msg := range days {
		params.Privmsgf(params.Target, "%s", msg)
		time.Sleep(30 * time.Second)
	}

	params.Privmsgf(params.Target, "according to %s", url)

	return nil
}

func NationalWeek(params bot.HandlerParams) error {
	weeks, err := fetchWeeks()
	if err != nil {
		return err
	}

	for _, msg := range weeks {
		params.Privmsgf(params.Target, "%s", msg)
		time.Sleep(30 * time.Second)
	}

	params.Privmsgf(params.Target, "according to %s", url)

	return nil
}

func NationalMonth(params bot.HandlerParams) error {
	months, err := fetchMonths()
	if err != nil {
		return err
	}

	for _, msg := range months {
		params.Privmsgf(params.Target, "%s", msg)
		time.Sleep(30 * time.Second)
	}

	params.Privmsgf(params.Target, "according to %s", url)

	return nil
}
