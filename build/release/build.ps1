param(
    [Parameter(Mandatory = $true)][string]$RepositoryRoot,
    [Parameter(Mandatory = $true)][string]$ExpectedToolingRevision,
    [Parameter(Mandatory = $true)][string]$AcquisitionDirectory,
    [Parameter(Mandatory = $true)][string]$OutputDirectory,
    [switch]$AllowCanonicalEolProjection,
    [ValidateSet('Candidate','Validation')][string]$Mode='Candidate'
)

$ErrorActionPreference = 'Stop'
$image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$root = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $RepositoryRoot).Path)
$output = [IO.Path]::GetFullPath($OutputDirectory)

function Require-Digest([string]$Path, [string]$Expected) {
    if (!(Test-Path -LiteralPath $Path -PathType Leaf)) { throw "Required file missing: $Path" }
    $actual = (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $Expected) { throw "Digest mismatch: $Path" }
}

function Write-Utf8([string]$Path, [string]$Value) {
    [IO.File]::WriteAllText($Path, $Value.Replace("`r`n", "`n") + $(if ($Value.EndsWith("`n")) { '' } else { "`n" }), [Text.UTF8Encoding]::new($false))
}

function Get-GitBlobSha256([string]$Repository, [string]$ObjectID) {
    if ($ObjectID -notmatch '^[0-9a-f]{40}$') { throw "Unsupported Git blob identity: $ObjectID" }
    $bytes = (Invoke-SourceTrustGit -Repository $Repository -Arguments @('cat-file','blob',$ObjectID)).Bytes
    $sha256 = [Security.Cryptography.SHA256]::Create()
    try {
        $digest = $sha256.ComputeHash($bytes)
    } finally {
        $sha256.Dispose()
    }
    return ([BitConverter]::ToString($digest)).Replace('-','').ToLowerInvariant()
}

function New-ExactGitTreeArchive(
    [string]$Repository,
    [string]$Revision,
    [string]$ExpectedTree,
    [string]$ArchivePath,
    [string]$Label
) {
    $objectType = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $Repository -Arguments @('cat-file','-t',$Revision)).Bytes -Label "$Label object type").Trim()
    $resolvedTree = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $Repository -Arguments @('rev-parse',"$Revision`^{tree}")).Bytes -Label "$Label tree").Trim()
    if ($objectType -ne 'commit' -or $resolvedTree -ne $ExpectedTree) {
        throw "$Label commit/tree identity cannot be proven"
    }

    $entries = @(Get-SourceTrustTreeEntries -Repository $Repository -Tree $ExpectedTree)
    if ($entries.Count -eq 0) { throw "$Label Git tree enumeration failed closed" }
    $hashLines = [Collections.Generic.List[string]]::new()
    $pathLines = [Collections.Generic.List[string]]::new()
    $modeLines = [Collections.Generic.List[string]]::new()
    foreach ($entry in $entries) {
        $path = $entry.Path
        $digest = Get-GitBlobSha256 -Repository $Repository -ObjectID $entry.Object
        $hashLines.Add("$digest  $path")
        $pathLines.Add($path)
        $modeLines.Add("$($entry.Mode)`t$path")
    }
    $paths = $pathLines.ToArray()
    [Array]::Sort($paths, [StringComparer]::Ordinal)
    $hashes = $hashLines.ToArray()
    [Array]::Sort($hashes, [StringComparer]::Ordinal)
    $modes = $modeLines.ToArray()
    [Array]::Sort($modes, [StringComparer]::Ordinal)

    $hashManifest = "$ArchivePath.blobs.sha256"
    $pathManifest = "$ArchivePath.paths"
    $modeManifest = "$ArchivePath.modes"
    Write-Utf8 $hashManifest ($hashes -join "`n")
    Write-Utf8 $pathManifest ($paths -join "`n")
    Write-Utf8 $modeManifest ($modes -join "`n")

    [void](Invoke-SourceTrustGit -Repository $Repository -Arguments @('-c','tar.umask=0022','archive','--format=tar','--output',$ArchivePath,$Revision))
    if (!(Test-Path -LiteralPath $ArchivePath -PathType Leaf) -or (Get-Item -LiteralPath $ArchivePath).Length -eq 0) {
        throw "$Label exact-tree archive materialization failed"
    }
    return [pscustomobject]@{
        Archive = $ArchivePath
        ArchiveSHA256 = (Get-FileHash -LiteralPath $ArchivePath -Algorithm SHA256).Hash.ToLowerInvariant()
        BlobManifest = $hashManifest
        PathManifest = $pathManifest
        ModeManifest = $modeManifest
        FileCount = $entries.Count
    }
}

