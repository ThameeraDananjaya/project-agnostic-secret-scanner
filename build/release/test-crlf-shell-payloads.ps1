param(
    [string]$WorkingDirectory,
    [string]$AcquisitionDirectory,
    [string]$BuildOutputDirectory,
    [string]$SourceRepository,
    [string]$SourceRevision
)

$ErrorActionPreference = 'Stop'
$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
if ([string]::IsNullOrWhiteSpace($SourceRepository)) {
    if ([string]::IsNullOrWhiteSpace($env:PSCAN_VERIFIED_REPOSITORY_ROOT)) { throw 'CRLF regression requires an explicit verified source repository' }
    $SourceRepository = $env:PSCAN_VERIFIED_REPOSITORY_ROOT
}
if ([string]::IsNullOrWhiteSpace($SourceRevision)) {
    if ([string]::IsNullOrWhiteSpace($env:PSCAN_TRUSTED_LAUNCHER_REVISION)) { throw 'CRLF regression requires an exact source revision' }
    $SourceRevision = $env:PSCAN_TRUSTED_LAUNCHER_REVISION
}
$sourceRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $SourceRepository).Path)
. (Join-Path $PSScriptRoot 'source-trust.ps1')
$sourceTrust = Assert-ExactGitSourceTrust -Repository $sourceRoot -ExpectedRevision $SourceRevision

if ([string]::IsNullOrWhiteSpace($WorkingDirectory)) {
    $WorkingDirectory = Join-Path ([IO.Path]::GetTempPath()) ("pscan-06-c1-crlf-" + [guid]::NewGuid().ToString('N'))
}
$testRoot = [IO.Path]::GetFullPath($WorkingDirectory)
if (Test-Path -LiteralPath $testRoot) {
    if (Get-ChildItem -LiteralPath $testRoot -Force | Select-Object -First 1) {
        throw 'CRLF regression working directory must be new or empty'
    }
} else {
    New-Item -ItemType Directory -Path $testRoot | Out-Null
}

function Get-EmbeddedPayload([string]$Path, [string]$Marker) {
    $tokens = $null
    $errors = $null
    $ast = [Management.Automation.Language.Parser]::ParseFile($Path, [ref]$tokens, [ref]$errors)
    if ($errors.Count -ne 0) { throw "CRLF fixture did not parse: $Path" }
    $matches = @($ast.FindAll({
        param($node)
        $node -is [Management.Automation.Language.StringConstantExpressionAst] -and
            $node.Value.Contains($Marker)
    }, $true))
    if ($matches.Count -ne 1) { throw "Expected one embedded shell payload containing '$Marker' in $Path" }
    return $matches[0].Value
}

function Require-Empty([string]$Path) {
    if (Test-Path -LiteralPath (Join-Path $Path 'acquisition-ledger.json')) {
        throw "Failed CRLF cache case emitted a success ledger: $Path"
    }
    if (Get-ChildItem -LiteralPath $Path -Force | Select-Object -First 1) {
        throw "CRLF cache case left canary residue: $Path"
    }
}

$fixtureRoot = Join-Path $testRoot 'checkout'
& git clone --no-hardlinks --no-checkout --quiet $sourceRoot $fixtureRoot
if ($LASTEXITCODE -ne 0) { throw 'Local CRLF fixture clone failed' }
& git -C $fixtureRoot config core.autocrlf true
& git -C $fixtureRoot checkout --quiet --detach $SourceRevision
if ($LASTEXITCODE -ne 0) { throw 'Exact CRLF fixture checkout failed' }
$fixtureTrust = Assert-ExactGitSourceTrust -Repository $fixtureRoot -ExpectedRevision $SourceRevision -AllowCanonicalEolProjection
$fixtureRevision = $fixtureTrust.Commit
$fixtureTree = $fixtureTrust.Tree
$crlfIntegrityAsset = Join-Path $fixtureRoot 'rules\generic\gitleaks-ignore-empty-v1.txt'
$crlfIntegrityBytes = [IO.File]::ReadAllBytes($crlfIntegrityAsset)
$crlfIntegrityEol = (& git -C $fixtureRoot ls-files --eol -- 'rules/generic/gitleaks-ignore-empty-v1.txt') -join "`n"
if ($crlfIntegrityEol -notmatch 'w/crlf' -or $crlfIntegrityBytes -notcontains 13) {
    throw 'Fixture did not create an actual CRLF checkout of the pinned integrity asset'
}

$fixtureReleaseRoot = Join-Path $fixtureRoot 'build\release'
. (Join-Path $PSScriptRoot 'shell-payload.ps1')
$payloadCases = @(
    [ordered]@{Name='cache-canary';Path=(Join-Path $fixtureReleaseRoot 'cache-canary.ps1');Marker='PSCAN-06-C1-CACHE-CANARY'},
    [ordered]@{Name='acquisition';Path=(Join-Path $fixtureReleaseRoot 'acquire.ps1');Marker='/work/runner/go/bin/go mod download'},
    [ordered]@{Name='build';Path=(Join-Path $fixtureReleaseRoot 'build.ps1');Marker='verifier_ldflags='}
)

