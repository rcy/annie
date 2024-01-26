package handlers

import (
	"goirc/bot"
	"os/exec"
	"strings"
)

func POM(params bot.HandlerParams) (string, error) {
	cmd := exec.Command("/usr/games/pom")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return out.String(), nil
}