function Copy-ReleaseFile([string]$Source, [string]$Name) {
    $destination = Join-Path $dist $Name
    Copy-Item -LiteralPath $Source -Destination $destination
    return $destination
}

if ($env:PSCAN_TRUSTED_LAUNCHER_REVISION -ne $ExpectedToolingRevision -or $env:PSCAN_TRUSTED_LAUNCHER_TREE -notmatch '^[0-9a-f]{40}$') {
    throw 'Release build must enter through the exact committed launcher'
}
. (Join-Path $PSScriptRoot 'source-trust.ps1')
$sourceTrust = Assert-ExactGitSourceTrust -Repository $root -ExpectedRevision $ExpectedToolingRevision -AllowCanonicalEolProjection:$AllowCanonicalEolProjection
if ($sourceTrust.Tree -ne $env:PSCAN_TRUSTED_LAUNCHER_TREE) { throw 'Launcher and build source-tree identities conflict' }
$toolingRevision = $sourceTrust.Commit
$toolingTree = $sourceTrust.Tree
$createdText = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('show','-s','--format=%cI',$toolingRevision)).Bytes -Label 'Tooling commit timestamp').Trim()
$created = ([DateTimeOffset]::Parse($createdText)).UtcDateTime.ToString('yyyy-MM-ddTHH:mm:ssZ')
$epoch = [DateTimeOffset]::Parse($created).ToUnixTimeSeconds()
$productTag = 'v1.0.0'
$productRevision = 'a13c28fe7273bc8dc6545f97966a02889524eb4c'
$productTree = '217b711ddea51fd0ea7e808edd2e27fdecef8427'
$toolingTag = 'release-tooling-v1.0.0-c2-linux-build-v2'
$workflow = '.github/workflows/release-build-unsigned.yml'
$workflowRef = 'refs/tags/release-tooling-v1.0.0-c2-linux-build-v2'
. (Join-Path $PSScriptRoot 'build-validation.ps1')
$validationInvocation=$null
if($Mode-ceq'Validation'){
    $validationInvocation=Assert-BuildValidationInvocation $toolingRevision
    $toolingTag=$null;$workflow=$validationInvocation.workflow;$workflowRef=$validationInvocation.ref
}else{
    if(![string]::IsNullOrEmpty($env:PSCAN_VALIDATION_WORKFLOW_SHA)){throw 'Validation workflow cannot enter candidate mode'}
    $resolvedToolingTag = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse',"$toolingTag`^{commit}")).Bytes -Label 'Immutable tooling commit').Trim()
    if ($resolvedToolingTag -cne $toolingRevision) { throw 'Build tag does not identify the exact corrected tooling commit' }
}
$resolvedProductTag = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse','v1.0.0^{commit}')).Bytes -Label 'Locked product commit').Trim()
$resolvedProductTree = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse','v1.0.0^{tree}')).Bytes -Label 'Locked product tree').Trim()
if ($resolvedProductTag -ne $productRevision -or $resolvedProductTree -ne $productTree) {
    throw 'Locked product tag, commit or tree identity does not match Correction C2 authority'
}
$acquisition = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $AcquisitionDirectory).Path)

if (Test-Path -LiteralPath $output) {
    if (Get-ChildItem -LiteralPath $output -Force | Select-Object -First 1) { throw 'Release output must be a new empty directory' }
} else {
    New-Item -ItemType Directory -Path $output | Out-Null
}
$raw = Join-Path $output '.raw'
$linuxStage = Join-Path $output '.stage-linux'
$windowsStage = Join-Path $output '.stage-windows'
$dist = Join-Path $output 'dist'
New-Item -ItemType Directory -Path $raw,$linuxStage,$windowsStage,$dist | Out-Null
$productArchive = Join-Path $raw 'product-source.tar'
$toolingArchive = Join-Path $raw 'release-tooling-source.tar'
$productMaterialization = New-ExactGitTreeArchive -Repository $root -Revision $productRevision -ExpectedTree $productTree -ArchivePath $productArchive -Label 'Locked product source'
$toolingMaterialization = New-ExactGitTreeArchive -Repository $root -Revision $toolingRevision -ExpectedTree $toolingTree -ArchivePath $toolingArchive -Label 'Correction tooling source'

