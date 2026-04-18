# Build a portable Windows amd64 moonBASIC *compiler* zip: single moonbasic.exe, no CGO, no raylib.dll.
#
# End users only need this exe to: compile .mb -> .mbc, --check, --lsp, --disasm.
# They do NOT need Python, GCC, or MSVC. (Building moonBASIC from source still requires Go.)
#
# Usage (PowerShell, from repo root):
#   .\scripts\release_compiler_windows.ps1
#   .\scripts\release_compiler_windows.ps1 -OutZip .\dist\custom.zip
#   $env:MOONBASIC_VERSION = "v1.2.18"; .\scripts\release_compiler_windows.ps1
#
# Requires: Go toolchain on the *build machine* only.

[CmdletBinding()]
param(
    [string]$RepoRoot = "",
    [string]$OutZip = "",
    [string]$Version = ""
)

$ErrorActionPreference = "Stop"

if (-not $RepoRoot) {
    $RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
}

if (-not $OutZip) {
    $dist = Join-Path $RepoRoot "dist"
    if (-not (Test-Path $dist)) {
        New-Item -ItemType Directory -Path $dist | Out-Null
    }
    $OutZip = Join-Path $dist "MoonBasic-compiler-windows-amd64.zip"
}

$staging = Join-Path $env:TEMP ("moonbasic_compiler_" + [Guid]::NewGuid().ToString("n"))
$bundle = Join-Path $staging "MoonBasic-compiler"
New-Item -ItemType Directory -Path $bundle -Force | Out-Null

try {
    Push-Location $RepoRoot

    $env:CGO_ENABLED = "0"
    $exeOut = Join-Path $bundle "moonbasic.exe"
    $ver = $Version
    if (-not $ver -and $env:MOONBASIC_VERSION) { $ver = $env:MOONBASIC_VERSION }
    $ldflags = "-s -w"
    if ($ver) {
        $ldflags += " -X moonbasic/internal/version.Version=$ver"
    }
    Write-Host "Building cmd/moonbasic with CGO_ENABLED=0 (no native game deps in the binary)..."
    go build -trimpath -ldflags="$ldflags" -o $exeOut ./cmd/moonbasic
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    $readme = @"
MoonBASIC compiler (Windows amd64)

This package is the toolchain only: compile .mb to .mbc, type-check, LSP, disassemble bytecode.
It is built with CGO disabled so you do not need raylib.dll or other native libraries next to this exe.

Commands (examples):
  moonbasic --version
  moonbasic game.mb
  moonbasic --check game.mb
  moonbasic --lsp
  moonbasic --disasm game.mbc

This build does NOT include the full game runtime (--run / moonrun). For running graphical games,
use a fullruntime moonrun build (see docs/BUILDING.md).

No Python or C compiler is required on the machine where you run this compiler.
"@
    Set-Content -Path (Join-Path $bundle "README-COMPILER.txt") -Value $readme -Encoding UTF8

    if (Test-Path -LiteralPath $OutZip) {
        Remove-Item -LiteralPath $OutZip -Force
    }
    $zipParent = Split-Path -Parent $OutZip
    if ($zipParent -and -not (Test-Path $zipParent)) {
        New-Item -ItemType Directory -Path $zipParent -Force | Out-Null
    }
    Compress-Archive -Path (Join-Path $staging "MoonBasic-compiler") -DestinationPath $OutZip -Force
    Write-Host "Created: $OutZip"
}
finally {
    Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Pop-Location
    if (Test-Path $staging) {
        Remove-Item -LiteralPath $staging -Recurse -Force -ErrorAction SilentlyContinue
    }
}
