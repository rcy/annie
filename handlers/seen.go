package handlers

import (
	"fmt"
	"goirc/bot"
	"goirc/model"
	"goirc/util"
)

func Seen(params bot.HandlerParams) (string, error) {
	nick := params.Matches[1]

	var channelNick model.ChannelNick

	err := model.DB.Get(&channelNick, "select * from channel_nicks where nick = ? and channel = ?", nick, params.Target)
	if err != nil {
		return "", err
	}

	if channelNick.Nick != "" {
		if channelNick.Present {
			return fmt.Sprintf("%s joined %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt)), nil
		} else {
			return fmt.Sprintf("%s left %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt)), nil
		}
	}
	return fmt.Sprintf("Never seen %s", nick), nil
}