$runnerGo = Join-Path $acquisition 'downloads\go1.27.1.linux-amd64.tar.gz'
$engineGo = Join-Path $acquisition 'downloads\go1.27.0.linux-amd64.tar.gz'
$gitleaksSource = Join-Path $acquisition 'downloads\gitleaks-83d9cd684c87d95d656c1458ef04895a7f1cbd8e.tar.gz'
$moduleCache = Join-Path $acquisition 'gomodcache'
$ledger = Join-Path $acquisition 'acquisition-ledger.json'
Require-Digest $runnerGo '63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445'
Require-Digest $engineGo '675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685'
Require-Digest $gitleaksSource '6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115'
if (!(Test-Path -LiteralPath $moduleCache -PathType Container) -or !(Test-Path -LiteralPath $ledger -PathType Leaf)) {
    throw 'Complete acquisition evidence is missing'
}
$ledgerValue = Get-Content -Raw -LiteralPath $ledger | ConvertFrom-Json
if ($ledgerValue.schemaVersion -ne '2.1' -or $ledgerValue.cacheCanary.semantics -ne 'host-and-container-write-atomic-rename-read-delete' -or
    $ledgerValue.cacheCanary.completedBeforeNetworkDependencyAcquisition -ne $true -or
    $ledgerValue.cacheCanary.completedBeforeImagePull -ne $true -or
    $ledgerValue.imageAdmission.canonicalReference -ne $image -or
    $ledgerValue.imageAdmission.executionBoundary -cne 'closed-private-contained' -or
    $ledgerValue.imageAdmission.postAdmissionRepoDigestProved -ne $true) {
    throw 'Acquisition ledger does not prove the Correction C2 host, image and container admission sequence'
}

$receiptPath = Join-Path $acquisition 'docker-admission.json'
if (!(Test-Path -LiteralPath $receiptPath -PathType Leaf)) { throw 'Release build requires the prior closed Docker admission receipt' }
. (Join-Path $PSScriptRoot 'execution-profile.ps1')
$imageAdmission = Read-ReleaseImageAdmission $receiptPath
if ($imageAdmission.sourceRevision -cne $toolingRevision -or $imageAdmission.image -cne $image -or
    $imageAdmission.dockerExecutableSHA256 -notmatch '^[0-9a-f]{64}$' -or $imageAdmission.containment -cne 'empty-after-every-operation') { throw 'Docker admission receipt does not bind this exact build' }
