package handler

import (
	"fmt"
	"net/http"
	"testing"

	Test "github.com/coderaiser/go-subscriber/internal/tape"
)

func TestWriteJSONEncodeError(t *testing.T) {
	Test.Test(t, "handler: writeJSON handles encode error", func(t *Test.T) {
		w := &failResponseWriter{}

		writeJSON(w, http.StatusOK, func() {})

		t.Ok(w.headerWritten)
		t.End()
	})
}

type failResponseWriter struct {
	headerWritten bool
}

func (w *failResponseWriter) Header() http.Header {
	return http.Header{}
}

func (w *failResponseWriter) WriteHeader(status int) {
	w.headerWritten = true
}

func (w *failResponseWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("forced write failure")
}
