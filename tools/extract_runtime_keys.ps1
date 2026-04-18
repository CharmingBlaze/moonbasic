$files = Get-ChildItem -Recurse -Include *.go -Path "$PSScriptRoot\..\runtime"
$keys = @()
foreach ($f in $files) {
    $content = Get-Content $f.FullName -Raw
    $matches = [regex]::Matches($content, '\.Register\("([^"]+)"')
    foreach ($m in $matches) {
        $keys += $m.Groups[1].Value.ToUpper()
    }
}
$unique = $keys | Sort-Object -Unique
$unique | Out-File -Encoding utf8 "$PSScriptRoot\..\docs\audit\runtime_keys.txt"
Write-Host "Runtime keys: $($unique.Count)"
