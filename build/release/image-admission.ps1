$script:ReleaseImage = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'

function ConvertTo-CanonicalReleaseImageReference([string]$Reference) {
    if ($Reference -notmatch '^(?<repository>[^@]+)@(?<digest>sha256:[0-9a-f]{64})$') {
        throw 'Build image reference must bind one lowercase SHA-256 digest'
    }
    $repository = $Matches.repository.ToLowerInvariant()
    if ($repository -eq 'golang') { $repository = 'docker.io/library/golang' }
    elseif ($repository -eq 'library/golang') { $repository = 'docker.io/library/golang' }
    if ($repository -ne 'docker.io/library/golang') {
        throw 'Build image repository identity is not the canonical Docker Hub library/golang repository'
    }
    return "$repository@$($Matches.digest)"
}

function Assert-AdmittedReleaseImage {
    [CmdletBinding()]
    param([string]$Image = $script:ReleaseImage)

    if ($Image -cne $script:ReleaseImage) { throw 'Build image must use the exact canonical Correction C2 repository and digest reference' }
    $canonical = ConvertTo-CanonicalReleaseImageReference -Reference $Image
    $json = & docker image inspect --format '{{json .RepoDigests}}' $canonical 2>$null
    if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace(($json -join ''))) {
        throw 'Pinned build image is absent'
    }
    try { $repoDigests = @(($json -join '') | ConvertFrom-Json) } catch { throw 'Docker image identity output is invalid' }
    return Assert-ReleaseImageIdentityEvidence -Image $canonical -RepoDigests $repoDigests
}

function Assert-ReleaseImageIdentityEvidence {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$Image,
        [Parameter(Mandatory = $true)][AllowEmptyCollection()][string[]]$RepoDigests
    )

    if ($Image -cne $script:ReleaseImage) { throw 'Build image must use the exact canonical Correction C2 repository and digest reference' }
    $canonical = ConvertTo-CanonicalReleaseImageReference -Reference $Image
    $matches = @($RepoDigests | ForEach-Object { ConvertTo-CanonicalReleaseImageReference -Reference $_ } | Where-Object { $_ -eq $canonical })
    if ($matches.Count -ne 1) { throw 'Docker image identity does not prove the exact canonical repository and digest' }
    return $canonical
}

function Invoke-ReleaseImageAdmission {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$CacheDirectory,
        [switch]$AllowImagePull
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
        $hostProof = Invoke-HostCacheCanary -CacheDirectory $path
    }
    try {
        $canonical = Assert-AdmittedReleaseImage
        $pulled = $false
    } catch {
        if (!$AllowImagePull -or $_.Exception.Message -ne 'Pinned build image is absent') { throw }
        & docker pull $script:ReleaseImage
        if ($LASTEXITCODE -ne 0) { throw 'Pinned build image pull failed' }
        $canonical = Assert-AdmittedReleaseImage
        $pulled = $true
    }
    return [pscustomobject]@{
        Image=$canonical; Pulled=$pulled; HostIdentityMode=$hostProof.IdentityMode
        HostUID=$hostProof.HostUID; HostGID=$hostProof.HostGID
        ProvedCachePaths=@($cacheRoot,$downloads,$moduleCache)
    }
}
