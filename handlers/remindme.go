package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model/reminders"
	"goirc/util"
	"log"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

func RemindMe(params bot.HandlerParams) bool {
	duration := params.Matches[1]
	what := params.Matches[2]

	when, err := remind(params.Nick, duration, what)
	if err != nil {
		params.Privmsgf(params.Target, "%s", err)
		return true
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		params.Privmsgf(params.Target, "%s: error: %s", err)
		return true
	}
	localFormat := when.In(loc).Format(time.RFC1123)

	params.Privmsgf(params.Target, "%s: reminder set for %s\n", params.Nick, localFormat)

	return true
}

func remind(nick string, dur string, what string) (*time.Time, error) {
	d, err := str2duration.ParseDuration(dur)
	if err != nil {
		return nil, err
	}

	at := time.Now().Add(d)

	err = reminders.Create(nick, at, what)
	if err != nil {
		return nil, err
	}

	return &at, nil
}

func DoRemind(params bot.HandlerParams) bool {
	row, err := reminders.NextDue(params.Target)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("DoRemind NextDue error: %s", err)
		}
		return false
	}

	ago := util.Ago(time.Now().Sub(row.CreatedAt).Round(time.Second))
	params.Privmsgf(params.Target, `%s: reminder (%s ago) "%s"`, row.Nick, ago, row.What)

	err = reminders.Delete(row.ID)
	if err != nil {
		log.Printf("DoRemind NextDue error: %s", err)
		return false
	}

	return true
}
