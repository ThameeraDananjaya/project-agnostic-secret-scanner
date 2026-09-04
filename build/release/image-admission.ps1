$script:ReleaseImage = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$script:ReleaseImageDigest = 'sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$script:ReleaseEngineRepoDigest = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$script:ReleaseDockerOutputLimit = 131072
$script:ReleaseDockerTimeoutMilliseconds = 15000
$script:ReleaseDockerInvoker = $null
$script:ReleaseHostCacheCanaryInvoker = $null
$script:ReleaseHostOnlyCrlfProofInvoker = $null

function Assert-ExactReleaseImageReference([string]$Reference) {
    if ($Reference -cne $script:ReleaseImage) {
        throw 'Build image must use the exact canonical Correction C2 repository and digest reference'
    }
    return $Reference
}

function Invoke-ReleaseDockerProcess {
    [CmdletBinding()]
    param([Parameter(Mandatory = $true)][string[]]$Arguments)

    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = 'docker'
    $startInfo.UseShellExecute = $false
    $startInfo.CreateNoWindow = $true
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    foreach ($argument in $Arguments) { [void]$startInfo.ArgumentList.Add($argument) }
    $process = [Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    try {
        if (!$process.Start()) { throw 'Docker command could not be started' }
        $stdoutTask = $process.StandardOutput.ReadToEndAsync()
        $stderrTask = $process.StandardError.ReadToEndAsync()
        $completed = $process.WaitForExit($script:ReleaseDockerTimeoutMilliseconds)
        if (!$completed) {
            try { $process.Kill($true) } catch {}
            $process.WaitForExit()
        }
        $stdout = $stdoutTask.GetAwaiter().GetResult()
        $stderr = $stderrTask.GetAwaiter().GetResult()
        return [pscustomobject]@{
            ExitCode = if ($completed) { $process.ExitCode } else { -1 }
            StdOut = $stdout
            StdErr = $stderr
            TimedOut = !$completed
        }
    } finally { $process.Dispose() }
}

function Invoke-ReleaseDockerBoundary {
    [CmdletBinding()]
    param([Parameter(Mandatory = $true)][string[]]$Arguments)

    $request = [pscustomobject]@{ Arguments = @($Arguments) }
    $result = if ($null -ne $script:ReleaseDockerInvoker) {
        & $script:ReleaseDockerInvoker $request
    } else {
        Invoke-ReleaseDockerProcess -Arguments $Arguments
    }
    if ($null -eq $result) { throw 'Docker command boundary returned no result' }
    foreach ($name in @('ExitCode','StdOut','StdErr','TimedOut')) {
        if ($name -notin $result.PSObject.Properties.Name) { throw 'Docker command boundary returned an invalid result' }
    }
    if ($result.ExitCode -isnot [sbyte] -and $result.ExitCode -isnot [byte] -and
        $result.ExitCode -isnot [int16] -and $result.ExitCode -isnot [uint16] -and
        $result.ExitCode -isnot [int32] -and $result.ExitCode -isnot [uint32] -and
        $result.ExitCode -isnot [int64] -and $result.ExitCode -isnot [uint64]) {
        throw 'Docker command boundary returned an invalid exit code'
    }
    if ($result.StdOut -isnot [string] -or $result.StdErr -isnot [string] -or $result.TimedOut -isnot [bool]) {
        throw 'Docker command boundary returned invalid stream or timeout evidence'
    }
    if ($result.StdOut.Length -gt $script:ReleaseDockerOutputLimit -or
        $result.StdErr.Length -gt $script:ReleaseDockerOutputLimit) {
        throw 'Docker command boundary exceeded its output limit'
    }
    return [pscustomobject]@{
        Arguments = @($Arguments)
        ExitCode = [int]$result.ExitCode
        StdOut = $result.StdOut
        StdErr = $result.StdErr
        TimedOut = $result.TimedOut
    }
}

function Assert-ReleaseDockerSuccess {
    param(
        [Parameter(Mandatory = $true)]$Result,
        [Parameter(Mandatory = $true)][string]$Operation
    )
    if ($Result.TimedOut) { throw "$Operation timed out and is untrusted" }
    if ($Result.ExitCode -ne 0) { throw "$Operation failed and is untrusted" }
    if (![string]::IsNullOrEmpty($Result.StdErr)) { throw "$Operation produced unexpected stderr and is untrusted" }
}

function ConvertFrom-ReleaseJsonDocument {
    param(
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json,
        [Parameter(Mandatory = $true)][string]$Operation
    )
    if ([string]::IsNullOrWhiteSpace($Json)) { throw "$Operation returned empty structured data" }
    try { return [Text.Json.JsonDocument]::Parse($Json) }
    catch { throw "$Operation returned malformed structured data" }
}

function Assert-ReleaseEngineResponsive {
    $result = Invoke-ReleaseDockerBoundary -Arguments @('version','--format','{{json .Server}}')
    Assert-ReleaseDockerSuccess -Result $result -Operation 'Docker engine reachability inspection'
    $document = ConvertFrom-ReleaseJsonDocument -Json $result.StdOut -Operation 'Docker engine reachability inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Object) {
            throw 'Docker engine reachability inspection returned a non-object result'
        }
        $properties = @{}
        foreach ($property in $document.RootElement.EnumerateObject()) {
            if ($properties.ContainsKey($property.Name)) { throw 'Docker engine reachability inspection returned duplicate fields' }
            $properties[$property.Name] = $property.Value
        }
        foreach ($name in @('Version','ApiVersion','Os','Arch')) {
            if (!$properties.ContainsKey($name) -or
                $properties[$name].ValueKind -ne [Text.Json.JsonValueKind]::String -or
                [string]::IsNullOrWhiteSpace($properties[$name].GetString())) {
                throw 'Docker engine reachability inspection returned incomplete structured data'
            }
        }
    } finally { $document.Dispose() }
}

