# Package a Windows folder zip using a *release-style* moonrun.exe (same link model as GitHub Releases).
#
# Does NOT compile anything. Provide moonrun.exe built like CI/release:
#   - Copy moonrun.exe (and optionally moonbasic.exe) from an official full-runtime zip, or
#   - Build locally with MSYS2 using the same steps as .github/workflows (Jolt from
#     third_party/jolt-go/scripts/build-libs-windows.ps1 + ldflags from
#     scripts/windows_fullruntime_go_ldflags.sh) — see docs/BUILDING.md.
#
# Layout (default: MoonBasic/ root inside the zip):
#   MoonBasic/
#     moonrun.exe
#     moonbasic.exe   (optional, -IncludeMoonbasic)
#     shaders/shd/...
#     assets/...
#     examples/...    (unless -NoExamples)
#     README-BUNDLE.txt
#
# Usage (PowerShell, from repo root):
#   .\scripts\package_release_style_zip.ps1
#   .\scripts\package_release_style_zip.ps1 -ExePath .\dist\moonrun.exe -MoonbasicPath .\dist\moonbasic.exe
#   .\scripts\package_release_style_zip.ps1 -NoExamples -Flat

[CmdletBinding()]
param(
    [string]$RepoRoot = "",
    [string]$OutZip = "",
    [string]$ExePath = "",
    [string]$MoonbasicPath = "",
    [switch]$IncludeMoonbasic,
    [switch]$NoExamples,
    [switch]$Flat
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
    $OutZip = Join-Path $dist "MoonBasic-release-style-windows-amd64.zip"
}

if (-not $ExePath) {
    $ExePath = Join-Path $RepoRoot "dist\moonrun.exe"
}
if (-not (Test-Path -LiteralPath $ExePath)) {
    Write-Error "moonrun.exe not found: $ExePath. Build or copy a release-style binary (see script header and docs/BUILDING.md)."
}

if ($IncludeMoonbasic) {
    if (-not $MoonbasicPath) {
        $MoonbasicPath = Join-Path $RepoRoot "dist\moonbasic.exe"
    }
    if (-not (Test-Path -LiteralPath $MoonbasicPath)) {
        Write-Error "moonbasic.exe not found: $MoonbasicPath (omit -IncludeMoonbasic or copy moonbasic.exe to dist/)."
    }
}

$staging = Join-Path $env:TEMP ("moonbasic_rel_pkg_" + [Guid]::NewGuid().ToString("n"))
New-Item -ItemType Directory -Path $staging -Force | Out-Null

try {
    Push-Location $RepoRoot

    $bundleRoot = if ($Flat) { $staging } else {
        $b = Join-Path $staging "MoonBasic"
        New-Item -ItemType Directory -Path $b -Force | Out-Null
        $b
    }

    Copy-Item -LiteralPath $ExePath -Destination (Join-Path $bundleRoot "moonrun.exe")
    if ($IncludeMoonbasic) {
        Copy-Item -LiteralPath $MoonbasicPath -Destination (Join-Path $bundleRoot "moonbasic.exe")
    }

    $shdSrc = Join-Path $RepoRoot "runtime\shaders\shd"
    if (-not (Test-Path $shdSrc)) {
        Write-Error "Missing shader tree: $shdSrc"
    }
    $shdDst = Join-Path $bundleRoot "shaders\shd"
    New-Item -ItemType Directory -Path $shdDst -Force | Out-Null
    Copy-Item -Path (Join-Path $shdSrc "*") -Destination $shdDst -Recurse -Force

    $assetsSrc = Join-Path $RepoRoot "assets"
    if (Test-Path $assetsSrc) {
        Copy-Item -LiteralPath $assetsSrc -Destination (Join-Path $bundleRoot "assets") -Recurse -Force
    }

    if (-not $NoExamples) {
        $exSrc = Join-Path $RepoRoot "examples"
        if (Test-Path $exSrc) {
            Copy-Item -LiteralPath $exSrc -Destination (Join-Path $bundleRoot "examples") -Recurse -Force
        }
    }

    $readme = @"
MoonBASIC bundle (Windows amd64) — release-style moonrun.exe

This zip was built with scripts/package_release_style_zip.ps1. moonrun.exe must match
official full-runtime builds (see docs/BUILDING.md, "Windows full-runtime PE link model").

Run from this folder:

  .\moonrun.exe examples\pong\main.mb

Keep shaders\, assets\, and examples\ next to moonrun.exe when samples use disk paths.
For players who only need to run your game, they can install the official full-runtime
zip from GitHub Releases instead of this bundle.

See docs/GETTING_STARTED.md — section "Ship your game (for authors)".
"@
    Set-Content -Path (Join-Path $bundleRoot "README-BUNDLE.txt") -Value $readme -Encoding UTF8

    if (Test-Path -LiteralPath $OutZip) {
        Remove-Item -LiteralPath $OutZip -Force
    }
    $zipParent = Split-Path -Parent $OutZip
    if ($zipParent -and -not (Test-Path $zipParent)) {
        New-Item -ItemType Directory -Path $zipParent -Force | Out-Null
    }

    $compressPath = if ($Flat) { Join-Path $staging "*" } else { Join-Path $staging "MoonBasic" }
    Compress-Archive -Path $compressPath -DestinationPath $OutZip -Force

    Write-Host "Created: $OutZip"
}
finally {
    Pop-Location
    if (Test-Path $staging) {
        Remove-Item -LiteralPath $staging -Recurse -Force -ErrorAction SilentlyContinue
    }
}
