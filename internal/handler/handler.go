package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/coderaiser/go-subscriber/internal/engine"
	"github.com/coderaiser/go-subscriber/internal/store"
)

func New(eng *engine.Engine, states *store.StateStore) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn       string `json:"msisdn"`
			ServiceID    string `json:"service_id"`
			Trial        bool   `json:"trial"`
			ChargeResult string `json:"charge_result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		if body.ChargeResult == "" {
			body.ChargeResult = engine.ResultSuccess
		}

		err := eng.Subscribe(body.Msisdn, body.ServiceID, body.Trial, body.ChargeResult)
		if err == engine.ErrAlreadySubscribed || err == engine.ErrCooloff || err == engine.ErrTrialAlreadyUsed {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if err == engine.ErrChargeFailed {
			writeJSON(w, http.StatusOK, map[string]string{"state": "charge_failed"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"state": "ok"})
	})

	mux.HandleFunc("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		state := eng.Unsubscribe(body.Msisdn, body.ServiceID)
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/charge-result", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
			Result    string `json:"result"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		state := eng.OnChargeResult(body.Msisdn, body.ServiceID, body.Result)
		if state == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/renew", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
			Success   bool   `json:"success"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		var result string
		if body.Success {
			result = engine.ResultSuccess
		} else {
			result = engine.ResultPermanent
		}
		state := eng.OnChargeResult(body.Msisdn, body.ServiceID, result)
		if state == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/expire-trial", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
			Success   bool   `json:"success"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		state := eng.ExpireTrial(body.Msisdn, body.ServiceID, body.Success)
		if state == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/retry", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
			Success   bool   `json:"success"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		state := eng.Retry(body.Msisdn, body.ServiceID, body.Success)
		if state == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/kick-out", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		var body struct {
			Msisdn    string `json:"msisdn"`
			ServiceID string `json:"service_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "bad request"})
			return
		}

		state := eng.KickOut(body.Msisdn, body.ServiceID)
		if state == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"state": state})
	})

	mux.HandleFunc("/state/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		msisdn := strings.TrimPrefix(r.URL.Path, "/state/")
		if msisdn == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing msisdn"})
			return
		}

		all := states.All()
		var services []map[string]string
		for k, v := range all {
			if strings.HasPrefix(k, msisdn+":") {
				serviceID := strings.TrimPrefix(k, msisdn+":")
				services = append(services, map[string]string{
					"service_id": serviceID,
					"state":      v,
				})
			}
		}

		resp := map[string]any{
			"msisdn":   msisdn,
			"services": services,
		}
		if services == nil {
			resp["services"] = []map[string]string{}
		}

		writeJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("/send-mt", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "stub": true})
	})

	mux.HandleFunc("/status", statusStubHandler)
	mux.HandleFunc("/services-history", servicesHistoryHandler)
	mux.HandleFunc("/request-pin", stubHandler("sent", true))
	mux.HandleFunc("/verify-pin", stubHandler("verified", true))
	mux.HandleFunc("/one-time-payment", stubHandler("charged", true))
	mux.HandleFunc("/forward-mo", stubHandler("forwarded", true))
	mux.HandleFunc("/sync", stubHandler("synced", true))
	mux.HandleFunc("/on-delivery-report", stubHandler("received", true))

	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Response headers are already written, so we cannot change the status.
		// Log the failure for diagnostics.
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// stubHandler returns an HTTP handler that responds with a static JSON stub.
func stubHandler(key string, value any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{key: value, "stub": true})
	}
}

func statusStubHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"state": "active", "stub": true})
}

func servicesHistoryHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"services": []any{}, "stub": true})
}
