#!/usr/bin/env bash
# Invoked by release-windows.ps1. Same build flags as .github/workflows/release.yml (Windows job).
set -euo pipefail
export PATH="/c/Progra~1/Go/bin:/mingw64/bin:/usr/bin:${PATH}"
export CGO_ENABLED=1
export CC=/mingw64/bin/gcc.exe
export CGO_LDFLAGS='-lraylib -lgdi32 -lwinmm -lws2_32'
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"
mkdir -p dist
go build -trimpath -ldflags='-s -w' -o dist/moonbasic.exe .
go build -trimpath -ldflags='-s -w' -tags fullruntime -o dist/moonrun.exe ./cmd/moonrun
cp packaging/README-RELEASE.txt dist/README-RELEASE.txt
ls -la dist/moonbasic.exe dist/moonrun.exe dist/README-RELEASE.txt
