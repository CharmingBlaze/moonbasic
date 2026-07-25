#!/usr/bin/env sh
# Launch moonBASIC IDE from a release extract (Linux / macOS).
set -e
cd "$(dirname "$0")"

chmod +x moonbasic moonrun START-IDE.sh 2>/dev/null || true

# Prefer a macOS .app bundle when present (Wails package).
for app in "moonbasic-ide.app" "moonBASIC IDE.app" "MoonBASIC IDE.app"; do
  if [ -d "./$app" ]; then
    open "./$app"
    exit 0
  fi
done

if [ ! -f ./moonbasic-ide ]; then
  echo "moonbasic-ide not found in this folder." >&2
  echo "Expected release layout: moonbasic-ide, moonbasic, moonrun next to this script." >&2
  exit 1
fi
chmod +x ./moonbasic-ide 2>/dev/null || true
exec ./moonbasic-ide
