param(
    [string]$WorkingDirectory,
    [string]$AcquisitionDirectory,
    [string]$BuildOutputDirectory,
    [string]$SourceRepository,
    [string]$SourceRevision,
    [ValidateSet('HostOnly','Container','Build','All')][string]$Phase = 'All',
    [ValidateSet('Candidate','Validation')][string]$Mode='Candidate'
)

$ErrorActionPreference = 'Stop'
$image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
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
    $WorkingDirectory = Join-Path ([IO.Path]::GetTempPath()) ("pscan-06-c2-crlf-" + [guid]::NewGuid().ToString('N'))
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
$dockerExecutionFixture = Join-Path $fixtureReleaseRoot 'docker-execution.ps1'
$payloadCases = @(
    [ordered]@{Name='cache-canary';Path=$dockerExecutionFixture;Marker='PSCAN-06-C2-CONTAINER-CACHE-CANARY'},
    [ordered]@{Name='acquisition';Path=$dockerExecutionFixture;Marker='/work/runner/go/bin/go mod download'},
    [ordered]@{Name='build';Path=$dockerExecutionFixture;Marker='verifier_ldflags='}
)

$normalizedCases = [Collections.Generic.List[object]]::new()
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
    $normalizedCases.Add([pscustomobject]@{Name=$case.Name;Payload=$normalized})
}

if ($Phase -eq 'HostOnly') {
    Write-Output 'CRLF host-only regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT docker=NOT_INVOKED'
    return
}

$engineJson = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation EngineInspection
$engine = ($engineJson -join "`n") | ConvertFrom-Json
if ($engine.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($engine.StdErr) -or !$engine.ContainmentEmpty) { throw 'Docker engine boundary is untrusted' }
. (Join-Path $PSScriptRoot 'image-admission.ps1')
[void](Assert-ReleaseEngineEvidence -Json $engine.StdOut)
$dockerSHA256 = $engine.DockerSHA256
foreach ($case in $normalizedCases) {
    $kind = switch ($case.Name) { 'cache-canary' {'CacheCanary'} 'acquisition' {'Acquisition'} 'build' {'Build'} }
    $parseJson = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation ContainerCrlfParse -PayloadKind $kind -ExpectedDockerSHA256 $dockerSHA256
    $parse = ($parseJson -join "`n") | ConvertFrom-Json
    if ($parse.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($parse.StdErr) -or !$parse.ContainmentEmpty -or $parse.DockerSHA256 -cne $dockerSHA256) { throw "$($case.Name) normalized payload failed pinned offline shell parsing" }
}

$boundaryContent = [IO.File]::ReadAllText($dockerExecutionFixture)
if (([regex]::Matches($boundaryContent, 'function Get-(Cache|Acquisition|Build|Package)Payload')).Count -ne 4) { throw 'Closed Docker operation table does not own every POSIX payload' }

. (Join-Path $PSScriptRoot 'cache-canary.ps1')
$positiveCache = Join-Path $testRoot 'positive-cache'
$readOnlyCache = Join-Path $testRoot 'read-only-cache'
New-Item -ItemType Directory -Path $positiveCache,$readOnlyCache | Out-Null
. (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
$cacheIdentity = Get-HostCacheIdentity -CacheDirectory $positiveCache
$cacheUser = @{}
if ($null -ne $cacheIdentity.HostUID) { $cacheUser.HostUID=$cacheIdentity.HostUID; $cacheUser.HostGID=$cacheIdentity.HostGID }
Invoke-ReleaseCacheCanary -Image $image -ModuleCache $positiveCache -ExpectedDockerSHA256 $dockerSHA256 @cacheUser
Require-Empty $positiveCache

$readOnlyRejected = $false
try {
    Invoke-ReleaseCacheCanary -Image $image -ModuleCache $readOnlyCache -ExpectedDockerSHA256 $dockerSHA256 -ReadOnlyCache @cacheUser
} catch {
    if($_.FullyQualifiedErrorId-notlike'PSCAN_CACHE_EXPECTED_DENIAL*'){throw}
    Write-Host ('PSCAN_CACHE_NEGATIVE '+($_.TargetObject|ConvertTo-Json -Compress))
    $readOnlyRejected = $true
}
if (!$readOnlyRejected) { throw 'CRLF read-only cache unexpectedly passed the pre-acquisition canary' }
Require-Empty $readOnlyCache

$completeBuild = 'SKIPPED'
if ([string]::IsNullOrWhiteSpace($AcquisitionDirectory) -xor [string]::IsNullOrWhiteSpace($BuildOutputDirectory)) {
    throw 'AcquisitionDirectory and BuildOutputDirectory must be supplied together'
}
if (![string]::IsNullOrWhiteSpace($AcquisitionDirectory)) {
    if ($Phase -notin @('Build','All')) { throw 'A complete build is only valid in the Build or All phase' }
    $buildLauncher = Join-Path $PSScriptRoot 'invoke-exact-build.ps1'
    & pwsh -NoProfile -File $buildLauncher -RepositoryRoot $fixtureRoot -ExpectedToolingRevision $SourceRevision -AcquisitionDirectory $AcquisitionDirectory -OutputDirectory $BuildOutputDirectory -AllowCanonicalEolProjection -Mode $Mode
    if ($LASTEXITCODE -ne 0) { throw 'Complete offline release build from the CRLF checkout failed' }
    $dist = Join-Path ([IO.Path]::GetFullPath($BuildOutputDirectory)) 'dist'
    $testSummary = Get-Content -LiteralPath (Join-Path $dist 'TEST-SUMMARY.json') -Raw | ConvertFrom-Json
    $manifest = if($Mode-ceq'Validation'){
        if(Test-Path -LiteralPath (Join-Path $dist 'release-manifest.json')){throw 'Validation emitted a release trust manifest'}
        $validation=Get-Content -LiteralPath (Join-Path $dist 'BUILD-VALIDATION.json') -Raw|ConvertFrom-Json
        if($validation.schema-cne'pscan-build-validation-v1'-or$validation.candidate-ne$false){throw 'Validation output identity is invalid'}
        [pscustomobject]@{releaseTooling=$validation.source}
    }else{Get-Content -LiteralPath (Join-Path $dist 'release-manifest.json') -Raw | ConvertFrom-Json}
    if ($testSummary.commands -notcontains "go test -p=1 -count=1 -run '^TestPinnedRuleAndCoverageIntegrityBindings$' ./tests/acceptance/gitleaks" -or
        $manifest.releaseTooling.commit -ne $SourceRevision -or $manifest.releaseTooling.tree -ne $sourceTrust.Tree) {
        throw 'CRLF build did not prove the pinned integrity test and exact tooling identity'
    }
    $completeBuild = 'PASS'
}

Write-Output "CRLF checkout regression PASS checkout-asset-CR=PROVED raw-shell-CR=REJECT normalized-shell-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT pinned-asset-integrity=$completeBuild complete-build=$completeBuild"
