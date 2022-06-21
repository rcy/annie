package main

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"text/template"
	"time"
	//"fmt"
	//	"database/sql"
	"github.com/BurntSushi/migration"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/thoj/go-ircevent"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"regexp"
)

type Note struct {
	CreatedAt string `db:"created_at"`
	Text      string
	Nick      string
	Kind      string
}

type NickWithNoteCount struct {
	Nick  string
	Count int
}

func getenv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s not set!", key)
	} else {
		log.Printf("%s=%s\n", key, val)
	}

	return val
}

func openDb(dbfile string) *sqlx.DB {
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

	db, err := migration.Open("sqlite", dbfile, migrations)
	if err != nil {
		log.Fatal("Failed to open with migration")
	}
	return sqlx.NewDb(db, "sqlite")
}

func main() {
	db := openDb(getenv("SQLITE_DB"))
	defer db.Close()

	conn, err := ircmain(db, getenv("IRC_NICK"), getenv("IRC_CHANNEL"), getenv("IRC_SERVER"))
	if err != nil {
		log.Fatal(err)
	}
	go conn.Loop()

	// // webserver
	// log.Printf("starting webserver on %s", os.Getenv("PORT"))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	var (
	// 		created_at string
	// 		nick       string
	// 		text       string
	// 		kind       string
	// 	)

	// 	rows, err := db.Query(`select created_at, nick, text, kind from notes order by created_at desc`)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "there was an error")
	// 	}
	// 	fmt.Fprintf(w, "<html><ul>")
	// 	for rows.Next() {
	// 		err := rows.Scan(&created_at, &nick, &text, &kind)

	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		if kind == "note" {
	// 			fmt.Fprintf(w, "<li>%s</li>", text)
	// 		} else if kind == "link" {
	// 			fmt.Fprintf(w, `<li><a href="%s">%s</a></li>`, text, text)
	// 		}
	// 	}
	// 	fmt.Fprintf(w, "</ul></html>")
	// })
	// http.HandleFunc("/links", func(w http.ResponseWriter, r *http.Request) {
	// 	var (
	// 		created_at string
	// 		nick       string
	// 		text       string
	// 		kind       string
	// 	)

	// 	rows, err := db.Query(`select * from notes where kind = 'link' order by created_at desc`)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "there was an error")
	// 	}
	// 	for rows.Next() {
	// 		err := rows.Scan(&created_at, &nick, &text, &kind)

	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Fprintf(w, "%s\n", text)
	// 	}
	// })
	// http.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
	// 	var (
	// 		created_at string
	// 		nick       string
	// 		text       string
	// 		kind       string
	// 	)

	// 	rows, err := db.Query(`select * from notes where kind = 'note' order by created_at desc`)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "there was an error")
	// 	}
	// 	for rows.Next() {
	// 		err := rows.Scan(&created_at, &nick, &text, &kind)

	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		fmt.Fprintf(w, "%s\n", text)
	// 	}
	// })
	// err = http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// new gin webserver
	webserver(db)
}

//go:embed "templates/index.gohtml"
var indexTemplate string

func webserver(db *sqlx.DB) {
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/test/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
	r.GET("/", func(c *gin.Context) {
		nick := c.Query("nick")

		notes, err := getNotes(db, nick)
		if err != nil {
			log.Fatal(err)
		}

		nicks, err := getNicks(db)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("name").Parse(indexTemplate)
		if err != nil {
			log.Fatal("error parsing template")
		}

		out := new(bytes.Buffer)
		err = tmpl.Execute(out, gin.H{
			"nicks": nicks,
			"notes": notes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", out.Bytes())
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getNicks(db *sqlx.DB) ([]NickWithNoteCount, error) {
	nicks := []NickWithNoteCount{}
	err := db.Select(&nicks, `select nick, count(nick) as count from notes group by nick`)
	return nicks, err
}

func getNotes(db *sqlx.DB, nick string) ([]Note, error) {
	notes := []Note{}
	var err error
	if nick == "" {
		err = db.Select(&notes, `select created_at, text, nick, kind from notes order by created_at desc limit 1000`)
	} else {
		err = db.Select(&notes, `select created_at, text, nick, kind from notes where nick = ? order by created_at desc limit 1000`, nick)
	}
	return notes, err
}

func ircmain(db *sqlx.DB, nick, channel, server string) (*irc.Connection, error) {
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
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick == nick {
			time.Sleep(10 * time.Second)
			irccon.Privmsg(channel, "don't worry devlan, I'm ok")
		}
	})
	err := irccon.Connect(server)

	return irccon, err
}

func matchNote(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) {
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

func matchLink(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) {
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
