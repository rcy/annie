package model

import (
	"goirc/db"
	"goirc/util"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func init() {
	DB = db.Open(util.Getenv("SQLITE_DB"))
}

func Close() {
	DB.Close()
}

type ChannelNick struct {
	Channel   string
	Nick      string
	Present   string
	UpdatedAt string `db:"updated_at"`
}
