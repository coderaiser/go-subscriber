package runner

import (
	"io"
	"os/exec"
)

// Executor runs a command and returns its stdout.
type Executor interface {
	Run(args ...string) (io.ReadCloser, error)
}

// OSExecutor runs real commands via os/exec.
// Command is injectable for testing — defaults to exec.Command.
type OSExecutor struct {
	Command func(name string, args ...string) *exec.Cmd
}

// NewOSExecutor returns an OSExecutor using the real exec.Command.
func NewOSExecutor() *OSExecutor {
	return &OSExecutor{Command: exec.Command}
}

func (e *OSExecutor) Run(args ...string) (io.ReadCloser, error) {
	cmd := e.Command("go", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return stdout, nil
}
