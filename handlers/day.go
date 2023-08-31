package handlers

import (
	"fmt"
	"goirc/bot"
	"os/exec"
	"strings"
)

func NationalDay(params bot.HandlerParams) error {
	r, err := shell(`curl -s https://nationaltoday.com/ | pup '.holiday-title-text json{}' | jq -r .[0].text`)
	if err != nil {
		params.Privmsgf(params.Target, "error: %v", err)
	}

	params.Privmsgf(params.Target, "%s", r)

	return nil
}

func shell(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)

	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s", stderr.String())
	}

	return stdout.String(), nil
}
