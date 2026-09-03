param(
    [Parameter(Mandatory = $true)][string]$AcquisitionDirectory,
    [Parameter(Mandatory = $true)][string]$OutputDirectory
)

$ErrorActionPreference = 'Stop'

. (Join-Path $PSScriptRoot 'shell-payload.ps1')

$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$acquisition = (Resolve-Path -LiteralPath $AcquisitionDirectory).Path
$output = [IO.Path]::GetFullPath($OutputDirectory)
if (Test-Path -LiteralPath $output) {
    if ((Get-ChildItem -LiteralPath $output -Force | Select-Object -First 1)) {
        throw 'Release output must be a new empty directory'
    }
} else {
    New-Item -ItemType Directory -Path $output | Out-Null
}
$raw = Join-Path $output '.raw'
$linuxStage = Join-Path $output '.stage-linux'
$windowsStage = Join-Path $output '.stage-windows'
$dist = Join-Path $output 'dist'
New-Item -ItemType Directory -Path $raw,$linuxStage,$windowsStage,$dist | Out-Null

function Require-Digest([string]$Path, [string]$Expected) {
    if (!(Test-Path -LiteralPath $Path -PathType Leaf)) { throw "Required file missing: $Path" }
    $actual = (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $Expected) { throw "Digest mismatch: $Path" }
}

function Write-Utf8([string]$Path, [string]$Value) {
    [IO.File]::WriteAllText($Path, $Value.Replace("`r`n", "`n") + $(if ($Value.EndsWith("`n")) { '' } else { "`n" }), [Text.UTF8Encoding]::new($false))
}

function Copy-ReleaseFile([string]$Source, [string]$Name) {
    $destination = Join-Path $dist $Name
    Copy-Item -LiteralPath $Source -Destination $destination
    return $destination
}

$status = & git -C $root status --porcelain=v1
if ($LASTEXITCODE -ne 0 -or $status) { throw 'Release build requires an exact clean Git worktree' }
$toolingRevision = (& git -C $root rev-parse HEAD).Trim()
$toolingTree = (& git -C $root rev-parse 'HEAD^{tree}').Trim()
$created = ([DateTimeOffset]::Parse((& git -C $root show -s --format=%cI HEAD).Trim())).UtcDateTime.ToString('yyyy-MM-ddTHH:mm:ssZ')
$epoch = [DateTimeOffset]::Parse($created).ToUnixTimeSeconds()
$productTag = 'v1.0.0'
$productRevision = 'a13c28fe7273bc8dc6545f97966a02889524eb4c'
$productTree = '217b711ddea51fd0ea7e808edd2e27fdecef8427'
$toolingTag = 'release-tooling-v1.0.0-c1'
$workflow = '.github/workflows/release-recovery-v1.0.0.yml'
$workflowRef = 'refs/tags/release-tooling-v1.0.0-c1'
if ($toolingRevision -notmatch '^[0-9a-f]{40}$' -or $toolingTree -notmatch '^[0-9a-f]{40}$') { throw 'Unsupported correction-tooling Git object identity' }
$resolvedProductTag = (& git -C $root rev-parse 'v1.0.0^{commit}').Trim()
$resolvedProductTree = (& git -C $root rev-parse 'v1.0.0^{tree}').Trim()
if ($LASTEXITCODE -ne 0 -or $resolvedProductTag -ne $productRevision -or $resolvedProductTree -ne $productTree) {
    throw 'Locked product tag, commit or tree identity does not match Correction C1 authority'
}
$productArchive = Join-Path $raw 'product-source.tar'
& git -C $root archive --format=tar --output=$productArchive $productRevision
if ($LASTEXITCODE -ne 0 -or !(Test-Path -LiteralPath $productArchive -PathType Leaf)) { throw 'Locked product-source materialization failed' }

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
if ($ledgerValue.schemaVersion -ne '2.0' -or $ledgerValue.cacheCanary.semantics -ne 'write-atomic-rename-read-delete' -or
    $ledgerValue.cacheCanary.completedBeforeNetworkDependencyAcquisition -ne $true) {
    throw 'Acquisition ledger does not prove the Correction C1 pre-network cache canary'
}

& docker image inspect $image *> $null
if ($LASTEXITCODE -ne 0) { throw 'Pinned build image is missing; image acquisition is forbidden during the offline build phase' }

