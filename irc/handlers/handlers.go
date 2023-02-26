package handers

import (
	"goirc/fin"
	"goirc/model"
	"goirc/trader"
	"goirc/util"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	irc "github.com/thoj/go-ircevent"
)

type HandlerParams struct {
	Irccon *irc.Connection
	Db     *sqlx.DB
	Msg    string
	Nick   string
	Target string
}

type Handler struct {
	Name     string
	Function func(HandlerParams) bool
}

var Handlers = []Handler{
	{
		Name: "Match Create Note",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`^,(.+)$`)
			matches := re.FindSubmatch([]byte(params.Msg))

			if len(matches) > 0 {
				if params.Target == params.Nick {
					params.Irccon.Privmsg(params.Target, "not your personal secretary")
					return false
				}

				text := string(matches[1])

				err := insertNote(params.Db, params.Target, params.Nick, "note", text)
				if err != nil {
					log.Print(err)
					params.Irccon.Privmsg(params.Target, err.Error())
				} else {
					params.Irccon.Privmsg(params.Target, "recorded note")
				}
				return true
			}
			return false
		},
	},
	{
		Name: "Match Deferred Delivery",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`^([^\s:]+): (.+)$`)
			matches := re.FindSubmatch([]byte(params.Msg))

			if len(matches) > 0 {
				if params.Target == params.Nick {
					params.Irccon.Privmsg(params.Target, "not your personal secretary")
					return false
				}

				prefix := string(matches[1])
				message := string(matches[2])

				// if the prefix matches a currently joined nick, we do nothing
				if model.PrefixMatchesJoinedNick(params.Db, params.Target, prefix) {
					return false
				}

				if model.PrefixMatchesKnownNick(params.Db, params.Target, prefix) {
					_, err := params.Db.Exec(`insert into laters values(datetime('now'), ?, ?, ?, ?)`, params.Nick, prefix, message, false)
					if err != nil {
						log.Fatal(err)
					}

					params.Irccon.Privmsgf(params.Target, "%s: will send to %s* later", params.Nick, prefix)
				}
				return true
			}
			return false
		},
	},
	{
		Name: "Match Link",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`(https?://\S+)`)
			matches := re.FindSubmatch([]byte(params.Msg))

			if len(matches) > 0 {
				if params.Target == params.Nick {
					params.Irccon.Privmsg(params.Target, "not your personal secretary")
					return false
				}

				url := string(matches[1])

				err := insertNote(params.Db, params.Target, params.Nick, "link", url)
				if err != nil {
					log.Print(err)
					params.Irccon.Privmsg(params.Target, err.Error())
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
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`^!feedme`)
			match := re.Find([]byte(params.Msg))

			if len(match) == 0 {
				return false
			}

			notes := []model.Note{}

			err := params.Db.Select(&notes, `select id, created_at, nick, text, kind from notes order by random() limit 1`)
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "%v", err)
			} else if len(notes) >= 1 {
				note := notes[0]
				params.Irccon.Privmsgf(params.Target, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
				err = markAsSeen(params.Db, note.Id, params.Target)
				if err != nil {
					params.Irccon.Privmsg(params.Target, err.Error())
				}
			}
			return true
		},
	},
	{
		Name: "Match Catchup Command",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`^!catchup`)
			match := re.Find([]byte(params.Msg))

			if len(match) == 0 {
				return false
			}

			notes := []model.Note{}

			// TODO: markAsSeen
			err := params.Db.Select(&notes, `select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc`)
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "%v", err)
				return false
			}
			if len(notes) >= 1 {
				for _, note := range notes {
					params.Irccon.Privmsgf(params.Nick, "%s (from %s %s ago)", note.Text, note.Nick, util.Since(note.CreatedAt))
					time.Sleep(1 * time.Second)
				}
			}
			params.Irccon.Privmsgf(params.Nick, "--- %d total from last 24 hours", len(notes))
			return true
		},
	},
	{
		Name: "Match Until Command",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile(`world.?cup`)
			match := re.Find([]byte(params.Msg))

			if len(match) == 0 {
				return false
			}

			end, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "error: %v", err)
				return true
			}
			until := time.Until(end)
			params.Irccon.Privmsgf(params.Target, "the world cup will start in %.0f days", math.Round(until.Hours()/24))
			return true
		},
	},
	{
		Name: "Match stock price",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile("^[$]([A-Za-z-]+)")
			matches := re.FindSubmatch([]byte(params.Msg))

			if len(matches) == 0 {
				return false
			}

			symbol := string(matches[1])

			data, err := fin.YahooFinanceFetch(symbol)
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "error: %s", err)
				return true
			}

			result := data.QuoteSummary.Result[0]
			params.Irccon.Privmsgf(params.Target, "%s %s %f", strings.ToUpper(symbol), util.BareDomain(result.SummaryProfile.Website), result.FinancialData.CurrentPrice.Raw)

			return true
		},
	},
	{
		Name: "Match quote",
		Function: func(params HandlerParams) bool {
			// match anything that starts with a quote and has no subsequent quotes
			re := regexp.MustCompile(`^("[^"]+)$`)
			matches := re.FindSubmatch([]byte(params.Msg))

			if len(matches) > 0 {
				if params.Target == params.Nick {
					params.Irccon.Privmsg(params.Target, "not your personal secretary")
					return false
				}

				text := string(matches[1])

				err := insertNote(params.Db, params.Target, params.Nick, "quote", text)
				if err != nil {
					log.Print(err)
				}

				// post to twitter
				nvurl := os.Getenv("NICHE_VOMIT_URL")
				if nvurl != "" {
					res, err := http.Post(nvurl, "text/plain", strings.NewReader(text))
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
		Name: "Trade stock",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile("^((buy|sell).*)$")
			matches := re.FindStringSubmatch(params.Msg)

			if len(matches) == 0 {
				return false
			}

			reply, err := trader.Trade(params.Nick, matches[1], params.Db)
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "error: %s", err)
				return true
			}

			if reply != "" {
				params.Irccon.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
				return true
			}

			return false
		},
	},
	{
		Name: "Trade: show holdings report",
		Function: func(params HandlerParams) bool {
			re := regexp.MustCompile("^((report).*)$")
			matches := re.FindStringSubmatch(params.Msg)

			if len(matches) == 0 {
				return false
			}

			reply, err := trader.Report(params.Nick, params.Db)
			if err != nil {
				params.Irccon.Privmsgf(params.Target, "error: %s", err)
				return true
			}

			if reply != "" {
				params.Irccon.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
				return true
			}

			return false
		},
	},
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

func markAsSeen(db *sqlx.DB, noteId int64, target string) error {
	//db.Select(`select * from channel_nicks where channel = ?`, target)
	channelNicks, err := model.JoinedNicks(db, target)
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
