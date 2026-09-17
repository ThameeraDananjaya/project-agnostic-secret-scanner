param(
    [Parameter(Mandatory)][ValidateSet('NativePrerequisites','ImageAdmission','ContainerCache','ContainerCRLF','Acquire','BuildA','BuildB','Compare')][string]$Stage,
    [Parameter(Mandatory)][string]$BaselinePath,
    [Parameter(Mandatory)][ValidatePattern('^[0-9a-f]{64}$')][string]$BaselineSHA256
)
$ErrorActionPreference = 'Stop'
$startEpoch = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$primaryFailure = $null
$afterFailure = $null
try {
    & /usr/bin/python3 -I -B "$PSScriptRoot/unsigned-build-profile-context.py" check $BaselinePath $BaselineSHA256 before
    if ($LASTEXITCODE -ne 0) { throw 'Named profile context admission failed' }
    & /usr/bin/python3 -I -B "$PSScriptRoot/native-host-facts.py" --phase before --since $startEpoch
    if ($LASTEXITCODE -ne 0) { throw 'Before-fixture host facts unavailable' }
    $repository=$env:GITHUB_WORKSPACE;$revision=$env:GITHUB_SHA;$temporary=$env:RUNNER_TEMP
    if($revision-notmatch'^[0-9a-f]{40}$'-or$env:GITHUB_REF-cne'refs/tags/release-tooling-v1.0.0-c2-linux-boundary'-or
        ![IO.Path]::IsPathFullyQualified($repository)-or![IO.Path]::IsPathFullyQualified($temporary)){throw 'Unsigned build invocation identity is invalid'}
    $cache=Join-Path $temporary 'pscan-acquisition'
    switch($Stage){
        'NativePrerequisites' { & "$PSScriptRoot/test-image-admission.ps1" }
        'ImageAdmission' { & "$PSScriptRoot/admit-image.ps1" -CacheDirectory $cache -SourceRepository $repository -SourceRevision $revision -WorkingDirectory (Join-Path $temporary 'pscan-crlf-regression') }
        'ContainerCache' { & "$PSScriptRoot/test-cache-boundary.ps1" -PositiveCacheDirectory (Join-Path $temporary 'cache-positive') -WrongOwnerCacheDirectory (Join-Path $temporary 'cache-wrong-owner') -ReadOnlyCacheDirectory (Join-Path $temporary 'cache-read-only') -Phase Container }
        'ContainerCRLF' { & "$PSScriptRoot/test-crlf-shell-payloads.ps1" -SourceRepository $repository -SourceRevision $revision -WorkingDirectory (Join-Path $temporary 'pscan-crlf-container') -Phase Container }
        'Acquire' { & "$PSScriptRoot/acquire.ps1" -CacheDirectory $cache }
        'BuildA' { & "$PSScriptRoot/test-crlf-shell-payloads.ps1" -SourceRepository $repository -SourceRevision $revision -WorkingDirectory (Join-Path $temporary 'pscan-crlf-build-a') -AcquisitionDirectory $cache -BuildOutputDirectory (Join-Path $temporary 'pscan-release-a') -Phase Build }
        'BuildB' { & "$PSScriptRoot/test-crlf-shell-payloads.ps1" -SourceRepository $repository -SourceRevision $revision -WorkingDirectory (Join-Path $temporary 'pscan-crlf-build-b') -AcquisitionDirectory $cache -BuildOutputDirectory (Join-Path $temporary 'pscan-release-b') -Phase Build }
        'Compare' { & "$PSScriptRoot/compare-builds.ps1" -FirstOutputDirectory (Join-Path $temporary 'pscan-release-a') -SecondOutputDirectory (Join-Path $temporary 'pscan-release-b') }
    }

} catch {
    $primaryFailure = $_
} finally {
    try {
        & /usr/bin/python3 -I -B "$PSScriptRoot/unsigned-build-profile-context.py" check $BaselinePath $BaselineSHA256 after
        if ($LASTEXITCODE -ne 0) { throw 'Named profile after-context mismatch' }
    } catch { $afterFailure = $_ }
    try {
        & /usr/bin/python3 -I -B "$PSScriptRoot/native-host-facts.py" --phase after --since $startEpoch
        if ($LASTEXITCODE -ne 0) { throw 'After-fixture host facts unavailable' }
    } catch {
        if ($afterFailure) { [Console]::Error.WriteLine("Additional host-facts failure: $_") }
        else { $afterFailure = $_ }
    }
}
if ($primaryFailure) {
    if ($afterFailure) { [Console]::Error.WriteLine("Additional after-check failure: $afterFailure") }
    throw $primaryFailure
}
if ($afterFailure) { throw $afterFailure }
