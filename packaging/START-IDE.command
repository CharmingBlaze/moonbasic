#!/bin/bash
# Double-click in Finder (macOS) to start the moonBASIC IDE.
cd "$(dirname "$0")"
chmod +x moonbasic moonrun START-IDE.sh START-IDE.command 2>/dev/null || true
for app in "moonbasic-ide.app" "moonBASIC IDE.app" "MoonBASIC IDE.app"; do
  if [ -d "./$app" ]; then
    open "./$app"
    exit 0
  fi
done
if [ ! -f ./moonbasic-ide ]; then
  osascript -e 'display alert "moonbasic-ide not found" message "Keep this script in the same folder as moonbasic-ide, moonbasic, and moonrun."' 2>/dev/null || \
    echo "moonbasic-ide not found in this folder." >&2
  exit 1
fi
chmod +x ./moonbasic-ide 2>/dev/null || true
# Clear quarantine on first launch when possible (Gatekeeper).
xattr -cr ./moonbasic-ide ./moonbasic ./moonrun 2>/dev/null || true
exec ./moonbasic-ide