# Carry the admitted ordinary host identity through every writable container.
. (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
$hostIdentity = Get-HostCacheIdentity -CacheDirectory $moduleCache
$containerUser = @{}
if ($IsLinux) {
    foreach ($directory in @($raw,$linuxStage,$windowsStage,$dist)) {
        $identity = Get-HostCacheIdentity -CacheDirectory $directory
        if ($identity.HostUID -ne $hostIdentity.HostUID -or $identity.HostGID -ne $hostIdentity.HostGID) { throw 'Build output ownership differs from the admitted host identity' }
        & chmod 0700 -- $directory
        if ($LASTEXITCODE -ne 0) { throw 'Private build directory mode preparation failed' }
    }
    $containerUser.HostUID=$hostIdentity.HostUID; $containerUser.HostGID=$hostIdentity.HostGID
}
$expectedHostMode = $hostIdentity.IdentityMode
if ($imageAdmission.hostIdentityMode -cne $expectedHostMode -or $imageAdmission.hostUID -ne $hostIdentity.HostUID -or $imageAdmission.hostGID -ne $hostIdentity.HostGID -or
    $ledgerValue.cacheCanary.hostIdentityMode -cne $expectedHostMode -or $ledgerValue.cacheCanary.hostUID -ne $hostIdentity.HostUID -or $ledgerValue.cacheCanary.hostGID -ne $hostIdentity.HostGID) { throw 'Build host identity differs from admission and acquisition evidence' }
$dockerSHA256 = $imageAdmission.dockerExecutableSHA256
if ($ledgerValue.imageAdmission.dockerExecutableSHA256 -cne $dockerSHA256) { throw 'Acquisition ledger and Docker admission receipt identities conflict' }
. (Join-Path $PSScriptRoot 'image-admission.ps1')
$inspectJson = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation RepositoryDigestInspection -ExpectedDockerSHA256 $dockerSHA256
$inspect = ($inspectJson -join "`n") | ConvertFrom-Json
if ($inspect.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($inspect.StdErr) -or !$inspect.ContainmentEmpty) { throw 'Pre-build image inspection failed closed' }
$repoDigests = Read-ReleaseRepoDigestsEvidence -Json $inspect.StdOut
[void](Assert-ReleaseImageIdentityEvidence -Image $image -RepoDigests $repoDigests)

$buildBoundary = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation ReleaseBuild @containerUser -ExpectedDockerSHA256 $dockerSHA256 -CacheDirectory $moduleCache -RunnerGoArchive $runnerGo -EngineGoArchive $engineGo -GitleaksArchive $gitleaksSource -ProductArchive $productArchive -ProductBlobManifest $productMaterialization.BlobManifest -ProductPathManifest $productMaterialization.PathManifest -ProductModeManifest $productMaterialization.ModeManifest -ToolingArchive $toolingArchive -ToolingBlobManifest $toolingMaterialization.BlobManifest -ToolingPathManifest $toolingMaterialization.PathManifest -ToolingModeManifest $toolingMaterialization.ModeManifest -RawOutputDirectory $raw -SourceDateEpoch ([string]$epoch) -ProductRevision $productRevision -ToolingRevision $toolingRevision -ToolingTree $toolingTree -ProductArchiveSHA256 $productMaterialization.ArchiveSHA256 -ToolingArchiveSHA256 $toolingMaterialization.ArchiveSHA256 -ProductFileCount $productMaterialization.FileCount -ToolingFileCount $toolingMaterialization.FileCount -Created $created
$buildResult = ($buildBoundary -join "`n") | ConvertFrom-Json
if ($buildResult.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($buildResult.StdErr) -or !$buildResult.ContainmentEmpty -or $buildResult.DockerSHA256 -cne $dockerSHA256) { throw 'Offline release build or validation failed closed at the Docker boundary' }

Require-Digest (Join-Path $raw 'gitleaks-linux-amd64') '657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586'
Require-Digest (Join-Path $raw 'gitleaks-windows-amd64.exe') 'b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178'
Require-Digest (Join-Path $raw 'scanner-runner-linux-amd64') '06043e9410e05a927356960f7b4ee64c9a073b7467e775452528a6813d97906e'
Require-Digest (Join-Path $raw 'scanner-runner-windows-amd64.exe') '1b74310e16e0df13cfd077c41338b93eabdc3d5959fb63dc256927ca51117543'

Copy-ReleaseFile (Join-Path $raw 'scanner-runner-linux-amd64') 'scanner-runner-linux-amd64' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'scanner-runner-windows-amd64.exe') 'scanner-runner-windows-amd64.exe' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'scanner-release-verifier-linux-amd64') 'scanner-release-verifier-linux-amd64' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'scanner-release-verifier-windows-amd64.exe') 'scanner-release-verifier-windows-amd64.exe' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'gitleaks-linux-amd64') 'gitleaks-linux-amd64' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'gitleaks-windows-amd64.exe') 'gitleaks-windows-amd64.exe' | Out-Null
$productFiles = Join-Path $raw 'product'
$toolingFiles = Join-Path $raw 'tooling-materialized'
Copy-ReleaseFile (Join-Path $productFiles 'rules\generic\gitleaks-v8.30.1.toml') 'rules-gitleaks-v8.30.1.toml' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'rules\generic\gitleaks-ignore-empty-v1.txt') 'rules-gitleaks-ignore-empty-v1.txt' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\scan-request\schema-1.1.json') 'schema-scan-request-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\scan-outcome\schema-1.0.json') 'schema-scan-outcome-1.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\release-manifest\schema-1.1.json') 'schema-release-manifest-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'contracts\release-manifest\schema-2.0.json') 'schema-release-manifest-2.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'contracts\release-manifest\schema-2.1.json') 'schema-release-manifest-2.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'contracts\release-manifest\schema-2.2.json') 'schema-release-manifest-2.2.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'contracts\release-manifest\schema-2.3.json') 'schema-release-manifest-2.3.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'contracts\release-manifest\schema-2.4.json') 'schema-release-manifest-2.4.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\global-revocation\schema-1.1.json') 'schema-global-revocation-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\rule-pack\schema-1.0.json') 'schema-rule-pack-1.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'LICENSE') 'LICENSE.txt' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'THIRD_PARTY_NOTICES.md') 'THIRD_PARTY_NOTICES.md' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'licenses\gitleaks\modules\manifest.json') 'GITLEAKS-LICENCE-MANIFEST.json' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'sbom.spdx.json') 'sbom.spdx.json' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'docs\release\OFFLINE-VERIFICATION-RUNBOOK.md') 'OFFLINE-VERIFICATION-RUNBOOK.md' | Out-Null
Copy-ReleaseFile (Join-Path $toolingFiles 'docs\release\SCANNER-IO-REFERENCE.md') 'SCANNER-IO-REFERENCE.md' | Out-Null

