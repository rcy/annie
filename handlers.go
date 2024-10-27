package main

import (
	"context"
	"errors"
	"goirc/bot"
	"goirc/db/model"
	"goirc/events"
	"goirc/handlers"
	"goirc/handlers/annie"
	"goirc/handlers/bedtime"
	"goirc/handlers/day"
	"goirc/handlers/election"
	"goirc/handlers/epigram"
	"goirc/handlers/gold"
	"goirc/handlers/hn"
	"goirc/handlers/kinfonet"
	"goirc/handlers/linkpool"
	"goirc/handlers/mlb"
	"goirc/handlers/weather"
	db "goirc/model"
	"goirc/web"
	"time"

	"github.com/robfig/cron"
)

func addHandlers(b *bot.Bot) {
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
	b.Handle(`(^annie\b)|(\bannie.?$)`, annie.Handle)

	b.Repeat(10*time.Second, handlers.DoRemind)
	b.IdleRepeatAfterReset(8*time.Hour, handlers.POM)

	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	c := cron.NewWithLocation(location)
	err = c.AddFunc("16 14 15 * * *", func() {
		q := model.New(db.DB.DB)
		note, err := q.RandomHistoricalTodayNote(context.TODO())
		if err != nil {
			// TODO test no rows
			b.Conn.Privmsg(b.Channel, err.Error())
			return
		}

		b.Conn.Privmsgf(b.Channel, "on this day in %d, %s posted: %s", note.CreatedAt.Year(), note.Nick.String, note.Text.String)
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
				if !errors.Is(err, linkpool.NoNoteFoundError) {
					b.Conn.Privmsg(b.Channel, "error: "+err.Error())
				}
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
				if !errors.Is(err, linkpool.NoNoteFoundError) {
					b.Conn.Privmsg(b.Channel, "error: "+err.Error())
				}
			}
		}()
	})

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, "https://github.com/rcy/annie/blob/main/handlers.go")
		return nil
	})
}
