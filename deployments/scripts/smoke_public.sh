#!/usr/bin/env bash
set -euo pipefail

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this smoke script" >&2
  exit 1
fi

API_BASE="${API_BASE:-}"
TENANT_NAME="${TENANT_NAME:-Acme Internet}"
TENANT_SLUG="${TENANT_SLUG:-acme-internet}"
ADMIN_EMAIL="${ADMIN_EMAIL:-admin@acme-internet.local}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-changeme123456}"
KEY_NAME="${KEY_NAME:-smoke-key}"
WHOAMI_ROUTE="${WHOAMI_ROUTE:-/api/consumer/whoami}"

if [[ -z "${API_BASE}" ]]; then
  echo "Set API_BASE, e.g. https://your-backend.onrender.com" >&2
  exit 1
fi

json_post() {
  local path="$1"
  local body="$2"
  shift 2
  curl -sS -X POST "${API_BASE}${path}" -H "Content-Type: application/json" "$@" -d "${body}"
}

echo "[1/8] health"
curl -fsS "${API_BASE}/health" >/dev/null

echo "[2/8] readiness"
curl -fsS "${API_BASE}/readyz" >/dev/null

echo "[3/8] register tenant/admin (idempotent best-effort)"
json_post "/api/admin/tenants/register" "{\"tenant_name\":\"${TENANT_NAME}\",\"tenant_slug\":\"${TENANT_SLUG}\",\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}" >/dev/null || true

echo "[4/8] login"
TOKEN="$(json_post "/api/admin/login" "{\"email\":\"${ADMIN_EMAIL}\",\"password\":\"${ADMIN_PASSWORD}\"}" | jq -r '.token')"
if [[ -z "${TOKEN}" || "${TOKEN}" == "null" ]]; then
  echo "failed to get JWT token" >&2
  exit 1
fi

echo "[5/8] create API key"
CREATE_KEY_RESP="$(json_post "/api/admin/api-keys" "{\"name\":\"${KEY_NAME}\"}" -H "Authorization: Bearer ${TOKEN}")"
API_KEY="$(echo "${CREATE_KEY_RESP}" | jq -r '.api_key')"
KEY_ID="$(echo "${CREATE_KEY_RESP}" | jq -r '.id')"
if [[ -z "${API_KEY}" || "${API_KEY}" == "null" ]]; then
  echo "failed to create API key" >&2
  exit 1
fi

echo "[6/8] consumer whoami"
curl -fsS "${API_BASE}${WHOAMI_ROUTE}" -H "X-API-Key: ${API_KEY}" >/dev/null

echo "[7/8] trigger rate-limit sample (expect at least one 429 eventually)"
RATE_HIT=0
for _ in {1..120}; do
  CODE="$(curl -sS -o /dev/null -w '%{http_code}' "${API_BASE}${WHOAMI_ROUTE}" -H "X-API-Key: ${API_KEY}")"
  if [[ "${CODE}" == "429" ]]; then
    RATE_HIT=1
    break
  fi
done
if [[ "${RATE_HIT}" -ne 1 ]]; then
  echo "warning: did not observe 429 in sample loop" >&2
fi

echo "[8/8] revoke key and confirm denied"
curl -fsS -X POST "${API_BASE}/api/admin/api-keys/${KEY_ID}/revoke" -H "Authorization: Bearer ${TOKEN}" >/dev/null
REVOKED_CODE="$(curl -sS -o /dev/null -w '%{http_code}' "${API_BASE}${WHOAMI_ROUTE}" -H "X-API-Key: ${API_KEY}")"
if [[ "${REVOKED_CODE}" != "401" ]]; then
  echo "expected 401 after revoke, got ${REVOKED_CODE}" >&2
  exit 1
fi

echo "smoke test completed"
