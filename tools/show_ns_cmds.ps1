param([string]$NS = "CHARCONTROLLER")
$raw  = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json = $raw | ConvertFrom-Json
$pat  = $NS + ".*"
$cmds = $json.commands | Where-Object { $_.key -like $pat } | Sort-Object key
$cmds | Select-Object key, description | Format-Table -AutoSize -Wrap
