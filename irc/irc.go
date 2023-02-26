package irc

import (
	"crypto/tls"
	handlers "goirc/irc/handlers"
	"goirc/model"
	"goirc/util"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	irc "github.com/thoj/go-ircevent"
)

func Connect(db *sqlx.DB, nick, channel, server string) (*irc.Connection, error) {
	ircnick1 := nick
	irccon := irc.IRC(ircnick1, "github.com/rcy/annie")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("353", func(e *irc.Event) {
		// clear the presence of all channel nicks
		_, err := db.Exec(`update channel_nicks set updated_at = current_timestamp, present = false`)
		if err != nil {
			log.Fatal(err)
		}

		// remove @ op markers from nick argument
		nickStr := strings.ReplaceAll(e.Arguments[len(e.Arguments)-1], "@", "")

		// mark nicks as present and record timestamp which can be intepreted as 'last seen', or 'online since'
		for _, nick := range strings.Split(nickStr, " ") {
			_, err = db.Exec(`
insert into channel_nicks(updated_at, channel, nick, present) values(current_timestamp, ?, ?, ?)
on conflict(channel, nick) do update set updated_at = current_timestamp, present=excluded.present`,
				channel, nick, true)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		go handlePrivmsg(irccon, db, e)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick != nick {
			sendLaters(irccon, db, channel, e.Nick)

			go sendMissed(irccon, db, channel, e.Nick)
		}

		// trigger NAMES to update the list of joined nicks
		irccon.SendRawf("NAMES %s", channel)
	})
	irccon.AddCallback("PART", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	irccon.AddCallback("QUIT", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	irccon.AddCallback("NICK", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	err := irccon.Connect(server)

	return irccon, err
}

func sendLaters(irccon *irc.Connection, db *sqlx.DB, channel string, nick string) {
	// loop through each later message and see if the prefix matches this nick
	laters, err := getLaters(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, later := range laters {
		if strings.Contains(nick, later.Target) {
			_, err := db.Exec(`delete from laters where rowid = ?`, later.RowId)
			if err != nil {
				log.Fatal(err)
			}
			irccon.Privmsgf(channel, "%s: %s (from %s %s ago)", nick, later.Message, later.Nick, util.Since(later.CreatedAt))
		}
	}
}

func handlePrivmsg(irccon *irc.Connection, db *sqlx.DB, e *irc.Event) {
	channel := e.Arguments[0]
	msg := e.Arguments[1]
	nick := e.Nick

	for _, f := range handlers.Handlers {
		var target string

		if channel == irccon.GetNick() {
			target = nick
		} else {
			target = channel
		}

		if f.Function(handlers.HandlerParams{Irccon: irccon, Db: db, Msg: msg, Nick: nick, Target: target}) {
			break
		}
	}
}

func isAltNick(nick string) bool {
	return strings.HasSuffix(nick, "`") || strings.HasSuffix(nick, "_")
}

func getLaters(db *sqlx.DB) ([]model.Later, error) {
	laters := []model.Later{}
	err := db.Select(&laters, `select rowid, created_at, nick, target, message, sent from laters limit 100`)
	return laters, err
}

func sendMissed(irccon *irc.Connection, db *sqlx.DB, channel string, nick string) {
	if isAltNick(nick) {
		return
	}

	channelNick := model.ChannelNick{}
	err := db.Get(&channelNick, `select * from channel_nicks where present = 0 and channel = ? and nick = ?`, channel, nick)
	if err != nil {
		return
	}

	notes := []model.Note{}
	db.Select(&notes, "select * from notes where created_at > ? order by created_at asc limit 69", channelNick.UpdatedAt)

	if len(notes) > 0 {
		irccon.Privmsgf(nick, "Hi %s, you missed %d thing(s) in %s since %s:",
			nick, len(notes), channel, channelNick.UpdatedAt)

		for _, note := range notes {
			irccon.Privmsgf(nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
}
