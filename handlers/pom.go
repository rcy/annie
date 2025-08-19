package handlers

import (
	"goirc/internal/responder"
	"os/exec"
	"strings"
)

func POM(params responder.Responder) error {
	cmd := exec.Command("/usr/games/pom")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	str := strings.TrimSpace(out.String())
	if strings.Contains(str, "69") {
		str = str + " (nice)"
	}

	params.Privmsgf(params.Target(), "%s", str)

	return nil
}
