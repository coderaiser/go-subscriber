package runner

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/coderaiser/go-subscriber/model"
)

func TestRunnerReturnsEventsFromFakeExecutor(t *testing.T) {
	lines := []string{
		`{"Action":"run","Package":"mypkg","Test":"TestFoo"}`,
		`{"Action":"pass","Package":"mypkg","Test":"TestFoo","Elapsed":0.1}`,
	}
	executor := &FakeExecutor{Lines: lines}
	r := New(executor)

	ch, err := r.Run("test", "-json")
	if err != nil {
		t.Fatal(err)
	}

	var events []model.Event
	for e := range ch {
		events = append(events, e)
	}

	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Action != "run" {
		t.Errorf("expected run, got %s", events[0].Action)
	}
	if events[1].Action != "pass" {
		t.Errorf("expected pass, got %s", events[1].Action)
	}
}

func TestRunnerSkipsInvalidJSON(t *testing.T) {
	lines := []string{
		`not json`,
	}
	executor := &FakeExecutor{Lines: lines}
	r := New(executor)

	ch, err := r.Run("test", "-json")
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	for range ch {
		count++
	}

	if count != 0 {
		t.Errorf("expected 0 events, got %d", count)
	}
}

func TestFakeExecutorReturnsLines(t *testing.T) {
	e := &FakeExecutor{Lines: []string{"line1", "line2"}}
	rc, _ := e.Run()
	buf := make([]byte, 100)
	n, _ := rc.Read(buf)
	got := string(buf[:n])
	if got != "line1\n" {
		t.Errorf("expected line1, got %s", got)
	}
}

func TestOSExecutorRun(t *testing.T) {
	// inject echo as the command — no real go test spawned
	e := &OSExecutor{
		Command: func(name string, args ...string) *exec.Cmd {
			return exec.Command("echo", `{"Action":"run","Test":"TestFoo"}`)
		},
	}
	rc, err := e.Run("test", "-json")
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	buf := make([]byte, 256)
	n, _ := rc.Read(buf)
	if n == 0 {
		t.Fatal("expected output")
	}
}

func TestRunnerClosesChannelOnEmptyInput(t *testing.T) {
	executor := &FakeExecutor{Lines: []string{}}
	r := New(executor)
	ch, err := r.Run("test", "-json")
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for range ch {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 events, got %d", count)
	}
}

func TestOSExecutorStartError(t *testing.T) {
	e := &OSExecutor{
		Command: func(name string, args ...string) *exec.Cmd {
			// nonexistent command — Start() will fail
			return exec.Command("nonexistent-command-xyz")
		},
	}
	_, err := e.Run("test")
	if err == nil {
		t.Fatal("expected error from failed Start()")
	}
}

func TestOSExecutorConformsToInterface(t *testing.T) {
	var e Executor = NewOSExecutor()
	_ = e
}

func TestOSExecutorStdoutPipeError(t *testing.T) {
	e := &OSExecutor{
		Command: func(name string, args ...string) *exec.Cmd {
			cmd := exec.Command("echo", "hello")
			// pre-assign stdout — StdoutPipe will fail
			cmd.Stdout = os.Stdout
			return cmd
		},
	}
	_, err := e.Run("test")
	if err == nil {
		t.Fatal("expected error when StdoutPipe fails")
	}
}

// errReader returns an error after the first read
type errReader struct{ done bool }

func (r *errReader) Read(p []byte) (int, error) {
	if r.done {
		return 0, errors.New("read error")
	}
	r.done = true
	copy(p, []byte("not json\n"))
	return 9, nil
}

func (r *errReader) Close() error { return nil }

// errExecutor returns an errReader
type errExecutor struct{}

func (e errExecutor) Run(args ...string) (io.ReadCloser, error) {
	return &errReader{}, nil
}

func TestRunnerHandlesScannerError(t *testing.T) {
	r := New(errExecutor{})
	ch, err := r.Run("test", "-json")
	if err != nil {
		t.Fatal(err)
	}
	// drain channel — scanner error causes goroutine to exit
	for range ch {
	}
}

// errExecutorRunError always fails
type errExecutorRunError struct{}

func (e errExecutorRunError) Run(args ...string) (io.ReadCloser, error) {
	return nil, errors.New("executor failed")
}

func TestRunnerExecutorError(t *testing.T) {
	r := New(errExecutorRunError{})
	_, err := r.Run("test")
	if err == nil {
		t.Fatal("expected executor error")
	}
}
