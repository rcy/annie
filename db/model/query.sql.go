// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package model

import (
	"context"
	"database/sql"
	"time"
)

const allNickNotes = `-- name: AllNickNotes :many
select id, created_at, nick, text, kind, target, anon from notes where target != nick and nick = ? order by created_at desc limit 10000
`

func (q *Queries) AllNickNotes(ctx context.Context, nick sql.NullString) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, allNickNotes, nick)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
			&i.Target,
			&i.Anon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const allNotes = `-- name: AllNotes :many
select id, created_at, nick, text, kind, target, anon from notes where target != nick order by created_at desc limit 10000
`

func (q *Queries) AllNotes(ctx context.Context) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, allNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
			&i.Target,
			&i.Anon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const channelNick = `-- name: ChannelNick :one
select channel, nick, present, updated_at from channel_nicks where present = ? and channel = ? and nick = ? collate nocase
`

type ChannelNickParams struct {
	Present bool
	Channel string
	Nick    string
}

func (q *Queries) ChannelNick(ctx context.Context, arg ChannelNickParams) (ChannelNick, error) {
	row := q.db.QueryRowContext(ctx, channelNick, arg.Present, arg.Channel, arg.Nick)
	var i ChannelNick
	err := row.Scan(
		&i.Channel,
		&i.Nick,
		&i.Present,
		&i.UpdatedAt,
	)
	return i, err
}

const channelNotesSince = `-- name: ChannelNotesSince :many
select id, created_at, nick, text, kind, target, anon from notes where target = ? and created_at > ? order by created_at asc limit 69
`

type ChannelNotesSinceParams struct {
	Target    string
	CreatedAt time.Time
}

func (q *Queries) ChannelNotesSince(ctx context.Context, arg ChannelNotesSinceParams) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, channelNotesSince, arg.Target, arg.CreatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
			&i.Target,
			&i.Anon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createGeneratedImage = `-- name: CreateGeneratedImage :one
insert into generated_images(filename, prompt, revised_prompt) values(?,?,?) returning id, created_at, filename, prompt, revised_prompt
`

type CreateGeneratedImageParams struct {
	Filename      string
	Prompt        string
	RevisedPrompt string
}

func (q *Queries) CreateGeneratedImage(ctx context.Context, arg CreateGeneratedImageParams) (GeneratedImage, error) {
	row := q.db.QueryRowContext(ctx, createGeneratedImage, arg.Filename, arg.Prompt, arg.RevisedPrompt)
	var i GeneratedImage
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Filename,
		&i.Prompt,
		&i.RevisedPrompt,
	)
	return i, err
}

const createNickSession = `-- name: CreateNickSession :exec
insert into nick_sessions(nick, session) values(?,?)
`

type CreateNickSessionParams struct {
	Nick    string
	Session string
}

func (q *Queries) CreateNickSession(ctx context.Context, arg CreateNickSessionParams) error {
	_, err := q.db.ExecContext(ctx, createNickSession, arg.Nick, arg.Session)
	return err
}

const deleteNickSessions = `-- name: DeleteNickSessions :exec
delete from nick_sessions where nick = ?
`

func (q *Queries) DeleteNickSessions(ctx context.Context, nick string) error {
	_, err := q.db.ExecContext(ctx, deleteNickSessions, nick)
	return err
}

const deleteNoteByID = `-- name: DeleteNoteByID :exec
delete from notes where id = ?
`

func (q *Queries) DeleteNoteByID(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteNoteByID, id)
	return err
}

const generatedImageByID = `-- name: GeneratedImageByID :one
select id, created_at, filename, prompt, revised_prompt from generated_images where id = ?
`

func (q *Queries) GeneratedImageByID(ctx context.Context, id int64) (GeneratedImage, error) {
	row := q.db.QueryRowContext(ctx, generatedImageByID, id)
	var i GeneratedImage
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Filename,
		&i.Prompt,
		&i.RevisedPrompt,
	)
	return i, err
}

