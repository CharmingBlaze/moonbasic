#!/usr/bin/env bash
# Dry-run both build-tag paths before push (compiler vs full runtime).
# Requires a working C toolchain when CGO_ENABLED=1 (fullruntime targets).
# From repo root: bash scripts/check_builds.sh   or: make check-builds
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
OUT="$ROOT/.check"
rm -rf "$OUT"
mkdir -p "$OUT"
cleanup() { rm -rf "$OUT"; }
trap cleanup EXIT

export CGO_ENABLED="${CGO_ENABLED:-1}"

echo "== cmd/moonbasic (default tags, compiler CLI, !fullruntime) =="
go build -o "$OUT/moonbasic-cli" ./cmd/moonbasic

echo "== root main.go path (default, compiler-only binary) =="
go build -o "$OUT/moonbasic-root" .

echo "== cmd/moonrun (-tags=fullruntime) =="
go build -tags=fullruntime -o "$OUT/moonrun" ./cmd/moonrun

echo "== root main_fullruntime.go path (-tags=fullruntime) =="
go build -tags=fullruntime -o "$OUT/moonbasic-full" .

echo "check_builds: OK (both tag axes compile)"
