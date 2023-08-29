package handlers

import (
	"goirc/bot"
	"goirc/model"
)

func DeferredDelivery(params bot.HandlerParams) error {
	if params.Target == params.Nick {
		params.Privmsgf(params.Target, "not your personal secretary")
		return nil
	}

	prefix := params.Matches[1]
	message := params.Matches[2]

	// if the prefix matches a currently joined nick, we do nothing
	if model.PrefixMatchesJoinedNick(model.DB, params.Target, prefix) {
		return nil
	}

	if model.PrefixMatchesKnownNick(model.DB, params.Target, prefix) {
		_, err := model.DB.Exec(`insert into laters values(datetime('now'), ?, ?, ?, ?)`, params.Nick, prefix, message, false)
		if err != nil {
			return err
		}

		params.Privmsgf(params.Target, "%s: will send to %s* later", params.Nick, prefix)
	}
	return nil
}
