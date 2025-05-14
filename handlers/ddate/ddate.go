package ddate

import (
	"goirc/bot"
	"os"
	"os/exec"
	"strings"
)

func Handle(params bot.HandlerParams) error {
	cmd := exec.Command("ddate")
	cmd.Env = append(os.Environ(), "TZ=America/Creston")

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return err
	}

	str := out.String()

	params.Privmsgf(params.Target, "%s", str)

	return nil
}
