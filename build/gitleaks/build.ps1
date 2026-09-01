param(
    [Parameter(Mandatory = $true)][string]$SourceRoot,
    [Parameter(Mandatory = $true)][string]$GoRoot,
    [Parameter(Mandatory = $true)][string]$OutputDirectory
)

$ErrorActionPreference = 'Stop'
$expected = @{
    'go.mod'                 = '607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3'
    'config/gitleaks.toml'   = 'e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf'
    'LICENSE'                = 'e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6'
    '.goreleaser.yml'        = '1e6a76e13378b4ad215b423411ee6b1732f04c82224fbed4dcdb20b99eab9717'
}
foreach ($entry in $expected.GetEnumerator()) {
    $path = Join-Path $SourceRoot $entry.Key
    if (!(Test-Path -LiteralPath $path -PathType Leaf)) { throw "Missing pinned source file: $($entry.Key)" }
    $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $entry.Value) { throw "Pinned source digest mismatch: $($entry.Key)" }
}

$goExe = Join-Path $GoRoot 'bin/go.exe'
if (!(Test-Path -LiteralPath $goExe -PathType Leaf)) { throw 'Pinned Go executable is missing' }
if ((& $goExe version) -notmatch '^go version go1\.27\.0 windows/amd64$') { throw 'Go toolchain binding mismatch' }
New-Item -ItemType Directory -Force -Path $OutputDirectory | Out-Null

$previousToolchain = $env:GOTOOLCHAIN
$previousCGO = $env:CGO_ENABLED
$previousFlags = $env:GOFLAGS
try {
    $env:GOTOOLCHAIN = 'local'
    $env:CGO_ENABLED = '0'
    $env:GOFLAGS = '-mod=readonly'
    $ldflags = '-s -w -buildid= -X=github.com/zricethezav/gitleaks/v8/version.Version=8.30.1'
    Push-Location $SourceRoot
    try {
        $env:GOOS = 'windows'; $env:GOARCH = 'amd64'
        & $goExe build -trimpath -buildvcs=false -ldflags $ldflags -o (Join-Path $OutputDirectory 'gitleaks-windows-amd64.exe') .
        if ($LASTEXITCODE -ne 0) { throw 'Windows build failed' }
        $env:GOOS = 'linux'; $env:GOARCH = 'amd64'
        & $goExe build -trimpath -buildvcs=false -ldflags $ldflags -o (Join-Path $OutputDirectory 'gitleaks-linux-amd64') .
        if ($LASTEXITCODE -ne 0) { throw 'Linux build failed' }
    } finally {
        Pop-Location
    }
} finally {
    $env:GOTOOLCHAIN = $previousToolchain
    $env:CGO_ENABLED = $previousCGO
    $env:GOFLAGS = $previousFlags
    Remove-Item Env:GOOS,Env:GOARCH -ErrorAction SilentlyContinue
}

$expectedOutputs = @{
    'gitleaks-windows-amd64.exe' = 'b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178'
    'gitleaks-linux-amd64'       = '657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586'
}
foreach ($entry in $expectedOutputs.GetEnumerator()) {
    $actual = (Get-FileHash -LiteralPath (Join-Path $OutputDirectory $entry.Key) -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($actual -ne $entry.Value) { throw "Non-reproducible output: $($entry.Key)" }
}
