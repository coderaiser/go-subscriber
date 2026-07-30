package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestHandleHealth(t *testing.T) {
	Test.Test(t, "server: GET /healthz returns 200", func(t *Test.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		w := httptest.NewRecorder()
		handleHealth(w, req)
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestHandleReady(t *testing.T) {
	Test.Test(t, "server: GET /readyz returns 200", func(t *Test.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		w := httptest.NewRecorder()
		handleReady(w, req)
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestNewLoggerJSON(t *testing.T) {
	Test.Test(t, "server: newLogger returns json logger when LOG_FORMAT=json", func(t *Test.T) {
		t.TB().Setenv("LOG_FORMAT", "json")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerText(t *testing.T) {
	Test.Test(t, "server: newLogger returns text logger when LOG_FORMAT=text", func(t *Test.T) {
		t.TB().Setenv("LOG_FORMAT", "text")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerDefaultDev(t *testing.T) {
	Test.Test(t, "server: newLogger defaults to text when no env set", func(t *Test.T) {
		os.Unsetenv("LOG_FORMAT")
		os.Unsetenv("PORT")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestNewLoggerDefaultProd(t *testing.T) {
	Test.Test(t, "server: newLogger defaults to json when PORT is set", func(t *Test.T) {
		os.Unsetenv("LOG_FORMAT")
		t.TB().Setenv("PORT", "9999")
		t.Ok(newLogger() != nil)
		t.End()
	})
}

func TestServeBadPort(t *testing.T) {
	Test.Test(t, "server: serve returns 1 for invalid port", func(t *Test.T) {
		code := serve("99999", newLogger())
		t.Ok(code == 1)
		t.End()
	})
}

func TestIoDiscardWrite(t *testing.T) {
	Test.Test(t, "server: ioDiscard.Write returns len and no error", func(t *Test.T) {
		d := ioDiscard{}
		n, err := d.Write([]byte("hello"))
		t.Ok(n == 5 && err == nil)
		t.End()
	})
}

func TestServerHealthzIntegration(t *testing.T) {
	Test.Test(t, "server: integration GET /healthz returns ok", func(t *Test.T) {
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

func TestRunVersionFlag(t *testing.T) {
	Test.Test(t, "server: run with --version returns 0", func(t *Test.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"subscriber", "--version"}
		code := run()
		t.Ok(code == 0)
		t.End()
	})
}

func TestRunHelpFlag(t *testing.T) {
	Test.Test(t, "server: run with --help returns 0", func(t *Test.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"subscriber", "--help"}
		code := run()
		t.Ok(code == 0)
		t.End()
	})
}

func TestRunUnknownFlag(t *testing.T) {
	Test.Test(t, "server: run with unknown flag returns 1", func(t *Test.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"subscriber", "--bad"}
		code := run()
		t.Ok(code == 1)
		t.End()
	})
}

func TestRunBadPort(t *testing.T) {
	Test.Test(t, "server: run with bad port returns 1", func(t *Test.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"subscriber"}
		t.TB().Setenv("PORT", "99999")
		code := run()
		t.Ok(code == 1)
		t.End()
	})
}
