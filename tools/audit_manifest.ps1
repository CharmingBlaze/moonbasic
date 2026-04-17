$j = Get-Content "$PSScriptRoot\..\compiler\builtinmanifest\commands.json" -Raw | ConvertFrom-Json
$manifestKeys = $j.commands | ForEach-Object { $_.key.ToUpper() } | Sort-Object -Unique
$manifestKeys | Out-File -Encoding utf8 "$PSScriptRoot\..\manifest_keys.txt"
Write-Host "Manifest keys: $($manifestKeys.Count)"
