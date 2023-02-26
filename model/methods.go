package model

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func JoinedNicks(channel string) ([]ChannelNick, error) {
	channelNicks := []ChannelNick{}
	err := DB.Select(&channelNicks, `select channel, nick, present from channel_nicks where present = true and channel = ?`, channel)
	return channelNicks, err
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
