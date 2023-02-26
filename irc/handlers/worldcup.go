package handlers

import (
	"math"
	"regexp"
	"time"
)

func worldcup(params Params) bool {
	re := regexp.MustCompile(`world.?cup`)
	match := re.Find([]byte(params.Msg))

	if len(match) == 0 {
		return false
	}

	end, err := time.Parse(time.RFC3339, "2026-06-01T15:00:00Z")
	if err != nil {
		params.Irccon.Privmsgf(params.Target, "error: %v", err)
		return true
	}
	until := time.Until(end)
	params.Irccon.Privmsgf(params.Target, "the world cup will start in %.0f days", math.Round(until.Hours()/24))
	return true
}
