# Run moonrun on examples/pong/main.mb for a few seconds to validate full-runtime + window + game loop.
# Intended for CI (Windows runner) and local release checks. Kills the process if it is still running
# after -Seconds (success = game loop stayed alive; immediate exit = failure).
#
# Usage (repo root, after building dist/moonrun.exe):
#   powershell -File scripts/smoke_moonrun_pong.ps1 -Exe dist/moonrun.exe -RepoRoot $PWD

param(
    [Parameter(Mandatory = $true)][string]$Exe,
    [Parameter(Mandatory = $true)][string]$RepoRoot,
    [int]$Seconds = 8
)

$ErrorActionPreference = "Stop"

$mb = Join-Path $RepoRoot "examples\pong\main.mb"
if (-not (Test-Path -LiteralPath $mb)) {
    Write-Error "Missing sample: $mb"
}
if (-not (Test-Path -LiteralPath $Exe)) {
    Write-Error "Missing moonrun binary: $Exe"
}

$pong = (Resolve-Path -LiteralPath $mb).Path
$exePath = (Resolve-Path -LiteralPath $Exe).Path

Write-Host "Smoke: $exePath $pong (${Seconds}s) ..."

$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $exePath
$psi.Arguments = "`"$pong`""
$psi.WorkingDirectory = $RepoRoot
$psi.UseShellExecute = $false
$psi.CreateNoWindow = $false

$p = New-Object System.Diagnostics.Process
$p.StartInfo = $psi
[void]$p.Start()

$deadline = (Get-Date).AddSeconds($Seconds)
while (-not $p.HasExited -and (Get-Date) -lt $deadline) {
    Start-Sleep -Milliseconds 150
}

if ($p.HasExited) {
    $code = $p.ExitCode
    Write-Error "moonrun exited early (exit code $code). Full-runtime window/game loop did not stay up."
}

try { $p.Kill($true) } catch { }
try { $p.Dispose() } catch { }

Write-Host "OK: pong stayed running for ${Seconds}s (smoke pass)."
