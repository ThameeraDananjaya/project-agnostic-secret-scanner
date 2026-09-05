param(
    [Parameter(Mandatory = $true)][string]$PositiveCacheDirectory,
    [Parameter(Mandatory = $true)][string]$WrongOwnerCacheDirectory,
    [Parameter(Mandatory = $true)][string]$ReadOnlyCacheDirectory,
    [ValidateSet('Host','Container','All')][string]$Phase = 'All'
)

$ErrorActionPreference = 'Stop'
if (!$IsLinux) { throw 'The Correction C2 ownership regression test requires an actual Linux host' }

$image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$uidText = (& id -u).Trim()
$gidText = (& id -g).Trim()
if ($LASTEXITCODE -ne 0 -or $uidText -notmatch '^\d+$' -or $gidText -notmatch '^\d+$') { throw 'Numeric Linux test identity is unavailable' }
$uid = [int]$uidText
$gid = [int]$gidText

. (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
. (Join-Path $PSScriptRoot 'cache-canary.ps1')

function Require-EmptyWithoutLedger([string]$Path) {
    if (Test-Path -LiteralPath (Join-Path $Path 'acquisition-ledger.json')) { throw "Failed cache test emitted a success ledger: $Path" }
    if (Get-ChildItem -LiteralPath $Path -Force | Select-Object -First 1) { throw "Cache test left a canary or acquired bytes: $Path" }
}

if ($Phase -in @('Host','All')) {
    [void](Invoke-HostCacheCanary -CacheDirectory $PositiveCacheDirectory)
    Require-EmptyWithoutLedger $PositiveCacheDirectory

    foreach ($case in @(
        [ordered]@{Name='wrong-owner';Path=$WrongOwnerCacheDirectory},
        [ordered]@{Name='read-only';Path=$ReadOnlyCacheDirectory}
    )) {
        $failedClosed = $false
        try { [void](Invoke-HostCacheCanary -CacheDirectory $case.Path) } catch { $failedClosed = $true }
        if (!$failedClosed) { throw "$($case.Name) cache unexpectedly admitted the host canary" }
        Require-EmptyWithoutLedger $case.Path
    }
    if ($Phase -eq 'Host') {
        Write-Output "Linux host cache boundary PASS uid=$uid gid=$gid positive=PASS wrong-owner=REJECT read-only=REJECT docker=NOT_INVOKED"
        return
    }
}

$engineJson = & (Join-Path $PSScriptRoot 'docker-execution.ps1') -Operation EngineInspection
$engine = ($engineJson -join "`n") | ConvertFrom-Json
if ($engine.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($engine.StdErr) -or !$engine.ContainmentEmpty) { throw 'Docker engine boundary is untrusted' }
. (Join-Path $PSScriptRoot 'image-admission.ps1')
[void](Assert-ReleaseEngineEvidence -Json $engine.StdOut)
$dockerSHA256 = $engine.DockerSHA256

Invoke-ReleaseCacheCanary -Image $image -ModuleCache $PositiveCacheDirectory -ExpectedDockerSHA256 $dockerSHA256 -HostUID $uid -HostGID $gid
Require-EmptyWithoutLedger $PositiveCacheDirectory

foreach ($case in @(
    [ordered]@{Name='wrong-owner';Path=$WrongOwnerCacheDirectory;ReadOnly=$false},
    [ordered]@{Name='read-only';Path=$ReadOnlyCacheDirectory;ReadOnly=$true}
)) {
    $failedClosed = $false
    try {
        Invoke-ReleaseCacheCanary -Image $image -ModuleCache $case.Path -ExpectedDockerSHA256 $dockerSHA256 -HostUID $uid -HostGID $gid -ReadOnlyCache:$case.ReadOnly
    } catch {
        $failedClosed = $true
    }
    if (!$failedClosed) { throw "$($case.Name) cache unexpectedly admitted acquisition" }
    Require-EmptyWithoutLedger $case.Path
}

Write-Output "Linux cache boundary PASS uid=$uid gid=$gid positive=PASS wrong-owner=REJECT read-only=REJECT"
