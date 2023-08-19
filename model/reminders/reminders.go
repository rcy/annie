package reminders

import (
	"goirc/model"
	"time"
)

type Reminder struct {
	ID        int64
	CreatedAt time.Time `db:"created_at"`
	Nick      string
	RemindAt  time.Time `db:"remind_at"`
	What      string
}

func Create(nick string, at time.Time, what string) error {
	_, err := model.DB.Exec(`insert into reminders(nick, remind_at, what) values(?, ?, ?)`, nick, at.UTC().Format(time.DateTime), what)
	return err
}

func All() (rows []Reminder, err error) {
	err = model.DB.Select(&rows, `select id, created_at, nick, remind_at, what from reminders limit 100`)
	return
}

// return the next due item for a nick that is in the channel
func NextDue(channel string) (row Reminder, err error) {
	query := `
select id, created_at, reminders.nick, remind_at, what
from reminders
join channel_nicks on channel_nicks.nick = reminders.nick
where
  remind_at < current_timestamp and
  channel_nicks.present = true   
limit 1`

	err = model.DB.Get(&row, query)
	return
}

func Delete(id int64) error {
	_, err := model.DB.Exec(`delete from reminders where id = ?`, id)
	return err
}
