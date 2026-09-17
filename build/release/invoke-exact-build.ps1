param(
    [Parameter(Mandatory = $true)][string]$RepositoryRoot,
    [Parameter(Mandatory = $true)][string]$ExpectedToolingRevision,
    [Parameter(Mandatory = $true)][string]$AcquisitionDirectory,
    [Parameter(Mandatory = $true)][string]$OutputDirectory,
    [switch]$AllowCanonicalEolProjection,
    [switch]$PreflightOnly,
    [ValidateSet('Candidate','Validation')][string]$Mode='Candidate'
)

$ErrorActionPreference = 'Stop'
if ($ExpectedToolingRevision -notmatch '^[0-9a-f]{40}$') { throw 'Exact launcher requires one lowercase SHA-1 tooling revision' }

. (Join-Path $PSScriptRoot 'source-trust.ps1')

$trust = Assert-ExactGitSourceTrust -Repository $RepositoryRoot -ExpectedRevision $ExpectedToolingRevision -AllowCanonicalEolProjection:$AllowCanonicalEolProjection
if ($PreflightOnly) {
    Write-Output "Exact source trust PASS commit=$($trust.Commit) tree=$($trust.Tree) files=$($trust.FileCount) raw_equal=$($trust.RawEqualCount) canonical_crlf=$($trust.CanonicalEolProjectionCount)"
    exit 0
}

$materializationRoot = Join-Path ([IO.Path]::GetTempPath()) ("pscan-06-exact-tooling-" + [guid]::NewGuid().ToString('N'))
try {
    $materialized = Export-ExactGitTreeMaterialization -Repository $trust.Repository -Revision $trust.Commit -Destination $materializationRoot
    $buildScript = Join-Path $materialized.Destination 'build\release\build.ps1'
    $sourceTrustScript = Join-Path $materialized.Destination 'build\release\source-trust.ps1'
    if (!(Test-Path -LiteralPath $buildScript -PathType Leaf) -or !(Test-Path -LiteralPath $sourceTrustScript -PathType Leaf)) {
        throw 'Exact tooling materialization omitted required build entry files'
    }
    $env:PSCAN_TRUSTED_LAUNCHER_REVISION = $trust.Commit
    $env:PSCAN_TRUSTED_LAUNCHER_TREE = $trust.Tree
    $env:PSCAN_VERIFIED_REPOSITORY_ROOT = $trust.Repository
    & pwsh -NoProfile -File $buildScript -RepositoryRoot $trust.Repository -ExpectedToolingRevision $trust.Commit -AcquisitionDirectory $AcquisitionDirectory -OutputDirectory $OutputDirectory -AllowCanonicalEolProjection:$AllowCanonicalEolProjection -Mode $Mode
    if ($LASTEXITCODE -ne 0) { throw 'Exact committed release build failed' }
} finally {
    Remove-Item -LiteralPath $materializationRoot -Recurse -Force -ErrorAction SilentlyContinue
}
