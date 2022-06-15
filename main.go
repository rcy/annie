package main

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/thoj/go-ircevent"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"regexp"
)

func getenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s not set!", key)
	} else {
		log.Printf("%s=%s\n", key, val)
	}

	return val
}

func main() {
	dbfile := getenv("SQLITE_DB")

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

	conn, err := ircmain(db, getenv("IRC_NICK"), getenv("IRC_CHANNEL"), getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	go conn.Loop()

	// webserver
	log.Printf("starting webserver on %s", os.Getenv("PORT"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var (
			created_at string
			nick       string
			text       string
		)

		rows, err := db.Query(`select * from notes order by created_at desc`)
		if err != nil {
			fmt.Fprintf(w, "there was an error")
		}
		for rows.Next() {
			err := rows.Scan(&created_at, &nick, &text)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n\n", text)
		}
	})
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
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
