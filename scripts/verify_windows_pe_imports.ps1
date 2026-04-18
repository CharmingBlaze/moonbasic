# Verify a PE built with MinGW does not import companion runtime / Raylib DLLs we intend to embed.
#
# Maintenance (contributors):
#   When you add a new Windows CGO dependency to moonbasic/moonrun, run this script on the
#   resulting .exe. If a NEW non-system DLL appears in the import table:
#   1) Prefer fixing the link (static archive, correct -l order) so the DLL is not required.
#   2) If a sidecar DLL is unavoidable, add an explicit allowlist entry BELOW with the DLL
#      name and a one-line rationale; update docs/BUILDING.md "Windows full-runtime PE link model".
#   3) Do not remove existing forbidden checks without maintainer review.
#
# Optional future allowlist (regex or literal), currently unused — example:
#   # $Allowlisted = @('^SomeCodec\.dll$')
#
# Usage (from repo root, MSYS2 mingw64 bin on PATH):
#   powershell -File scripts/verify_windows_pe_imports.ps1 -Exe dist/moonrun.exe -MingwBin "C:\msys64\mingw64\bin"

param(
    [Parameter(Mandatory = $true)][string]$Exe,
    [Parameter(Mandatory = $true)][string]$MingwBin
)

$ErrorActionPreference = "Stop"

$objdump = Join-Path $MingwBin "objdump.exe"
if (-not (Test-Path -LiteralPath $objdump)) {
    Write-Error "objdump not found: $objdump (install mingw-w64-x86_64-binutils or full toolchain)"
}

$out = & $objdump -p $Exe 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error "objdump failed for $Exe"
}

$dlls = @()
foreach ($line in $out) {
    if ($line -match "^\s*DLL Name:\s+(\S+)") {
        $dlls += $Matches[1]
    }
}

# Names we do not want next to the exe for "fully static MinGW runtime" distro builds.
$forbidden = @()
foreach ($d in $dlls) {
    $lower = $d.ToLowerInvariant()
    if ($lower -eq "raylib.dll") { $forbidden += $d; continue }
    if ($lower -eq "libstdc++-6.dll") { $forbidden += $d; continue }
    if ($lower -like "libwinpthread*.dll") { $forbidden += $d; continue }
    if ($lower -like "libgcc_s_*.dll") { $forbidden += $d; continue }
}

if ($forbidden.Count -gt 0) {
    Write-Error ("PE imports unexpected DLL(s): {0}. Full import list: {1}" -f `
            ($forbidden -join ", "), ($dlls -join ", "))
}

Write-Host ('OK: {0} - no forbidden MinGW/Raylib DLL imports ({1} DLL import name(s) in PE).' -f $Exe, $dlls.Count)