const generatedImages = `-- name: GeneratedImages :many
select id, created_at, filename, prompt, revised_prompt from generated_images order by created_at desc
`

func (q *Queries) GeneratedImages(ctx context.Context) ([]GeneratedImage, error) {
	rows, err := q.db.QueryContext(ctx, generatedImages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GeneratedImage
	for rows.Next() {
		var i GeneratedImage
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Filename,
			&i.Prompt,
			&i.RevisedPrompt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertNickWeatherRequest = `-- name: InsertNickWeatherRequest :exec
insert into nick_weather_requests(nick, query, city, country) values(?,?,?,?)
`

type InsertNickWeatherRequestParams struct {
	Nick    string
	Query   string
	City    string
	Country string
}

func (q *Queries) InsertNickWeatherRequest(ctx context.Context, arg InsertNickWeatherRequestParams) error {
	_, err := q.db.ExecContext(ctx, insertNickWeatherRequest,
		arg.Nick,
		arg.Query,
		arg.City,
		arg.Country,
	)
	return err
}

const insertNote = `-- name: InsertNote :one
insert into notes(target, nick, kind, text, anon) values(?,?,?,?,?) returning id, created_at, nick, text, kind, target, anon
`

type InsertNoteParams struct {
	Target string
	Nick   sql.NullString
	Kind   string
	Text   sql.NullString
	Anon   bool
}

func (q *Queries) InsertNote(ctx context.Context, arg InsertNoteParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, insertNote,
		arg.Target,
		arg.Nick,
		arg.Kind,
		arg.Text,
		arg.Anon,
	)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const insertVisit = `-- name: InsertVisit :exec
insert into visits(session, note_id) values(?,?)
`

type InsertVisitParams struct {
	Session string
	NoteID  int64
}

func (q *Queries) InsertVisit(ctx context.Context, arg InsertVisitParams) error {
	_, err := q.db.ExecContext(ctx, insertVisit, arg.Session, arg.NoteID)
	return err
}

const lastDaysNotes = `-- name: LastDaysNotes :many
select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc
`

type LastDaysNotesRow struct {
	CreatedAt time.Time
	Nick      sql.NullString
	Text      sql.NullString
	Kind      string
}

func (q *Queries) LastDaysNotes(ctx context.Context) ([]LastDaysNotesRow, error) {
	rows, err := q.db.QueryContext(ctx, lastDaysNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []LastDaysNotesRow
	for rows.Next() {
		var i LastDaysNotesRow
		if err := rows.Scan(
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const lastNickWeatherRequest = `-- name: LastNickWeatherRequest :one
select id, created_at, nick, "query", city, country from nick_weather_requests where nick = ? order by created_at desc limit 1
`

func (q *Queries) LastNickWeatherRequest(ctx context.Context, nick string) (NickWeatherRequest, error) {
	row := q.db.QueryRowContext(ctx, lastNickWeatherRequest, nick)
	var i NickWeatherRequest
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Query,
		&i.City,
		&i.Country,
	)
	return i, err
}

const lastWeatherRequestByPrefix = `-- name: LastWeatherRequestByPrefix :one
select id, created_at, nick, "query", city, country from nick_weather_requests where city like ? || '%' order by created_at desc limit 1
`

func (q *Queries) LastWeatherRequestByPrefix(ctx context.Context, dollar_1 sql.NullString) (NickWeatherRequest, error) {
	row := q.db.QueryRowContext(ctx, lastWeatherRequestByPrefix, dollar_1)
	var i NickWeatherRequest
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Query,
		&i.City,
		&i.Country,
	)
	return i, err
}

const link = `-- name: Link :one
select id, created_at, nick, text, kind, target, anon from notes where id = ? and kind = 'link'
`

func (q *Queries) Link(ctx context.Context, id int64) (Note, error) {
	row := q.db.QueryRowContext(ctx, link, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const markAnonymousNoteDelivered = `-- name: MarkAnonymousNoteDelivered :one
update notes set target = ?, created_at = current_timestamp where id = ? returning id, created_at, nick, text, kind, target, anon
`

type MarkAnonymousNoteDeliveredParams struct {
	Target string
	ID     int64
}

func (q *Queries) MarkAnonymousNoteDelivered(ctx context.Context, arg MarkAnonymousNoteDeliveredParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, markAnonymousNoteDelivered, arg.Target, arg.ID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const nickBySession = `-- name: NickBySession :one
select id, created_at, nick, session from nick_sessions where session = ?
`

func (q *Queries) NickBySession(ctx context.Context, session string) (NickSession, error) {
	row := q.db.QueryRowContext(ctx, nickBySession, session)
	var i NickSession
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Session,
	)
	return i, err
}

const nicksWithNoteCount = `-- name: NicksWithNoteCount :many
select nick, count(nick) as count from notes group by nick
`

type NicksWithNoteCountRow struct {
	Nick  sql.NullString
	Count int64
}

func (q *Queries) NicksWithNoteCount(ctx context.Context) ([]NicksWithNoteCountRow, error) {
	rows, err := q.db.QueryContext(ctx, nicksWithNoteCount)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []NicksWithNoteCountRow
	for rows.Next() {
		var i NicksWithNoteCountRow
		if err := rows.Scan(&i.Nick, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const noteByID = `-- name: NoteByID :one
select id, created_at, nick, text, kind, target, anon from notes where id = ?
`

func (q *Queries) NoteByID(ctx context.Context, id int64) (Note, error) {
	row := q.db.QueryRowContext(ctx, noteByID, id)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const randomHistoricalTodayNote = `-- name: RandomHistoricalTodayNote :one
select id, created_at, nick, text, kind, target, anon from notes
where
  strftime('%m-%d', created_at) = strftime('%m-%d', 'now')
and
  strftime('%Y', created_at) != strftime('%Y', 'now')
order by random()
limit 1
`

func (q *Queries) RandomHistoricalTodayNote(ctx context.Context) (Note, error) {
	row := q.db.QueryRowContext(ctx, randomHistoricalTodayNote)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const unsentAnonymousNotes = `-- name: UnsentAnonymousNotes :many
select id, created_at, nick, text, kind, target, anon from notes where created_at <= ? and kind = ? and nick = target order by id asc limit 420
`

type UnsentAnonymousNotesParams struct {
	CreatedAt time.Time
	Kind      string
}

func (q *Queries) UnsentAnonymousNotes(ctx context.Context, arg UnsentAnonymousNotesParams) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, unsentAnonymousNotes, arg.CreatedAt, arg.Kind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
			&i.Target,
			&i.Anon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateNoteTextByID = `-- name: UpdateNoteTextByID :one
update notes set text = ? where id = ? returning id, created_at, nick, text, kind, target, anon
`

type UpdateNoteTextByIDParams struct {
	Text sql.NullString
	ID   int64
}

func (q *Queries) UpdateNoteTextByID(ctx context.Context, arg UpdateNoteTextByIDParams) (Note, error) {
	row := q.db.QueryRowContext(ctx, updateNoteTextByID, arg.Text, arg.ID)
	var i Note
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Nick,
		&i.Text,
		&i.Kind,
		&i.Target,
		&i.Anon,
	)
	return i, err
}

const youtubeLinks = `-- name: YoutubeLinks :many
select id, created_at, nick, text, kind, target, anon from notes where kind = 'link' and text like '%youtube.com%' or text like '%youtu.be%'
`

func (q *Queries) YoutubeLinks(ctx context.Context) ([]Note, error) {
	rows, err := q.db.QueryContext(ctx, youtubeLinks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Note
	for rows.Next() {
		var i Note
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Nick,
			&i.Text,
			&i.Kind,
			&i.Target,
			&i.Anon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
