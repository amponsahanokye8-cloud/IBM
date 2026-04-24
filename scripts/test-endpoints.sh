#!/usr/bin/env bash
set -euo pipefail

failures=0

check() {
  local name="$1"
  local cmd="$2"
  echo "== $name =="
  if bash -c "$cmd"; then
    echo "[OK] $name"
  else
    echo "[FAIL] $name"
    failures=$((failures + 1))
  fi
  echo
}

check "Auth health" "curl -fsS http://localhost:3001/health"
check "Trip health" "curl -fsS http://localhost:3002/health"
check "Auth request OTP" "curl -fsS -X POST http://localhost:3001/v1/auth/request-otp -H 'Content-Type: application/json' -d '{\"phone\":\"+233200000000\"}'"
check "Trip request" "curl -fsS -X POST http://localhost:3002/v1/trips/request -H 'Content-Type: application/json' -d '{\"rider_id\":\"rider_1\",\"pickup\":{\"lat\":5.6037,\"lng\":-0.1870},\"dropoff\":{\"lat\":5.5600,\"lng\":-0.2050}}'"

if [ "$failures" -gt 0 ]; then
  echo "$failures endpoint check(s) failed. Start services first, then rerun this script."
  exit 1
fi

echo "All endpoint checks passed."
