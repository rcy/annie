package db

import (
	"log"

	"github.com/BurntSushi/migration"
	"github.com/jmoiron/sqlx"

	_ "modernc.org/sqlite"
)

func Open(dbfile string) *sqlx.DB {
	log.Printf("Opening db: %s", dbfile)

	migrations := []migration.Migrator{
		func(tx migration.LimitedTx) error {
			_, err := tx.Exec(`create table if not exists notes(created_at text, nick text, text text)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			_, err := tx.Exec(`create table if not exists links(created_at text, nick text, text text)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding kind column to notes")
			_, err := tx.Exec(`alter table notes add column kind string not null default "note"`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding laters table")
			_, err := tx.Exec(`create table laters(created_at text, nick text, target text, message text, sent boolean default false)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: adding channel_nicks table")
			_, err := tx.Exec(`create table channel_nicks(channel text not null, nick text not null, present bool not null default false)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add unique constrant to channel_nicks table")

			// delete duplicates, keeping oldest records
			_, err := tx.Exec(`delete from channel_nicks where rowid not in (select min(rowid) from channel_nicks group by nick, channel)`)
			if err != nil {
				return err
			}

			// add unique constraint
			_, err = tx.Exec(`create unique index channel_nick_unique_index on channel_nicks(channel, nick)`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add primary key to notes")
			_, err := tx.Exec(`
pragma foreign_key = off;

alter table notes rename to old_notes;

create table notes(
  id INTEGER not null primary key,
  created_at datetime not null default current_timestamp,
  nick text,
  text text,
  kind string not null default "note"
);

insert into notes select rowid, * from old_notes;

drop table old_notes;

pragma foreign_key = on;
`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add seen table")
			_, err := tx.Exec(`
create table seen_by(
  created_at datetime not null default current_timestamp,
  note_id references notes not null,
  nick text not null
);`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: add updated_at to channel_nicks")
			_, err := tx.Exec(`alter table channel_nicks add column updated_at text`)
			return err
		},
		func(tx migration.LimitedTx) error {
			log.Println("MIGRATE: transactions table")
			_, err := tx.Exec(`
create table transactions(
  created_at datetime not null default current_timestamp,
  nick text not null,
  verb text not null,
  symbol text not null,
  shares number not null,
  price number not null
);`)
			return err
		},
	}

	db, err := migration.Open("sqlite", dbfile, migrations)
	if err != nil {
		log.Fatalf("MIGRATION: %v", err)
	}
	return sqlx.NewDb(db, "sqlite")
}
