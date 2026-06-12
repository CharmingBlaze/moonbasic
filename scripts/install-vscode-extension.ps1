# Install the moonBASIC VS Code / Cursor extension (one command).
# Usage: powershell -ExecutionPolicy Bypass -File scripts/install-vscode-extension.ps1
$ErrorActionPreference = "Stop"
$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
$mb = Join-Path $root "moonbasic.exe"
if (Test-Path $mb) {
  & $mb install-vscode @args
} elseif (Get-Command moonbasic -ErrorAction SilentlyContinue) {
  moonbasic install-vscode @args
} else {
  Write-Host "Building moonbasic from source …"
  Push-Location $root
  go run . install-vscode @args
  Pop-Location
}
