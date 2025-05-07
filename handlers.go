package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/events"
	"goirc/handlers"
	"goirc/handlers/annie"
	"goirc/handlers/bedtime"
	"goirc/handlers/bible"
	"goirc/handlers/day"
	"goirc/handlers/election"
	"goirc/handlers/epigram"
	"goirc/handlers/fed25"
	"goirc/handlers/gold"
	"goirc/handlers/hn"
	"goirc/handlers/kinfonet"
	"goirc/handlers/linkpool"
	"goirc/handlers/mlb"
	"goirc/handlers/tip"
	"goirc/handlers/weather"
	"goirc/internal/ai"
	db "goirc/model"
	"goirc/web"
	"regexp"
	"time"

	"github.com/robfig/cron"
	"github.com/sashabaranov/go-openai"
)

func addHandlers(b *bot.Bot) {
	nick := regexp.QuoteMeta(b.Conn.GetNick())

	b.Handle(`^!catchup`, handlers.Catchup)
	b.Handle(`^,(.+)$`, handlers.CreateNote)
	b.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	b.Handle(`^!feedme`, handlers.AnonLink)
	b.Handle(`^!pipehealth\b`, handlers.AnonStatus)
	b.Handle(`(https?://\S+)`, handlers.Link)
	b.Handle(`^!day\b`, day.NationalDay)
	b.Handle(`^!dayi\b`, day.Dayi)
	b.Handle(`^!week\b`, day.NationalWeek)
	b.Handle(`^!weeki\b`, day.Weeki)
	b.Handle(`^!month\b`, day.NationalMonth)
	b.Handle(`^!monthi\b`, day.Monthi)
	b.Handle(`^!refs`, day.NationalRefs)
	b.Handle(`^!img (.+)$`, day.Image)
	b.Handle(`\b69[^0-9]*\b`, handlers.Nice)
	b.Handle(`^!odds`, mlb.PlayoffOdds)
	b.Handle(`^!godds`, mlb.GameOdds)
	b.Handle(`^!pom`, handlers.POM)
	b.Handle(`^("[^"]+)$`, handlers.Quote)
	b.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	b.Handle(`^\?(\S+)`, handlers.Seen)
	b.Handle(`world.?cup`, handlers.Worldcup)
	b.Handle(`^!left`, handlers.TimeLeft)
	b.Handle(`^!epi`, epigram.Handle)
	b.Handle(`^!weather (.*)$`, weather.Handle)
	b.Handle(`^!weather$`, weather.Handle)
	b.Handle(`^!w (.*)$`, weather.Handle)
	b.Handle(`^!w$`, weather.Handle)
	b.Handle(`^!f (.*)$`, weather.HandleForecast)
	b.Handle(`^!f$`, weather.HandleForecast)
	b.Handle(`^!wf (.*)$`, weather.HandleWeatherForecast)
	b.Handle(`^!wf$`, weather.HandleWeatherForecast)
	b.Handle(`^!xweather (.+)$`, weather.XHandle)
	b.Handle(`^!k`, kinfonet.TodaysQuoteHandler)
	b.Handle(`^!gold`, gold.Handle)
	b.Handle(`^!hn`, hn.Handle)
	b.Handle(`^!auth$`, web.HandleAuth)
	b.Handle(`^!deauth$`, web.HandleDeauth)
	b.Handle(`night`, bedtime.Handle)
	b.Handle(`^!election`, election.Handle)
	b.Handle(fmt.Sprintf(`^%s:?(.+)$`, nick), annie.Handle)
	b.Handle(fmt.Sprintf(`^(.+),? %s.?$`, nick), annie.Handle)
	b.Handle(`^!bible (.+)$`, bible.Handle)
	b.Handle(`^tip$`, tip.Handle)

	// show summary
	b.Handle(`^v$`, fed25.Handler)
	// if it starts with a letter, look up ridings
	b.Handle(`^v ([a-zA-Z].*)$`, fed25.FindRidingsByNameHandler)
	// if it starts with a number, pull results from riding id
	b.Handle(`^v ([0-9]+)$`, fed25.RidingHandler)
	b.Handle(`^vc ([a-zA-Z]+)$`, fed25.FindCandidatesHandler)
	b.Repeat(time.Minute, fed25.LeaderboardHandler)

	b.Repeat(10*time.Second, handlers.DoRemind)
	b.IdleRepeatAfterReset(8*time.Hour, handlers.POM)

	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	q := model.New(db.DB.DB)

	c := cron.NewWithLocation(location)
	err = c.AddFunc("16 14 15 * * 1,2,3,4,5,6", func() {
		note, err := q.RandomHistoricalTodayNote(context.TODO())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return
			}
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}

		b.Conn.Privmsgf(b.Channel, "on this day in %d, %s posted: %s", note.CreatedAt.Year(), note.Nick.String, note.Text.String)
	})
	if err != nil {
		panic(err)
	}
	err = c.AddFunc("0 0 0 * * 0", func() {
		signoff, err := ai.Complete(context.TODO(), openai.GPT4o, "you are annie, and have been hanging out with friends in irc all week. now its time for you to rest for a day. respond with a short goodbye. use all lowercase, minimal punctuation, no emojis", "see you later, annie")
		if err != nil {
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}
		b.Conn.Privmsg(b.Channel, signoff)
		b.Conn.Part(b.Channel)
	})
	if err != nil {
		panic(err)
	}

	err = c.AddFunc("0 0 0 * * 1", func() {
		b.Conn.Join(b.Channel)
	})
	if err != nil {
		panic(err)
	}

	err = c.AddFunc("57 * * * * *", func() {
		ctx := context.TODO()
		msg, err := q.ReadyFutureMessage(ctx, handlers.FutureMessageInterval)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return
			}
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}
		err = q.DeleteFutureMessage(ctx, msg.ID)
		if err != nil {
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}

		// send anonymous note
		switch msg.Kind {
		case "link":
			err = handlers.AnonLink(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
		case "quote":
			err = handlers.AnonQuote(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
		default:
			b.Conn.Privmsgf(b.Channel, "unhandled msg.Kind: %s", msg.Kind)
		}
		if err != nil {
			if errors.Is(err, ai.ErrBilling) {
				// the quote was sent, but no generated image, this is fine
				return
			}
			if errors.Is(err, linkpool.NoNoteFoundError) {
				// didn't find a note, reschedule
				_, scheduleErr := q.ScheduleFutureMessage(ctx, msg.Kind)
				if scheduleErr != nil {
					b.Conn.Privmsg(b.Channel, "error rescheduling: "+scheduleErr.Error())
				}
				return
			}
			// something else happened, spam the channel
			b.Conn.Privmsg(b.Channel, "error: "+err.Error())
		}
	})
	if err != nil {
		panic(err)
	}

	c.Start()

	events.Subscribe("anonnoteposted", func(note any) {
		go func() {
			err := handlers.AnonLink(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
			if err != nil {
				if errors.Is(err, ai.ErrBilling) {
					return
				}
				if errors.Is(err, linkpool.NoNoteFoundError) {
					return
				}
				b.Conn.Privmsg(b.Channel, "error: "+err.Error())
			}
		}()
	})

	events.Subscribe("anonquoteposted", func(note any) {
		go func() {
			err := handlers.AnonQuote(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
			if err != nil {
				if errors.Is(err, ai.ErrBilling) {
					return
				}
				if errors.Is(err, linkpool.NoNoteFoundError) {
					return
				}
				b.Conn.Privmsg(b.Channel, "error: "+err.Error())
			}
		}()
	})

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, "https://github.com/rcy/annie/blob/main/handlers.go")
		return nil
	})
}