$testSummary = [ordered]@{
    schemaVersion = '2.4'; productSourceRevision = $productRevision; releaseToolingRevision = $toolingRevision; createdAt = $created
    authoritativeEnvironment = 'pinned-network-disabled-linux-container'
    goToolchain = 'go1.27.1'; engineToolchain = 'go1.27.0'
    commands = @("go test -p=1 -count=1 -run '^TestPinnedRuleAndCoverageIntegrityBindings$' ./tests/acceptance/gitleaks", 'go test -p=1 -count=1 ./...', 'go vet -p=1 ./...', 'GOOS=windows GOARCH=amd64 go test -p=1 -exec /bin/true ./...')
    linuxExecution = 'PASS'; windowsCompilation = 'PASS'; windowsNativeExecution = 'UNPROVEN_SMART_APP_CONTROL'
    signing = 'NOT_PERFORMED_OWNER_GATE'; remoteWorkflow = 'UNSIGNED_BUILD_ONLY'
}
if($Mode-ceq'Validation'){$testSummary.schemaVersion='pscan-validation-test-summary-v1';$testSummary.remoteWorkflow='VALIDATION_ONLY';$testSummary.candidate=$false}
Write-Utf8 (Join-Path $dist 'TEST-SUMMARY.json') ($testSummary | ConvertTo-Json -Depth 6)
$compatibility = [ordered]@{
    schemaVersion = '2.4'; releaseVersion = 'v1.0.0'; productSourceRevision = $productRevision; releaseToolingRevision = $toolingRevision
    platforms = @([ordered]@{os='linux';arch='amd64'},[ordered]@{os='windows';arch='amd64'})
    requestSchemas = @('1.0','1.1'); outcomeSchemas = @('1.0'); releaseManifestSchemas = @('1.0','1.1','2.0','2.1','2.2','2.3','2.4')
    engine = [ordered]@{name='gitleaks';version='8.30.1';adapterVersion='2.0.0'}
}
Write-Utf8 (Join-Path $dist 'COMPATIBILITY.json') ($compatibility | ConvertTo-Json -Depth 6)
$limitations = @"
# v1.0.0 limitations

