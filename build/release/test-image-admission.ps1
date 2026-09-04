$ErrorActionPreference = 'Stop'

foreach ($scriptName in @('image-admission.ps1', 'admit-image.ps1', 'test-image-admission.ps1')) {
    [void][scriptblock]::Create((Get-Content -Raw (Join-Path $PSScriptRoot $scriptName)))
}

. (Join-Path $PSScriptRoot 'image-admission.ps1')

$exact = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineDigest = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$wrongDigest = 'golang@sha256:1111111111111111111111111111111111111111111111111111111111111111'
$wrongRepository = 'example.invalid/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineJson = '{"Version":"29.7.2","ApiVersion":"1.52","Os":"linux","Arch":"amd64"}'
$listJson = '{"Repository":"golang","Digest":"sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452","ID":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}'
$commandVersion = 'docker|version|--format|{{json .Server}}'
$commandList = "docker|image|ls|--all|--no-trunc|--digests|--filter|reference=$exact|--format|{{json .}}"
$commandInspect = "docker|image|inspect|--format|{{json .RepoDigests}}|$exact"
$commandPull = "docker|pull|$exact"

function New-FakeResult {
    param(
        [int]$ExitCode = 0,
        [string]$StdOut = '',
        [string]$StdErr = '',
        [bool]$TimedOut = $false
    )
    return [pscustomobject]@{ ExitCode=$ExitCode; StdOut=$StdOut; StdErr=$StdErr; TimedOut=$TimedOut }
}

function Assert-Equal {
    param($Actual,$Expected,[string]$Message)
    if ($Actual -cne $Expected) { throw "$Message expected=[$Expected] actual=[$Actual]" }
}

function Assert-Sequence {
    param([string[]]$Actual,[string[]]$Expected,[string]$Name)
    if ($Actual.Count -ne $Expected.Count) {
        throw "$Name command/event count mismatch expected=$($Expected.Count) actual=$($Actual.Count): $($Actual -join '; ')"
    }
    for ($index = 0; $index -lt $Expected.Count; $index++) {
        if ($Actual[$index] -cne $Expected[$index]) {
            throw "$Name order mismatch at $index expected=[$($Expected[$index])] actual=[$($Actual[$index])]"
        }
    }
}

$temporaryRoot = Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-002-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $temporaryRoot)
$cache = Join-Path $temporaryRoot 'cache'
[void](New-Item -ItemType Directory -Path $cache)
[void](New-Item -ItemType Directory -Path (Join-Path $cache 'downloads'))
[void](New-Item -ItemType Directory -Path (Join-Path $cache 'gomodcache'))

$script:FakeResults = $null
$script:InvocationEvents = $null
$script:ReleaseHostCacheCanaryInvoker = {
    param([string]$Path)
    [void]$script:InvocationEvents.Add('host-cache')
    Invoke-HostCacheCanary -CacheDirectory $Path
}
$script:ReleaseHostOnlyCrlfProofInvoker = {
    param([string]$SourceRepository,[string]$SourceRevision,[string]$WorkingDirectory)
    [void]$script:InvocationEvents.Add('host-crlf')
}
$script:ReleaseDockerInvoker = {
    param($Request)
    [void]$script:InvocationEvents.Add(('docker|' + ($Request.Arguments -join '|')))
    if ($script:FakeResults.Count -eq 0) { throw 'Fake engine received an unexpected Docker command' }
    return $script:FakeResults.Dequeue()
}

function Invoke-FakeScenario {
    param(
        [Parameter(Mandatory = $true)][string]$Name,
        [Parameter(Mandatory = $true)][object[]]$Results,
        [Parameter(Mandatory = $true)][string[]]$ExpectedEvents,
        [bool]$ExpectSuccess = $false,
        [bool]$ExpectPulled = $false
    )
    $script:FakeResults = [Collections.Generic.Queue[object]]::new()
    foreach ($result in $Results) { $script:FakeResults.Enqueue($result) }
    $script:InvocationEvents = [Collections.Generic.List[string]]::new()
    $laterActions = 0
    $succeeded = $false
    $admission = $null
    try {
        $admission = Invoke-ReleaseImageBootstrap `
            -CacheDirectory $cache `
            -SourceRepository $PSScriptRoot `
            -SourceRevision '0000000000000000000000000000000000000000' `
            -WorkingDirectory (Join-Path $temporaryRoot 'crlf-proof')
        $succeeded = $true
        $laterActions++
    } catch {
        if ($ExpectSuccess) { throw "$Name unexpectedly rejected: $($_.Exception.Message)" }
    }
    if (!$ExpectSuccess -and $succeeded) { throw "$Name unexpectedly admitted" }
    if ($ExpectSuccess) {
        Assert-Equal -Actual $admission.Image -Expected $exact -Message "$Name image"
        Assert-Equal -Actual $admission.Pulled -Expected $ExpectPulled -Message "$Name pull state"
        Assert-Equal -Actual $laterActions -Expected 1 -Message "$Name later-action gate"
    } else {
        Assert-Equal -Actual $laterActions -Expected 0 -Message "$Name later-action gate"
    }
    Assert-Sequence -Actual $script:InvocationEvents.ToArray() -Expected $ExpectedEvents -Name $Name
}

