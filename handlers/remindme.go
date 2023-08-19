package handlers

import (
	"goirc/bot"
	"goirc/model/reminders"
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

	params.Privmsgf(params.Target, "I will remind %s on %s at %s \"%s\"\n",
		params.Nick,
		when.Format(time.DateOnly),
		when.Format("15:04 MST"),
		string(matches[2]))

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
