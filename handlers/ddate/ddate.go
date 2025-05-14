package ddate

import (
	"goirc/bot"
	"os/exec"
	"strings"
)

func Handle(params bot.HandlerParams) error {
	cmd := exec.Command("ddate")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	str := out.String()

	params.Privmsgf(params.Target, "%s", str)

	return nil
}
