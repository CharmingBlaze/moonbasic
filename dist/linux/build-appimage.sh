#!/usr/bin/env bash
set -euo pipefail
# Minimal AppImage layout — requires appimagetool on PATH.
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
APP="$ROOT/dist/AppDir"
rm -rf "$APP"
mkdir -p "$APP/usr/bin" "$APP/usr/share/moonbasic"
cp "$ROOT/dist/stage/bin/moonbasic" "$APP/usr/bin/moonbasic"
cp -r "$ROOT/examples" "$APP/usr/share/moonbasic/" 2>/dev/null || true
cp -r "$ROOT/assets" "$APP/usr/share/moonbasic/" 2>/dev/null || true
cat >"$APP/AppRun" <<'EOF'
#!/bin/sh
HERE="$(dirname "$(readlink -f "$0")")"
exec "$HERE/usr/bin/moonbasic" "$@"
EOF
chmod +x "$APP/AppRun"
cat >"$APP/moonbasic.desktop" <<EOF
[Desktop Entry]
Type=Application
Name=moonBASIC
Exec=moonbasic
Categories=Development;
EOF
appimagetool "$APP" moonbasic-x86_64.AppImage