function Get-ReleaseImageListState {
    $result = Invoke-ReleaseDockerBoundary -Arguments @(
        'image','ls','--all','--no-trunc','--digests','--filter',"reference=$script:ReleaseImage",'--format','{{json .}}'
    )
    Assert-ReleaseDockerSuccess -Result $result -Operation 'Docker exact-image inventory inspection'
    if ([string]::IsNullOrWhiteSpace($result.StdOut)) { return 'ConclusiveAbsent' }
    $lines = @([regex]::Split($result.StdOut.Trim(), '\r?\n'))
    if ($lines.Count -ne 1 -or [string]::IsNullOrWhiteSpace($lines[0])) {
        throw 'Docker exact-image inventory inspection returned ambiguous structured data'
    }
    $document = ConvertFrom-ReleaseJsonDocument -Json $lines[0] -Operation 'Docker exact-image inventory inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Object) {
            throw 'Docker exact-image inventory inspection returned a non-object result'
        }
        $properties = @{}
        foreach ($property in $document.RootElement.EnumerateObject()) {
            if ($properties.ContainsKey($property.Name)) { throw 'Docker exact-image inventory inspection returned duplicate fields' }
            $properties[$property.Name] = $property.Value
        }
        foreach ($name in @('Repository','Digest','ID')) {
            if (!$properties.ContainsKey($name) -or $properties[$name].ValueKind -ne [Text.Json.JsonValueKind]::String) {
                throw 'Docker exact-image inventory inspection returned incomplete structured data'
            }
        }
        if ($properties.Repository.GetString() -cne 'golang' -or
            $properties.Digest.GetString() -cne $script:ReleaseImageDigest -or
            $properties.ID.GetString() -cnotmatch '^sha256:[0-9a-f]{64}$') {
            throw 'Docker exact-image inventory inspection returned unexpected identity data'
        }
    } finally { $document.Dispose() }
    return 'CandidatePresent'
}

function Read-ReleaseRepoDigestsEvidence {
    param([Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json)
    $document = ConvertFrom-ReleaseJsonDocument -Json $Json -Operation 'Docker repository-digest inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Array) {
            throw 'Docker repository-digest inspection returned a non-array result'
        }
        $values = [Collections.Generic.List[string]]::new()
        foreach ($element in $document.RootElement.EnumerateArray()) {
            if ($element.ValueKind -ne [Text.Json.JsonValueKind]::String) {
                throw 'Docker repository-digest inspection returned a non-string item'
            }
            $values.Add($element.GetString())
        }
        return ,$values.ToArray()
    } finally { $document.Dispose() }
}

function Assert-ReleaseImageIdentityEvidence {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$Image,
        [Parameter(Mandatory = $true)][AllowNull()]$RepoDigests
    )
    [void](Assert-ExactReleaseImageReference -Reference $Image)
    if ($null -eq $RepoDigests -or $RepoDigests -isnot [array]) {
        throw 'Docker repository-digest evidence must be one structured array'
    }
    $values = @($RepoDigests)
    if ($values.Count -ne 1 -or $values[0] -isnot [string] -or
        $values[0] -cne $script:ReleaseEngineRepoDigest) {
        throw 'Docker repository-digest evidence must contain exactly one canonical engine identity'
    }
    return $Image
}

function Assert-ReleaseImageByInspection {
    param([Parameter(Mandatory = $true)][string]$Image)
    [void](Assert-ExactReleaseImageReference -Reference $Image)
    $result = Invoke-ReleaseDockerBoundary -Arguments @('image','inspect','--format','{{json .RepoDigests}}',$Image)
    Assert-ReleaseDockerSuccess -Result $result -Operation 'Docker repository-digest inspection'
    $repoDigests = Read-ReleaseRepoDigestsEvidence -Json $result.StdOut
    return Assert-ReleaseImageIdentityEvidence -Image $Image -RepoDigests $repoDigests
}

function Get-ReleaseImagePrePullState {
    param([string]$Image = $script:ReleaseImage)
    [void](Assert-ExactReleaseImageReference -Reference $Image)
    Assert-ReleaseEngineResponsive
    $inventoryState = Get-ReleaseImageListState
    if ($inventoryState -eq 'ConclusiveAbsent') {
        return [pscustomobject]@{ State = 'ConclusiveAbsent'; Image = $Image }
    }
    $proved = Assert-ReleaseImageByInspection -Image $Image
    return [pscustomobject]@{ State = 'ExactPresent'; Image = $proved }
}

