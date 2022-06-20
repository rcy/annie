package main

import (
	"github.com/BurntSushi/migration"
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

func openDb(dbfile string) (*sql.DB, error) {
	log.Printf("Opening db: %s", dbfile)

	migrations := []migration.Migrator{
		func(tx migration.LimitedTx) error {
			_, err := tx.Exec(`create table if not exists notes(created_at text, nick text, text text)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			_, err := tx.Exec(`create table if not exists links(created_at text, nick text, text text)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding kind column to notes")
			_, err := tx.Exec(`alter table notes add column kind string not null default "note"`)
			return err
		},
	}

	return migration.Open("sqlite", dbfile, migrations)
}

func main() {
	db, err := openDb(getenv("SQLITE_DB"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
			kind       string
		)

		rows, err := db.Query(`select created_at, nick, text, kind from notes order by created_at desc`)
		if err != nil {
			fmt.Fprintf(w, "there was an error")
		}
		fmt.Fprintf(w, "<html><ul>")
		for rows.Next() {
			err := rows.Scan(&created_at, &nick, &text, &kind)

			if err != nil {
				log.Fatal(err)
			}
			if kind == "note" {
				fmt.Fprintf(w, "<li>%s</li>", text)
			} else if kind == "link" {
				fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, text, text)
			}
		}
		fmt.Fprintf(w, "</ul></html>")
	})
	http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
		var (
			created_at string
			nick       string
			text       string
			kind       string
		)

		rows, err := db.Query(`select * from notes where kind = 'link' order by created_at desc`)
		if err != nil {
			fmt.Fprintf(w, "there was an error")
		}
		for rows.Next() {
			err := rows.Scan(&created_at, &nick, &text, &kind)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "%s\n", text)
		}
	})
	http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		var (
			created_at string
			nick       string
			text       string
			kind       string
		)

		rows, err := db.Query(`select * from notes where kind = 'note' order by created_at desc`)
		if err != nil {
			fmt.Fprintf(w, "there was an error")
		}
		for rows.Next() {
			err := rows.Scan(&created_at, &nick, &text, &kind)

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
		_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?, 'note')`, nick, note)
		if err != nil {
			log.Print(err)
			irccon.Privmsg(channel, err.Error())
		} else {
			irccon.Privmsg(channel, "recorded note")
		}
	}
}

func matchLink(irccon *irc.Connection, db *sql.DB, msg, nick, channel string) {
	re := regexp.MustCompile(`^.*(https?://\S+)$`)
	matches := re.FindSubmatch([]byte(msg))

	if len(matches) > 0 {
		url := string(matches[1])
		_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?, 'link')`, nick, url)
		if err != nil {
			log.Print(err)
			irccon.Privmsg(channel, err.Error())
		} else {
			log.Printf("recorded url %s", url)
		}
	}
}
