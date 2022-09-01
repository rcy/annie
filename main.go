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
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add unique constrant to channel_nicks table")

			// delete duplicates, keeping oldest records
			_, err := tx.Exec(`delete from channel_nicks where rowid not in (select min(rowid) from channel_nicks group by nick, channel)`)
			if err != nil {
				return err
			}

			// add unique constraint
			_, err = tx.Exec(`create unique index channel_nick_unique_index on channel_nicks(channel, nick)`)
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

	go webserver(db)

	conn.Loop()
}

//go:embed "templates/index.gohtml"
var indexTemplate string

//go:embed "templates/rss.gohtml"
var rssTemplate string

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

	r.GET("/rss.xml", func(c *gin.Context) {
		nick := c.Query("nick")

		notes, err := getNotes(db, nick)
		if err != nil {
			log.Fatal(err)
		}

		nicks, err := getNicks(db)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("name").Parse(rssTemplate)
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

		c.Data(http.StatusOK, "text/xml; charset=utf-8", out.Bytes())
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
			_, err = db.Exec(`
insert into channel_nicks(channel, nick, present) values(?, ?, ?)
on conflict(channel, nick) do update set present=excluded.present`,
				channel, nick, true)
			if err != nil {
				log.Fatal(err)
			}
		}
	})
	irccon.AddCallback("366", func(e *irc.Event) {})
	irccon.AddCallback("PRIVMSG", func(e *irc.Event) {
		go handlePrivmsg(irccon, db, e)
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

func handlePrivmsg(irccon *irc.Connection, db *sqlx.DB, e *irc.Event) {
	channel := e.Arguments[0]
	msg := e.Arguments[1]
	nick := e.Nick

	for _, f := range matchHandlers {
		var target string

		if channel == irccon.GetNick() {
			target = nick
		} else {
			target = channel
		}

		if f.Function(irccon, db, msg, nick, target) {
			break
		}
	}
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

type MatchHandler struct {
	Name     string
	Function func(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) bool
}

var matchHandlers = []MatchHandler{
	{
		Name: "Match Create Note",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`^,(.+)$`)
			matches := re.FindSubmatch([]byte(msg))

			if len(matches) > 0 {
				if target == nick {
					irccon.Privmsg(target, "not your personal secretary")
					return false
				}

				note := string(matches[1])

				_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?, 'note')`, nick, note)
				if err != nil {
					log.Print(err)
					irccon.Privmsg(target, err.Error())
				} else {
					irccon.Privmsg(target, "recorded note")
				}
				return true
			}
			return false
		},
	},
	{
		Name: "Match Deferred Delivery",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`^([^\s:]+): (.+)$`)
			matches := re.FindSubmatch([]byte(msg))

			if len(matches) > 0 {
				if target == nick {
					irccon.Privmsg(target, "not your personal secretary")
					return false
				}

				prefix := string(matches[1])
				message := string(matches[2])

				// if the prefix matches a currently joined nick, we do nothing
				if prefixMatchesJoinedNick(db, target, prefix) {
					return false
				}

				if prefixMatchesKnownNick(db, target, prefix) {
					_, err := db.Exec(`insert into laters values(datetime('now'), ?, ?, ?, ?)`, nick, prefix, message, false)
					if err != nil {
						log.Fatal(err)
					}

					irccon.Privmsgf(target, "%s: will send to %s* later", nick, prefix)
				} else {
					irccon.Privmsgf(target, "%s: %s* doesn't match any known nick", nick, prefix)
				}
				return true
			}
			return false
		},
	},
	{
		Name: "Match Link",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`(https?://\S+)`)
			matches := re.FindSubmatch([]byte(msg))

			if len(matches) > 0 {
				if target == nick {
					irccon.Privmsg(target, "not your personal secretary")
					return false
				}

				url := string(matches[1])
				_, err := db.Exec(`insert into notes values(datetime('now'), ?, ?, 'link')`, nick, url)
				if err != nil {
					log.Print(err)
					irccon.Privmsg(target, err.Error())
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

				return true
			}
			return false
		},
	},
	{
		Name: "Match Feedme Command",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`^!feedme`)
			match := re.Find([]byte(msg))

			if len(match) == 0 {
				return false
			}

			notes := []Note{}

			err := db.Select(&notes, `select created_at, nick, text, kind from notes order by random() limit 1`)
			if err != nil {
				irccon.Privmsgf(target, "%v", err)
				return false
			}
			if len(notes) >= 1 {
				note := notes[0]
				irccon.Privmsgf(target, "%s (from %s %s ago)", note.Text, note.Nick, since(note.CreatedAt))
				return true
			}
			return false
		},
	},
	{
		Name: "Match Catchup Command",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`^!catchup`)
			match := re.Find([]byte(msg))

			if len(match) == 0 {
				return false
			}

			notes := []Note{}

			err := db.Select(&notes, `select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc`)
			if err != nil {
				irccon.Privmsgf(target, "%v", err)
				return false
			}
			if len(notes) >= 1 {
				for _, note := range notes {
					irccon.Privmsgf(nick, "%s (from %s %s ago)", note.Text, note.Nick, since(note.CreatedAt))
					time.Sleep(1 * time.Second)
				}
			}
			irccon.Privmsgf(nick, "--- %d total from last 24 hours", len(notes))
			return true
		},
	},
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
