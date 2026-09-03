param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [switch]$AllowImagePull
)

$ErrorActionPreference = 'Stop'

$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$cache = [IO.Path]::GetFullPath($CacheDirectory)
if (Test-Path -LiteralPath $cache) {
    if (Test-Path -LiteralPath (Join-Path $cache 'acquisition-ledger.json')) {
        throw 'Acquisition cache is already complete and is not mutable'
    }
    $unexpected = Get-ChildItem -LiteralPath $cache -Force | Where-Object { $_.Name -notin @('downloads','gomodcache') }
    if ($unexpected) {
        throw 'Incomplete acquisition cache contains unexpected paths'
    }
} else {
    New-Item -ItemType Directory -Path $cache | Out-Null
}
$downloads = Join-Path $cache 'downloads'
$moduleCache = Join-Path $cache 'gomodcache'
New-Item -ItemType Directory -Path $downloads,$moduleCache -Force | Out-Null

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

& docker image inspect $image *> $null
if ($LASTEXITCODE -ne 0) {
    if (!$AllowImagePull) {
        throw 'Pinned build image is absent; rerun only with explicit -AllowImagePull during the networked acquisition phase'
    }
    & docker pull $image
    if ($LASTEXITCODE -ne 0) { throw 'Pinned build image pull failed' }
}

$runnerGo = (Resolve-Path (Join-Path $downloads 'go1.27.1.linux-amd64.tar.gz')).Path
$engineGo = (Resolve-Path (Join-Path $downloads 'go1.27.0.linux-amd64.tar.gz')).Path
$gitleaks = (Resolve-Path (Join-Path $downloads 'gitleaks-83d9cd684c87d95d656c1458ef04895a7f1cbd8e.tar.gz')).Path

$arguments = @(
    'run', '--rm', '--pull=never',
    '--cap-drop', 'ALL', '--security-opt', 'no-new-privileges',
    '--pids-limit', '256', '--memory', '4g', '--memory-swap', '4g', '--cpus', '2',
    '--tmpfs', '/work:rw,exec,nosuid,nodev,size=1g',
    '--mount', "type=bind,src=$root,dst=/src,readonly",
    '--mount', "type=bind,src=$runnerGo,dst=/input/runner-go.tar.gz,readonly",
    '--mount', "type=bind,src=$engineGo,dst=/input/engine-go.tar.gz,readonly",
    '--mount', "type=bind,src=$gitleaks,dst=/input/gitleaks.tar.gz,readonly",
    '--mount', "type=bind,src=$moduleCache,dst=/gomodcache",
    '--workdir', '/src',
    $image,
    '/usr/bin/env', '-i', 'PATH=/usr/bin:/bin', 'HOME=/work',
    '/bin/sh', '-ceu', @'
echo "63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445  /input/runner-go.tar.gz" | sha256sum -c -
echo "675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685  /input/engine-go.tar.gz" | sha256sum -c -
echo "6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115  /input/gitleaks.tar.gz" | sha256sum -c -
mkdir -p /work/runner /work/engine /work/gitleaks
tar -xzf /input/runner-go.tar.gz -C /work/runner
tar -xzf /input/engine-go.tar.gz -C /work/engine
tar -xzf /input/gitleaks.tar.gz -C /work/gitleaks --strip-components=1
test "$(/work/runner/go/bin/go version)" = "go version go1.27.1 linux/amd64"
test "$(/work/engine/go/bin/go version)" = "go version go1.27.0 linux/amd64"
cd /src
env GOTOOLCHAIN=local GOMODCACHE=/gomodcache GOCACHE=/work/runner-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/runner/go/bin/go mod download
cd /work/gitleaks
env GOTOOLCHAIN=local GOMODCACHE=/gomodcache GOCACHE=/work/engine-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/engine/go/bin/go mod download
'@
)

& docker @arguments
if ($LASTEXITCODE -ne 0) { throw 'Pinned dependency acquisition failed' }

$ledger = [ordered]@{
    schemaVersion = '1.0'
    acquiredAt = [DateTime]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ssZ')
    image = $image
    artifacts = $artifacts
    moduleCache = 'gomodcache'
    networkBoundary = 'Network enabled only in this acquisition phase; builds require --network none and read-only cache mounts.'
}
$json = $ledger | ConvertTo-Json -Depth 8
[IO.File]::WriteAllText((Join-Path $cache 'acquisition-ledger.json'), $json + "`n", [Text.UTF8Encoding]::new($false))
Write-Output "Acquisition complete: $cache"
