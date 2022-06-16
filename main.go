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

	_, err = db.Exec(`create table if not exists links(created_at text, nick text, text text);`)
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
			fmt.Fprintf(w, "%s\n", text)
		}
	})
	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		var (
			created_at string
			nick       string
			text       string
		)

		rows, err := db.Query(`select * from links order by created_at desc`)
		if err != nil {
			fmt.Fprintf(w, "there was an error")
		}
		for rows.Next() {
			err := rows.Scan(&created_at, &nick, &text)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", text)
		}
	})
	err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ircmain(db *sql.DB, nick, channel, server string) (*irc.Connection, error) {
	ircnick1 := nick
	irccon := irc.IRC(ircnick1, "github.com/rcy/annie")
	irccon.VerboseCallbackHandler = true
	irccon.Debug = true
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		channel := e.Arguments[0]
		msg := e.Arguments[1]
		nick := e.Nick

		matchNote(irccon, db, msg, nick, channel)
		matchLink(irccon, db, msg, nick, channel)
	})
	err := irccon.Connect(server)

	return irccon, err
}

func matchNote(irccon *irc.Connection, db *sql.DB, msg, nick, channel string) {
	re := regexp.MustCompile(`^,(.+)$`)
	matches := re.FindSubmatch([]byte(msg))

	if len(matches) > 0 {
		note := string(matches[1])
		_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?)`, nick, note)
		if err != nil {
			log.Print(err)
			irccon.Privmsg(channel, err.Error())
		} else {
			irccon.Privmsg(channel, "recorded note")
		}
	}
}

func matchLink(irccon *irc.Connection, db *sql.DB, msg, nick, channel string) {
	re := regexp.MustCompile(`^.*(http://\S+)$`)
	matches := re.FindSubmatch([]byte(msg))

	if len(matches) > 0 {
		url := string(matches[1])
		_, err := db.Exec(`insert into links values(datetime('now'), ?, ?)`, nick, url)
		if err != nil {
			log.Print(err)
			irccon.Privmsg(channel, err.Error())
		} else {
			log.Printf("recorded url %s", url)
		}
	}
}
