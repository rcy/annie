// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package model

import (
	"database/sql"
	"time"
)

type Bedtime struct {
	ID        int64
	CreatedAt time.Time
	Nick      string
	Message   sql.NullString
}

type ChannelNick struct {
	Channel   string
	Nick      string
	Present   bool
	UpdatedAt time.Time
}

type FutureMessage struct {
	ID        int64
	CreatedAt time.Time
	Kind      string
}

type GeneratedImage struct {
	ID            int64
	CreatedAt     time.Time
	Filename      string
	Prompt        string
	RevisedPrompt string
}

type Later struct {
	CreatedAt sql.NullString
	Nick      sql.NullString
	Target    sql.NullString
	Message   sql.NullString
	Sent      sql.NullBool
}

type Link struct {
	CreatedAt sql.NullString
	Nick      sql.NullString
	Text      sql.NullString
}

type MigrationVersion struct {
	Version sql.NullInt64
}

type NickSession struct {
	ID        int64
	CreatedAt time.Time
	Nick      string
	Session   string
}

type NickWeatherRequest struct {
	ID        int64
	CreatedAt time.Time
	Nick      string
	Query     string
	City      string
	Country   string
}

type Note struct {
	ID        int64
	CreatedAt time.Time
	Nick      sql.NullString
	Text      sql.NullString
	Kind      string
	Target    string
	Anon      bool
}

type Reminder struct {
	ID        int64
	CreatedAt time.Time
	Nick      string
	RemindAt  time.Time
	What      string
}

type Rev struct {
	ID        int64
	CreatedAt time.Time
	Sha       string
}

type Visit struct {
	ID        int64
	CreatedAt time.Time
	Session   string
	NoteID    int64
}
