#!/usr/bin/env bash
# Build a portable Linux amd64 moonBASIC *compiler* tarball: single moonbasic binary, CGO off.
# End users need only this binary to compile, --check, --lsp, --disasm (no Python/C compiler).
# Usage: from repo root: bash scripts/release_compiler_linux.sh
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT="${ROOT}/dist/MoonBasic-compiler-linux-amd64.tar.gz"
mkdir -p "${ROOT}/dist"
STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
DIR="${STAGE}/MoonBasic-compiler"
mkdir -p "$DIR"
export CGO_ENABLED=0
VER="${MOONBASIC_VERSION:-}"
LDFLAGS="-s -w"
if [ -n "$VER" ]; then
  LDFLAGS="$LDFLAGS -X moonbasic/internal/version.Version=${VER}"
fi
go build -trimpath -ldflags="$LDFLAGS" -o "${DIR}/moonbasic" "${ROOT}/cmd/moonbasic"
cat > "${DIR}/README-COMPILER.txt" <<'EOF'
MoonBASIC compiler (Linux amd64)

Toolchain only: .mb -> .mbc, --check, --lsp, --disasm. Built with CGO_ENABLED=0.

  ./moonbasic --version
  ./moonbasic game.mb
  ./moonbasic --check game.mb

Full game runtime (moonrun) is a separate fullruntime build; see docs/BUILDING.md.
EOF
tar -C "$STAGE" -czf "$OUT" MoonBasic-compiler
echo "Created: $OUT"
