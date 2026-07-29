package runner

import (
	"bufio"

	"github.com/coderaiser/go-subscriber/model"
	"github.com/coderaiser/go-subscriber/parser"
)

// Runner runs go test -json and returns parsed events.
type Runner struct {
	executor Executor
}

func New(executor Executor) *Runner {
	return &Runner{executor: executor}
}

// Run starts go test -json and returns a channel of parsed events.
func (r *Runner) Run(args ...string) (<-chan model.Event, error) {
	rc, err := r.executor.Run(args...)
	if err != nil {
		return nil, err
	}

	ch := make(chan model.Event)
	go func() {
		defer close(ch)
		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			line := scanner.Text()
			e, err := parser.Parse(line)
			if err != nil {
				continue
			}
			ch <- e
		}
		if err := scanner.Err(); err != nil {
			return
		}
	}()

	return ch, nil
}
