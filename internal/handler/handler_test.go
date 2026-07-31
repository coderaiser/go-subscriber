package handler_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/handler"
	"github.com/coderaiser/go-subscriber/internal/store"
	Test "github.com/coderaiser/go-subscriber/internal/tape"
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
	Test.Test(t, "handler: POST /subscribe returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestSubscribeDuplicate(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe duplicate returns 409", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		t.Ok(w.Code == http.StatusConflict)
		t.End()
	})
}

func TestSubscribeBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/subscribe", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestUnsubscribeOK(t *testing.T) {
	Test.Test(t, "handler: POST /unsubscribe returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/unsubscribe", map[string]any{"msisdn": "111", "service_id": "svc1"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestChargeResultOK(t *testing.T) {
	Test.Test(t, "handler: POST /charge-result returns 200", func(t *Test.T) {
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
	Test.Test(t, "handler: GET /state/{msisdn} returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := get(h, "/state/111")
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

// ── /subscribe ─────────────────────────────────────────────────────────────

func TestSubscribeMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /subscribe returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/subscribe")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestSubscribeCooloff(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe during cooloff returns 409", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		post(h, "/unsubscribe", map[string]any{"msisdn": "111", "service_id": "svc1"})
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		t.Ok(w.Code == http.StatusConflict)
		t.End()
	})
}

func TestRenewOK(t *testing.T) {
	Test.Test(t, "handler: POST /renew returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/renew", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestExpireTrialOK(t *testing.T) {
	Test.Test(t, "handler: POST /expire-trial returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": true})
		w := post(h, "/expire-trial", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestRetryOK(t *testing.T) {
	Test.Test(t, "handler: POST /retry returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		post(h, "/charge-result", map[string]any{
			"msisdn": "111", "service_id": "svc1", "result": engine.ResultSubscriberState,
		})
		w := post(h, "/retry", map[string]any{"msisdn": "111", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestKickOutOK(t *testing.T) {
	Test.Test(t, "handler: POST /kick-out returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		post(h, "/charge-result", map[string]any{
			"msisdn": "111", "service_id": "svc1", "result": engine.ResultSubscriberState,
		})
		w := post(h, "/kick-out", map[string]any{"msisdn": "111", "service_id": "svc1"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubSendMT(t *testing.T) {
	Test.Test(t, "handler: POST /send-mt stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/send-mt", map[string]any{"msisdn": "111"})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

// ── /unsubscribe ───────────────────────────────────────────────────────────

func TestUnsubscribeMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /unsubscribe returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/unsubscribe")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestUnsubscribeBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /unsubscribe bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/unsubscribe", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

// ── /charge-result ─────────────────────────────────────────────────────────

func TestChargeResultMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /charge-result returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/charge-result")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestChargeResultBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /charge-result bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/charge-result", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestChargeResultUnknownSub(t *testing.T) {
	Test.Test(t, "handler: POST /charge-result unknown sub returns 404", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/charge-result", map[string]any{
			"msisdn": "nobody", "service_id": "svc1", "result": engine.ResultSuccess,
		})
		t.Ok(w.Code == http.StatusNotFound)
		t.End()
	})
}

// ── /renew ─────────────────────────────────────────────────────────────────

func TestRenewMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /renew returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/renew")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestRenewBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /renew bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/renew", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestRenewUnknownSub(t *testing.T) {
	Test.Test(t, "handler: POST /renew unknown sub returns 404", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/renew", map[string]any{"msisdn": "nobody", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusNotFound)
		t.End()
	})
}

// ── /expire-trial ──────────────────────────────────────────────────────────

func TestExpireTrialMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /expire-trial returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/expire-trial")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestExpireTrialBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /expire-trial bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/expire-trial", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestExpireTrialUnknownSub(t *testing.T) {
	Test.Test(t, "handler: POST /expire-trial unknown sub returns 404", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/expire-trial", map[string]any{"msisdn": "nobody", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusNotFound)
		t.End()
	})
}

// ── /retry ─────────────────────────────────────────────────────────────────

func TestRetryMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /retry returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/retry")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestRetryBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /retry bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/retry", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestRetryUnknownSub(t *testing.T) {
	Test.Test(t, "handler: POST /retry unknown sub returns 404", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/retry", map[string]any{"msisdn": "nobody", "service_id": "svc1", "success": true})
		t.Ok(w.Code == http.StatusNotFound)
		t.End()
	})
}

// ── /kick-out ──────────────────────────────────────────────────────────────

func TestKickOutMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: GET /kick-out returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/kick-out")
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestKickOutBadBody(t *testing.T) {
	Test.Test(t, "handler: POST /kick-out bad JSON returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		req := httptest.NewRequest(http.MethodPost, "/kick-out", bytes.NewReader([]byte("bad")))
		w := httptest.NewRecorder()
		h.ServeHTTP(w, req)
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestKickOutUnknownSub(t *testing.T) {
	Test.Test(t, "handler: POST /kick-out unknown sub returns 404", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/kick-out", map[string]any{"msisdn": "nobody", "service_id": "svc1"})
		t.Ok(w.Code == http.StatusNotFound)
		t.End()
	})
}

// ── /state/ ────────────────────────────────────────────────────────────────

func TestStateMethodNotAllowed(t *testing.T) {
	Test.Test(t, "handler: POST /state/ returns 405", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/state/", map[string]any{})
		t.Ok(w.Code == http.StatusMethodNotAllowed)
		t.End()
	})
}

func TestStateMissingMsisdn(t *testing.T) {
	Test.Test(t, "handler: GET /state/ with empty msisdn returns 400", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/state/")
		t.Ok(w.Code == http.StatusBadRequest)
		t.End()
	})
}

func TestStateNoSubscriptions(t *testing.T) {
	Test.Test(t, "handler: GET /state/{msisdn} with no subscriptions returns empty array", func(t *Test.T) {
		h := newHandler(t.TB())
		w := get(h, "/state/nobody")
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestRenewFail(t *testing.T) {
	Test.Test(t, "handler: POST /renew with success=false returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": false})
		w := post(h, "/renew", map[string]any{"msisdn": "111", "service_id": "svc1", "success": false})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

// ── /subscribe charge_result ─────────────────────────────────────────────

func TestSubscribeChargeResultSuccess(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe with charge_result=success returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/subscribe", map[string]any{
			"msisdn": "111", "service_id": "svc1", "trial": false,
			"charge_result": engine.ResultSuccess,
		})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestSubscribeChargeResultLowBalance(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe with charge_result=low_balance returns 200 with charge_failed", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/subscribe", map[string]any{
			"msisdn": "111", "service_id": "svc1", "trial": false,
			"charge_result": engine.ResultLowBalance,
		})
		t.Ok(w.Code == http.StatusOK)
		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		t.Equal(resp["state"], "charge_failed")
		t.End()
	})
}

func TestSubscribeChargeResultPermanent(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe with charge_result=permanent returns 200 with charge_failed", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/subscribe", map[string]any{
			"msisdn": "111", "service_id": "svc1", "trial": false,
			"charge_result": engine.ResultPermanent,
		})
		t.Ok(w.Code == http.StatusOK)
		var resp map[string]string
		json.NewDecoder(w.Body).Decode(&resp)
		t.Equal(resp["state"], "charge_failed")
		t.End()
	})
}

