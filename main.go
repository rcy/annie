package main

import (
	"goirc/bot"
	db "goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
)

func main() {
	go web.Serve(db.DB)

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
