param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [Parameter(Mandatory = $true)][string]$SourceRepository,
    [Parameter(Mandatory = $true)][string]$SourceRevision,
    [Parameter(Mandatory = $true)][string]$WorkingDirectory
)

$ErrorActionPreference = 'Stop'
if ($MyInvocation.InvocationName -eq '.') { throw 'The release image admission entrypoint cannot be dot-sourced' }

$releaseImage = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$dockerBoundary = [IO.Path]::Combine($PSScriptRoot, 'docker-execution.ps1')
$receiptPath = Join-Path ([IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)) 'docker-admission.json'
if (Test-Path -LiteralPath $receiptPath) { throw 'Docker admission receipt path must be absent before any admission operation' }
. ([IO.Path]::Combine($PSScriptRoot, 'image-admission.ps1'))

function Invoke-ClosedOperation([string]$Operation, [string]$ExpectedDockerSHA256) {
    $parameters = @{ Operation = $Operation }
    if (![string]::IsNullOrEmpty($ExpectedDockerSHA256)) { $parameters.ExpectedDockerSHA256 = $ExpectedDockerSHA256 }
    $json = & $dockerBoundary @parameters
    if ($LASTEXITCODE -ne 0 -or $null -eq $json) { throw "Docker $Operation boundary failed" }
    try { $result = ($json -join "`n") | ConvertFrom-Json } catch { throw "Docker $Operation returned malformed lifecycle evidence" }
    if ($result.Operation -cne $Operation -or $result.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($result.StdErr) -or !$result.ContainmentEmpty -or $result.DockerSHA256 -notmatch '^[0-9a-f]{64}$') {
        throw "Docker $Operation failed and is terminal untrusted evidence"
    }
    $result
}

function Invoke-GenuineHostPrerequisites {
    $cacheRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)
    $downloads = [IO.Path]::Combine($cacheRoot, 'downloads')
    $moduleCache = [IO.Path]::Combine($cacheRoot, 'gomodcache')
    foreach ($path in @($downloads, $moduleCache)) { if (!(Test-Path -LiteralPath $path -PathType Container)) { throw "Exact acquisition cache path is absent before image admission: $path" } }
    $hostProof = & {
        . ([IO.Path]::Combine($PSScriptRoot, 'host-cache-canary.ps1'))
        $proof = $null
        foreach ($path in @($cacheRoot, $downloads, $moduleCache)) { $proof = Invoke-HostCacheCanary -CacheDirectory $path; if ($null -eq $proof) { throw 'Host cache prerequisite returned no proof' } }
        $proof
    }
    [void](& ([IO.Path]::Combine($PSScriptRoot, 'test-crlf-shell-payloads.ps1')) -SourceRepository $SourceRepository -SourceRevision $SourceRevision -WorkingDirectory $WorkingDirectory -Phase HostOnly)
    [pscustomobject]@{ CacheRoot=$cacheRoot; Downloads=$downloads; ModuleCache=$moduleCache; HostProof=$hostProof }
}

$prerequisites = Invoke-GenuineHostPrerequisites
$engine = Invoke-ClosedOperation -Operation EngineInspection
[void](Assert-ReleaseEngineEvidence -Json $engine.StdOut)
$dockerSHA256 = $engine.DockerSHA256
$inventory = Invoke-ClosedOperation -Operation ExactImageInventory -ExpectedDockerSHA256 $dockerSHA256
$state = Resolve-ReleaseImageListEvidence -Json $inventory.StdOut
$pulled = $false
if ($state -eq 'ConclusiveAbsent') {
    [void](Invoke-ClosedOperation -Operation ApprovedImagePull -ExpectedDockerSHA256 $dockerSHA256)
    $pulled = $true
}
$inspect = Invoke-ClosedOperation -Operation RepositoryDigestInspection -ExpectedDockerSHA256 $dockerSHA256
$repoDigests = Read-ReleaseRepoDigestsEvidence -Json $inspect.StdOut
[void](Assert-ReleaseImageIdentityEvidence -Image $releaseImage -RepoDigests $repoDigests)

$cacheParameters = @{
    Operation = 'ContainerCacheProof'; ExpectedDockerSHA256 = $dockerSHA256
    CacheDirectory = $prerequisites.ModuleCache
}
if ($null -ne $prerequisites.HostProof.HostUID) { $cacheParameters.HostUID = [int]$prerequisites.HostProof.HostUID; $cacheParameters.HostGID = [int]$prerequisites.HostProof.HostGID }
$cacheJson = & $dockerBoundary @cacheParameters
if ($LASTEXITCODE -ne 0) { throw 'Post-admission container cache proof failed' }
$cacheResult = ($cacheJson -join "`n") | ConvertFrom-Json
if ($cacheResult.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($cacheResult.StdErr) -or !$cacheResult.ContainmentEmpty -or $cacheResult.DockerSHA256 -cne $dockerSHA256) { throw 'Post-admission container cache proof is terminal untrusted evidence' }

$result = [pscustomobject]@{
    Image=$releaseImage; Pulled=$pulled; DockerExecutableSHA256=$dockerSHA256
    HostIdentityMode=$prerequisites.HostProof.IdentityMode; HostUID=$prerequisites.HostProof.HostUID; HostGID=$prerequisites.HostProof.HostGID
    ProvedCachePaths=@($prerequisites.CacheRoot,$prerequisites.Downloads,$prerequisites.ModuleCache)
}
$receipt = [ordered]@{
    schemaVersion='1.0'; sourceRevision=$SourceRevision; image=$releaseImage
    dockerExecutableSHA256=$dockerSHA256; pulledDuringAdmission=$pulled
    hostIdentityMode=$result.HostIdentityMode; hostUID=$result.HostUID; hostGID=$result.HostGID
    containment='empty-after-every-operation'; streamLimitBytes=131072; commandBudgetMilliseconds=15000; cleanupGraceMilliseconds=2000
}
[IO.File]::WriteAllText($receiptPath, ($receipt | ConvertTo-Json -Depth 4) + "`n", [Text.UTF8Encoding]::new($false))
$result
