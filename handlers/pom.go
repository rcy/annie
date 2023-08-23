package handlers

import (
	"goirc/bot"

	"github.com/liamg/gomoon"
)

func POM(params bot.HandlerParams) bool {
	params.Privmsgf(params.Target, "%s", desc(gomoon.PhaseNow()))
	return true
}

func desc(phase gomoon.MoonPhase) string {
	switch phase {
	case gomoon.FULL_MOON:
		return "The moon is full"
	default:
		return "Not full"
	}
}