& docker image inspect $image *> $null
if ($LASTEXITCODE -ne 0) { throw 'Pinned image is required for the CRLF shell-payload regression' }

foreach ($case in $payloadCases) {
    $rawPayload = Get-EmbeddedPayload -Path $case.Path -Marker $case.Marker
    if ($rawPayload.IndexOf([char]13) -lt 0) {
        throw "$($case.Name) fixture did not reproduce CRLF shell-payload transport"
    }
    $rejected = $false
    try {
        Assert-LFPosixShellPayload -Payload $rawPayload
    } catch {
        $rejected = $true
    }
    if (!$rejected) { throw "$($case.Name) raw CRLF payload did not fail closed" }

    $normalized = ConvertTo-LFPosixShellPayload -Payload $rawPayload
    Assert-LFPosixShellPayload -Payload $normalized
    $dockerArguments = @(
        'run','--rm','--pull=never','--network','none','--read-only',
        '--cap-drop','ALL','--security-opt','no-new-privileges',
        '--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1',
        '--tmpfs','/work:rw,noexec,nosuid,nodev,size=16m,mode=0700',
        $image,'/bin/sh','-n','-c',$normalized
    )
    & docker @dockerArguments
    if ($LASTEXITCODE -ne 0) { throw "$($case.Name) normalized payload failed pinned offline shell parsing" }
}

$normalizationCounts = [ordered]@{
    'cache-canary.ps1' = 1
    'acquire.ps1' = 1
    'build.ps1' = 2
}
foreach ($entry in $normalizationCounts.GetEnumerator()) {
    $content = [IO.File]::ReadAllText((Join-Path $fixtureReleaseRoot $entry.Key))
    $count = ([regex]::Matches($content, 'ConvertTo-LFPosixShellPayload\s+-Payload')).Count
    if ($count -ne $entry.Value) {
        throw "$($entry.Key) does not normalize every Docker POSIX shell payload immediately before invocation"
    }
}

. (Join-Path $PSScriptRoot 'cache-canary.ps1')
$positiveCache = Join-Path $testRoot 'positive-cache'
$readOnlyCache = Join-Path $testRoot 'read-only-cache'
New-Item -ItemType Directory -Path $positiveCache,$readOnlyCache | Out-Null
Invoke-ReleaseCacheCanary -Image $image -ModuleCache $positiveCache
Require-Empty $positiveCache

$readOnlyRejected = $false
try {
    Invoke-ReleaseCacheCanary -Image $image -ModuleCache $readOnlyCache -ReadOnlyCache
} catch {
    $readOnlyRejected = $true
}
if (!$readOnlyRejected) { throw 'CRLF read-only cache unexpectedly passed the pre-acquisition canary' }
Require-Empty $readOnlyCache

$completeBuild = 'SKIPPED'
if ([string]::IsNullOrWhiteSpace($AcquisitionDirectory) -xor [string]::IsNullOrWhiteSpace($BuildOutputDirectory)) {
    throw 'AcquisitionDirectory and BuildOutputDirectory must be supplied together'
}
if (![string]::IsNullOrWhiteSpace($AcquisitionDirectory)) {
    $buildLauncher = Join-Path $PSScriptRoot 'invoke-exact-build.ps1'
    & pwsh -NoProfile -File $buildLauncher -RepositoryRoot $fixtureRoot -ExpectedToolingRevision $SourceRevision -AcquisitionDirectory $AcquisitionDirectory -OutputDirectory $BuildOutputDirectory -AllowCanonicalEolProjection
    if ($LASTEXITCODE -ne 0) { throw 'Complete offline release build from the CRLF checkout failed' }
    $dist = Join-Path ([IO.Path]::GetFullPath($BuildOutputDirectory)) 'dist'
    $testSummary = Get-Content -LiteralPath (Join-Path $dist 'TEST-SUMMARY.json') -Raw | ConvertFrom-Json
    $manifest = Get-Content -LiteralPath (Join-Path $dist 'release-manifest.json') -Raw | ConvertFrom-Json
    if ($testSummary.commands -notcontains "go test -p=1 -count=1 -run '^TestPinnedRuleAndCoverageIntegrityBindings$' ./tests/acceptance/gitleaks" -or
        $manifest.releaseTooling.commit -ne $SourceRevision -or $manifest.releaseTooling.tree -ne $sourceTrust.Tree) {
        throw 'CRLF build did not prove the pinned integrity test and exact tooling identity'
    }
    $completeBuild = 'PASS'
}

Write-Output "CRLF checkout regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT pinned-asset-integrity=$completeBuild complete-build=$completeBuild"
