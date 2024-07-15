package day

import (
	"context"
	"database/sql"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	db "goirc/model"
	"goirc/shell"
	"strings"
)

var url = "https://www.daysoftheyear.com/today/"

var dayCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Day$'`
var weekCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Week$'`
var monthCmd = `curl -s https://www.daysoftheyear.com/today/ | pup 'body img json{}' | jq -r .[].alt | grep -E ' Month$'`

var dayCache = NewCache(dayCmd)
var weekCache = NewCache(weekCmd)
var monthCache = NewCache(monthCmd)

func NationalDay(params bot.HandlerParams) error {
	str, err := dayCache.Pop()
	if err != nil {
		return err
	}

	if str == "EOF" {
		link, err := dayImage(params, dayCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "Today's image: %s", link)
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalWeek(params bot.HandlerParams) error {
	str, err := weekCache.Pop()
	if err != nil {
		return err
	}

	if str == "EOF" {
		link, err := dayImage(params, weekCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "This week's image: %s", link)
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalMonth(params bot.HandlerParams) error {
	str, err := monthCache.Pop()
	if err != nil {
		return err
	}

	if str == "EOF" {
		link, err := dayImage(params, monthCmd)
		if err != nil {
			return err
		}
		params.Privmsgf(params.Target, "This month's image: %s", link)
	} else {
		params.Privmsgf(params.Target, "%s", str)
	}

	return nil
}

func NationalRefs(params bot.HandlerParams) error {
	params.Privmsgf(params.Target, "%s", url)

	return nil
}

func dayImage(params bot.HandlerParams, cmd string) (string, error) {
	r, err := shell.Command(cmd)
	if err != nil {
		return "", err
	}

	days := strings.Split(strings.TrimSpace(r), "\n")
	prompt := fmt.Sprintf("a scene incorporating themes from: %s", strings.Join(days, ","))
	url, err := generateImage(context.Background(), prompt)
	if err != nil {
		return "", err
	}

	q := model.New(db.DB)
	note, err := q.InsertNote(context.TODO(), model.InsertNoteParams{
		Target: params.Target,
		Nick:   sql.NullString{String: params.Nick, Valid: true},
		Kind:   "link",
		Text:   sql.NullString{String: url, Valid: true},
		Anon:   true,
	})
	if err != nil {
		return "", err
	}

	link, err := note.Link()
	if err != nil {
		return "", err
	}

	return link, nil
}
