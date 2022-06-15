package main

import (
	"net/http"
	"fmt"
	"crypto/tls"
	"database/sql"
	"github.com/thoj/go-ircevent"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"regexp"
)

func main() {
	dbfile := os.Getenv("ANNIE_DB")

	if dbfile == "" {
		log.Fatal("ANNIE_DB not set")
	}

	log.Printf("Opening db: %s", dbfile)

	db, err := sql.Open("sqlite", dbfile)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`create table if not exists notes(created_at text, nick text, text text);`)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ircmain(db, "annie", "#embx", "irc.libera.chat:6697")
	if err != nil {
		log.Fatal(err)
	}
	go conn.Loop()

	// webserver
	log.Printf("starting webserver on %s", os.Getenv("PORT"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to new server!")
	})
	err = http.ListenAndServe(":" + os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ircmain(db *sql.DB, nick, channel, server string) (*irc.Connection, error) {
	ircnick1 := nick
	irccon := irc.IRC(ircnick1, "IRCTestSSL")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := e.Arguments[1]
		re := regexp.MustCompile(`^,(.+)$`)
		matches := re.FindSubmatch([]byte(msg))

		if len(matches) > 0 {
			note := string(matches[1])
			_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?)`, e.Nick, note)
			if err != nil {
				log.Print(err)
				irccon.Privmsg(e.Arguments[0], err.Error())
			} else {
				irccon.Privmsg(e.Arguments[0], "recorded note")
			}
		}
	})
	err := irccon.Connect(server)

	return irccon, err
}
