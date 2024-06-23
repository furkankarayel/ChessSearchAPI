package helper

import (
	"os/exec"
)

// Runner struct to hold the path to executable
type Runner struct {
	Executable string
}

// NewRunner creates a new instance of Runner
func NewRunner(executable string) *Runner {
	return &Runner{Executable: executable}
}

// Run executes the command with given arguments/flags
func (cw *Runner) Run(args ...string) (string, error) {
	cmd := exec.Command(cw.Executable, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
