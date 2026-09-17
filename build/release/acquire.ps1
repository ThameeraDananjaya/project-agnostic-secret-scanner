param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [Parameter(Mandatory = $true)][string]$SourceRepository,
    [Parameter(Mandatory = $true)][ValidatePattern('^[0-9a-f]{40}$')][string]$SourceRevision
)

$ErrorActionPreference = 'Stop'

$image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
if ($SourceRevision -cne $env:PSCAN_TRUSTED_LAUNCHER_REVISION) { throw 'Acquisition source and launcher revisions conflict' }
. (Join-Path $PSScriptRoot 'source-trust.ps1')
$sourceTrust = Assert-ExactGitSourceTrust -Repository $SourceRepository -ExpectedRevision $SourceRevision
$root = $sourceTrust.Repository
# Acquisition consumes these files from the admitted complete repository, never
# from the release-entrypoint-only materialization or an ambient fallback.
foreach ($name in @('go.mod','go.sum')) {
    $blob = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse',"${SourceRevision}:$name")).Bytes -Label 'Acquisition module source').Trim()
    if ($blob -notmatch '^[0-9a-f]{40}$' -or (Get-RawFileGitBlobID -Path (Join-Path $root $name)) -cne $blob) { throw 'Acquisition module source bytes differ from the exact commit' }
}
$cache = [IO.Path]::GetFullPath($CacheDirectory)
if (Test-Path -LiteralPath $cache) {
    if (Test-Path -LiteralPath (Join-Path $cache 'acquisition-ledger.json')) {
        throw 'Acquisition cache is already complete and is not mutable'
    }
    $unexpected = Get-ChildItem -LiteralPath $cache -Force | Where-Object { $_.Name -notin @('downloads','gomodcache','docker-admission.json') }
    if ($unexpected) {
        throw 'Incomplete acquisition cache contains unexpected paths'
    }
} else {
    New-Item -ItemType Directory -Path $cache | Out-Null
}
$downloads = Join-Path $cache 'downloads'
$moduleCache = Join-Path $cache 'gomodcache'
New-Item -ItemType Directory -Path $downloads,$moduleCache -Force | Out-Null

$hostUID = $null
$hostGID = $null
$platformMode = 'windows-invoking-host'
if ($IsLinux) {
    $hostUIDText = (& id -u).Trim()
    if ($LASTEXITCODE -ne 0 -or $hostUIDText -notmatch '^\d+$') { throw 'Linux numeric UID discovery failed' }
    $hostGIDText = (& id -g).Trim()
    if ($LASTEXITCODE -ne 0 -or $hostUIDText -notmatch '^\d+$' -or $hostGIDText -notmatch '^\d+$') {
        throw 'Linux host numeric UID/GID discovery failed before dependency acquisition'
    }
    $hostUID = [int]$hostUIDText
    $hostGID = [int]$hostGIDText
    & chmod 0700 -- $cache $downloads $moduleCache
    if ($LASTEXITCODE -ne 0) { throw 'Linux cache mode preparation failed before dependency acquisition' }
    foreach ($path in @($cache,$downloads,$moduleCache)) {
        $owner = (& stat -c '%u:%g' -- $path).Trim()
        if ($LASTEXITCODE -ne 0 -or $owner -ne "${hostUID}:${hostGID}") {
            throw "Linux cache ownership does not match the invoking host identity: $path"
        }
    }
    $platformMode = 'linux-host-numeric-uid-gid'
}
$sourceRevision = $env:PSCAN_TRUSTED_LAUNCHER_REVISION
if ($sourceRevision -notmatch '^[0-9a-f]{40}$') { throw 'Acquisition requires the exact committed launcher revision' }
$receiptPath = Join-Path $cache 'docker-admission.json'
if (!(Test-Path -LiteralPath $receiptPath -PathType Leaf)) { throw 'Dependency acquisition requires the prior closed Docker admission receipt' }
. (Join-Path $PSScriptRoot 'execution-profile.ps1')
$imageAdmission = Read-ReleaseImageAdmission $receiptPath
if ($imageAdmission.sourceRevision -cne $sourceRevision -or $imageAdmission.image -cne $image -or
    $imageAdmission.dockerExecutableSHA256 -notmatch '^[0-9a-f]{64}$' -or $imageAdmission.containment -cne 'empty-after-every-operation' -or
    $imageAdmission.hostIdentityMode -cne $platformMode -or $imageAdmission.hostUID -ne $hostUID -or $imageAdmission.hostGID -ne $hostGID) { throw 'Docker admission receipt does not bind this exact acquisition' }
