package main

import (
	"goirc/bot"
	"goirc/handlers"
	"goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
	"strings"
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

	b.Handle(`^!help`, func(params bot.HandlerParams) (string, error) {
		lines := []string{}
		for _, h := range b.Handlers {
			lines = append(lines, h.String())
		}
		return strings.Join(lines, "\n"), nil
	})
	b.Handle(`^!catchup`, handlers.Catchup)
	b.Handle(`^,(.+)$`, handlers.CreateNote)
	b.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	b.Handle(`^!feedme`, handlers.FeedMe)
	b.Handle(`^!pipehealth`, handlers.PipeHealth)
	b.Handle(`(https?://\S+)`, handlers.Link)
	b.Handle(`^!day`, handlers.NationalDay)
	b.Handle(`\b69[^0-9]*\b`, handlers.Nice)
	b.Handle(`^!odds`, handlers.MLBOddsHandler)
	b.Handle(`^!pom`, handlers.POM)
	b.Handle(`^("[^"]+)$`, handlers.Quote)
	b.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	b.Handle(`^\?(\S+)`, handlers.Seen)
	b.Handle(`world.?cup`, handlers.Worldcup)
	b.Handle(`^!left`, handlers.TimeLeft)
	b.Handle(`^!epi`, handlers.EpigramHandler)
	b.Handle(`^!weather (.*)$`, handlers.WeatherHandler)
	b.Handle(`^!weather$`, handlers.WeatherHandler)
	b.Handle(`^!w (.*)$`, handlers.WeatherHandler)
	b.Handle(`^!w$`, handlers.WeatherHandler)
}
