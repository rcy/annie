package handlers

import (
	"goirc/model"

	"github.com/jmoiron/sqlx"
)

type replyFunction func()

type Params struct {
	Privmsgf func(string, string, ...interface{})
	Db       *sqlx.DB
	Msg      string
	Nick     string
	Target   string
}

type handlerFunction func(Params) bool

var Functions = []handlerFunction{
	createNote,
	deferredDelivery,
	link,
	feedMe,
	catchup,
	worldcup,
	ticker,
	quote,
	trade,
	report,
}

func insertNote(db *sqlx.DB, target string, nick string, kind string, text string) error {
	result, err := db.Exec(`insert into notes(nick, text, kind) values(?, ?, ?) returning id`, nick, text, kind)
	if err != nil {
		return err
	} else {
		noteId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		err = markAsSeen(db, noteId, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func markAsSeen(db *sqlx.DB, noteId int64, target string) error {
	//db.Select(`select * from channel_nicks where channel = ?`, target)
	channelNicks, err := model.JoinedNicks(db, target)
	if err != nil {
		return err
	}
	// for each channelNick insert a seen_by record
	for _, nick := range channelNicks {
		_, err := db.Exec(`insert into seen_by(note_id, nick) values(?, ?)`, noteId, nick.Nick)
		if err != nil {
			return err
		}
	}
	return nil
}
