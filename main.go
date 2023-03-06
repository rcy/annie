package main

import (
	"goirc/bot"
	"goirc/handlers"
	"goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
)

func main() {
	var functions = []bot.HandlerFunction{
		handlers.Catchup,
		handlers.CreateNote,
		handlers.DeferredDelivery,
		handlers.FeedMe,
		handlers.Link,
		handlers.Nice,
		handlers.OpeningDay,
		handlers.Quote,
		handlers.Report,
		handlers.Ticker,
		handlers.Trade,
		handlers.Worldcup,
	}

	conn, err := bot.Connect(util.Getenv("IRC_NICK"), util.Getenv("IRC_CHANNEL"), util.Getenv("IRC_SERVER"), functions)
	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(model.DB)

	conn.Loop()
}
