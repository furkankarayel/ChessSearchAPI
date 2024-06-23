package helper

import (
	"os/exec"
)

// Runner struct to hold any state or configuration
type Runner struct {
	Executable string
}

// NewRunner creates a new instance of Runner
func NewRunner(executable string) *Runner {
	return &Runner{Executable: executable}
}

// Run executes the command with given arguments
func (cw *Runner) Run(args ...string) (string, error) {
	cmd := exec.Command(cw.Executable, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// Example function to run echo with a message
func (cw *Runner) Echo(message string) (string, error) {
	return cw.Run(message)
}

// Example function to run echo with additional flags
func (cw *Runner) EchoWithFlags(flags ...string) (string, error) {
	return cw.Run(flags...)
}
