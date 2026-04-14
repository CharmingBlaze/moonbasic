# Experimental: static-linked moonrun (CGO + Zig). Intended to produce a single .exe without
# Raylib or Jolt DLLs when the toolchain and Jolt static libraries are available.
#
# Prerequisites:
#   - Go, Zig on PATH (unless you override CC/CXX to MinGW gcc/g++).
#   - CGO_ENABLED=1 (this script sets it).
#   - Raylib is compiled from source via vendored raylib-go CGO (no separate raylib.dll).
#   - Native Jolt on Windows: build static libs first — see
#     third_party/jolt-go/jolt/lib/windows_amd64/README.md
#     (third_party/jolt-go/scripts/build-libs-windows.ps1 with JPH_SRC set).
#   - Purego / CGO_ENABLED=0 builds are out of scope here (they load raylib.dll at runtime).
#
# Usage (PowerShell, from repo root):
#   .\scripts\build_static.ps1
#
# Optional overrides:
#   $env:CC = "zig cc -target x86_64-windows-gnu -static"
#   $env:CXX = "zig c++ -target x86_64-windows-gnu -static"
#   $env:MOONBASIC_SKIP_STATIC_EXTLDFLAGS = "1"   # omit -ldflags -extldflags (if link fails)
#   $env:OUTPUT = "moonrun_static.exe"

$ErrorActionPreference = "Stop"

$ZigCC = "zig cc -target x86_64-windows-gnu -static"
$ZigCXX = "zig c++ -target x86_64-windows-gnu -static"

if (-not $env:CC) {
    if (-not (Get-Command zig -ErrorAction SilentlyContinue)) {
        Write-Error "zig not found on PATH. Install Zig from https://ziglang.org/ or set CC/CXX to your MinGW toolchain (e.g. gcc/g++ from MSYS2)."
    }
    $env:CC = $ZigCC
}
if (-not $env:CXX) {
    if (Get-Command zig -ErrorAction SilentlyContinue) {
        $env:CXX = $ZigCXX
    }
}

$env:CGO_ENABLED = "1"

$out = if ($env:OUTPUT) { $env:OUTPUT } else { "moonrun_static.exe" }

Write-Host "CC=$($env:CC)"
Write-Host "CXX=$($env:CXX)"
Write-Host "CGO_ENABLED=$($env:CGO_ENABLED)"
Write-Host "Building cmd/moonrun with -tags fullruntime -> $out ..."

# Use -ldflags=-... as a single flag value so PowerShell/go do not split on spaces (see `go help build`).
$goArgs = @(
    "build",
    "-tags", "fullruntime",
    "-o", $out,
    "./cmd/moonrun"
)
if (-not $env:MOONBASIC_SKIP_STATIC_EXTLDFLAGS) {
    # CGO final link: static libgcc/libstdc++ where possible (in addition to zig -static on CC/CXX).
    $goArgs = @(
        "build",
        "-tags", "fullruntime",
        "-ldflags=-linkmode external -extldflags=-static",
        "-o", $out,
        "./cmd/moonrun"
    )
}

& go @goArgs
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}

Write-Host "Output: $out"

$dumpbin = Get-Command dumpbin -ErrorAction SilentlyContinue
if ($dumpbin) {
    Write-Host "--- dumpbin /dependents (verify no unexpected DLLs) ---"
    & dumpbin /dependents $out
} else {
    Write-Host "Tip: In a VS Developer shell, run: dumpbin /dependents $out"
    Write-Host "     Non-system DLLs (e.g. raylib.dll, jolt.dll) should be absent for a static build."
}
