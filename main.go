package main

import (
	"goirc/bot"
	"goirc/handlers"
	"goirc/handlers/epigram"
	"goirc/handlers/mlb"
	"goirc/handlers/weather"
	"goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
	"time"
)

func main() {
	go web.Serve(model.DB)

	b, err := bot.Connect(
		util.Getenv("IRC_NICK"),
		util.Getenv("IRC_CHANNEL"),
		util.Getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}

	addHandlers(b)

	b.Loop()
}

func addHandlers(b *bot.Bot) {
	b.Repeat(10*time.Second, handlers.DoRemind)
	b.IdleRepeatAfterReset(8*time.Hour, handlers.POM)
	b.IdleRepeatAfterReset(1*time.Minute, handlers.FeedMe)

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		for _, h := range b.Handlers {
			params.Privmsgf(params.Target, "%s", h.String())
			time.Sleep(200 * time.Millisecond) // prevent flooding
		}
		return nil
	})
	b.Handle(`^!catchup`, handlers.Catchup)
	b.Handle(`^,(.+)$`, handlers.CreateNote)
	b.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	b.Handle(`^!feedme`, handlers.FeedMe)
	b.Handle(`^!pipehealth`, handlers.PipeHealth)
	b.Handle(`(https?://\S+)`, handlers.Link)
	b.Handle(`^!day`, handlers.NationalDay)
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
	b.Handle(`^!xweather (.+)$`, weather.XHandle)
}
