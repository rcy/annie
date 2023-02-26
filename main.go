package main

import (
	"goirc/db"
	"goirc/irc"
	"goirc/util"
	"goirc/web"
	"log"
)

func main() {
	db := db.Open(util.Getenv("SQLITE_DB"))
	defer db.Close()

	conn, err := irc.Connect(db, util.Getenv("IRC_NICK"), util.Getenv("IRC_CHANNEL"), util.Getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}

	go web.Serve(db)

	conn.Loop()
}