$hostPrefix = @('host-cache','host-cache','host-cache','host-crlf')
$presentEvents = @($hostPrefix + @($commandVersion,$commandList,$commandInspect))
$absentEvents = @($hostPrefix + @($commandVersion,$commandList,$commandPull,$commandInspect))
$versionOnly = @($hostPrefix + @($commandVersion))
$listReached = @($hostPrefix + @($commandVersion,$commandList))

try {
    Invoke-FakeScenario -Name 'pre-existing-exact' -ExpectSuccess $true -ExpectPulled $false `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut $listJson),
            (New-FakeResult -StdOut ('["' + $engineDigest + '"]'))
        ) -ExpectedEvents $presentEvents

    Invoke-FakeScenario -Name 'conclusive-absence-single-pull' -ExpectSuccess $true -ExpectPulled $true `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut ''),
            (New-FakeResult -StdOut 'pulled'),
            (New-FakeResult -StdOut ('["' + $engineDigest + '"]'))
        ) -ExpectedEvents $absentEvents

    Invoke-FakeScenario -Name 'daemon-unavailable' `
        -Results @((New-FakeResult -ExitCode 1 -StdErr 'daemon unavailable')) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'permission-denial' `
        -Results @((New-FakeResult -StdOut $engineJson),(New-FakeResult -ExitCode 1 -StdErr 'permission denied')) -ExpectedEvents $listReached
    Invoke-FakeScenario -Name 'timeout' `
        -Results @((New-FakeResult -ExitCode -1 -TimedOut $true)) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'invalid-invocation' `
        -Results @((New-FakeResult -ExitCode 125 -StdErr 'invalid invocation')) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'invalid-command-boundary-protocol' `
        -Results @([pscustomobject]@{ ExitCode=0; StdOut=$engineJson; StdErr='' }) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'unexpected-success-stderr' `
        -Results @((New-FakeResult -StdOut $engineJson -StdErr 'warning')) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'over-limit-command-output' `
        -Results @((New-FakeResult -StdOut ('x' * ($script:ReleaseDockerOutputLimit + 1)))) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'malformed-engine-data' `
        -Results @((New-FakeResult -StdOut '{')) -ExpectedEvents $versionOnly
    Invoke-FakeScenario -Name 'malformed-image-list-data' `
        -Results @((New-FakeResult -StdOut $engineJson),(New-FakeResult -StdOut '{')) -ExpectedEvents $listReached
    Invoke-FakeScenario -Name 'ambiguous-image-list-data' `
        -Results @((New-FakeResult -StdOut $engineJson),(New-FakeResult -StdOut ($listJson + "`n" + $listJson))) -ExpectedEvents $listReached
    Invoke-FakeScenario -Name 'deceptive-absence-text' `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut $listJson),
            (New-FakeResult -ExitCode 1 -StdErr 'No such image: absent')
        ) -ExpectedEvents $presentEvents

    $identityCases = @(
        [ordered]@{ Name='empty-identity-output'; Json='' },
        [ordered]@{ Name='null-identity'; Json='null' },
        [ordered]@{ Name='scalar-identity'; Json=('"' + $engineDigest + '"') },
        [ordered]@{ Name='empty-identity-array'; Json='[]' },
        [ordered]@{ Name='malformed-identity'; Json='[' },
        [ordered]@{ Name='exact-plus-wrong'; Json=('["' + $engineDigest + '","' + $wrongDigest + '"]') },
        [ordered]@{ Name='exact-plus-alias'; Json=('["' + $engineDigest + '","' + $exact + '"]') },
        [ordered]@{ Name='duplicate-canonical'; Json=('["' + $engineDigest + '","' + $engineDigest + '"]') },
        [ordered]@{ Name='alias-only'; Json=('["' + $exact + '"]') },
        [ordered]@{ Name='wrong-repository'; Json=('["' + $wrongRepository + '"]') },
        [ordered]@{ Name='wrong-digest'; Json=('["' + $wrongDigest + '"]') }
    )
    foreach ($case in $identityCases) {
        Invoke-FakeScenario -Name $case.Name `
            -Results @(
                (New-FakeResult -StdOut $engineJson),
                (New-FakeResult -StdOut $listJson),
                (New-FakeResult -StdOut $case.Json)
            ) -ExpectedEvents $presentEvents
    }

    Invoke-FakeScenario -Name 'pull-failure-no-retry' `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut ''),
            (New-FakeResult -ExitCode 1 -StdErr 'pull failed')
        ) -ExpectedEvents @($hostPrefix + @($commandVersion,$commandList,$commandPull))
    Invoke-FakeScenario -Name 'post-pull-inspect-failure' `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut ''),
            (New-FakeResult -StdOut 'pulled'),
            (New-FakeResult -ExitCode 1 -StdErr 'inspect failed')
        ) -ExpectedEvents $absentEvents
    Invoke-FakeScenario -Name 'post-pull-mixed-identity' `
        -Results @(
            (New-FakeResult -StdOut $engineJson),
            (New-FakeResult -StdOut ''),
            (New-FakeResult -StdOut 'pulled'),
            (New-FakeResult -StdOut ('["' + $engineDigest + '","' + $wrongDigest + '"]'))
        ) -ExpectedEvents $absentEvents

    $script:InvocationEvents = [Collections.Generic.List[string]]::new()
    $legacyRejected = $false
    try { [void](Invoke-ReleaseImageAdmission -CacheDirectory $cache -AllowImagePull) }
    catch { $legacyRejected = $true }
    if (!$legacyRejected) { throw 'Legacy AllowImagePull request unexpectedly retained pull authority' }
    Assert-Sequence -Actual $script:InvocationEvents.ToArray() -Expected @() -Name 'legacy-pull-authority'

    $workflow = Get-Content -Raw (Join-Path $PSScriptRoot '..\..\.github\workflows\release-recovery-v1.0.0.yml')
    $admitScript = Get-Content -Raw (Join-Path $PSScriptRoot 'admit-image.ps1')
    $imageScript = Get-Content -Raw (Join-Path $PSScriptRoot 'image-admission.ps1')
    if ($workflow.Contains('AllowImagePull') -or $admitScript.Contains('AllowImagePull')) {
        throw 'Recovery workflow or mandatory orchestrator still exposes an image-pull bypass switch'
    }
    foreach ($required in @('-SourceRepository','-SourceRevision','-WorkingDirectory')) {
        if (!$workflow.Contains($required)) { throw "Recovery workflow omits mandatory host prerequisite argument $required" }
    }
    if (($imageScript.Split("'pull',",[StringSplitOptions]::None).Count - 1) -ne 1) {
        throw 'Image admission implementation does not contain exactly one pull-capable command boundary'
    }
    $order = @(
        $workflow.IndexOf('Prove image identity admission cases without Docker'),
        $workflow.IndexOf('Prove host cache ownership failures before any image operation'),
        $workflow.IndexOf('Run mandatory host prerequisites and admit the pinned image'),
        $workflow.IndexOf('Prove offline container cache rejection after image admission'),
        $workflow.IndexOf('Acquire pinned public inputs after cache canary'),
        $workflow.IndexOf('Build and validate from an actual CRLF checkout without network')
    )
    if ($order -contains -1) { throw 'Recovery workflow is missing an ordered admission boundary' }
    for ($index=1; $index -lt $order.Count; $index++) {
        if ($order[$index] -le $order[$index-1]) { throw 'Recovery workflow admission ordering is invalid' }
    }

    Write-Output 'Image admission iteration-002 PASS present=NO-PULL absent=ONE-PULL failures=ZERO-PULL identity-set=UNIQUE ordering=BOUND'
} finally {
    $script:ReleaseDockerInvoker = $null
    $script:ReleaseHostCacheCanaryInvoker = $null
    $script:ReleaseHostOnlyCrlfProofInvoker = $null
    $resolvedTemporaryRoot = [IO.Path]::GetFullPath($temporaryRoot)
    $resolvedSystemTemp = [IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if (!$resolvedTemporaryRoot.StartsWith($resolvedSystemTemp,[StringComparison]::OrdinalIgnoreCase)) {
        throw "Unsafe test cleanup path: $resolvedTemporaryRoot"
    }
    if (Test-Path -LiteralPath $resolvedTemporaryRoot) {
        Remove-Item -LiteralPath $resolvedTemporaryRoot -Recurse -Force
    }
}
