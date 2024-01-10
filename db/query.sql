-- name: InsertVisit :exec
insert into visits(session, note_id) values(?,?);

-- name: Link :one
select * from notes where id = ? and kind = 'link';
