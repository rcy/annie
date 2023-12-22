package web

import (
	"bytes"
	"database/sql"
	_ "embed"
	"errors"
	"goirc/internal/idstr"
	"goirc/model/notes"
	"goirc/util"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/kkdai/youtube/v2"

	"github.com/go-chi/chi/v5"
)

type NickWithNoteCount struct {
	Nick  string
	Count int
}

//go:embed "templates/index.gohtml"
var indexTemplate string

//go:embed "templates/rss.gohtml"
var rssTemplate string

//go:embed "templates/player.gohtml"
var playerTemplateContent string
var playerTemplate = template.Must(template.New("").Parse(playerTemplateContent))

func Serve(db *sqlx.DB) {
	r := chi.NewRouter()

	r.Get("/snapshot.db", func(w http.ResponseWriter, r *http.Request) {
		os.Remove("/tmp/snapshot.db")
		if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeFile(w, r, "/tmp/snapshot.db")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		nick := r.URL.Query().Get("nick")

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
		err = tmpl.Execute(out, map[string]any{
			"nicks": nicks,
			"notes": notes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		w.Write(out.Bytes())
	})

	r.Get("/rss.xml", func(w http.ResponseWriter, r *http.Request) {
		nick := r.URL.Query().Get("nick")

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
		err = tmpl.Execute(out, map[string]any{
			"notes": fnotes,
		})
		if err != nil {
			log.Fatal("error executing template on data")
		}

		w.Write(out.Bytes())
	})

	r.Get("/player", func(w http.ResponseWriter, r *http.Request) {
		var youtubeLinks []notes.Note
		err := db.Select(&youtubeLinks, "select * from notes where kind = 'link' and text like '%youtube.com%' or text like '%youtu.be%'")
		if err != nil {
			log.Fatal("could not select links")
		}

		var videoIDs []string
		for _, link := range youtubeLinks {
			id, err := youtube.ExtractVideoID(link.Text)
			if err != nil {
				log.Fatalf("error extracting video id %s", link.Text)
			}
			videoIDs = append(videoIDs, id)
		}

		out := new(bytes.Buffer)
		err = playerTemplate.Execute(out, map[string]any{"VideoIDs": videoIDs})
		if err != nil {
			log.Fatalf("error executing template: %s", err)
		}

		w.Write(out.Bytes())
	})

	r.Get("/{sqid}", func(w http.ResponseWriter, r *http.Request) {
		sqid := chi.URLParam(r, "sqid")
		id, err := idstr.Decode(sqid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var text string
		err = db.Get(&text, `select text from notes where id = ? and kind = 'link'`, id)
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, text, http.StatusSeeOther)
	})

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}

func getNotes(db *sqlx.DB, nick string) ([]notes.Note, error) {
	result := []notes.Note{}
	var err error
	if nick == "" {
		err = db.Select(&result, `select created_at, text, nick, kind from notes where target != nick order by created_at desc limit 10000`)
	} else {
		err = db.Select(&result, `select created_at, text, nick, kind from notes where target != nick and nick = ? order by created_at desc limit 10000`, nick)
	}
	return result, err
}

func getNicks(db *sqlx.DB) ([]NickWithNoteCount, error) {
	nicks := []NickWithNoteCount{}
	err := db.Select(&nicks, `select nick, count(nick) as count from notes group by nick`)
	return nicks, err
}

func formatNotesDates(narr []notes.Note) ([]notes.Note, error) {
	result := []notes.Note{}
	for _, n := range narr {
		newNote := n

		createdAt, err := util.ParseTime(n.CreatedAt)
		if err != nil {
			return nil, err
		}

		newNote.CreatedAt = createdAt.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		result = append(result, newNote)
	}
	return result, nil
}
