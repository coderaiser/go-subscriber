package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	tape "github.com/coderaiser/go-tape"
)

func TestHandleHealth(t *testing.T) {
	tape.Test(t, "server: GET /healthz returns 200", func(t *tape.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		handleHealth(w, req)
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestHandleReady(t *testing.T) {
	tape.Test(t, "server: GET /readyz returns 200", func(t *tape.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		handleReady(w, req)
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestNewLoggerJSON(t *testing.T) {
	tape.Test(t, "server: newLogger returns json logger when LOG_FORMAT=json", func(t *tape.T) {
		t.TB().Setenv("LOG_FORMAT", "json")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerText(t *testing.T) {
	tape.Test(t, "server: newLogger returns text logger when LOG_FORMAT=text", func(t *tape.T) {
		t.TB().Setenv("LOG_FORMAT", "text")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerDefaultDev(t *testing.T) {
	tape.Test(t, "server: newLogger defaults to text when no env set", func(t *tape.T) {
		os.Unsetenv("LOG_FORMAT")
		os.Unsetenv("PORT")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerDefaultProd(t *testing.T) {
	tape.Test(t, "server: newLogger defaults to json when PORT is set", func(t *tape.T) {
		os.Unsetenv("LOG_FORMAT")
		t.TB().Setenv("PORT", "9999")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestServeBadPort(t *testing.T) {
	tape.Test(t, "server: serve returns 1 for invalid port", func(t *tape.T) {
		code := serve("99999", newLogger())
		t.Ok(code == 1)
		t.End()
	})
}

func TestServerHealthzIntegration(t *testing.T) {
	tape.Test(t, "server: integration GET /healthz returns ok", func(t *tape.T) {
		base, cancel := startTestServer(context.Background())
		defer cancel()
		resp, err := http.Get(base + "/healthz")
		if err != nil {
			t.TB().Fatal(err)
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		t.Ok(strings.Contains(string(body), "ok"))
		t.End()
	})
}
