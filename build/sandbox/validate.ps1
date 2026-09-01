param(
    [Parameter(Mandatory = $true)]
    [string]$ModuleCache,
    [Parameter(Mandatory = $false)]
    [string]$GitleaksBinary = "",
    [ValidateSet('pr', 'release')]
    [string]$Profile = 'pr'
)

$ErrorActionPreference = 'Stop'

$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$root = (Resolve-Path (Join-Path $PSScriptRoot '..\..')).Path
$cache = (Resolve-Path -LiteralPath $ModuleCache).Path
$memory = if ($Profile -eq 'release') { '8g' } else { '4g' }
$memoryBytes = if ($Profile -eq 'release') { '8589934592' } else { '4294967296' }

$arguments = @(
    'run', '--rm', '--pull=never',
    '--network', 'none',
    '--read-only',
    '--cap-drop', 'ALL',
    '--security-opt', 'no-new-privileges',
    '--pids-limit', '256',
    '--memory', $memory,
    '--memory-swap', $memory,
    '--cpus', '2',
    '--tmpfs', '/tmp:rw,noexec,nosuid,nodev,size=64m',
    '--tmpfs', '/work:rw,exec,nosuid,nodev,size=1g',
    '--env', 'PSCAN_AUTHORITATIVE_DOCKER=1',
    '--env', "PSCAN_SANDBOX_MEMORY_BYTES=$memoryBytes",
    '--env', 'PSCAN_SANDBOX_PROCESS_LIMIT=1',
    '--mount', "type=bind,src=$root,dst=/src,readonly",
    '--mount', "type=bind,src=$cache,dst=/modcache,readonly",
    '--workdir', '/src'
)

if ($GitleaksBinary -ne '') {
    $binary = (Resolve-Path -LiteralPath $GitleaksBinary).Path
    $arguments += @('--mount', "type=bind,src=$binary,dst=/engine/gitleaks,readonly")
}

$arguments += @(
    $image,
    '/usr/bin/env', '-i',
    'PATH=/usr/local/go/bin:/usr/bin:/bin',
    'HOME=/work',
    'GOMODCACHE=/modcache',
    'GOCACHE=/work/cache',
    'GOTMPDIR=/work/tmp',
    'TMPDIR=/work/tmp',
    'GOTOOLCHAIN=local',
    'GOPROXY=off',
    'GOSUMDB=off',
    'PSCAN_AUTHORITATIVE_DOCKER=1',
    "PSCAN_SANDBOX_MEMORY_BYTES=$memoryBytes",
    'PSCAN_SANDBOX_PROCESS_LIMIT=1',
    $(if ($GitleaksBinary -ne '') { 'PSCAN_GITLEAKS_BINARY=/engine/gitleaks' } else { 'PSCAN_GITLEAKS_BINARY=' }),
    'PSCAN_GITLEAKS_CONFIG=/src/rules/generic/gitleaks-v8.30.1.toml',
    'PSCAN_GITLEAKS_IGNORE=/src/rules/generic/gitleaks-ignore-empty-v1.txt',
    '/bin/sh', '-ceu',
    'mkdir -p /work/tmp /work/cache; go version; go test -count=1 ./...; go vet ./...'
)

& docker @arguments
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
