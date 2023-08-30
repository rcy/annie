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

	bot, err := bot.Connect(
		util.Getenv("IRC_NICK"),
		util.Getenv("IRC_CHANNEL"),
		util.Getenv("IRC_SERVER"),
		idleParam,
		repeatParam)

	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(`^!catchup`, handlers.Catchup)
	bot.Handle(`^,(.+)$`, handlers.CreateNote)
	bot.Handle(`^([^\s:]+): (.+)$`, handlers.DeferredDelivery)
	bot.Handle(`^!feedme`, handlers.FeedMe)
	bot.Handle(`(https?://\S+)`, handlers.Link)
	bot.Handle(`\b69\b`, handlers.Nice)
	bot.Handle(`^!pom`, handlers.POM)
	bot.Handle(`^("[^"]+)$`, handlers.Quote)
	bot.Handle(`^((report).*)$`, handlers.Report) // broken
	bot.Handle(`^!remindme ([^\s]+) (.+)$`, handlers.RemindMe)
	bot.Handle(`^\?(\S+)`, handlers.Seen)
	bot.Handle(`^[$]([A-Za-z-]+)`, handlers.Ticker) // broken
	bot.Handle(`^((buy|sell).*)$`, handlers.Trade)  // broken
	bot.Handle(`world.?cup`, handlers.Worldcup)

	bot.Loop()
}
