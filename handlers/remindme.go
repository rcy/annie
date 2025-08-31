package handlers

import (
	"database/sql"
	"goirc/internal/responder"
	"goirc/model/reminders"
	"goirc/util"
	"strings"
	"time"

	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
)

func RemindMe(params responder.Responder) error {
	input := params.Match(1)

	at, what, err := parseTimeAndTask(input)
	if err != nil {
		return err
	}

	result, err := reminders.Create(params.Nick(), at, what)
	if err != nil {
		return err
	}

	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return err
	}
	localFormat := at.In(loc).Format(time.RFC1123)

	id, _ := result.LastInsertId()
	params.Privmsgf(params.Target(), "%s: '%s' reminder set for %s [%d]\n", params.Nick(), what, localFormat, id)

	return err
}

func DoRemind(params responder.Responder) error {
	row, err := reminders.NextDue(params.Target())
	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
		return nil
	}

	ago := util.Ago(time.Since(row.CreatedAt).Round(time.Second))
	params.Privmsgf(params.Target(), `%s: reminder (set %s ago) "%s"`, row.Nick, ago, row.What)

	err = reminders.Delete(row.ID)
	if err != nil {
		return err
	}

	return nil
}

func parseTimeAndTask(input string) (time.Time, string, error) {
	w := when.New(nil)
	w.Add(en.All...)

	now := time.Now()
	result, err := w.Parse(input, now)
	if err != nil {
		return time.Time{}, "", err
	}
	if result == nil {
		return time.Time{}, input, nil
	}

	remaining := strings.TrimSpace(strings.Replace(input, result.Text, "", 1))

	return result.Time, remaining, nil
}
