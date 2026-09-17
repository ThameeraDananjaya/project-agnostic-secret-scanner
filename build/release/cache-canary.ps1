function Invoke-ReleaseCacheCanary {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$Image,
        [Parameter(Mandatory = $true)][string]$ModuleCache,
        [Parameter(Mandatory = $true)][ValidatePattern('^[0-9a-f]{64}$')][string]$ExpectedDockerSHA256,
        [Nullable[int]]$HostUID,
        [Nullable[int]]$HostGID,
        [switch]$ReadOnlyCache
    )

    $exact = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
    if ($Image -cne $exact) { throw 'Cache proof requires the exact admitted image identity' }
    $parameters = @{
        Operation='ContainerCacheProof'; ExpectedDockerSHA256=$ExpectedDockerSHA256
        CacheDirectory=$ModuleCache; ReadOnlyCache=$ReadOnlyCache
    }
    if ($HostUID.HasValue -or $HostGID.HasValue) {
        if (!$HostUID.HasValue -or !$HostGID.HasValue) { throw 'Cache proof requires a complete UID/GID pair' }
        $parameters.HostUID=$HostUID.Value; $parameters.HostGID=$HostGID.Value
    }
    $json = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') @parameters
    $result = ($json -join "`n") | ConvertFrom-Json
    if ($result.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($result.StdErr) -or !$result.ContainmentEmpty -or $result.DockerSHA256 -cne $ExpectedDockerSHA256) {
        throw 'Module-cache canary failed closed at the Docker boundary'
    }
}
