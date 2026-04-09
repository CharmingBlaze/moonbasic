#!/usr/bin/env bash
# moonBASIC dev helpers (Unix / Git Bash / WSL). Repo root: ./scripts/dev.sh <cmd>
set -euo pipefail
cd "$(dirname "$0")/.."
cmd="${1:-help}"
case "$cmd" in
  build-compiler) go build -o moonbasic . ;;
  build-moonrun)  go build -tags fullruntime -o moonrun ./cmd/moonrun ;;
  test)           go test ./... ;;
  check)          go run . --check examples/mario64/main_entities.mb ;;
  run-spin-cube)  CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb ;;
  help|*)
    echo "Usage: scripts/dev.sh <build-compiler|build-moonrun|test|check|run-spin-cube>"
    ;;
esac
