package main

import (
	"bytes"
	"crypto/tls"
	_ "embed"
	"fmt"
	"goirc/fin"
	"goirc/trader"
	"strings"
	"text/template"
	"time"

	//	"database/sql"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/BurntSushi/migration"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	irc "github.com/thoj/go-ircevent"
	_ "modernc.org/sqlite"
)

type ChannelNick struct {
	Channel   string
	Nick      string
	Present   string
	UpdatedAt string `db:"updated_at"`
}

type Note struct {
	Id        int64
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
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add primary key to notes")
			_, err := tx.Exec(`
pragma foreign_key = off;

alter table notes rename to old_notes;

create table notes(
  id INTEGER not null primary key,
  created_at datetime not null default current_timestamp,
  nick text,
  text text,
  kind string not null default "note"
);

insert into notes select rowid, * from old_notes;

drop table old_notes;

pragma foreign_key = on;
`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add seen table")
			_, err := tx.Exec(`
create table seen_by(
  created_at datetime not null default current_timestamp,
  note_id references notes not null,
  nick text not null
);`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add updated_at to channel_nicks")
			_, err := tx.Exec(`alter table channel_nicks add column updated_at text`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: transactions table")
			_, err := tx.Exec(`
create table transactions(
  created_at datetime not null default current_timestamp,
  nick text not null,
  verb text not null,
  symbol text not null,
  shares number not null,
  price number not null
);`)
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

		tmpl, err := template.New("name").Parse(rssTemplate)
		if err != nil {
			log.Fatal("error parsing template")
		}

		fnotes, err := formatNotesDates(notes)
		if err != nil {
			log.Fatalf("error formatting notes: %v", err)
		}

		out := new(bytes.Buffer)
		err = tmpl.Execute(out, gin.H{
			"notes": fnotes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		c.Data(http.StatusOK, "text/xml; charset=utf-8", out.Bytes())
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func parseTime(str string) (time.Time, error) {
	result, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		result, err = time.Parse("2006-01-02T15:04:05Z", str)
	}
	return result, err
}

func formatNotesDates(notes []Note) ([]Note, error) {
	result := []Note{}
	for _, n := range notes {
		newNote := n

		createdAt, err := parseTime(n.CreatedAt)
		if err != nil {
			return nil, err
		}

		newNote.CreatedAt = createdAt.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		result = append(result, newNote)
	}
	return result, nil
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
	irccon.VerboseCallbackHandler = false
	irccon.Debug = false
	irccon.UseTLS = true
	irccon.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	irccon.AddCallback("001", func(e *irc.Event) { irccon.Join(channel) })
	irccon.AddCallback("353", func(e *irc.Event) {
		// clear the presence of all channel nicks
		_, err := db.Exec(`update channel_nicks set updated_at = current_timestamp, present = false`)
		if err != nil {
			log.Fatal(err)
		}

		// remove @ op markers from nick argument
		nickStr := strings.ReplaceAll(e.Arguments[len(e.Arguments)-1], "@", "")

		// mark nicks as present and record timestamp which can be intepreted as 'last seen', or 'online since'
		for _, nick := range strings.Split(nickStr, " ") {
			_, err = db.Exec(`
insert into channel_nicks(updated_at, channel, nick, present) values(current_timestamp, ?, ?, ?)
on conflict(channel, nick) do update set updated_at = current_timestamp, present=excluded.present`,
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
			sendLaters(irccon, db, channel, e.Nick)

			go sendMissed(irccon, db, channel, e.Nick)
		}

		// trigger NAMES to update the list of joined nicks
		irccon.SendRawf("NAMES %s", channel)
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

func isAltNick(nick string) bool {
	return strings.HasSuffix(nick, "`") || strings.HasSuffix(nick, "_")
}

func sendMissed(irccon *irc.Connection, db *sqlx.DB, channel string, nick string) {
	if isAltNick(nick) {
		return
	}

	channelNick := ChannelNick{}
	err := db.Get(&channelNick, `select * from channel_nicks where present = 0 and channel = ? and nick = ?`, channel, nick)
	if err != nil {
		return
	}

	notes := []Note{}
	db.Select(&notes, "select * from notes where created_at > ? order by created_at asc limit 69", channelNick.UpdatedAt)

	if len(notes) > 0 {
		irccon.Privmsgf(nick, "Hi %s, you missed %d thing(s) in %s since %s:",
			nick, len(notes), channel, channelNick.UpdatedAt)

		for _, note := range notes {
			irccon.Privmsgf(nick, "%s (from %s %s ago)", note.Text, note.Nick, since(note.CreatedAt))
			time.Sleep(1 * time.Second)
		}
	}
}

type MatchHandler struct {
	Name     string
	Function func(irccon *irc.Connection, db *sqlx.DB, msg, nick, channel string) bool
}

func markAsSeen(db *sqlx.DB, noteId int64, target string) error {
	//db.Select(`select * from channel_nicks where channel = ?`, target)
	channelNicks, err := joinedNicks(db, target)
	if err != nil {
		return err
	}
	// for each channelNick insert a seen_by record
	for _, nick := range channelNicks {
		_, err := db.Exec(`insert into seen_by(note_id, nick) values(?, ?)`, noteId, nick.Nick)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertNote(db *sqlx.DB, target string, nick string, kind string, text string) error {
	result, err := db.Exec(`insert into notes(nick, text, kind) values(?, ?, ?) returning id`, nick, text, kind)
	if err != nil {
		return err
	} else {
		noteId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		err = markAsSeen(db, noteId, target)
		if err != nil {
			return err
		}
	}
	return nil
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

				text := string(matches[1])

				err := insertNote(db, target, nick, "note", text)
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

				err := insertNote(db, target, nick, "link", url)
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

			err := db.Select(&notes, `select id, created_at, nick, text, kind from notes order by random() limit 1`)
			if err != nil {
				irccon.Privmsgf(target, "%v", err)
			} else if len(notes) >= 1 {
				note := notes[0]
				irccon.Privmsgf(target, "%s (from %s %s ago)", note.Text, note.Nick, since(note.CreatedAt))
				err = markAsSeen(db, note.Id, target)
				if err != nil {
					irccon.Privmsg(target, err.Error())
				}
			}
			return true
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

			// TODO: markAsSeen
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
	{
		Name: "Match Until Command",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`world.?cup`)
			match := re.Find([]byte(msg))

			if len(match) == 0 {
				return false
			}

			end, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
			if err != nil {
				irccon.Privmsgf(target, "error: %v", err)
				return true
			}
			until := time.Until(end)
			irccon.Privmsgf(target, "the world cup will start in %.0f days", math.Round(until.Hours()/24))
			return true
		},
	},
	{
		Name: "Match stock price",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile("^[$]([A-Za-z-]+)")
			matches := re.FindSubmatch([]byte(msg))

			if len(matches) == 0 {
				return false
			}

			symbol := string(matches[1])

			data, err := fin.YahooFinanceFetch(symbol)
			if err != nil {
				irccon.Privmsgf(target, "error: %s", err)
				return true
			}

			result := data.QuoteSummary.Result[0]
			irccon.Privmsgf(target, "%s %s %f", strings.ToUpper(symbol), bareDomain(result.SummaryProfile.Website), result.FinancialData.CurrentPrice.Raw)

			return true
		},
	},
	{
		Name: "Match quote",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile(`^"(.+)$`)
			matches := re.FindSubmatch([]byte(msg))

			if len(matches) > 0 {
				if target == nick {
					irccon.Privmsg(target, "not your personal secretary")
					return false
				}

				text := string(matches[1])

				err := insertNote(db, target, nick, "quote", text)
				if err != nil {
					log.Print(err)
				}
				return true
			}
			return false
		},
	},
	{
		Name: "Trade stock",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile("^((buy|sell).*)$")
			matches := re.FindStringSubmatch(msg)

			if len(matches) == 0 {
				return false
			}

			reply, err := trader.Trade(nick, matches[1], db)
			if err != nil {
				irccon.Privmsgf(target, "error: %s", err)
				return true
			}

			if reply != "" {
				irccon.Privmsgf(target, "%s: %s", nick, reply)
				return true
			}

			return false
		},
	},
	{
		Name: "Trade: show holdings report",
		Function: func(irccon *irc.Connection, db *sqlx.DB, msg, nick, target string) bool {
			re := regexp.MustCompile("^((report).*)$")
			matches := re.FindStringSubmatch(msg)

			if len(matches) == 0 {
				return false
			}

			reply, err := trader.Report(nick, db)
			if err != nil {
				irccon.Privmsgf(target, "error: %s", err)
				return true
			}

			if reply != "" {
				irccon.Privmsgf(target, "%s: %s", nick, reply)
				return true
			}

			return false
		},
	},
}

// from a uri like https://www.google.com/abc?def=123 return google.com
func bareDomain(uri string) string {
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		// just punt and return the original uri
		return uri
	}
	return strings.Replace(parsedUrl.Host, "www.", "", 1)
}

func joinedNicks(db *sqlx.DB, channel string) ([]ChannelNick, error) {
	channelNicks := []ChannelNick{}
	err := db.Select(&channelNicks, `select channel, nick, present from channel_nicks where present = true and channel = ?`, channel)
	return channelNicks, err
}

func prefixMatchesJoinedNick(db *sqlx.DB, channel, prefix string) bool {
	channelNicks, err := joinedNicks(db, channel)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range channelNicks {
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
	t, err := parseTime(tstr)
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