$arguments = @(
    'run', '--rm', '--pull=never', '--network', 'none', '--read-only',
    '--cap-drop', 'ALL', '--security-opt', 'no-new-privileges',
    '--pids-limit', '512', '--memory', '8g', '--memory-swap', '8g', '--cpus', '2',
    '--tmpfs', '/work:rw,exec,nosuid,nodev,size=4g',
    '--mount', "type=bind,src=$root,dst=/src,readonly",
    '--mount', "type=bind,src=$productArchive,dst=/input/product-source.tar,readonly",
    '--mount', "type=bind,src=$runnerGo,dst=/input/runner-go.tar.gz,readonly",
    '--mount', "type=bind,src=$engineGo,dst=/input/engine-go.tar.gz,readonly",
    '--mount', "type=bind,src=$gitleaksSource,dst=/input/gitleaks.tar.gz,readonly",
    '--mount', "type=bind,src=$moduleCache,dst=/gomodcache,readonly",
    '--mount', "type=bind,src=$raw,dst=/out",
    '--workdir', '/src',
    $image,
    '/usr/bin/env', '-i', 'PATH=/usr/bin:/bin', 'HOME=/work', "SOURCE_DATE_EPOCH=$epoch",
    '/bin/sh', '-ceu', @'
echo "63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445  /input/runner-go.tar.gz" | sha256sum -c -
echo "675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685  /input/engine-go.tar.gz" | sha256sum -c -
echo "6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115  /input/gitleaks.tar.gz" | sha256sum -c -
mkdir -p /work/runner /work/engine /work/gitleaks /work/product /work/cache /work/tmp /out/product
tar -xzf /input/runner-go.tar.gz -C /work/runner
tar -xzf /input/engine-go.tar.gz -C /work/engine
tar -xzf /input/gitleaks.tar.gz -C /work/gitleaks --strip-components=1
tar -xf /input/product-source.tar -C /work/product
test "$(/work/runner/go/bin/go version)" = "go version go1.27.1 linux/amd64"
test "$(/work/engine/go/bin/go version)" = "go version go1.27.0 linux/amd64"
test "$(sha256sum /work/gitleaks/go.mod | cut -d ' ' -f 1)" = "607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3"
test "$(sha256sum /work/gitleaks/config/gitleaks.toml | cut -d ' ' -f 1)" = "e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf"
test "$(sha256sum /work/gitleaks/LICENSE | cut -d ' ' -f 1)" = "e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6"
test "$(sha256sum /work/gitleaks/.goreleaser.yml | cut -d ' ' -f 1)" = "1e6a76e13378b4ad215b423411ee6b1732f04c82224fbed4dcdb20b99eab9717"
export GOTOOLCHAIN=local GOMODCACHE=/gomodcache GOPROXY=off GOSUMDB=off CGO_ENABLED=0
export GOCACHE=/work/cache GOTMPDIR=/work/tmp TMPDIR=/work/tmp
cd /work/gitleaks
GOOS=linux GOARCH=amd64 /work/engine/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid= -X=github.com/zricethezav/gitleaks/v8/version.Version=8.30.1' -o /out/gitleaks-linux-amd64 .
GOOS=windows GOARCH=amd64 /work/engine/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid= -X=github.com/zricethezav/gitleaks/v8/version.Version=8.30.1' -o /out/gitleaks-windows-amd64.exe .
test "$(sha256sum /out/gitleaks-linux-amd64 | cut -d ' ' -f 1)" = "657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586"
test "$(sha256sum /out/gitleaks-windows-amd64.exe | cut -d ' ' -f 1)" = "b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178"
cd /work/product
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/scanner-runner-linux-amd64 ./cmd/scanner-runner
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/scanner-runner-windows-amd64.exe ./cmd/scanner-runner
cp -R rules contracts licenses /out/product/
cp LICENSE THIRD_PARTY_NOTICES.md /out/product/
cd /src
verifier_ldflags="-s -w -buildid= -X=main.releaseToolingCommit=$PSCAN_TOOLING_REVISION -X=main.releaseToolingTree=$PSCAN_TOOLING_TREE"
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags "$verifier_ldflags" -o /out/scanner-release-verifier-linux-amd64 ./build/release/cmd/release-verifier
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags "$verifier_ldflags" -o /out/scanner-release-verifier-windows-amd64.exe ./build/release/cmd/release-verifier
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/packager-linux-amd64 ./build/release/cmd/packager
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/sbom-linux-amd64 ./build/release/cmd/sbom
/work/runner/go/bin/go list -mod=readonly -m -json all > /out/modules.json
PSCAN_GITLEAKS_BINARY=/out/gitleaks-linux-amd64 PSCAN_GITLEAKS_CONFIG=/work/product/rules/generic/gitleaks-v8.30.1.toml PSCAN_GITLEAKS_IGNORE=/work/product/rules/generic/gitleaks-ignore-empty-v1.txt /work/runner/go/bin/go test -p=1 -count=1 ./...
/work/runner/go/bin/go vet -p=1 ./...
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go test -p=1 -exec /bin/true ./...
/out/sbom-linux-amd64 -input /out/modules.json -output /out/sbom.spdx.json -revision "$PSCAN_PRODUCT_REVISION" -created "$PSCAN_CREATED"
'@
)
$revisionIndex = $arguments.IndexOf('/usr/bin/env')
$arguments = $arguments[0..($revisionIndex)] + @('-i','PATH=/usr/bin:/bin','HOME=/work',"SOURCE_DATE_EPOCH=$epoch","PSCAN_PRODUCT_REVISION=$productRevision","PSCAN_TOOLING_REVISION=$toolingRevision","PSCAN_TOOLING_TREE=$toolingTree","PSCAN_CREATED=$created") + $arguments[($revisionIndex + 5)..($arguments.Count - 1)]

