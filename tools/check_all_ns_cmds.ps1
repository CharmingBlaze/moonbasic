$raw  = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json = $raw | ConvertFrom-Json

# Get all unique namespaces
$namespaces = $json.commands | ForEach-Object { $_.key.Split('.')[0] } | Sort-Object -Unique

$results = @()
foreach ($ns in $namespaces) {
    $pat  = $ns + ".*"
    $cmds = $json.commands | Where-Object { $_.key -like $pat }
    $cnt  = $cmds.Count
    $docPath = "docs\reference\$ns.md"
    if (-not (Test-Path $docPath)) { continue }
    $content  = [System.IO.File]::ReadAllText($docPath)
    $missing  = 0
    foreach ($c in $cmds) {
        $shortKey = $c.key.Split('.')[1]
        $inDoc = $content -match [regex]::Escape($c.key) -or $content -match [regex]::Escape($shortKey)
        if (-not $inDoc) { $missing++ }
    }
    if ($missing -gt 0) {
        $results += [PSCustomObject]@{ NS=$ns; Total=$cnt; Missing=$missing }
    }
}

if ($results.Count -eq 0) {
    Write-Output "ALL CLEAR - every namespace doc covers all manifest commands"
} else {
    Write-Output "Namespaces with undocumented commands: $($results.Count)"
    $results | Sort-Object -Descending Missing | Format-Table -AutoSize
}
