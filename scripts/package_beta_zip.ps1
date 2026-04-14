# Package a Windows "Beta" distribution zip: static-linked moonrun + loose folders.
#
# Layout (default: one root folder inside the zip so "Extract all" stays tidy):
#   MoonBasic/
#     moonrun.exe          (from build_static.ps1 output, renamed for clarity)
#     shaders/shd/*.fs     (disk mirror of runtime/shaders/shd for path-based loads / overrides)
#     assets/              (repo assets/, e.g. fonts)
#     examples/            (optional, whole examples tree)
#     README-BETA.txt      (how to run, cwd, failure modes)
#
# Prerequisites: same as scripts/build_static.ps1 (Go, Zig or CC/CXX, Jolt static libs for full physics).
#
# Usage (PowerShell, from repo root):
#   .\scripts\package_beta_zip.ps1
#   .\scripts\package_beta_zip.ps1 -SkipBuild -ExePath .\moonrun_static.exe
#   .\scripts\package_beta_zip.ps1 -Flat          # zip root = files (no MoonBasic/ wrapper)
#   .\scripts\package_beta_zip.ps1 -NoExamples    # smaller archive

[CmdletBinding()]
param(
    [string]$RepoRoot = "",
    [string]$OutZip = "",
    [switch]$SkipBuild,
    [string]$ExePath = "",
    [switch]$NoExamples,
    [switch]$Flat
)

$ErrorActionPreference = "Stop"

if (-not $RepoRoot) {
    $here = $PSScriptRoot
    if (-not $here -and $MyInvocation.MyCommand.Path) {
        $here = Split-Path -Parent $MyInvocation.MyCommand.Path
    }
    if (-not $here) {
        $here = (Get-Location).Path
    }
    $RepoRoot = (Resolve-Path (Join-Path $here "..")).Path
}

if (-not $OutZip) {
    $dist = Join-Path $RepoRoot "dist"
    if (-not (Test-Path $dist)) {
        New-Item -ItemType Directory -Path $dist | Out-Null
    }
    $OutZip = Join-Path $dist "MoonBasic-beta-windows-amd64.zip"
}

$staging = Join-Path $env:TEMP ("moonbasic_pkg_" + [Guid]::NewGuid().ToString("n"))
New-Item -ItemType Directory -Path $staging -Force | Out-Null

try {
    Push-Location $RepoRoot
    if (-not $SkipBuild) {
        Write-Host "Building static moonrun (see scripts/build_static.ps1)..."
        & (Join-Path $RepoRoot "scripts\build_static.ps1")
        if ($LASTEXITCODE -ne 0) {
            exit $LASTEXITCODE
        }
    }

    if (-not $ExePath) {
        $ExePath = Join-Path $RepoRoot "moonrun_static.exe"
    }
    if (-not (Test-Path -LiteralPath $ExePath)) {
        Write-Error "Executable not found: $ExePath (build first or pass -ExePath)"
    }

    $bundleRoot = if ($Flat) { $staging } else {
        $b = Join-Path $staging "MoonBasic"
        New-Item -ItemType Directory -Path $b -Force | Out-Null
        $b
    }

    Copy-Item -LiteralPath $ExePath -Destination (Join-Path $bundleRoot "moonrun.exe")

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
MoonBASIC Beta bundle (Windows amd64)

Run from this folder (or the unzipped root that contains moonrun.exe, shaders, assets, examples):

  .\moonrun.exe examples\sphere_drop\main.mb

RES.PATH and paths relative to the executable resolve against the directory containing moonrun.exe.
Keep shaders\, assets\, and examples\ next to moonrun.exe.

Failure modes:
  - File not found: a script references a missing file; restore the full zip tree or fix the path.
  - Wrong working directory: prefer running commands from this bundle root so relative paths in scripts match.
  - Missing DLL errors: this build should be static (no raylib.dll / jolt.dll). If you see DLL load errors,
    you are not running the static-linked exe from scripts/build_static.ps1 / package_beta_zip.ps1.

Future: embed.FS single-file distribution is optional and documented separately.
"@
    Set-Content -Path (Join-Path $bundleRoot "README-BETA.txt") -Value $readme -Encoding UTF8

    if (Test-Path -LiteralPath $OutZip) {
        Remove-Item -LiteralPath $OutZip -Force
    }
    $zipParent = Split-Path -Parent $OutZip
    if ($zipParent -and -not (Test-Path $zipParent)) {
        New-Item -ItemType Directory -Path $zipParent -Force | Out-Null
    }

    # Single top-level entry in archive: either MoonBasic\... or flat files under staging
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