func TestSubscribeTrialAfterUnsubConflict(t *testing.T) {
	Test.Test(t, "handler: POST /subscribe trial after unsub returns 409", func(t *Test.T) {
		h := newHandler(t.TB())
		post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": true})
		post(h, "/unsubscribe", map[string]any{"msisdn": "111", "service_id": "svc1"})
		w := post(h, "/subscribe", map[string]any{"msisdn": "111", "service_id": "svc1", "trial": true})
		t.Ok(w.Code == http.StatusConflict)
		t.End()
	})
}

// ── Group A stubs ───────────────────────────────────────────────────────

func TestStubStatus(t *testing.T) {
	Test.Test(t, "handler: POST /status stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/status", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubServicesHistory(t *testing.T) {
	Test.Test(t, "handler: POST /services-history stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/services-history", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubRequestPin(t *testing.T) {
	Test.Test(t, "handler: POST /request-pin stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/request-pin", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubVerifyPin(t *testing.T) {
	Test.Test(t, "handler: POST /verify-pin stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/verify-pin", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubOneTimePayment(t *testing.T) {
	Test.Test(t, "handler: POST /one-time-payment stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/one-time-payment", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubForwardMO(t *testing.T) {
	Test.Test(t, "handler: POST /forward-mo stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/forward-mo", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubSync(t *testing.T) {
	Test.Test(t, "handler: POST /sync stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/sync", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}

func TestStubOnDeliveryReport(t *testing.T) {
	Test.Test(t, "handler: POST /on-delivery-report stub returns 200", func(t *Test.T) {
		h := newHandler(t.TB())
		w := post(h, "/on-delivery-report", map[string]any{})
		t.Ok(w.Code == http.StatusOK)
		t.End()
	})
}
