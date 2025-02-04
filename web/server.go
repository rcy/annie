package web

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"goirc/bot"
	"goirc/db/model"
	"goirc/image"
	"goirc/internal/idstr"
	"goirc/internal/summary"
	db "goirc/model"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kkdai/youtube/v2"

	"github.com/go-chi/chi/v5"
)

//go:embed "templates/index.gohtml"
var indexTemplate string

//go:embed "templates/login.gohtml"
var loginTemplate string

//go:embed "templates/note.gohtml"
var noteTemplate string

//go:embed "templates/rss.gohtml"
var rssTemplate string

//go:embed "templates/player.gohtml"
var playerTemplateContent string
var playerTemplate = template.Must(template.New("").Parse(playerTemplateContent))

//go:embed "templates/generatedimage.gohtml"
var generatedImageTemplateContent string
var generatedImageTemplate = template.Must(template.New("").Parse(generatedImageTemplateContent))

//go:embed "templates/generatedimages.gohtml"
var generatedImagesTemplateContent string
var generatedImagesTemplate = template.Must(template.New("").Parse(generatedImagesTemplateContent))

var pacific = func(name string) *time.Location {
	result, err := time.LoadLocation(name)
	if err != nil {
		log.Fatal(err)
	}
	return result
}("America/Los_Angeles")

type keyType int

const (
	sessionKey keyType = iota
	nickKey
)

const sessionCookieKey = "annie.session"
const fromCookieKey = "annie.from"

type code string

type oneTimeCode struct {
	session string
	nick    string
}

var codes = make(map[code]oneTimeCode)

func HandleAuth(params bot.HandlerParams) error {
	if params.Nick == params.Target {
		params.Privmsgf(params.Nick, "cannot !auth privately, do it in channel")
		return nil
	}
	var c = code(strings.Split(uuid.Must(uuid.NewV4()).String(), "-")[0])
	codes[c] = oneTimeCode{nick: params.Nick}
	params.Privmsgf(params.Nick, "hi %s, login with this link: %s/login/code/%s", params.Nick, os.Getenv("ROOT_URL"), c)
	return nil
}

func HandleDeauth(params bot.HandlerParams) error {
	q := model.New(db.DB.DB)

	err := q.DeleteNickSessions(context.Background(), params.Nick)
	if err != nil {
		return err
	}
	params.Privmsgf(params.Nick, "%s: all your sessions have been destroyed on %s", params.Nick, os.Getenv("ROOT_URL"))
	return nil
}

