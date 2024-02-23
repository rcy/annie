package handlers

import (
	"goirc/bot"
	"testing"
)

func TestLink(t *testing.T) {
	err := Link(bot.HandlerParams{
		Matches: []string{"", "https://www.youtube.com/watch?v=8nrYCOZZ9sg"},
		Target:  "#channel",
		Nick:    "theguy",
	})

	if err != nil {
		panic(err)
	}
}
