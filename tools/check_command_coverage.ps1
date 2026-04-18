$raw  = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json = $raw | ConvertFrom-Json

$check = @('DRAW','GAME','PHYSICS','CHAR','CHARCONTROLLER','NAV','ENTITY','PLAYER','ACTION','GAMEPAD','BOX2D','JOLT','CHARACTERREF','BODYREF','NET')
$results = @()
foreach ($ns in $check) {
    $pat  = $ns + ".*"
    $cmds = $json.commands | Where-Object { $_.key -like $pat }
    $cnt  = $cmds.Count
    $docPath = "docs\reference\$ns.md"
    if (Test-Path $docPath) {
        $content  = [System.IO.File]::ReadAllText($docPath)
        $mentions = ([regex]::Matches($content, [regex]::Escape($ns + "."))).Count
    } else { $mentions = 0 }
    $results += [PSCustomObject]@{ NS=$ns; ManifestCmds=$cnt; DocMentions=$mentions; DocExists=(Test-Path $docPath) }
}
$results | Format-Table -AutoSize
