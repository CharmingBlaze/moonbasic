# Build moonbasic.exe + moonrun.exe into ./dist (same flags as .github/workflows/release.yml).
# Requires MSYS2 MinGW64 with gcc and raylib (pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-raylib).
# From repo root: .\scripts\release-windows.ps1

$ErrorActionPreference = "Stop"
$MsysBash = "C:\msys64\usr\bin\bash.exe"
$Cygpath  = "C:\msys64\usr\bin\cygpath.exe"
if (-not (Test-Path $MsysBash)) {
    Write-Error "MSYS2 not found at $MsysBash. Install MSYS2 or edit this script."
}

$Sh = Join-Path $PSScriptRoot "release-windows.sh"
if (-not (Test-Path $Sh)) {
    Write-Error "Missing $Sh"
}

$ShUnix = if (Test-Path $Cygpath) { & $Cygpath -u $Sh } else { $Sh -replace '\\', '/' -replace '^C:', '/c' -replace '^c:', '/c' }

& $MsysBash -lc "export PATH=/c/Progra~1/Go/bin:/mingw64/bin:/usr/bin:`$PATH && exec bash '$ShUnix'"
