package main

import (
	"goirc/bot"
	db "goirc/model"
	"goirc/util"
	"goirc/web"
	"log"
)

//go:generate go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate --file db/sqlc.yaml

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
