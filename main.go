package main

import (
	"goirc/db"
	"goirc/handlers"
	"goirc/irc"
	"goirc/util"
	"goirc/web"
	"log"
)

func main() {
	db := db.Open(util.Getenv("SQLITE_DB"))
	defer db.Close()

	var functions = []irc.HandlerFunction{
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

	conn, err := irc.Connect(db, util.Getenv("IRC_NICK"), util.Getenv("IRC_CHANNEL"), util.Getenv("IRC_SERVER"), functions)
	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(db)

	conn.Loop()
}
