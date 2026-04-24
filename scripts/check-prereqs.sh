#!/usr/bin/env bash
set -euo pipefail

check() {
  local cmd="$1"
  if command -v "$cmd" >/dev/null 2>&1; then
    echo "[OK] $cmd"
  else
    echo "[MISSING] $cmd"
  fi
}

echo "Checking required tools..."
check docker
check node
check npm
check go
check curl
