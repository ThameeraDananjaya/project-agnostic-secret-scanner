param(
    [Parameter(Mandatory = $true)][string]$FirstOutputDirectory,
    [Parameter(Mandatory = $true)][string]$SecondOutputDirectory
)

$ErrorActionPreference = 'Stop'
$first = (Resolve-Path -LiteralPath (Join-Path $FirstOutputDirectory 'dist')).Path
$second = (Resolve-Path -LiteralPath (Join-Path $SecondOutputDirectory 'dist')).Path
$firstFiles = Get-ChildItem -LiteralPath $first -File | Sort-Object Name
$secondFiles = Get-ChildItem -LiteralPath $second -File | Sort-Object Name
if ($firstFiles.Count -ne $secondFiles.Count) { throw 'Independent build file counts differ' }
for ($index = 0; $index -lt $firstFiles.Count; $index++) {
    if ($firstFiles[$index].Name -ne $secondFiles[$index].Name) { throw 'Independent build file sets differ' }
    $firstDigest = (Get-FileHash -LiteralPath $firstFiles[$index].FullName -Algorithm SHA256).Hash
    $secondDigest = (Get-FileHash -LiteralPath $secondFiles[$index].FullName -Algorithm SHA256).Hash
    if ($firstDigest -ne $secondDigest) { throw "Independent build bytes differ: $($firstFiles[$index].Name)" }
}
Write-Output "Independent network-disabled builds are byte-identical: $($firstFiles.Count) files"
