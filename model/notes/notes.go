package notes

import "log"

type Note struct {
	Id        int64
	CreatedAt string `db:"created_at"`
	Text      string
	Nick      string
	Kind      string
}

func CreateNote(url string) {
	log.Printf("created %s", url)
}
