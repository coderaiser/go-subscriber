package runner

import "io"

// FakeExecutor is a fake for testing.
type FakeExecutor struct {
	Lines []string
}

func (e *FakeExecutor) Run(args ...string) (io.ReadCloser, error) {
	r, w := io.Pipe()
	go func() {
		for _, line := range e.Lines {
			_, _ = w.Write([]byte(line + "\n"))
		}
		_ = w.Close()
	}()
	return r, nil
}
