package laters

import "goirc/model"

type Later struct {
	RowId     int    `db:"rowid"`
	CreatedAt string `db:"created_at"`
	Nick      string
	Target    string
	Message   string
	Sent      bool
}

func Get() (rows []Later, err error) {
	err = model.DB.Select(&rows, `select rowid, created_at, nick, target, message, sent from laters limit 100`)
	return
}