function Assert-AdmittedReleaseImage {
    [CmdletBinding()]
    param([string]$Image = $script:ReleaseImage)
    $state = Get-ReleaseImagePrePullState -Image $Image
    if ($state.State -ne 'ExactPresent') { throw 'Pinned build image is conclusively absent' }
    return $state.Image
}

function Invoke-ReleaseHostPrerequisites {
    param(
        [Parameter(Mandatory = $true)][string]$CacheDirectory,
        [Parameter(Mandatory = $true)][string]$SourceRepository,
        [Parameter(Mandatory = $true)][string]$SourceRevision,
        [Parameter(Mandatory = $true)][string]$WorkingDirectory
    )
    . (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
    $cacheRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)
    $downloads = Join-Path $cacheRoot 'downloads'
    $moduleCache = Join-Path $cacheRoot 'gomodcache'
    foreach ($path in @($downloads,$moduleCache)) {
        if (!(Test-Path -LiteralPath $path -PathType Container)) {
            throw "Exact acquisition cache path is absent before image admission: $path"
        }
    }
    $hostProof = $null
    foreach ($path in @($cacheRoot,$downloads,$moduleCache)) {
        $hostProof = if ($null -ne $script:ReleaseHostCacheCanaryInvoker) {
            & $script:ReleaseHostCacheCanaryInvoker $path
        } else {
            Invoke-HostCacheCanary -CacheDirectory $path
        }
        if ($null -eq $hostProof) { throw 'Host cache prerequisite returned no proof' }
    }
    if ($null -ne $script:ReleaseHostOnlyCrlfProofInvoker) {
        & $script:ReleaseHostOnlyCrlfProofInvoker $SourceRepository $SourceRevision $WorkingDirectory
    } else {
        & (Join-Path $PSScriptRoot 'test-crlf-shell-payloads.ps1') `
            -SourceRepository $SourceRepository `
            -SourceRevision $SourceRevision `
            -WorkingDirectory $WorkingDirectory `
            -Phase HostOnly
    }
    return [pscustomobject]@{
        CacheRoot = $cacheRoot
        Downloads = $downloads
        ModuleCache = $moduleCache
        HostProof = $hostProof
    }
}

function Invoke-ReleaseImageBootstrap {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$CacheDirectory,
        [Parameter(Mandatory = $true)][string]$SourceRepository,
        [Parameter(Mandatory = $true)][string]$SourceRevision,
        [Parameter(Mandatory = $true)][string]$WorkingDirectory
    )
    $prerequisites = Invoke-ReleaseHostPrerequisites @PSBoundParameters
    $state = Get-ReleaseImagePrePullState -Image $script:ReleaseImage
    $pulled = $false
    if ($state.State -eq 'ConclusiveAbsent') {
        $pull = Invoke-ReleaseDockerBoundary -Arguments @('pull',$script:ReleaseImage)
        Assert-ReleaseDockerSuccess -Result $pull -Operation 'Pinned build image pull'
        $state = [pscustomobject]@{
            State = 'ExactPresent'
            Image = Assert-ReleaseImageByInspection -Image $script:ReleaseImage
        }
        $pulled = $true
    }
    if ($state.State -ne 'ExactPresent') { throw 'Release image admission reached an invalid state' }
    return [pscustomobject]@{
        Image = $state.Image
        Pulled = $pulled
        HostIdentityMode = $prerequisites.HostProof.IdentityMode
        HostUID = $prerequisites.HostProof.HostUID
        HostGID = $prerequisites.HostProof.HostGID
        ProvedCachePaths = @($prerequisites.CacheRoot,$prerequisites.Downloads,$prerequisites.ModuleCache)
    }
}

function Invoke-ReleaseImageAdmission {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$CacheDirectory,
        [Alias('AllowImagePull')][switch]$LegacyImagePullRequest
    )
    if ($LegacyImagePullRequest) {
        throw 'Image pull authority is restricted to the mandatory release-image bootstrap orchestrator'
    }
    . (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
    $cacheRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)
    $downloads = Join-Path $cacheRoot 'downloads'
    $moduleCache = Join-Path $cacheRoot 'gomodcache'
    foreach ($path in @($downloads,$moduleCache)) {
        if (!(Test-Path -LiteralPath $path -PathType Container)) {
            throw "Exact acquisition cache path is absent before image verification: $path"
        }
    }
    $hostProof = $null
    foreach ($path in @($cacheRoot,$downloads,$moduleCache)) {
        $hostProof = Invoke-HostCacheCanary -CacheDirectory $path
    }
    $canonical = Assert-AdmittedReleaseImage -Image $script:ReleaseImage
    return [pscustomobject]@{
        Image = $canonical; Pulled = $false; HostIdentityMode = $hostProof.IdentityMode
        HostUID = $hostProof.HostUID; HostGID = $hostProof.HostGID
        ProvedCachePaths = @($cacheRoot,$downloads,$moduleCache)
    }
}