$arguments[$arguments.Count - 1] = ConvertTo-LFPosixShellPayload -Payload $arguments[$arguments.Count - 1]
Assert-LFPosixShellPayload -Payload $arguments[$arguments.Count - 1]
& docker @arguments
if ($LASTEXITCODE -ne 0) { throw 'Offline release build or validation failed' }

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
Copy-ReleaseFile (Join-Path $productFiles 'rules\generic\gitleaks-v8.30.1.toml') 'rules-gitleaks-v8.30.1.toml' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'rules\generic\gitleaks-ignore-empty-v1.txt') 'rules-gitleaks-ignore-empty-v1.txt' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\scan-request\schema-1.1.json') 'schema-scan-request-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\scan-outcome\schema-1.0.json') 'schema-scan-outcome-1.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\release-manifest\schema-1.1.json') 'schema-release-manifest-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $root 'contracts\release-manifest\schema-2.0.json') 'schema-release-manifest-2.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\global-revocation\schema-1.1.json') 'schema-global-revocation-1.1.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'contracts\rule-pack\schema-1.0.json') 'schema-rule-pack-1.0.json' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'LICENSE') 'LICENSE.txt' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'THIRD_PARTY_NOTICES.md') 'THIRD_PARTY_NOTICES.md' | Out-Null
Copy-ReleaseFile (Join-Path $productFiles 'licenses\gitleaks\modules\manifest.json') 'GITLEAKS-LICENCE-MANIFEST.json' | Out-Null
Copy-ReleaseFile (Join-Path $raw 'sbom.spdx.json') 'sbom.spdx.json' | Out-Null
Copy-ReleaseFile (Join-Path $root 'docs\release\OFFLINE-VERIFICATION-RUNBOOK.md') 'OFFLINE-VERIFICATION-RUNBOOK.md' | Out-Null
Copy-ReleaseFile (Join-Path $root 'docs\release\SCANNER-IO-REFERENCE.md') 'SCANNER-IO-REFERENCE.md' | Out-Null

$testSummary = [ordered]@{
    schemaVersion = '2.0'; productSourceRevision = $productRevision; releaseToolingRevision = $toolingRevision; createdAt = $created
    authoritativeEnvironment = 'pinned-network-disabled-linux-container'
    goToolchain = 'go1.27.1'; engineToolchain = 'go1.27.0'
    commands = @('go test -count=1 ./...', 'go vet ./...', 'GOOS=windows GOARCH=amd64 go test -exec /bin/true ./...')
    linuxExecution = 'PASS'; windowsCompilation = 'PASS'; windowsNativeExecution = 'UNPROVEN_SMART_APP_CONTROL'
    signing = 'NOT_PERFORMED_OWNER_GATE'; remoteWorkflow = 'NOT_PERFORMED_OWNER_GATE'
}
Write-Utf8 (Join-Path $dist 'TEST-SUMMARY.json') ($testSummary | ConvertTo-Json -Depth 6)
$compatibility = [ordered]@{
    schemaVersion = '2.0'; releaseVersion = 'v1.0.0'; productSourceRevision = $productRevision; releaseToolingRevision = $toolingRevision
    platforms = @([ordered]@{os='linux';arch='amd64'},[ordered]@{os='windows';arch='amd64'})
    requestSchemas = @('1.0','1.1'); outcomeSchemas = @('1.0'); releaseManifestSchemas = @('1.0','1.1','2.0')
    engine = [ordered]@{name='gitleaks';version='8.30.1';adapterVersion='2.0.0'}
}
Write-Utf8 (Join-Path $dist 'COMPATIBILITY.json') ($compatibility | ConvertTo-Json -Depth 6)
$limitations = @"
# v1.0.0 limitations

