package notes

import (
	"errors"
	"goirc/model"
)

type Note struct {
	Id        int64
	CreatedAt string `db:"created_at"`
	Target    string
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
	query := `insert into notes(target, nick, text, kind) values(?, ?, ?, ?) returning *`

	if p.Target == "" {
		return nil, errors.New("target cannot be empty")
	}

	err := model.DB.Get(&note, query, p.Target, p.Nick, p.Text, p.Kind)
	if err != nil {
		return nil, err
	}
	return &note, nil
}
