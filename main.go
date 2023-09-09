package main

import (
	"goirc/bot"
	"goirc/handlers"
	"goirc/handlers/mlb"
	"goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
	"time"
)

func main() {
	go web.Serve(model.DB)

	var idleParam = bot.IdleParam{
		Duration: 24 * time.Hour,
		Handler:  handlers.FeedMe,
	}

	var repeatParam = bot.RepeatParam{
		Duration: 10 * time.Second,
		Handler:  handlers.DoRemind,
	}

	b, err := bot.Connect(
		util.Getenv("IRC_NICK"),
		util.Getenv("IRC_CHANNEL"),
		util.Getenv("IRC_SERVER"),
		idleParam,
		repeatParam)

	if err != nil {
		log.Fatal(err)
	}

	b.Handle(`^!help`, func(params bot.HandlerParams) error {
		for _, h := range b.Handlers {
			params.Privmsgf(params.Target, "%s", h.String())
		}
		return nil
	})
	b.Handle(`^!catchup`, handlers.Catchup)
	b.Handle(`^,(.+)$`, handlers.CreateNote)
	b.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	b.Handle(`^!feedme`, handlers.FeedMe)
	b.Handle(`(https?://\S+)`, handlers.Link)
	b.Handle(`^!day`, handlers.NationalDay)
	b.Handle(`\b69\b`, handlers.Nice)
	b.Handle(`^mlb odds (...)?$`, mlb.SingleTeamOdds)
	b.Handle(`^mlb teams$`, mlb.Teams)
	b.Handle(`^!pom`, handlers.POM)
	b.Handle(`^("[^"]+)$`, handlers.Quote)
	b.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	b.Handle(`^\?(\S+)`, handlers.Seen)
	b.Handle(`world.?cup`, handlers.Worldcup)

	b.Loop()
}
