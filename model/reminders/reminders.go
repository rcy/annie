package reminders

import (
	"goirc/model"
	"time"
)

type Reminder struct {
	Id        int64
	CreatedAt string `db:"created_at"`
	Nick      string
	RemindAt  time.Time `db:"remind_at"`
	What      string
}

func Create(nick string, at time.Time, what string) error {
	_, err := model.DB.Exec(`insert into reminders(nick, remind_at, what) values(?, ?, ?)`, nick, at, what)
	return err
}

func All() (rows []Reminder, err error) {
	err = model.DB.Select(&rows, `select id, created_at, nick, remind_at, what from reminders limit 100`)
	return
}
