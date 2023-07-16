package handlers

import (
	"goirc/bot"
	"regexp"

	"github.com/liamg/gomoon"
)

func MatchPOM(params bot.HandlerParams) bool {
	re := regexp.MustCompile(`^!pom`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	return POM(params)
}

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
