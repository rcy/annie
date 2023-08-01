package web

import (
	"bytes"
	_ "embed"
	"fmt"
	"goirc/model/notes"
	"goirc/util"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/kkdai/youtube/v2"
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
	r := gin.Default()
	//r.LoadHTMLGlob("templates/*")

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"Origin"},
	}))

	r.GET("/snapshot.db", func(c *gin.Context) {
		os.Remove("/tmp/snapshot.db")
		if _, err := db.Exec(`vacuum into '/tmp/snapshot.db'`); err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			return
		}
		c.File("/tmp/snapshot.db")
	})

	r.HEAD("/snapshot.db", func(c *gin.Context) {
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

	r.GET("/player", func(c *gin.Context) {
		var youtubeLinks []notes.Note
		err := db.Select(&youtubeLinks, "select * from notes where kind = 'link' and text like '%youtube%'")
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
		err = playerTemplate.Execute(out, gin.H{"VideoIDs": videoIDs})
		if err != nil {
			log.Fatalf("error executing template: %s", err)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", out.Bytes())
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func getNotes(db *sqlx.DB, nick string) ([]notes.Note, error) {
	result := []notes.Note{}
	var err error
	if nick == "" {
		err = db.Select(&result, `select created_at, text, nick, kind from notes order by created_at desc limit 1000`)
	} else {
		err = db.Select(&result, `select created_at, text, nick, kind from notes where nick = ? order by created_at desc limit 1000`, nick)
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
