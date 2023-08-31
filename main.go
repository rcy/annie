package main

import (
	"goirc/bot"
	"goirc/commit"
	"goirc/handlers"
	"goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
	"time"
)

func main() {
	log.Printf("VERSION %s", commit.URL())

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
	b.Handle(`\b69\b`, handlers.Nice)
	b.Handle(`^!pom`, handlers.POM)
	b.Handle(`^("[^"]+)$`, handlers.Quote)
	b.Handle(`^((report).*)$`, handlers.Report) // broken
	b.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	b.Handle(`^\?(\S+)`, handlers.Seen)
	b.Handle(`^[$]([A-Za-z-]+)`, handlers.Ticker) // broken
	b.Handle(`^((buy|sell).*)$`, handlers.Trade)  // broken
	b.Handle(`world.?cup`, handlers.Worldcup)

	b.Loop()
}
