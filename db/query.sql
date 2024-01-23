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
