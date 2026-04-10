# Dry-run both build-tag paths (same as scripts/check_builds.sh).
# From repo root: powershell -File scripts/check_builds.ps1
# If fullruntime fails with runtime/cgo in plain PowerShell, use Git Bash / MSYS2:
#   bash scripts/check_builds.sh
# (MinGW gcc on PATH; same idea as scripts/release-windows.sh.)
$ErrorActionPreference = "Stop"
function Invoke-GoBuild {
    param([string[]]$GoArgs)
    & go @GoArgs
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root
$Out = Join-Path $Root ".check"
if (Test-Path $Out) { Remove-Item -Recurse -Force $Out }
New-Item -ItemType Directory -Path $Out | Out-Null
try {
    if (-not $env:CGO_ENABLED) { $env:CGO_ENABLED = "1" }
    $MsysGcc = "C:\msys64\mingw64\bin\gcc.exe"
    if (-not $env:CC -and (Test-Path $MsysGcc)) {
        $env:CC = $MsysGcc
    }

    Write-Host "== cmd/moonbasic (default tags, compiler CLI) =="
    Invoke-GoBuild -GoArgs @("build", "-o", (Join-Path $Out "moonbasic-cli"), "./cmd/moonbasic")

    Write-Host "== root (default, compiler-only) =="
    Invoke-GoBuild -GoArgs @("build", "-o", (Join-Path $Out "moonbasic-root"), ".")

    Write-Host "== cmd/moonrun (-tags=fullruntime) =="
    Invoke-GoBuild -GoArgs @("build", "-tags=fullruntime", "-o", (Join-Path $Out "moonrun"), "./cmd/moonrun")

    Write-Host "== root (-tags=fullruntime) =="
    Invoke-GoBuild -GoArgs @("build", "-tags=fullruntime", "-o", (Join-Path $Out "moonbasic-full"), ".")

    Write-Host "check_builds: OK (both tag axes compile)"
}
finally {
    if (Test-Path $Out) { Remove-Item -Recurse -Force $Out }
}
