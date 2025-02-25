package handlers

import (
	"context"
	"database/sql"
	"errors"
	"goirc/bot"
	"goirc/model"
	"goirc/util"
)

func Seen(ctx context.Context, params bot.HandlerParams) error {
	nick := params.Matches[1]

	var channelNick model.ChannelNick

	err := model.DB.Get(&channelNick, "select * from channel_nicks where nick = ? and channel = ?", nick, params.Target)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	if channelNick.Nick != "" {
		if channelNick.Present {
			params.Privmsgf(params.Target, "%s joined %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt))
		} else {
			params.Privmsgf(params.Target, "%s left %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt))
		}
	} else {
		params.Privmsgf(params.Target, "Never seen %s", nick)
	}
	return nil
}
