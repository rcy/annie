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

	var privmsgHandlers = []bot.HandlerFunction{
		handlers.Catchup,
		handlers.CreateNote,
		handlers.DeferredDelivery,
		handlers.MatchFeedMe,
		handlers.Link,
		handlers.Nice,
		handlers.MatchPOM,
		handlers.Quote,
		handlers.Report,
		handlers.RemindMe,
		handlers.Seen,
		handlers.Ticker,
		handlers.Trade,
		handlers.Worldcup,
	}

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
		privmsgHandlers,
		idleParam,
		repeatParam)

	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(model.DB)

	bot.Loop()
}
