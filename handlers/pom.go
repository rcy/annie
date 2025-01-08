package handlers

import (
	"goirc/bot"
	"os/exec"
	"strings"
)

func POM(params bot.HandlerParams) error {
	cmd := exec.Command("/usr/games/pom")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	str := out.String()
	if strings.Contains(str, "69") {
		str = str + " (nice)"
	}

	params.Privmsgf(params.Target, "%s", str)

	return nil
}
