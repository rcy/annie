package ddate

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func NowInZone(tz string) (string, error) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return "", fmt.Errorf("LoadLocation: %w", err)
	}
	now := time.Now().In(location)
	return On(now)
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
