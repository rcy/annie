package irc

import (
	"github.com/jmoiron/sqlx"
)

type Params struct {
	Privmsgf func(string, string, ...interface{})
	Db       *sqlx.DB
	Msg      string
	Nick     string
	Target   string
}

type HandlerFunction func(Params) bool
