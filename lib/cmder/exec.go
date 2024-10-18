package cmder

import (
	"fmt"
	"os/exec"
)

func ExecuteCommand(program string, args ...string) (string, error) {
	cmd := exec.Command(program, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("running %s: %w", program, err)
	}
	return string(out), nil
}
