param(
    [Parameter(Mandatory)][string]$BaselinePath,
    [Parameter(Mandatory)][ValidatePattern('^[0-9a-f]{64}$')][string]$BaselineSHA256
)
$ErrorActionPreference = 'Stop'
$startEpoch = [DateTimeOffset]::UtcNow.ToUnixTimeSeconds()
$primaryFailure = $null
$afterFailure = $null
try {
    & /usr/bin/python3 -I -B "$PSScriptRoot/native-profile-context.py" check $BaselinePath $BaselineSHA256 before
    if ($LASTEXITCODE -ne 0) { throw 'Named profile context admission failed' }
    & /usr/bin/python3 -I -B "$PSScriptRoot/native-host-facts.py" --phase before --since $startEpoch
    if ($LASTEXITCODE -ne 0) { throw 'Before-fixture host facts unavailable' }
    & "$PSScriptRoot/test-docker-execution.ps1" -InternalMode CleanNative
} catch {
    $primaryFailure = $_
} finally {
    try {
        & /usr/bin/python3 -I -B "$PSScriptRoot/native-profile-context.py" check $BaselinePath $BaselineSHA256 after
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