$dockerSHA256 = $imageAdmission.dockerExecutableSHA256
. (Join-Path $PSScriptRoot 'image-admission.ps1')
$inspectJson = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation RepositoryDigestInspection -ExpectedDockerSHA256 $dockerSHA256
$inspect = ($inspectJson -join "`n") | ConvertFrom-Json
if ($inspect.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($inspect.StdErr) -or !$inspect.ContainmentEmpty) { throw 'Pre-acquisition image inspection failed closed' }
$repoDigests = Read-ReleaseRepoDigestsEvidence -Json $inspect.StdOut
[void](Assert-ReleaseImageIdentityEvidence -Image $image -RepoDigests $repoDigests)

$artifacts = @(
    [ordered]@{
        Name = 'go1.27.1.linux-amd64.tar.gz'
        Uri = 'https://go.dev/dl/go1.27.1.linux-amd64.tar.gz'
        Sha256 = '63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445'
    },
    [ordered]@{
        Name = 'go1.27.0.linux-amd64.tar.gz'
        Uri = 'https://go.dev/dl/go1.27.0.linux-amd64.tar.gz'
        Sha256 = '675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685'
    },
    [ordered]@{
        Name = 'gitleaks-83d9cd684c87d95d656c1458ef04895a7f1cbd8e.tar.gz'
        Uri = 'https://codeload.github.com/gitleaks/gitleaks/tar.gz/83d9cd684c87d95d656c1458ef04895a7f1cbd8e'
        Sha256 = '6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115'
    }
)

foreach ($artifact in $artifacts) {
    $destination = Join-Path $downloads $artifact.Name
    if (Test-Path -LiteralPath $destination -PathType Leaf) {
        $existing = (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($existing -ne $artifact.Sha256) {
            throw "Existing partial acquisition has an invalid digest: $($artifact.Name)"
        }
    } else {
        Invoke-WebRequest -Uri $artifact.Uri -OutFile $destination -MaximumRedirection 5
    }
    $actual = (Get-FileHash -LiteralPath $destination -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $artifact.Sha256) {
        throw "Acquired file digest mismatch: $($artifact.Name)"
    }
}

$runnerGo = (Resolve-Path (Join-Path $downloads 'go1.27.1.linux-amd64.tar.gz')).Path
$engineGo = (Resolve-Path (Join-Path $downloads 'go1.27.0.linux-amd64.tar.gz')).Path
$gitleaks = (Resolve-Path (Join-Path $downloads 'gitleaks-83d9cd684c87d95d656c1458ef04895a7f1cbd8e.tar.gz')).Path

$boundaryParameters = @{
    Operation='DependencyAcquisition'; ExpectedDockerSHA256=$dockerSHA256
    SourceRoot=$root; RunnerGoArchive=$runnerGo; EngineGoArchive=$engineGo
    GitleaksArchive=$gitleaks; CacheDirectory=$moduleCache
}
if ($null -ne $hostUID) { $boundaryParameters.HostUID=[int]$hostUID; $boundaryParameters.HostGID=[int]$hostGID }
$boundaryJson = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') @boundaryParameters
$boundaryResult = ($boundaryJson -join "`n") | ConvertFrom-Json
if ($boundaryResult.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($boundaryResult.StdErr) -or !$boundaryResult.ContainmentEmpty -or $boundaryResult.DockerSHA256 -cne $dockerSHA256) { throw 'Pinned dependency acquisition failed closed at the Docker boundary' }

$ledger = [ordered]@{
    schemaVersion = '2.1'
    acquiredAt = [DateTime]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ssZ')
    image = $image
    artifacts = $artifacts
    moduleCache = 'gomodcache'
    cacheCanary = [ordered]@{
        semantics = 'host-and-container-write-atomic-rename-read-delete'
        completedBeforeNetworkDependencyAcquisition = $true
        completedBeforeImagePull = $true
        hostIdentityMode = $platformMode
        hostUID = $hostUID
        hostGID = $hostGID
    }
    imageAdmission = [ordered]@{
        canonicalReference = $image
        repository = 'docker.io/library/golang'
        digest = 'sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
        dockerExecutableSHA256 = $dockerSHA256
        executionBoundary = 'closed-private-contained'
        pulledDuringAdmission = $imageAdmission.pulledDuringAdmission
        postAdmissionRepoDigestProved = $true
    }
    networkBoundary = 'Network enabled only in this acquisition phase; builds require --network none and read-only cache mounts.'
}
$json = $ledger | ConvertTo-Json -Depth 8
[IO.File]::WriteAllText((Join-Path $cache 'acquisition-ledger.json'), $json + "`n", [Text.UTF8Encoding]::new($false))
Write-Output "Acquisition complete: $cache"
