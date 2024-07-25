package main

import (
	"goirc/bot"
	db "goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
)

func main() {
	b, err := bot.Connect(
		util.Getenv("IRC_NICK"),
		util.Getenv("IRC_CHANNEL"),
		util.Getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(db.DB, b)

	addHandlers(b)

	b.Loop()
}
