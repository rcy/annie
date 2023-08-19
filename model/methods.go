package model

import (
	"database/sql"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func JoinedNicks(channel string) ([]ChannelNick, error) {
	channelNicks := []ChannelNick{}
	err := DB.Select(&channelNicks, `select channel, nick, present from channel_nicks where present = true and channel = ?`, channel)
	return channelNicks, err
}

func IsJoined(channel string, nick string) (*ChannelNick, error) {
	channelNick := ChannelNick{}

	query := `
select channel, nick, present
from channel_nicks
where
  present = true and
  channel = ? and
  nick = ?`

	log.Printf("query %s %s", channel, nick)

	err := DB.Get(&channelNick, query, channel, nick)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &channelNick, err
}

func PrefixMatchesJoinedNick(db *sqlx.DB, channel, prefix string) bool {
	channelNicks, err := JoinedNicks(channel)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range channelNicks {
		if strings.HasPrefix(row.Nick, prefix) {
			return true
		}
	}
	return false
}

func PrefixMatchesKnownNick(db *sqlx.DB, channel, prefix string) bool {
	channel_nicks := []ChannelNick{}
	err := db.Select(&channel_nicks, `select channel, nick, present from channel_nicks`)
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range channel_nicks {
		if strings.HasPrefix(row.Nick, prefix) {
			return true
		}
	}
	return false
}
