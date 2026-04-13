# Experimental: build moonrun with CGO and a Zig-provided Windows C compiler (static-friendly).
# Requires: Go, Zig on PATH, and Raylib + deps available to the C compiler (e.g. MSYS2 MinGW + raylib).
# Usage (PowerShell, from repo root):
#   .\scripts\build_static.ps1
#
# Override targets or flags as needed:
#   $env:CC = "zig cc -target x86_64-windows-gnu"
#   go build -tags fullruntime -o moonrun_static.exe ./cmd/moonrun

$ErrorActionPreference = "Stop"

if (-not $env:CC) {
    $env:CC = "zig cc -target x86_64-windows-gnu"
}

$env:CGO_ENABLED = "1"

Write-Host "CC=$($env:CC)"
Write-Host "CGO_ENABLED=$($env:CGO_ENABLED)"
Write-Host "Building cmd/moonrun with -tags fullruntime ..."

go build -tags fullruntime -o moonrun_static.exe ./cmd/moonrun

Write-Host "Output: moonrun_static.exe"
