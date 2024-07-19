package main

import (
	"goirc/bot"
	"goirc/events"
	"goirc/handlers"
	"goirc/handlers/day"
	"goirc/handlers/epigram"
	"goirc/handlers/gold"
	"goirc/handlers/hn"
	"goirc/handlers/kinfonet"
	"goirc/handlers/mlb"
	"goirc/handlers/weather"
	"time"
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
	b.Handle(`^!week`, day.NationalWeek)
	b.Handle(`^!month`, day.NationalMonth)
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
	b.Handle(`^!f (.*)$`, weather.HandleForecast)
	b.Handle(`^!w$`, weather.Handle)
	b.Handle(`^!xweather (.+)$`, weather.XHandle)
	b.Handle(`^!k`, kinfonet.TodaysQuoteHandler)
	b.Handle(`^!gold`, gold.Handle)
	b.Handle(`^!hn`, hn.Handle)

	b.Repeat(10*time.Second, handlers.DoRemind)
	b.IdleRepeatAfterReset(8*time.Hour, handlers.POM)

	events.Subscribe("anonnoteposted", func(note any) {
		go func() {
			err := handlers.AnonLink(bot.HandlerParams{
				Target:   b.Channel,
				Privmsgf: b.MakePrivmsgf(),
			})
			if err != nil {
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
				b.Conn.Privmsg(b.Channel, "error: "+err.Error())
			}
		}()
	})

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, "https://github.com/rcy/annie/blob/main/handlers.go")
		return nil
	})
}