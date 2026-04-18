$manifest = Get-Content "$PSScriptRoot\..\docs\audit\manifest_keys.txt" | Where-Object { $_.Trim() -ne "" } | ForEach-Object { $_.Trim().ToUpper() }
$runtime  = Get-Content "$PSScriptRoot\..\docs\audit\runtime_keys.txt" | Where-Object { $_.Trim() -ne "" } | ForEach-Object { $_.Trim().ToUpper() }

$manifestSet = @{}
foreach ($k in $manifest) { $manifestSet[$k] = $true }

$runtimeSet = @{}
foreach ($k in $runtime) { $runtimeSet[$k] = $true }

# Commands in runtime but NOT in manifest (missing from compiler)
$missingFromManifest = @()
foreach ($k in $runtime) {
    if (-not $manifestSet.ContainsKey($k)) {
        $missingFromManifest += $k
    }
}

# Commands in manifest but NOT in runtime (orphaned in compiler)
$missingFromRuntime = @()
foreach ($k in $manifest) {
    if (-not $runtimeSet.ContainsKey($k)) {
        $missingFromRuntime += $k
    }
}

$missingFromManifest = $missingFromManifest | Sort-Object -Unique
$missingFromRuntime  = $missingFromRuntime  | Sort-Object -Unique

Write-Host "=== IN RUNTIME BUT MISSING FROM MANIFEST ($($missingFromManifest.Count)) ==="
$missingFromManifest | ForEach-Object { Write-Host "  $_" }

Write-Host ""
Write-Host "=== IN MANIFEST BUT MISSING FROM RUNTIME ($($missingFromRuntime.Count)) ==="
$missingFromRuntime | ForEach-Object { Write-Host "  $_" }

# Also write to file
$out = @()
$out += "# Missing Commands Audit"
$out += ""
$out += "## In Runtime but Missing from Manifest ($($missingFromManifest.Count))"
$out += "These commands are registered in Go runtime code but have no entry in commands.json."
$out += "The compiler will reject .mb scripts that try to use them."
$out += ""
foreach ($k in $missingFromManifest) { $out += "- ``$k``" }
$out += ""
$out += "## In Manifest but Missing from Runtime ($($missingFromRuntime.Count))"
$out += "These commands are declared in commands.json but have no runtime registration."
$out += "Scripts compile but will fail at runtime with 'unknown command'."
$out += ""
foreach ($k in $missingFromRuntime) { $out += "- ``$k``" }

$out | Out-File -Encoding utf8 "$PSScriptRoot\..\docs\MISSING_COMMANDS_AUDIT.md"
Write-Host ""
Write-Host "Written to docs/MISSING_COMMANDS_AUDIT.md"