- Native Windows amd64 execution is unproven on the author host because Windows Smart App Control blocks unsigned locally built executables. Windows amd64 cross-compilation and byte verification pass; this is not substituted for native execution.
- This candidate was built by the exact unsigned build workflow recorded in BUILD-PROVENANCE.json. It is not a trusted signed release. Signing, artifact attestation, draft upload and publication remain separately gated.
- The offline verifier requires separately acquired Cosign v3.1.3 at its documented exact SHA-256 and rejects absent, stale, conflicting or untrusted evidence.
- Gitleaks is the only primary detector. TruffleHog is not assessed, downloaded, integrated, distributed or enabled.
- Passing validation placeholders do not establish any consuming-project integration, policy, receipt, deployment or production result.
"@
if($Mode-ceq'Validation'){$limitations=$limitations.Replace('This candidate was built by the exact unsigned build workflow','This validation-only package was built by the exact validation workflow');$limitations="# Validation only; not a release candidate`n`n"+$limitations}
Write-Utf8 (Join-Path $dist 'LIMITATIONS.md') $limitations
$revocationLocation = 'https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100'
$revocations = [ordered]@{schemaFamily='global-scanner-revocation-snapshot';schemaVersion='1.0';capturedAt=$created;records=@()}
Write-Utf8 (Join-Path $dist 'global-revocations.json') ($revocations | ConvertTo-Json -Depth 8)
$checkpoint = [ordered]@{schemaFamily='global-scanner-revocation-checkpoint';schemaVersion='1.0';sequence=0;digest=$null;capturedAt=$created;discoveryLocation=$revocationLocation}
Write-Utf8 (Join-Path $dist 'global-revocation-checkpoint.json') ($checkpoint | ConvertTo-Json -Depth 6)
$provenance = [ordered]@{
    schemaVersion='2.4';createdAt=$created
    productSource=[ordered]@{tag=$productTag;commit=$productRevision;tree=$productTree}
    releaseTooling=[ordered]@{tag=$toolingTag;commit=$toolingRevision;tree=$toolingTree;workflow=$workflow;workflowRef=$workflowRef;workflowSha=$toolingRevision;trigger='workflow_dispatch'}
    sourceTrust=[ordered]@{trackedFiles=$sourceTrust.FileCount;rawEqual=$sourceTrust.RawEqualCount;canonicalCrlfProjection=$sourceTrust.CanonicalEolProjectionCount;workingTreeInputsUsed=$false;buildDriver='exact-git-object-materialization'}
    buildImage=$image;runnerGo='go1.27.1';engineGo='go1.27.0'
    gitleaksSourceRevision='83d9cd684c87d95d656c1458ef04895a7f1cbd8e'
    network='disabled';moduleCache='read-only';cgo=$false;trimpath=$true;buildVCS=$false;buildId='empty'
}
if($Mode-ceq'Validation'){
    $provenance.Remove('releaseTooling');$provenance.schemaVersion='pscan-validation-provenance-v1'
    $provenance.purpose='validation-only';$provenance.candidate=$false
    $provenance.source=[ordered]@{commit=$toolingRevision;tree=$toolingTree};$provenance.invocation=$validationInvocation
}
Write-Utf8 (Join-Path $dist 'BUILD-PROVENANCE.json') ($provenance | ConvertTo-Json -Depth 6)

foreach ($stage in @($linuxStage,$windowsStage)) {
    New-Item -ItemType Directory -Path (Join-Path $stage 'bin'),(Join-Path $stage 'rules'),(Join-Path $stage 'schemas'),(Join-Path $stage 'licenses\gitleaks') -Force | Out-Null
    Copy-Item -LiteralPath (Join-Path $dist 'rules-gitleaks-v8.30.1.toml') -Destination (Join-Path $stage 'rules\gitleaks-v8.30.1.toml')
    Copy-Item -LiteralPath (Join-Path $dist 'rules-gitleaks-ignore-empty-v1.txt') -Destination (Join-Path $stage 'rules\gitleaks-ignore-empty-v1.txt')
    Copy-Item -Path (Join-Path $dist 'schema-*.json') -Destination (Join-Path $stage 'schemas')
    Copy-Item -Path (Join-Path $productFiles 'licenses\gitleaks\*') -Destination (Join-Path $stage 'licenses\gitleaks') -Recurse
    Copy-Item -LiteralPath (Join-Path $dist 'LICENSE.txt'),(Join-Path $dist 'THIRD_PARTY_NOTICES.md'),(Join-Path $dist 'sbom.spdx.json'),(Join-Path $dist 'COMPATIBILITY.json'),(Join-Path $dist 'LIMITATIONS.md'),(Join-Path $dist 'OFFLINE-VERIFICATION-RUNBOOK.md'),(Join-Path $dist 'SCANNER-IO-REFERENCE.md') -Destination $stage
}
Copy-Item -LiteralPath (Join-Path $dist 'scanner-runner-linux-amd64') -Destination (Join-Path $linuxStage 'bin\scanner-runner')
Copy-Item -LiteralPath (Join-Path $dist 'scanner-release-verifier-linux-amd64') -Destination (Join-Path $linuxStage 'bin\scanner-release-verifier')
Copy-Item -LiteralPath (Join-Path $dist 'gitleaks-linux-amd64') -Destination (Join-Path $linuxStage 'bin\gitleaks')
Copy-Item -LiteralPath (Join-Path $dist 'scanner-runner-windows-amd64.exe') -Destination (Join-Path $windowsStage 'bin\scanner-runner.exe')
Copy-Item -LiteralPath (Join-Path $dist 'scanner-release-verifier-windows-amd64.exe') -Destination (Join-Path $windowsStage 'bin\scanner-release-verifier.exe')
Copy-Item -LiteralPath (Join-Path $dist 'gitleaks-windows-amd64.exe') -Destination (Join-Path $windowsStage 'bin\gitleaks.exe')

