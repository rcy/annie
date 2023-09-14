package notes

import (
	"goirc/model"
)

type Note struct {
	Id        int64
	CreatedAt string `db:"created_at"`
	Text      string
	Nick      string
	Kind      string
}

type CreateParams struct {
	Target string
	Nick   string
	Kind   string
	Text   string
}

func Create(p CreateParams) (*Note, error) {
	var note Note
	err := model.DB.Get(&note, `insert into notes(nick, text, kind) values(?, ?, ?) returning *`, p.Nick, p.Text, p.Kind)
	if err != nil {
		return nil, err
	}
	return &note, nil
}
