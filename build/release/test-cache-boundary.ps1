param(
    [Parameter(Mandatory = $true)][string]$PositiveCacheDirectory,
    [Parameter(Mandatory = $true)][string]$WrongOwnerCacheDirectory,
    [Parameter(Mandatory = $true)][string]$ReadOnlyCacheDirectory
)

$ErrorActionPreference = 'Stop'
if (!$IsLinux) { throw 'The Correction C1 ownership regression test requires an actual Linux host' }

$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$uidText = (& id -u).Trim()
$gidText = (& id -g).Trim()
if ($LASTEXITCODE -ne 0 -or $uidText -notmatch '^\d+$' -or $gidText -notmatch '^\d+$') { throw 'Numeric Linux test identity is unavailable' }
$uid = [int]$uidText
$gid = [int]$gidText

. (Join-Path $PSScriptRoot 'cache-canary.ps1')

function Require-EmptyWithoutLedger([string]$Path) {
    if (Test-Path -LiteralPath (Join-Path $Path 'acquisition-ledger.json')) { throw "Failed cache test emitted a success ledger: $Path" }
    if (Get-ChildItem -LiteralPath $Path -Force | Select-Object -First 1) { throw "Cache test left a canary or acquired bytes: $Path" }
}

Invoke-ReleaseCacheCanary -Image $image -ModuleCache $PositiveCacheDirectory -HostUID $uid -HostGID $gid
Require-EmptyWithoutLedger $PositiveCacheDirectory

foreach ($case in @(
    [ordered]@{Name='wrong-owner';Path=$WrongOwnerCacheDirectory;ReadOnly=$false},
    [ordered]@{Name='read-only';Path=$ReadOnlyCacheDirectory;ReadOnly=$true}
)) {
    $failedClosed = $false
    try {
        Invoke-ReleaseCacheCanary -Image $image -ModuleCache $case.Path -HostUID $uid -HostGID $gid -ReadOnlyCache:$case.ReadOnly
    } catch {
        $failedClosed = $true
    }
    if (!$failedClosed) { throw "$($case.Name) cache unexpectedly admitted acquisition" }
    Require-EmptyWithoutLedger $case.Path
}

Write-Output "Linux cache boundary PASS uid=$uid gid=$gid positive=PASS wrong-owner=REJECT read-only=REJECT"
