function Invoke-HostCacheCanary {
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
        $gidText = (& id -g).Trim()
        if ($LASTEXITCODE -ne 0 -or $uidText -notmatch '^\d+$' -or $gidText -notmatch '^\d+$') {
            throw 'Linux host numeric UID/GID discovery failed before image or dependency acquisition'
        }
        $hostUID = [int]$uidText
        $hostGID = [int]$gidText
        $owner = (& stat -c '%u:%g' -- $resolvedCache).Trim()
        if ($LASTEXITCODE -ne 0 -or $owner -ne "${hostUID}:${hostGID}") {
            throw "Linux cache ownership does not match the invoking host identity: $resolvedCache"
        }
        $identityMode = 'linux-host-numeric-uid-gid'
    }

    $token = [guid]::NewGuid().ToString('N')
    $canary = Join-Path $resolvedCache ".pscan-host-cache-canary-$token"
    $renamed = Join-Path $resolvedCache ".pscan-host-cache-canary-ready-$token"
    try {
        [IO.File]::WriteAllText($canary, "PSCAN-06-C2-HOST-CACHE-CANARY`n", [Text.UTF8Encoding]::new($false))
        Move-Item -LiteralPath $canary -Destination $renamed
        $value = [IO.File]::ReadAllText($renamed)
        if ($value -ne "PSCAN-06-C2-HOST-CACHE-CANARY`n") {
            throw 'Host cache canary read did not reproduce the exact bytes'
        }
        Remove-Item -LiteralPath $renamed
        if ((Test-Path -LiteralPath $canary) -or (Test-Path -LiteralPath $renamed)) {
            throw 'Host cache canary cleanup was incomplete'
        }
    } catch {
        Remove-Item -LiteralPath $canary,$renamed -Force -ErrorAction SilentlyContinue
        throw 'Host cache write, atomic rename, read and delete canary failed before image or dependency acquisition'
    }

    return [pscustomobject]@{ IdentityMode=$identityMode; HostUID=$hostUID; HostGID=$hostGID }
}
