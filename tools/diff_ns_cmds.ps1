param([string]$NS = "PLAYER")

$raw  = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json = $raw | ConvertFrom-Json
$pat  = $NS + ".*"
$cmds = $json.commands | Where-Object { $_.key -like $pat } | Sort-Object key

$docPath = "docs\reference\$NS.md"
if (-not (Test-Path $docPath)) { Write-Output "No doc file: $docPath"; exit 1 }
$content = [System.IO.File]::ReadAllText($docPath)

$missing = @()
foreach ($c in $cmds) {
    $shortKey = $c.key.Split('.')[1]   # e.g. MOVE from PLAYER.MOVE
    $inDoc = $content -match [regex]::Escape($c.key) -or $content -match [regex]::Escape($shortKey)
    if (-not $inDoc) {
        $missing += $c.key
    }
}

Write-Output "Namespace: $NS  |  Manifest commands: $($cmds.Count)  |  Undocumented: $($missing.Count)"
if ($missing.Count -gt 0) {
    $missing | ForEach-Object { Write-Output "  MISSING: $_" }
}
