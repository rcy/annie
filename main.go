package main

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
	"time"

	//	"database/sql"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"

	"github.com/BurntSushi/migration"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	irc "github.com/thoj/go-ircevent"
	_ "modernc.org/sqlite"
)

type ChannelNick struct {
	Channel string
	Nick    string
	Present string
}

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

type Later struct {
	RowId     int    `db:"rowid"`
	CreatedAt string `db:"created_at"`
	Nick      string
	Target    string
	Message   string
	Sent      bool
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
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding laters table")
			_, err := tx.Exec(`create table laters(created_at text, nick text, target text, message text, sent boolean default false)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding channel_nicks table")
			_, err := tx.Exec(`create table channel_nicks(channel text not null, nick text not null, present bool not null default false)`)
			return err
		},
	}

	db, err := migration.Open("sqlite", dbfile, migrations)
	if err != nil {
		log.Fatalf("MIGRATION: %v", err)
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

	webserver(db)
}

//go:embed "templates/index.gohtml"
var indexTemplate string

func webserver(db *sqlx.DB) {
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")

	r.GET("/snapshot.db", func(c *gin.Context) {
		os.Remove("/tmp/snapshot.db")
		if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
		c.File("/tmp/snapshot.db")
	})
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

func getLaters(db *sqlx.DB) ([]Later, error) {
	laters := []Later{}
	err := db.Select(&laters, `select rowid, created_at, nick, target, message, sent from laters limit 100`)
	return laters, err
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
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("353", func(e *irc.Event) {
		// clear the presence of all channel nicks
		_, err := db.Exec(`update channel_nicks set present = false`)
		if err != nil {
			log.Fatal(err)
		}

		// remove @ op markers from nick argument
		nickStr := strings.ReplaceAll(e.Arguments[len(e.Arguments)-1], "@", "")

		// mark nicks as present
		for _, nick := range strings.Split(nickStr, " ") {
			_, err := db.Exec(`insert into channel_nicks(channel, nick, present) values(?, ?, ?)`, channel, nick, true)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		channel := e.Arguments[0]
		msg := e.Arguments[1]
		nick := e.Nick

		matchNote(irccon, db, msg, nick, channel)
		matchLink(irccon, db, msg, nick, channel)
		matchLater(irccon, db, msg, nick, channel)
		matchCommand(irccon, db, msg, nick, channel)
	})
	irccon.AddCallback("JOIN", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
			sendLaters(irccon, db, channel, e.Nick)
		}
	})
	irccon.AddCallback("PART", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	irccon.AddCallback("QUIT", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	irccon.AddCallback("NICK", func(e *irc.Event) {
		if e.Nick != nick {
			// trigger NAMES to update the list of joined nicks
			irccon.SendRawf("NAMES %s", channel)
		}
	})
	err := irccon.Connect(server)

	return irccon, err
}

func sendLaters(irccon *irc.Connection, db *sqlx.DB, channel string, nick string) {
	// loop through each later message and see if the prefix matches this nick
	laters, err := getLaters(db)
	if err != nil {
		log.Fatal(err)
	}
	for _, later := range laters {
		if strings.Contains(nick, later.Target) {
			_, err := db.Exec(`delete from laters where rowid = ?`, later.RowId)
			if err != nil {
				log.Fatal(err)
			}
			irccon.Privmsgf(channel, "%s: %s (from %s %s ago)", nick, later.Message, later.Nick, since(later.CreatedAt))
		}
	}
}

func matchLater(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) {
	re := regexp.MustCompile(`^([^\s:]+): (.+)$`)
	matches := re.FindSubmatch([]byte(msg))

	if len(matches) > 0 {
		prefix := string(matches[1])
		message := string(matches[2])

		// if the prefix matches a currently joined nick, we do nothing
		if prefixMatchesJoinedNick(db, channel, prefix) {
			return
		}

		if prefixMatchesKnownNick(db, channel, prefix) {
			_, err := db.Exec(`insert into laters values(datetime('now'), ?, ?, ?, ?)`, nick, prefix, message, false)
			if err != nil {
				log.Fatal(err)
			}

			irccon.Privmsgf(channel, "%s: will send to %s* later", nick, prefix)
		} else {
			irccon.Privmsgf(channel, "%s: %s* doesn't match any known nick", nick, prefix)
		}
	}
}

func prefixMatchesJoinedNick(db *sqlx.DB, channel, prefix string) bool {
	channel_nicks := []ChannelNick{}
	err := db.Select(&channel_nicks, `select channel, nick, present from channel_nicks where present = true`)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range channel_nicks {
		if strings.HasPrefix(row.Nick, prefix) {
			return true
		}
	}
	return false
}

func prefixMatchesKnownNick(db *sqlx.DB, channel, prefix string) bool {
	channel_nicks := []ChannelNick{}
	err := db.Select(&channel_nicks, `select channel, nick, present from channel_nicks`)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range channel_nicks {
		if strings.HasPrefix(row.Nick, prefix) {
			return true
		}
	}
	return false
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
	re := regexp.MustCompile(`(https?://\S+)`)
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

		// post to twitter
		nvurl := os.Getenv("NICHE_VOMIT_URL")
		if nvurl != "" {
			res, err := http.Post(nvurl, "text/plain", strings.NewReader(url))
			if res.StatusCode >= 300 || err != nil {
				log.Printf("error posting to twitter %d %v\n", res.StatusCode, err)
			}
		}
	}
}

func since(tstr string) string {
	t, err := time.Parse("2006-01-02 15:04:05", tstr)
	if err != nil {
		log.Fatal(err)
	}
	return ago(time.Now().Sub(t).Round(time.Second))
}

func ago(d time.Duration) string {
	if d.Hours() >= 48.0 {
		return fmt.Sprintf("%dd", int(math.Round(d.Hours()/24)))
	} else {
		return d.String()
	}
}

func matchCommand(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) {
	re := regexp.MustCompile(`^!feedme`)
	match := re.Find([]byte(msg))

	if len(match) == 0 {
		return
	}

	notes := []Note{}

	err := db.Select(&notes, `select created_at, nick, text, kind from notes order by random() limit 1`)
	if err != nil {
		irccon.Privmsgf(channel, "%v", err)
		return
	}
	if len(notes) >= 1 {
		note := notes[0]
		irccon.Privmsgf(channel, "%s (from %s %s ago)", note.Text, note.Nick, since(note.CreatedAt))
	}
}
