package handler_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	tape "github.com/coderaiser/go-tape"
	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/handler"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func newHandler(t *testing.T) http.Handler {
	t.Helper()
	ss := store.NewStateStore()
	eng := engine.New(
		ss,
		store.NewFactsStore(),
		time.Now,
		slog.Default(),
	)
	return handler.New(eng, ss)
}

func post(h http.Handler, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func get(h http.Handler, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
	return w
}

func TestSubscribeOK(t *testing.T) {
	tape.Test(t, "handler: POST /subscribe returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestSubscribeDuplicate(t *testing.T) {
	tape.Test(t, "handler: POST /subscribe duplicate returns 409", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		t.Ok(w.Code == http.StatusConflict)
		t.End()
	})
}

func TestSubscribeBadBody(t *testing.T) {
	tape.Test(t, "handler: POST /subscribe bad JSON returns 400", func(t *tape.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestUnsubscribeOK(t *testing.T) {
	tape.Test(t, "handler: POST /unsubscribe returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/unsubscribe", map[string]any{"msisdn": "111", "service_id": "svc1"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestChargeResultOK(t *testing.T) {
	tape.Test(t, "handler: POST /charge-result returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/charge-result", map[string]any{
			"msisdn": "111", "service_id": "svc1", "result": engine.ResultSuccess,
		})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStateOK(t *testing.T) {
	tape.Test(t, "handler: GET /state/{msisdn} returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := get(h, "/state/111")
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestRenewOK(t *testing.T) {
	tape.Test(t, "handler: POST /renew returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/renew", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestExpireTrialOK(t *testing.T) {
	tape.Test(t, "handler: POST /expire-trial returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": true})
		w := post(h, "/expire-trial", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestRetryOK(t *testing.T) {
	tape.Test(t, "handler: POST /retry returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		post(h, "/charge-result", map[string]any{
			"msisdn": "111", "service_id": "svc1", "result": engine.ResultPermanent,
		})
		w := post(h, "/retry", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestKickOutOK(t *testing.T) {
	tape.Test(t, "handler: POST /kick-out returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		post(h, "/charge-result", map[string]any{
			"msisdn": "111", "service_id": "svc1", "result": engine.ResultPermanent,
		})
		w := post(h, "/kick-out", map[string]any{"msisdn": "111", "service_id": "svc1"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubSendMT(t *testing.T) {
	tape.Test(t, "handler: POST /send-mt stub returns 200", func(t *tape.T) {
		h := newHandler(t.TB())
		w := post(h, "/send-mt", map[string]any{"msisdn": "111"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}
