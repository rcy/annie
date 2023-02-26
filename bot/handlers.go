package bot

import (
	"github.com/jmoiron/sqlx"
)

type HandlerParams struct {
	Privmsgf func(string, string, ...interface{})
	Db       *sqlx.DB
	Msg      string
	Nick     string
	Target   string
}

type HandlerFunction func(HandlerParams) bool