$packageBoundary = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') -Operation ReleasePackage @containerUser -ExpectedDockerSHA256 $dockerSHA256 -RawOutputDirectory $raw -LinuxStageDirectory $linuxStage -WindowsStageDirectory $windowsStage -DistributionDirectory $dist -SourceDateEpoch ([string]$epoch)
$packageResult = ($packageBoundary -join "`n") | ConvertFrom-Json
if ($packageResult.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($packageResult.StdErr) -or !$packageResult.ContainmentEmpty -or $packageResult.DockerSHA256 -cne $dockerSHA256) { throw 'Deterministic packaging failed closed at the Docker boundary' }

$checksumTargets = Get-ChildItem -LiteralPath $dist -File | Sort-Object Name
$checksumLines = foreach ($file in $checksumTargets) {
    "$((Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant())  $($file.Name)"
}
Write-Utf8 (Join-Path $dist 'CHECKSUMS.sha256') ($checksumLines -join "`n")
if($Mode-ceq'Validation'){
    Complete-BuildValidation -Distribution $dist -SourceTrust $sourceTrust -Invocation $validationInvocation -Created $created
    return
}

function Asset([string]$Name, [string]$Kind, [string]$OS='none', [string]$Arch='none') {
    $file = Get-Item -LiteralPath (Join-Path $dist $Name)
    return [ordered]@{path=$Name;kind=$Kind;os=$OS;arch=$Arch;size=$file.Length;sha256=(Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()}
}
$assets = @(
    (Asset 'scanner-runner-linux-amd64' 'runner' 'linux' 'amd64'),
    (Asset 'scanner-runner-windows-amd64.exe' 'runner' 'windows' 'amd64'),
    (Asset 'scanner-release-verifier-linux-amd64' 'verifier' 'linux' 'amd64'),
    (Asset 'scanner-release-verifier-windows-amd64.exe' 'verifier' 'windows' 'amd64'),
    (Asset 'gitleaks-linux-amd64' 'engine' 'linux' 'amd64'),
    (Asset 'gitleaks-windows-amd64.exe' 'engine' 'windows' 'amd64'),
    (Asset 'project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz' 'platform-bundle' 'linux' 'amd64'),
    (Asset 'project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip' 'platform-bundle' 'windows' 'amd64'),
    (Asset 'rules-gitleaks-v8.30.1.toml' 'rules'),
    (Asset 'rules-gitleaks-ignore-empty-v1.txt' 'rules'),
    (Asset 'schema-scan-request-1.1.json' 'schema'),
    (Asset 'schema-scan-outcome-1.0.json' 'schema'),
    (Asset 'schema-release-manifest-1.1.json' 'schema'),
    (Asset 'schema-release-manifest-2.0.json' 'schema'),
    (Asset 'schema-release-manifest-2.1.json' 'schema'),
    (Asset 'schema-release-manifest-2.2.json' 'schema'),
    (Asset 'schema-global-revocation-1.1.json' 'schema'),
    (Asset 'schema-rule-pack-1.0.json' 'schema'),
    (Asset 'LICENSE.txt' 'licence'),
    (Asset 'THIRD_PARTY_NOTICES.md' 'licence'),
    (Asset 'GITLEAKS-LICENCE-MANIFEST.json' 'licence-manifest'),
    (Asset 'sbom.spdx.json' 'sbom'),
    (Asset 'TEST-SUMMARY.json' 'test-summary'),
    (Asset 'LIMITATIONS.md' 'limitations'),
    (Asset 'COMPATIBILITY.json' 'compatibility'),
    (Asset 'global-revocations.json' 'revocation-snapshot'),
    (Asset 'global-revocation-checkpoint.json' 'revocation-checkpoint'),
    (Asset 'BUILD-PROVENANCE.json' 'documentation'),
    (Asset 'OFFLINE-VERIFICATION-RUNBOOK.md' 'documentation'),
    (Asset 'SCANNER-IO-REFERENCE.md' 'documentation'),
    (Asset 'schema-release-manifest-2.3.json' 'schema'),
    (Asset 'schema-release-manifest-2.4.json' 'schema'),
    (Asset 'CHECKSUMS.sha256' 'checksums')
)
$manifest = [ordered]@{
    schemaFamily='scanner-release-manifest';manifestSchemaVersion='2.4';releaseVersion='v1.0.0'
    productSource=[ordered]@{tag=$productTag;commit=$productRevision;tree=$productTree}
    releaseTooling=[ordered]@{tag=$toolingTag;commit=$toolingRevision;tree=$toolingTree;workflow=$workflow;workflowRef=$workflowRef;workflowSha=$toolingRevision;trigger='workflow_dispatch'}
    runnerVersion='1.0.0';goToolchainVersion='go1.27.1'
    runnerBindings=@(
        [ordered]@{os='linux';arch='amd64';path='scanner-runner-linux-amd64';sha256=$assets[0].sha256},
        [ordered]@{os='windows';arch='amd64';path='scanner-runner-windows-amd64.exe';sha256=$assets[1].sha256}
    )
    engineBindings=@(
        [ordered]@{name='gitleaks';version='8.30.1';sourceRevision='83d9cd684c87d95d656c1458ef04895a7f1cbd8e';os='linux';arch='amd64';path='gitleaks-linux-amd64';sha256='657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586'},
        [ordered]@{name='gitleaks';version='8.30.1';sourceRevision='83d9cd684c87d95d656c1458ef04895a7f1cbd8e';os='windows';arch='amd64';path='gitleaks-windows-amd64.exe';sha256='b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178'}
    )
    rulePack=[ordered]@{path='rules-gitleaks-v8.30.1.toml';sha256='cebfe007ae88a55540e42ccb8e042848d089fe2f59bcf604c5fd4904d6460792'}
    schemaBindings=@(
        [ordered]@{family='scan-request';version='1.1';path='schema-scan-request-1.1.json';sha256=$assets[10].sha256},
        [ordered]@{family='scan-outcome';version='1.0';path='schema-scan-outcome-1.0.json';sha256=$assets[11].sha256},
        [ordered]@{family='scanner-release-manifest';version='1.1';path='schema-release-manifest-1.1.json';sha256=$assets[12].sha256},
        [ordered]@{family='scanner-release-manifest';version='2.0';path='schema-release-manifest-2.0.json';sha256=$assets[13].sha256},
        [ordered]@{family='scanner-release-manifest';version='2.1';path='schema-release-manifest-2.1.json';sha256=$assets[14].sha256},
        [ordered]@{family='scanner-release-manifest';version='2.2';path='schema-release-manifest-2.2.json';sha256=$assets[15].sha256},
        [ordered]@{family='global-scanner-revocation';version='1.1';path='schema-global-revocation-1.1.json';sha256=$assets[16].sha256},
        [ordered]@{family='rule-pack';version='1.0';path='schema-rule-pack-1.0.json';sha256=$assets[17].sha256},
        [ordered]@{family='scanner-release-manifest';version='2.3';path='schema-release-manifest-2.3.json';sha256=$assets[30].sha256},
        [ordered]@{family='scanner-release-manifest';version='2.4';path='schema-release-manifest-2.4.json';sha256=$assets[31].sha256}
    )
    releaseState='unsigned-candidate'
    releaseIdentity=$null
    buildIdentity=[ordered]@{
        repository='ThameeraDananjaya/project-agnostic-secret-scanner';repositoryOwnerId=50274860;workflow=$workflow;ref=$workflowRef;workflowSha=$toolingRevision;trigger='workflow_dispatch'
    }
    compatibility=[ordered]@{minimumRunnerVersion='1.0.0';supportedOperatingSystems=@('linux','windows');supportedArchitectures=@('amd64');testSummaryAsset='TEST-SUMMARY.json';limitationsAsset='LIMITATIONS.md'}
    revocation=[ordered]@{discoveryLocation=$revocationLocation;snapshotAsset='global-revocations.json';checkpointAsset='global-revocation-checkpoint.json';schemaVersion='1.1';maximumSnapshotAgeHours=24}
    assets=$assets;createdAt=$created
}
Write-Utf8 (Join-Path $dist 'release-manifest.json') ($manifest | ConvertTo-Json -Depth 10)

Write-Output "Release build complete: $dist"
Write-Output "product_source_revision=$productRevision"
Write-Output "product_source_tree=$productTree"
Write-Output "release_tooling_revision=$toolingRevision"
Write-Output "release_tooling_tree=$toolingTree"
Write-Output "manifest_sha256=$((Get-FileHash -LiteralPath (Join-Path $dist 'release-manifest.json') -Algorithm SHA256).Hash.ToLowerInvariant())"
