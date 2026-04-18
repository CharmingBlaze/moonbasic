$raw    = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json   = $raw | ConvertFrom-Json
$nsSet  = New-Object System.Collections.Generic.HashSet[string]
foreach ($c in $json.commands) { $null = $nsSet.Add($c.key.Split('.')[0]) }
$docSet = New-Object System.Collections.Generic.HashSet[string]
foreach ($f in Get-ChildItem 'docs\reference\*.md') { $null = $docSet.Add($f.BaseName.ToUpper()) }

$missing = @()
foreach ($ns in ($nsSet | Sort-Object)) {
    $pat = $ns + ".*"
    $cnt = 0
    foreach ($c in $json.commands) { if ($c.key -like $pat) { $cnt++ } }
    if ($cnt -gt 0 -and -not $docSet.Contains($ns)) {
        $missing += "$ns ($cnt commands)"
    }
}

Write-Output "Total doc files: $($docSet.Count)"
Write-Output "Missing namespace docs: $($missing.Count)"
if ($missing.Count -gt 0) {
    $missing | ForEach-Object { Write-Output "  MISSING: $_" }
} else {
    Write-Output "ALL CLEAR - every real-command namespace has a doc file!"
}
