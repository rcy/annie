package shell

import (
	"fmt"
	"os/exec"
	"strings"
)

func Command(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)

	var stdout, stderr strings.Builder

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: stderr=%s", err.Error(), stderr.String())
	}

	return stdout.String(), nil
}
