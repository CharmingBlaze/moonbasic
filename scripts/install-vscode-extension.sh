#!/usr/bin/env sh
# Install the moonBASIC VS Code / Cursor extension (one command).
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
if [ -x "$ROOT/moonbasic" ]; then
  exec "$ROOT/moonbasic" install-vscode "$@"
fi
if command -v moonbasic >/dev/null 2>&1; then
  exec moonbasic install-vscode "$@"
fi
cd "$ROOT"
exec go run . install-vscode "$@"
