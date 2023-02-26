package handlers

import (
	"goirc/bot"
	"goirc/model"
	"log"
	"regexp"
)

func DeferredDelivery(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^([^\s:]+): (.+)$`)
	matches := re.FindSubmatch([]byte(params.Msg))

	if len(matches) > 0 {
		if params.Target == params.Nick {
			params.Privmsgf(params.Target, "not your personal secretary")
			return false
		}

		prefix := string(matches[1])
		message := string(matches[2])

		// if the prefix matches a currently joined nick, we do nothing
		if model.PrefixMatchesJoinedNick(params.Db, params.Target, prefix) {
			return false
		}

		if model.PrefixMatchesKnownNick(params.Db, params.Target, prefix) {
			_, err := params.Db.Exec(`insert into laters values(datetime('now'), ?, ?, ?, ?)`, params.Nick, prefix, message, false)
			if err != nil {
				log.Fatal(err)
			}

			params.Privmsgf(params.Target, "%s: will send to %s* later", params.Nick, prefix)
		}
		return true
	}
	return false
}
