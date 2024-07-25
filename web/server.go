package web

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"goirc/db/model"
	"goirc/image"
	"goirc/internal/idstr"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jmoiron/sqlx"
	"github.com/kkdai/youtube/v2"

	"github.com/go-chi/chi/v5"
)

//go:embed "templates/index.gohtml"
var indexTemplate string

//go:embed "templates/note.gohtml"
var noteTemplate string

//go:embed "templates/rss.gohtml"
var rssTemplate string

//go:embed "templates/player.gohtml"
var playerTemplateContent string
var playerTemplate = template.Must(template.New("").Parse(playerTemplateContent))

type keyType int

var sessionKey keyType

const cookieKey = "annie"

func Serve(db *sqlx.DB) {
	r := chi.NewRouter()

	q := model.New(db.DB)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var value string
			c, err := r.Cookie(cookieKey)
			if err != nil {
				value = uuid.Must(uuid.NewV7()).String()
				http.SetCookie(w, &http.Cookie{
					Name:     cookieKey,
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

	r.Get("/snapshot.db", func(w http.ResponseWriter, r *http.Request) {
		os.Remove("/tmp/snapshot.db")
		if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, "/tmp/snapshot.db")
	})

	pacific, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		log.Fatal(err)
	}

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

	fs := http.FileServer(http.Dir(image.ImageFileBase))
	r.Handle("/images/*", http.StripPrefix("/images/", fs))

	addr := ":" + os.Getenv("PORT")
	log.Printf("web server listening on %s", addr)
	err = http.ListenAndServe(addr, r)
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
