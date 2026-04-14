# Verifies MinGW link inputs for Jolt on Windows x64 (fullruntime + CGO).
# Exit 0 if both archives exist; exit 1 with a short hint otherwise.
$ErrorActionPreference = "Stop"
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$libDir = Join-Path $repoRoot "third_party\jolt-go\jolt\lib\windows_amd64"
$jolt = Join-Path $libDir "libJolt.a"
$wrap = Join-Path $libDir "libjolt_wrapper.a"
if (-not (Test-Path $jolt) -or -not (Test-Path $wrap)) {
    Write-Host "Missing Jolt static libraries under:" -ForegroundColor Yellow
    Write-Host "  $libDir" -ForegroundColor Yellow
    Write-Host "Expected: libJolt.a, libjolt_wrapper.a" -ForegroundColor Yellow
    Write-Host "See third_party/jolt-go/jolt/lib/windows_amd64/README.md and third_party/jolt-go/scripts/build-libs-windows.ps1 (set JPH_SRC)." -ForegroundColor Gray
    exit 1
}
Write-Host "OK: $jolt" -ForegroundColor Green
Write-Host "OK: $wrap" -ForegroundColor Green
