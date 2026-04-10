# moonBASIC dev helpers (PowerShell). From repo root: .\scripts\dev.ps1 <target>
param(
    [Parameter(Position = 0)]
    [ValidateSet("build-compiler", "build-moonrun", "test", "check", "check-builds", "run-spin-cube", "help")]
    [string]$Target = "help"
)

$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

switch ($Target) {
    "build-compiler" { go build -o moonbasic.exe . }
    "build-moonrun"  { go build -tags fullruntime -o moonrun.exe ./cmd/moonrun }
    "test"           { go test ./... }
    "check"          { go run . --check examples/mario64/main_entities.mb }
    "check-builds"   { & powershell -File (Join-Path $PSScriptRoot "check_builds.ps1") }
    "run-spin-cube" {
        $env:CGO_ENABLED = "1"
        go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
    }
    default {
        Write-Host "Usage: .\scripts\dev.ps1 <build-compiler|build-moonrun|test|check|check-builds|run-spin-cube>"
    }
}
