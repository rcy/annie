package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model/reminders"
	"goirc/util"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

func RemindMe(params bot.HandlerParams) error {
	duration := params.Matches[1]
	what := params.Matches[2]

	when, err := remind(params.Nick, duration, what)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	localFormat := when.In(loc).Format(time.RFC1123)

	params.Privmsgf(params.Target, "%s: reminder set for %s\n", params.Nick, localFormat)

	return err
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

func DoRemind(params bot.HandlerParams) error {
	row, err := reminders.NextDue(params.Target)
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return nil
	}

	ago := util.Ago(time.Since(row.CreatedAt).Round(time.Second))
	params.Privmsgf(params.Target, `%s: reminder (%s ago) "%s"`, row.Nick, ago, row.What)

	err = reminders.Delete(row.ID)
	if err != nil {
		return err
	}

	return nil
}