func Serve(db *sqlx.DB, b *bot.Bot) {
	r := chi.NewRouter()

	q := model.New(db.DB)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var value string
			c, err := r.Cookie(sessionCookieKey)
			if err != nil {
				value = uuid.Must(uuid.NewV7()).String()
				http.SetCookie(w, &http.Cookie{
					Name:     sessionCookieKey,
					Value:    value,
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
					Expires:  time.Now().Add(time.Hour * 24 * 400),
				})
			} else {
				value = c.Value
			}

			ctx := context.WithValue(r.Context(), sessionKey, value)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	r.Route("/login", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			tmpl, err := template.New("").Parse(loginTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "promptNick", map[string]string{
				"botNick": b.Conn.GetNick(),
				"channel": b.Channel,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		r.Post("/nick", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			nick := r.FormValue("nick")

			skey := r.Context().Value(sessionKey).(string)

			var c = code(strings.Split(uuid.Must(uuid.NewV4()).String(), "-")[0])
			codes[c] = oneTimeCode{session: skey, nick: nick}

			_, err := q.ChannelNick(ctx, model.ChannelNickParams{Nick: nick, Channel: b.Channel, Present: true})
			if err != nil {
				http.Error(w, fmt.Sprintf("couldn't find %s in %s: %s", nick, b.Channel, err.Error()), http.StatusForbidden)
				return
			}

			b.Conn.Privmsgf(nick, "hi %s, login with this link: %s/login/code/%s", nick, os.Getenv("ROOT_URL"), c)

			tmpl, err := template.New("").Parse(loginTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = tmpl.ExecuteTemplate(w, "promptCode", map[string]string{
				"nick":    nick,
				"botNick": b.Conn.GetNick(),
				"channel": b.Channel,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		r.Get("/code/{code}", func(w http.ResponseWriter, r *http.Request) {
			c := code(chi.URLParam(r, "code"))
			otc, ok := codes[c]
			if !ok {
				http.Error(w, "invalid code", http.StatusBadRequest)
				return
			}

			delete(codes, c)

			sess := r.Context().Value(sessionKey).(string)

			err := q.CreateNickSession(r.Context(), model.CreateNickSessionParams{
				Session: sess,
				Nick:    otc.nick,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cookie, err := r.Cookie(fromCookieKey)
			if err != nil {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			// expire "from" cookie
			http.SetCookie(w, &http.Cookie{
				Name:     fromCookieKey,
				Value:    r.URL.Path,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				Expires:  time.Unix(0, 0),
			})
			http.Redirect(w, r, cookie.Value, http.StatusSeeOther)
		})
	})

	r.Get("/{sqid}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		sqid := chi.URLParam(r, "sqid")
		id, err := idstr.Decode(sqid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sess := r.Context().Value(sessionKey).(string)
		m := model.New(db.DB)

		note, err := m.Link(ctx, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = m.InsertVisit(r.Context(), model.InsertVisitParams{Session: sess, NoteID: note.ID})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, note.Text.String, http.StatusSeeOther)
	})

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := r.Context()
				key := r.Context().Value(sessionKey).(string)
				session, err := q.NickBySession(ctx, key)
				if err != nil {
					http.SetCookie(w, &http.Cookie{
						Name:     fromCookieKey,
						Value:    r.URL.Path,
						Path:     "/",
						Secure:   true,
						HttpOnly: true,
						Expires:  time.Now().Add(time.Hour),
					})
					http.Redirect(w, r, "/login", http.StatusFound)
					return
				}
				ctx = context.WithValue(ctx, nickKey, session.Nick)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		})

		r.Get("/snapshot.db", func(w http.ResponseWriter, r *http.Request) {
			os.Remove("/tmp/snapshot.db")
			if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.ServeFile(w, r, "/tmp/snapshot.db")
		})

		funcMap := template.FuncMap{
			"time": func(t time.Time) string {
				return t.In(pacific).Format("2006-01-02 15:04:05")
			},
		}

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			nick := r.URL.Query().Get("nick")

			notes, err := getNotes(r.Context(), q, nick)
			if err != nil {
				log.Fatal(err)
			}

			nicks, err := q.NicksWithNoteCount(r.Context())
			if err != nil {
				log.Fatal(err)
			}

			tmpl, err := template.New("name").Funcs(funcMap).Parse(indexTemplate)
			if err != nil {
				log.Fatal("error parsing template")
			}

			out := new(bytes.Buffer)
			err = tmpl.Execute(out, map[string]any{
				"nicks": nicks,
				"notes": notes,
			})
			if err != nil {
				log.Fatal("error executing template on data")
			}

			_, _ = w.Write(out.Bytes())
		})

		r.Get("/note/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id, _ := strconv.Atoi(chi.URLParam(r, "id"))

			note, err := q.NoteByID(ctx, int64(id))
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.New("name").Funcs(funcMap).Parse(noteTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = tmpl.Execute(w, map[string]any{
				"note": note,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})

		r.Post("/note/{id}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			id, _ := strconv.Atoi(chi.URLParam(r, "id"))
			text := r.FormValue("text")

			if text == "" {
				err := q.DeleteNoteByID(ctx, int64(id))
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				_, err := q.UpdateNoteTextByID(ctx, model.UpdateNoteTextByIDParams{
					ID:   int64(id),
					Text: sql.NullString{String: text, Valid: true},
				})
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, r.URL.String(), http.StatusSeeOther)
		})

		r.Get("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
			nick := r.URL.Query().Get("nick")

			notes, err := getNotes(r.Context(), q, nick)
			if err != nil {
				log.Fatal(err)
			}

			tmpl, err := template.New("name").Parse(rssTemplate)
			if err != nil {
				log.Fatal("error parsing template")
			}

			out := new(bytes.Buffer)
			err = tmpl.Execute(out, map[string]any{
				"notes": notes,
			})
			if err != nil {
				log.Fatal("error executing template on data")
			}

			_, _ = w.Write(out.Bytes())
		})

		r.Get("/player", func(w http.ResponseWriter, r *http.Request) {
			youtubeLinks, err := q.YoutubeLinks(r.Context())
			if err != nil {
				log.Fatal("could not select links")
			}

			var videoIDs []string
			for _, link := range youtubeLinks {
				id, err := youtube.ExtractVideoID(link.Text.String)
				if err != nil {
					log.Fatalf("error extracting video id %s", link.Text.String)
				}
				videoIDs = append(videoIDs, id)
			}

			out := new(bytes.Buffer)
			err = playerTemplate.Execute(out, map[string]any{"VideoIDs": videoIDs})
			if err != nil {
				log.Fatalf("error executing template: %s", err)
			}

			_, _ = w.Write(out.Bytes())
		})

		r.Get("/week/{date}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			q := model.New(db.DB)

			date, err := time.ParseInLocation(time.DateOnly, chi.URLParam(r, "date"), pacific)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			start := summary.WeekStart(date, pacific)
			end := start.Add(time.Hour * 24 * 7)

			s := summary.New(q, start, end)
			b, err := s.Cache(ctx, s.WeeklyNewsletter)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(b)
		})
	})

	r.Get("/generated_images/{id}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/i/"+chi.URLParam(r, "id"), http.StatusSeeOther)
	})

	r.Get("/i", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := model.New(db.DB)
		images, err := q.GeneratedImages(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = generatedImagesTemplate.Execute(w, map[string]any{
			"images": images,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	r.Get("/i/{id}", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := model.New(db.DB)
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		image, err := q.GeneratedImageByID(ctx, int64(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = generatedImageTemplate.Execute(w, map[string]any{
			"image":            image,
			"absoluteImageURL": fmt.Sprintf("%s/images/%d.png", os.Getenv("ROOT_URL"), image.ID),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fs := http.FileServer(http.Dir(image.ImageFileBase))
	r.Handle("/images/*",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%.0f", (time.Hour*24*365).Seconds()))
			http.StripPrefix("/images/", fs).ServeHTTP(w, r)
		}),
	)

	addr := ":" + os.Getenv("PORT")
	log.Printf("web server listening on %s", addr)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func getNotes(ctx context.Context, q *model.Queries, nick string) ([]model.Note, error) {
	if nick == "" {
		return q.AllNotes(ctx)
	}
	return q.AllNickNotes(ctx, sql.NullString{String: nick, Valid: true})
}
