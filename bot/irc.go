package bot

import (
	"crypto/tls"
	"fmt"
	"goirc/bot/idle"
	"goirc/commit"
	"goirc/model"
	"goirc/model/laters"
	"goirc/model/notes"
	"goirc/util"
	"log"
	"log/slog"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	irc "github.com/thoj/go-ircevent"
)

type Handler struct {
	pattern string
	regexp  regexp.Regexp
	action  HandlerFunction
}

func (h Handler) String() string {
	strs := strings.Split(runtime.FuncForPC(reflect.ValueOf(h.action).Pointer()).Name(), ".")
	return fmt.Sprintf("%-32s %s", h.pattern, strs[len(strs)-1])
}

type Bot struct {
	Conn               *irc.Connection
	Channel            string
	Handlers           []Handler
	idleResetFunctions []func()
}

func (b *Bot) Handle(pattern string, action HandlerFunction) {
	h := Handler{
		pattern,
		*regexp.MustCompile(pattern),
		action,
	}

	b.Handlers = append(b.Handlers, h)
}

func (b *Bot) Repeat(timeout time.Duration, action HandlerFunction) {
	go func() {
		for {
			time.Sleep(timeout)
			err := action(HandlerParams{
				Privmsgf: b.MakePrivmsgf(),
				Target:   b.Channel,
			})
			if err != nil {
				slog.Warn("Repeat", "err", err)
			}
		}
	}()
}

func (b *Bot) IdleRepeat(timeout time.Duration, action HandlerFunction) {
	reset := idle.Repeat(timeout, func() {
		err := action(HandlerParams{
			Privmsgf: b.MakePrivmsgf(),
			Target:   b.Channel,
		})
		if err != nil {
			slog.Warn("IdleRepeat", "err", err)
		}
	})

	b.idleResetFunctions = append(b.idleResetFunctions, reset)
}

func (b *Bot) IdleRepeatAfterReset(timeout time.Duration, action HandlerFunction) {
	reset := idle.RepeatAfterReset(timeout, func() {
		err := action(HandlerParams{
			Privmsgf: b.MakePrivmsgf(),
			Target:   b.Channel,
		})
		if err != nil {
			slog.Warn("IdleRepeatAfterReset", "err", err)
		}
	})

	b.idleResetFunctions = append(b.idleResetFunctions, reset)
}

func (b *Bot) resetIdle() {
	for _, fn := range b.idleResetFunctions {
		fn()
	}
}

func (b *Bot) Loop() {
	b.Conn.Loop()
}

func Connect(nick string, channel string, server string) (*Bot, error) {
	initialized := make(chan bool)
	var bot Bot
	bot.Channel = channel
	bot.Conn = irc.IRC(nick, "github.com/rcy/annie")
	bot.Conn.VerboseCallbackHandler = false
	bot.Conn.Debug = false
	bot.Conn.UseTLS = true
	bot.Conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	bot.Conn.AddCallback("001", func(e *irc.Event) { bot.Conn.Join(channel) })
	bot.Conn.AddCallback("353", func(e *irc.Event) {
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
	bot.Conn.AddCallback("366", func(e *irc.Event) {})
	bot.Conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		bot.resetIdle()
		go bot.RunHandlers(e)
	})
	bot.Conn.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick != nick {
			go func() {
				time.Sleep(10 * time.Second)
				bot.SendLaters(channel, e.Nick)
			}()

			go bot.SendMissed(channel, e.Nick)
		} else {
			go func() {
				time.Sleep(1 * time.Second)
				url, err := commit.URL()
				if err != nil {
					bot.Conn.Privmsgf(channel, "error: %s", err)
					return
				}
				if url != "" {
					bot.Conn.Privmsgf(channel, url)
				}
				initialized <- true
			}()
		}

		// trigger NAMES to update the list of joined nicks
		bot.Conn.SendRawf("NAMES %s", channel)
	})
	bot.Conn.AddCallback("PART", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			bot.Conn.SendRawf("NAMES %s", channel)
		}
	})
	bot.Conn.AddCallback("QUIT", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			bot.Conn.SendRawf("NAMES %s", channel)
		}
	})
	bot.Conn.AddCallback("NICK", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			bot.Conn.SendRawf("NAMES %s", channel)
		}
	})
	err := bot.Conn.Connect(server)

	<-initialized

	return &bot, err
}

func (bot *Bot) SendLaters(channel string, nick string) {
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
			bot.Conn.Privmsgf(channel, "%s: %s (from %s %s ago)", nick, row.Message, row.Nick, util.Since(row.CreatedAt))
		}
	}
}

func (bot *Bot) MakePrivmsgf() func(string, string, ...interface{}) {
	return func(target, message string, a ...interface{}) {
		bot.Conn.Privmsgf(target, message, a...)
	}
}

func (bot *Bot) RunHandlers(e *irc.Event) {
	channel := e.Arguments[0]
	msg := e.Arguments[1]
	nick := e.Nick

	var target string
	if channel == bot.Conn.GetNick() {
		target = nick
	} else {
		target = channel
	}

	for _, handler := range bot.Handlers {
		matches := handler.regexp.FindStringSubmatch(msg)
		if len(matches) > 0 {
			err := handler.action(HandlerParams{
				Privmsgf: bot.MakePrivmsgf(),
				Msg:      msg,
				Nick:     nick,
				Target:   target,
				Matches:  matches,
			})

			if err != nil {
				bot.Conn.Privmsgf(target, "error: %s", err)
				return
			}
		}
	}
}

func isAltNick(nick string) bool {
	return strings.HasSuffix(nick, "`") || strings.HasSuffix(nick, "_")
}

func (bot *Bot) SendMissed(channel string, nick string) {
	if isAltNick(nick) {
		return
	}

	channelNick := model.ChannelNick{}
	err := model.DB.Get(&channelNick, `select * from channel_nicks where present = 0 and channel = ? and nick = ?`, channel, nick)
	if err != nil {
		return
	}

	notes := []notes.Note{}
	model.DB.Select(&notes, "select * from notes where created_at > ? order and nick <> target by created_at asc limit 69", channelNick.UpdatedAt)

	if len(notes) > 0 {
		bot.Conn.Privmsgf(nick, "Hi %s, you missed %d thing(s) in %s since %s:",
			nick, len(notes), channel, channelNick.UpdatedAt)

		for _, note := range notes {
			bot.Conn.Privmsgf(nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
}
