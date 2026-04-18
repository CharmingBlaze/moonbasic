# Known redirect/stub docs that intentionally omit Core Workflow and Full Example
# because they point to a parent doc for full documentation.
$stubs = @(
    "ANGLE.md","AXIS.md","CHAR.md","CHECK.md","CLIENT.md","CLIPBOARD.md",
    "DB.md","DRAWTEX2.md","DRAWTEXPRO.md","DRAWTEXREC.md","ENEMY.md","ENET.md",
    "FREE.md","JOINT.md","KEY.md","KINEMATIC.md","KINEMATICREF.md","MATRIX.md",
    "MODERN_BLITZ_COMMANDS.md","MOVE.md","MUSIC.md","PACKET.md","PEER.md",
    "RES.md","RPC.md","SERVER.md","SHAPEREF.md","SPAWNER.md","SPRITEGROUP.md",
    "SPRITELAYER.md","SPRITEUI.md","STATIC.md","TEXTEXOBJ.md",
    "API_CONVENTIONS.md"
)

$files = Get-ChildItem 'docs\reference\*.md'
$results = @()
foreach ($f in $files) {
    if ($stubs -contains $f.Name) { continue }
    $c = [System.IO.File]::ReadAllText($f.FullName)
    $hasWorkflow = $c -match '## Core Workflow'
    $hasExample  = $c -match '## Full Example'
    $hasSeeAlso  = $c -match '## See also'
    if (-not $hasWorkflow -or -not $hasExample) {
        $results += [PSCustomObject]@{
            File      = $f.Name
            Workflow  = $hasWorkflow
            Example   = $hasExample
            SeeAlso   = $hasSeeAlso
        }
    }
}
Write-Output "Non-stub files missing Core Workflow or Full Example: $($results.Count)"
if ($results.Count -gt 0) {
    $results | Sort-Object File | Format-Table -AutoSize
} else {
    Write-Output "ALL WAVE PATTERN CHECKS PASS"
}
