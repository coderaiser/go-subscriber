package debug

import (
	"context"
	"log/slog"
	"os"
)

// Logger returns a *slog.Logger for the given namespace.
// Output is controlled by two env vars:
//
//   DEBUG=subscriber:*           → all namespaces log to stderr
//   DEBUG=subscriber:engine      → only that namespace logs
//   DEBUG=*                      → everything logs
//   LOG_FORMAT=json              → JSON output (default text)
func Logger(namespace string) *slog.Logger {
	if !enabled(namespace) {
		return slog.New(noopHandler{})
	}

	opts := &slog.HandlerOptions{}
	var h slog.Handler

	if os.Getenv("LOG_FORMAT") == "json" {
		h = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		h = slog.NewTextHandler(os.Stderr, opts)
	}

	return slog.New(h).With("ns", namespace)
}

func enabled(namespace string) bool {
	v := os.Getenv("DEBUG")
	return v == "*" || v == "subscriber:*" || v == namespace
}

// noopHandler discards all log records with zero allocations.
type noopHandler struct{}

func (noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (noopHandler) Handle(context.Context, slog.Record) error  { return nil }
func (noopHandler) WithAttrs([]slog.Attr) slog.Handler         { return noopHandler{} }
func (noopHandler) WithGroup(string) slog.Handler              { return noopHandler{} }
