param(
    [string]$WorkingDirectory
)

$ErrorActionPreference = 'Stop'
$image = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$releaseRoot = $PSScriptRoot

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

function Write-CrlfCopy([string]$Source, [string]$Destination) {
    $content = [IO.File]::ReadAllText($Source).Replace("`r`n", "`n").Replace("`r", "`n")
    [IO.File]::WriteAllText($Destination, $content.Replace("`n", "`r`n"), [Text.UTF8Encoding]::new($false))
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
New-Item -ItemType Directory -Path $fixtureRoot | Out-Null
foreach ($name in @('shell-payload.ps1','cache-canary.ps1','acquire.ps1','build.ps1')) {
    Write-CrlfCopy (Join-Path $releaseRoot $name) (Join-Path $fixtureRoot $name)
}

. (Join-Path $fixtureRoot 'shell-payload.ps1')
$payloadCases = @(
    [ordered]@{Name='cache-canary';Path=(Join-Path $fixtureRoot 'cache-canary.ps1');Marker='PSCAN-06-C1-CACHE-CANARY'},
    [ordered]@{Name='acquisition';Path=(Join-Path $fixtureRoot 'acquire.ps1');Marker='/work/runner/go/bin/go mod download'},
    [ordered]@{Name='build';Path=(Join-Path $fixtureRoot 'build.ps1');Marker='verifier_ldflags='}
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
    $content = [IO.File]::ReadAllText((Join-Path $fixtureRoot $entry.Key))
    $count = ([regex]::Matches($content, 'ConvertTo-LFPosixShellPayload\s+-Payload')).Count
    if ($count -ne $entry.Value) {
        throw "$($entry.Key) does not normalize every Docker POSIX shell payload immediately before invocation"
    }
}

. (Join-Path $fixtureRoot 'cache-canary.ps1')
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

Write-Output 'CRLF shell-payload regression PASS raw-CR=REJECT normalized-CR=ABSENT canary=PASS acquisition-syntax=PASS build-syntax=PASS read-only-cache=REJECT'