- Native Windows amd64 execution is unproven on the author host because Windows Smart App Control blocks unsigned locally built executables. Windows amd64 cross-compilation and byte verification pass; this is not substituted for native execution.
- Remote repository identity, settings, rules, environment, immutable releases, workflow execution, keyless signature, artifact attestation, draft upload and publication are not performed without separate exact owner approval.
- The offline verifier requires separately acquired Cosign v3.1.3 at its documented exact SHA-256 and rejects absent, stale, conflicting or untrusted evidence.
- Gitleaks is the only primary detector. TruffleHog is not assessed, downloaded, integrated, distributed or enabled.
- Passing validation placeholders do not establish any consuming-project integration, policy, receipt, deployment or production result.
"@
Write-Utf8 (Join-Path $dist 'LIMITATIONS.md') $limitations
$revocationLocation = 'https://api.github.com/repos/ThameeraDananjaya/project-agnostic-secret-scanner/releases?per_page=100'
$revocations = [ordered]@{schemaFamily='global-scanner-revocation-snapshot';schemaVersion='1.0';capturedAt=$created;records=@()}
Write-Utf8 (Join-Path $dist 'global-revocations.json') ($revocations | ConvertTo-Json -Depth 8)
$checkpoint = [ordered]@{schemaFamily='global-scanner-revocation-checkpoint';schemaVersion='1.0';sequence=0;digest=$null;capturedAt=$created;discoveryLocation=$revocationLocation}
Write-Utf8 (Join-Path $dist 'global-revocation-checkpoint.json') ($checkpoint | ConvertTo-Json -Depth 6)
$provenance = [ordered]@{
    schemaVersion='2.0';createdAt=$created
    productSource=[ordered]@{tag=$productTag;commit=$productRevision;tree=$productTree}
    releaseTooling=[ordered]@{tag=$toolingTag;commit=$toolingRevision;tree=$toolingTree;workflow=$workflow;workflowRef=$workflowRef;workflowSha=$toolingRevision;trigger='workflow_dispatch'}
    buildImage=$image;runnerGo='go1.27.1';engineGo='go1.27.0'
    gitleaksSourceRevision='83d9cd684c87d95d656c1458ef04895a7f1cbd8e'
    network='disabled';moduleCache='read-only';cgo=$false;trimpath=$true;buildVCS=$false;buildId='empty'
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

$packageArguments = @(
    'run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','64','--memory','1g','--memory-swap','1g','--cpus','1',
    '--tmpfs','/tmp:rw,exec,nosuid,nodev,size=64m',
    '--mount',"type=bind,src=$raw,dst=/tools,readonly",
    '--mount',"type=bind,src=$linuxStage,dst=/input/linux,readonly",
    '--mount',"type=bind,src=$windowsStage,dst=/input/windows,readonly",
    '--mount',"type=bind,src=$dist,dst=/dist",
    $image,'/bin/sh','-ceu',
    "cp /tools/packager-linux-amd64 /tmp/packager; chmod 0755 /tmp/packager; /tmp/packager -root /input/linux -output /dist/project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz -format tar.gz -epoch $epoch; /tmp/packager -root /input/windows -output /dist/project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip -format zip -epoch $epoch"
)
$packageArguments[$packageArguments.Count - 1] = ConvertTo-LFPosixShellPayload -Payload $packageArguments[$packageArguments.Count - 1]
Assert-LFPosixShellPayload -Payload $packageArguments[$packageArguments.Count - 1]
& docker @packageArguments
if ($LASTEXITCODE -ne 0) { throw 'Deterministic packaging failed' }

$checksumTargets = Get-ChildItem -LiteralPath $dist -File | Sort-Object Name
$checksumLines = foreach ($file in $checksumTargets) {
    "$((Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant())  $($file.Name)"
}
Write-Utf8 (Join-Path $dist 'CHECKSUMS.sha256') ($checksumLines -join "`n")

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
    (Asset 'CHECKSUMS.sha256' 'checksums')
)
$manifest = [ordered]@{
    schemaFamily='scanner-release-manifest';manifestSchemaVersion='2.0';releaseVersion='v1.0.0'
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
        [ordered]@{family='global-scanner-revocation';version='1.1';path='schema-global-revocation-1.1.json';sha256=$assets[14].sha256},
        [ordered]@{family='rule-pack';version='1.0';path='schema-rule-pack-1.0.json';sha256=$assets[15].sha256}
    )
    releaseIdentity=[ordered]@{
        repository='ThameeraDananjaya/project-agnostic-secret-scanner';repositoryOwnerId=50274860;workflow=$workflow;ref=$workflowRef;workflowSha=$toolingRevision;trigger='workflow_dispatch'
        oidcIssuer='https://token.actions.githubusercontent.com';certificateIdentity='https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-recovery-v1.0.0.yml@refs/tags/release-tooling-v1.0.0-c1'
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
