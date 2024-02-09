CREATE TABLE migration_version (
			version INTEGER
		);
CREATE TABLE links(created_at text, nick text, text text);
CREATE TABLE laters(created_at text, nick text, target text, message text, sent boolean default false);
CREATE TABLE channel_nicks(channel text not null, nick text not null, present bool not null default false, updated_at datetime not null);
CREATE UNIQUE INDEX channel_nick_unique_index on channel_nicks(channel, nick);
CREATE TABLE notes(
  id INTEGER not null primary key,
  created_at datetime not null default current_timestamp,
  nick text,
  text text,
  kind text not null default 'note', target text not null default '');
CREATE TABLE reminders(
  id integer not null primary key,
  created_at datetime not null default current_timestamp,
  nick text not null,
  remind_at datetime not null,
  what text not null
);
CREATE TABLE revs(
  id integer not null primary key,
  created_at datetime not null default current_timestamp,
  sha text not null
);
CREATE TABLE visits(
  id integer not null primary key,
  created_at datetime not null default current_timestamp,
  session text not null,
  note_id integer references notes not null
);
CREATE TABLE nick_weather_requests(
  id integer not null primary key,
  created_at datetime not null default current_timestamp,
  nick text not null,
  query text not null,
  city text not null,
  country text not null
);
