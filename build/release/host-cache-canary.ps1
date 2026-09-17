function Get-HostCacheIdentity {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$CacheDirectory
    )

    $resolvedCache = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)
    $identityMode = 'windows-invoking-host'
    $hostUID = $null
    $hostGID = $null

    if ($IsLinux) {
        $uidText = (& id -u).Trim()
        if ($LASTEXITCODE -ne 0 -or $uidText -notmatch '^\d+$') {
            throw 'Linux host numeric UID/GID discovery failed before image or dependency acquisition'
        }
        $gidText = (& id -g).Trim()
        if ($LASTEXITCODE -ne 0 -or $gidText -notmatch '^\d+$') { throw 'Linux host numeric GID discovery failed' }
        $hostUID = [int]$uidText
        $hostGID = [int]$gidText
        $owner = (& stat -c '%u:%g' -- $resolvedCache).Trim()
        if ($LASTEXITCODE -ne 0-or$owner-notmatch'^\d+:\d+$') { throw 'Linux cache owner inspection failed' }
        if ($owner -ne "${hostUID}:${hostGID}") {
            throw [Management.Automation.ErrorRecord]::new([IO.IOException]::new('Linux cache ownership differs from invoking host'),'PSCAN_HOST_CACHE_OWNER_MISMATCH',[Management.Automation.ErrorCategory]::PermissionDenied,$null)
        }
        $identityMode = 'linux-host-numeric-uid-gid'
    }

    return [pscustomobject]@{ IdentityMode=$identityMode; HostUID=$hostUID; HostGID=$hostGID }
}

function Invoke-HostCacheCanary {
    [CmdletBinding()]
    param([Parameter(Mandatory = $true)][string]$CacheDirectory)
    $resolvedCache = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $CacheDirectory).Path)
    $identity = Get-HostCacheIdentity -CacheDirectory $resolvedCache
    $token = [guid]::NewGuid().ToString('N')
    $canary = Join-Path $resolvedCache ".pscan-host-cache-canary-$token"
    $renamed = Join-Path $resolvedCache ".pscan-host-cache-canary-ready-$token"
    $phase = 'create-write'
    $ownedPath = $null
    try {
        $stream = [IO.File]::Open($canary, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write, [IO.FileShare]::None)
        $ownedPath = $canary
        try {
            $bytes = [Text.UTF8Encoding]::new($false).GetBytes("PSCAN-06-C2-HOST-CACHE-CANARY`n")
            $stream.Write($bytes, 0, $bytes.Length)
        } finally { $stream.Dispose() }
        $phase = 'atomic-rename'
        [IO.File]::Move($canary, $renamed, $false)
        $ownedPath = $renamed
        $phase = 'read'
        $value = [IO.File]::ReadAllText($renamed)
        if ($value -ne "PSCAN-06-C2-HOST-CACHE-CANARY`n") {
            throw 'Host cache canary read did not reproduce the exact bytes'
        }
        $phase = 'delete'
        [IO.File]::Delete($renamed)
        $phase = 'verify-cleanup'
        if ((Test-Path -LiteralPath $canary) -or (Test-Path -LiteralPath $renamed)) {
            throw 'Host cache canary cleanup was incomplete'
        }
    } catch {
        $failure = $_.Exception
        $cleanup = 'not-needed'
        if ($null -ne $ownedPath) {
            try { [IO.File]::Delete($ownedPath); $cleanup = 'completed' }
            catch { $cleanup = 'failed' }
        }
        $cause = $failure
        while ($null -ne $cause.InnerException) { $cause = $cause.InnerException }
        $message = $cause.Message
        if ($message.Length -gt 512) { $message = $message.Substring(0, 512) }
        $diagnostic = [ordered]@{
            schema='pscan-host-cache-failure-v1'; phase=$phase
            exceptionType=$cause.GetType().FullName; hresult=$cause.HResult
            message=$message; cleanup=$cleanup
        } | ConvertTo-Json -Compress
        $canaryFailure=[IO.IOException]::new("Host cache canary failed before acquisition: $diagnostic", $failure)
        $canaryFailure.Data['PSCANCachePhase']=$phase;$canaryFailure.Data['PSCANCacheCleanup']=$cleanup
        throw $canaryFailure
    }

    return $identity
}
