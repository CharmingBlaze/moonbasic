# Run both checks and summarise
Write-Output "=== NAMESPACE COVERAGE ==="
& "$PSScriptRoot\check_doc_coverage.ps1"

Write-Output ""
Write-Output "=== WAVE PATTERN ==="
& "$PSScriptRoot\check_wave_pattern.ps1"
