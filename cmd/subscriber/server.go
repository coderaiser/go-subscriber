package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/flags"
	"github.com/coderaiser/go-subscriber/internal/flags/meta"
	"github.com/coderaiser/go-subscriber/internal/handler"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func run() int {
	meta.Set(VersionLine, Help)

	result := flags.Parse(os.Args[1:])
	if result.ExitCode >= 0 {
		fmt.Print(result.Output)
		return result.ExitCode
	}

	log := newLogger()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return serve(port, log)
}

func serve(port string, log *slog.Logger) int {
	ss := store.NewStateStore()
	fs := store.NewFactsStore()
	eng := engine.New(ss, fs, time.Now, log)

	h := handler.New(eng, ss)

	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(handleHealth))
	mux.Handle("/readyz", http.HandlerFunc(handleReady))
	mux.Handle("/", h)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Info("starting server", "port", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Error("server error", "error", err)
		return 1
	}
	return 0
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// At this point headers have already been sent, so we can't
		// change the status code. Just log the error.
		slog.Error("failed to encode JSON response", "error", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func newLogger() *slog.Logger {
	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		if os.Getenv("PORT") != "" {
			logFormat = "json"
		} else {
			logFormat = "text"
		}
	}

	var h slog.Handler
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	switch logFormat {
	case "json":
		h = slog.NewJSONHandler(os.Stderr, opts)
	default:
		h = slog.NewTextHandler(os.Stderr, opts)
	}
	return slog.New(h)
}

func startTestServer(ctx context.Context) (string, context.CancelFunc) {
	log := slog.New(slog.NewTextHandler(ioDiscard{}, nil))
	st := store.NewStateStore()
	if err := store.Seed("users.json", st); err != nil {
        panic("subscriber: failed to load users.json: " + err.Error())
    }
	fs := store.NewFactsStore()
	eng := engine.New(st, fs, time.Now, log)

	h := handler.New(eng, st)
	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(handleHealth))
	mux.Handle("/readyz", http.HandlerFunc(handleReady))
	mux.Handle("/", h)

	srv := &http.Server{Handler: mux}
	lis, _ := net.Listen("tcp", "127.0.0.1:0")

	ctx, cancel := context.WithCancel(ctx)
	go srv.Serve(lis)
	go func() {
		<-ctx.Done()
		srv.Close()
	}()

	return "http://" + lis.Addr().String(), cancel
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }