package handlers

import (
	"goirc/model"

	"github.com/jmoiron/sqlx"
	irc "github.com/thoj/go-ircevent"
)

type HandlerParams struct {
	Irccon *irc.Connection
	Db     *sqlx.DB
	Msg    string
	Nick   string
	Target string
}

type Handler struct {
	Name     string
	Function func(HandlerParams) bool
}

var Handlers = []Handler{
	{
		Name:     "Match Create Note",
		Function: createNote,
	},
	{
		Name:     "Match Deferred Delivery",
		Function: deferredDelivery,
	},
	{
		Name:     "Match Link",
		Function: link,
	},
	{
		Name:     "Match Feedme Command",
		Function: feedMe,
	},
	{
		Name:     "Match Catchup Command",
		Function: catchup,
	},
	{
		Name:     "Match Until Command",
		Function: worldcup,
	},
	{
		Name:     "Match stock price",
		Function: ticker,
	},
	{
		Name:     "Match quote",
		Function: quote,
	},
	{
		Name:     "Trade stock",
		Function: trade,
	},
	{
		Name:     "Trade: show holdings report",
		Function: report,
	},
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
