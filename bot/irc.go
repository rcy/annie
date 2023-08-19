package bot

import (
	"crypto/tls"
	"goirc/bot/idle"
	"goirc/bot/repeat"
	"goirc/commit"
	"goirc/model"
	"goirc/model/laters"
	"goirc/model/notes"
	"goirc/util"
	"log"
	"strings"
	"time"

	irc "github.com/thoj/go-ircevent"
)

type IdleParam struct {
	Duration time.Duration
	Handler  HandlerFunction
}

type RepeatParam struct {
	Duration time.Duration
	Handler  HandlerFunction
}

func Connect(nick string, channel string, server string, privmsgHandlers []HandlerFunction, idleParam IdleParam, repeatParam RepeatParam) (*irc.Connection, error) {
	ircnick1 := nick
	irccon := irc.IRC(ircnick1, "github.com/rcy/annie")
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("353", func(e *irc.Event) {
		// clear the presence of all channel nicks
		_, err := model.DB.Exec(`update channel_nicks set updated_at = current_timestamp, present = false where present = true`)
		if err != nil {
			log.Fatal(err)
		}

		// remove @ op markers from nick argument
		nickStr := strings.ReplaceAll(e.Arguments[len(e.Arguments)-1], "@", "")

		// mark nicks as present and record timestamp which can be intepreted as 'last seen', or 'online since'
		for _, nick := range strings.Split(nickStr, " ") {
			_, err = model.DB.Exec(`
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
		idle.Reset()
		go handlePrivmsg(irccon, e, privmsgHandlers)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick != nick {
			go func() {
				time.Sleep(10 * time.Second)
				sendLaters(irccon, channel, e.Nick)
			}()

			go sendMissed(irccon, channel, e.Nick)
		} else {
			go func() {
				time.Sleep(1 * time.Second)
				irccon.Privmsgf(channel, commit.URL())
			}()
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

	go idle.Every(idleParam.Duration, func() {
		idleParam.Handler(HandlerParams{
			Privmsgf: makePrivmsgf(irccon),
			Target:   channel,
		})
	})

	go repeat.Every(repeatParam.Duration, func() {
		repeatParam.Handler(HandlerParams{
			Privmsgf: makePrivmsgf(irccon),
			Target:   channel,
		})
	})

	return irccon, err
}

func sendLaters(irccon *irc.Connection, channel string, nick string) {
	// loop through each later message and see if the prefix matches this nick
	rows, err := laters.Get()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range rows {
		if strings.Contains(nick, row.Target) {
			_, err := model.DB.Exec(`delete from laters where rowid = ?`, row.RowId)
			if err != nil {
				log.Fatal(err)
			}
			irccon.Privmsgf(channel, "%s: %s (from %s %s ago)", nick, row.Message, row.Nick, util.Since(row.CreatedAt))
		}
	}
}

func makePrivmsgf(irccon *irc.Connection) func(string, string, ...interface{}) {
	return func(target, message string, a ...interface{}) {
		irccon.Privmsgf(target, message, a...)
	}
}

func handlePrivmsg(irccon *irc.Connection, e *irc.Event, handlers []HandlerFunction) {
	channel := e.Arguments[0]
	msg := e.Arguments[1]
	nick := e.Nick

	for _, f := range handlers {
		var target string

		if channel == irccon.GetNick() {
			target = nick
		} else {
			target = channel
		}

		f(HandlerParams{
			Privmsgf: makePrivmsgf(irccon),
			Msg:      msg,
			Nick:     nick,
			Target:   target,
		})
	}
}

func isAltNick(nick string) bool {
	return strings.HasSuffix(nick, "`") || strings.HasSuffix(nick, "_")
}

func sendMissed(irccon *irc.Connection, channel string, nick string) {
	if isAltNick(nick) {
		return
	}

	channelNick := model.ChannelNick{}
	err := model.DB.Get(&channelNick, `select * from channel_nicks where present = 0 and channel = ? and nick = ?`, channel, nick)
	if err != nil {
		return
	}

	notes := []notes.Note{}
	model.DB.Select(&notes, "select * from notes where created_at > ? order by created_at asc limit 69", channelNick.UpdatedAt)

	if len(notes) > 0 {
		irccon.Privmsgf(nick, "Hi %s, you missed %d thing(s) in %s since %s:",
			nick, len(notes), channel, channelNick.UpdatedAt)

		for _, note := range notes {
			irccon.Privmsgf(nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
}
