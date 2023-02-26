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
		handlers.CreateNote,
		handlers.DeferredDelivery,
		handlers.Link,
		handlers.FeedMe,
		handlers.Catchup,
		handlers.Worldcup,
		handlers.Ticker,
		handlers.Quote,
		handlers.Trade,
		handlers.Report,
	}

	conn, err := bot.Connect(model.DB, util.Getenv("IRC_NICK"), util.Getenv("IRC_CHANNEL"), util.Getenv("IRC_SERVER"), functions)
	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(model.DB)

	conn.Loop()
}
