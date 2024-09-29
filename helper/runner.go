package helper

import (
	"log"
	"os/exec"
)

type Runner struct {
	Executable string
}

func NewRunner(executable string) *Runner {
	return &Runner{Executable: executable}
}

func (cw *Runner) Run(args ...string) ([]byte, error) {
	cmd := exec.Command(cw.Executable, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Command execution failed: %v\nOutput: %s", err, output)
		return output, err
	}
	return output, nil
}
