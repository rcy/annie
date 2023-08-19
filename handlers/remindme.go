package handlers

import (
	"database/sql"
	"goirc/bot"
	"goirc/model/reminders"
	"goirc/util"
	"log"
	"regexp"
	"time"

	"github.com/xhit/go-str2duration/v2"
)

func RemindMe(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^!remindme ([^\s]+) (.+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) == 0 {
		return false
	}

	duration := string(matches[1])
	what := string(matches[2])

	when, err := remind(params.Nick, duration, what)
	if err != nil {
		params.Privmsgf(params.Target, "%s", err)
		return true
	}

	params.Privmsgf(params.Target, "%s: reminding at %s %s\n", params.Nick, when.Format(time.DateTime), what)

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
