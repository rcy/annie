package ddate

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func NowInZone(tz string) (string, error) {
	cmd := exec.Command("ddate")
	cmd.Env = append(os.Environ(), "TZ="+tz)

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func On(day time.Time) (string, error) {
	cmd := exec.Command("ddate", strconv.Itoa(day.Day()), strconv.Itoa(int(day.Month())), strconv.Itoa(day.Year()))

	var out strings.Builder
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}
