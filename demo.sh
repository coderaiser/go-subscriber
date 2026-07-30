#!/usr/bin/env bash
set -euo pipefail

BASE="${BASE:-http://localhost:8080}"

echo "=== subscriber demo ==="
echo ""

# ── scenario 1: trial → expire → low_balance × 3 → retry → unsubscribe ──

echo "--- scenario 1: trial lifecycle ---"

echo "[1] subscribe (trial)"
curl -sf -X POST "$BASE/subscribe" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","trial":true}' | jq .

echo "[2] expire trial (success) → active"
curl -sf -X POST "$BASE/expire-trial" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","success":true}' | jq .

echo "[3] low_balance attempt 1 (ladder not exhausted)"
curl -sf -X POST "$BASE/charge-result" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","result":"low_balance"}' | jq .

echo "[4] low_balance attempt 2 (ladder not exhausted)"
curl -sf -X POST "$BASE/charge-result" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","result":"low_balance"}' | jq .

echo "[5] low_balance attempt 3 → suspended"
curl -sf -X POST "$BASE/charge-result" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","result":"low_balance"}' | jq .

echo "[6] retry (success) → active"
curl -sf -X POST "$BASE/retry" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","success":true}' | jq .

echo "[7] check state"
curl -sf "$BASE/state/96478001" | jq .

echo "[8] unsubscribe → terminated + 30d cooloff"
curl -sf -X POST "$BASE/unsubscribe" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily"}' | jq .

echo ""

# ── scenario 2: re-subscribe during cooloff → 409 ────────────────────────

echo "--- scenario 2: cooloff rejection ---"

echo "[9] subscribe during cooloff → 409"
curl -s -o /dev/null -w "%{http_code}\n" -X POST "$BASE/subscribe" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","service_id":"quiz-daily","trial":false}'

echo ""

# ── scenario 3: paid → permanent → kick_out ──────────────────────────────

echo "--- scenario 3: paid subscribe → kick out ---"

echo "[10] subscribe paid → active"
curl -sf -X POST "$BASE/subscribe" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478002","service_id":"sports-alerts","trial":false}' | jq .

echo "[11] permanent failure → suspended"
curl -sf -X POST "$BASE/charge-result" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478002","service_id":"sports-alerts","result":"permanent"}' | jq .

echo "[12] kick out → removed"
curl -sf -X POST "$BASE/kick-out" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478002","service_id":"sports-alerts"}' | jq .

echo "[13] check final state"
curl -sf "$BASE/state/96478002" | jq .

echo ""

# ── scenario 4: stubs ─────────────────────────────────────────────────────

echo "--- scenario 4: stub endpoints ---"

echo "[14] send-mt stub"
curl -sf -X POST "$BASE/send-mt" \
  -H "Content-Type: application/json" \
  -d '{"msisdn":"96478001","message":"Welcome!"}' | jq .

echo ""
echo "--- done ---"
