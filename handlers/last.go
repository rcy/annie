package handlers

import (
	"goirc/bot"
	"goirc/model"
	"goirc/util"
	"regexp"
)

func Seen(params bot.HandlerParams) bool {
	re := regexp.MustCompile("^\\?(\\S+)")
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) == 0 {
		return false
	}

	nick := string(matches[1])

	var channelNick model.ChannelNick

	model.DB.Get(&channelNick, "select * from channel_nicks where nick = ? and channel = ?", nick, params.Target)

	if channelNick.Nick != "" {
		if channelNick.Present {
			params.Privmsgf(params.Target, "%s joined %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt))
		} else {
			params.Privmsgf(params.Target, "%s left %s ago", channelNick.Nick, util.Since(channelNick.UpdatedAt))
		}
	} else {
		params.Privmsgf(params.Target, "Never seen %s", nick)
	}
	return true
}
