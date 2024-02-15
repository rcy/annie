-- name: InsertVisit :exec
insert into visits(session, note_id) values(?,?);

-- name: Link :one
select * from notes where id = ? and kind = 'link';

-- name: InsertNickWeatherRequest :exec
insert into nick_weather_requests(nick, query, city, country) values(?,?,?,?);

-- name: LastNickWeatherRequest :one
select * from nick_weather_requests where nick = ? order by created_at desc limit 1;

-- name: LastWeatherRequestByPrefix :one
select * from nick_weather_requests where city like ? || '%' order by created_at desc limit 1;

-- name: InsertNote :one
insert into notes(target, nick, kind, text) values(?,?,?,?) returning *;

-- name: LastDaysNotes :many
select created_at, nick, text, kind from notes where created_at > datetime('now', '-1 day') order by created_at asc;

-- name: UnsentAnonymousNotes :many
select * from notes where created_at <= ? and nick = target order by id asc limit 420;

-- name: MarkAnonymousNoteDelivered :one
update notes set target = ? where id = ? returning *;

-- name: YoutubeLinks :many
select * from notes where kind = 'link' and text like '%youtube.com%' or text like '%youtu.be%';

-- name: AllNotes :many
select * from notes where target != nick order by created_at desc limit 10000;

-- name: AllNickNotes :many
select * from notes where target != nick and nick = ? order by created_at desc limit 10000;

-- name: NicksWithNoteCount :many
select nick, count(nick) as count from notes group by nick;

-- name: ChannelNick :one
select * from channel_nicks where present = 0 and channel = ? and nick = ? collate nocase;

-- name: ChannelNotesSince :many
select * from notes where target = ? and created_at > ? order by created_at asc limit 69;
