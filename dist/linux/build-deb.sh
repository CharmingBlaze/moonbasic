#!/usr/bin/env bash
set -euo pipefail
# Usage: from repo root, after filling dist/stage/: ./dist/linux/build-deb.sh VERSION
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
VER="${1:-0.0.0}"
STAGE="$ROOT/dist/stage"
mkdir -p "$STAGE/DEBIAN" "$STAGE/usr/bin" "$STAGE/usr/share/moonbasic"
cat >"$STAGE/DEBIAN/control" <<EOF
Package: moonbasic
Version: $VER
Section: devel
Priority: optional
Architecture: amd64
Maintainer: moonBASIC <https://github.com>
Description: moonBASIC compiler and runtime
Depends: libgl1, libx11-6
EOF
cp "$STAGE/bin/moonbasic" "$STAGE/usr/bin/moonbasic" 2>/dev/null || { echo "copy moonbasic to dist/stage/bin/moonbasic first"; exit 1; }
cp -r "$ROOT/examples" "$STAGE/usr/share/moonbasic/" 2>/dev/null || true
cp -r "$ROOT/assets" "$STAGE/usr/share/moonbasic/" 2>/dev/null || true
dpkg-deb --build "$STAGE" "moonbasic_${VER}_amd64.deb"
echo "wrote moonbasic_${VER}_amd64.deb"
