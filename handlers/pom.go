package handlers

import (
	"goirc/bot"
	"os/exec"
	"strings"
)

func POM(params bot.HandlerParams) bool {
	cmd := exec.Command("/usr/games/pom")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		params.Privmsgf(params.Target, "%s", err)
	}

	params.Privmsgf(params.Target, "%s", out.String())

	return true
}
