package handlers

import (
	"goirc/irc"
	"goirc/trader"
	"regexp"
)

func Report(params irc.Params) bool {
	re := regexp.MustCompile("^((report).*)$")
	matches := re.FindStringSubmatch(params.Msg)

	if len(matches) == 0 {
		return false
	}

	reply, err := trader.Report(params.Nick, params.Db)
	if err != nil {
		params.Privmsgf(params.Target, "error: %s", err)
		return true
	}

	if reply != "" {
		params.Privmsgf(params.Target, "%s: %s", params.Nick, reply)
		return true
	}

	return false
}
