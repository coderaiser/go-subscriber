package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/handler"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func run() int {
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
	if eng == nil {
		log.Error("failed to create engine")
		return 1
	}

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

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
	ss := store.NewStateStore()
	fs := store.NewFactsStore()
	eng := engine.New(ss, fs, time.Now, log)

	h := handler.New(eng, ss)
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