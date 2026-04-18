$raw  = [System.IO.File]::ReadAllText('compiler\builtinmanifest\commands.json')
$json = $raw | ConvertFrom-Json
$json.commands | Where-Object { $_.key -like '*.MAKE' } | Sort-Object key | Select-Object key, description | Format-Table -AutoSize -Wrap
