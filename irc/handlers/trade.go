package handlers

import (
	"goirc/trader"
	"regexp"
)

func trade(params HandlerParams) bool {
	re := regexp.MustCompile("^((buy|sell).*)$")
	matches := re.FindStringSubmatch(params.Msg)

	if len(matches) == 0 {
		return false
	}

	reply, err := trader.Trade(params.Nick, matches[1], params.Db)
	if err != nil {
		params.Irccon.Privmsgf(params.Target, "error: %s", err)
		return true
	}

	if reply != "" {
		params.Irccon.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
		return true
	}

	return false
}
