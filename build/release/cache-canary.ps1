function Assert-ReleaseCacheResult($Result,[string]$ExpectedDockerSHA256) {
    if ($Result.ContainmentEmpty-isnot[bool]-or!$Result.ContainmentEmpty-or$Result.DaemonContainerRemoved-isnot[bool]-or!$Result.DaemonContainerRemoved-or$Result.DaemonContainerID-cnotmatch'^[0-9a-f]{64}$'-or$Result.DockerSHA256-cne$ExpectedDockerSHA256-or
        $Result.Operation-cne'ContainerCacheProof'-or($Result.ExitCode-isnot[int]-and$Result.ExitCode-isnot[long])-or$Result.StdOut-isnot[string]-or$Result.StdErr-isnot[string]) { throw 'Cache proof lifecycle or executable identity is untrusted' }
    if ($Result.ExitCode-ne 0-or![string]::IsNullOrEmpty($Result.StdErr)) {
        # The fixed /bin/sh payload's create-write denial is the only admitted
        # negative control. Timeout, malformed evidence and cleanup uncertainty
        # must propagate, never become an expected cache rejection.
        $writeFailure='{"schema":"pscan-cache-write-v1","phase":"create-write","outcome":"failed"}'+"`n"
        if ($Result.ExitCode-eq 74-and$Result.StdOut-ceq$writeFailure-and$Result.StdErr-cmatch'^/bin/sh: [0-9]+: cannot create /gomodcache/\.pscan-cache-canary-[0-9]+: (Permission denied|Read-only file system)\r?\n?$') {
            $outcome=[pscustomobject]@{schema='pscan-cache-negative-v1';phase='create-write';outcome='expected-denial';container=$Result.DaemonContainerID;nativeEmpty=$true;daemonRemoved=$true;exitCode=74}
            $record=[Management.Automation.ErrorRecord]::new([IO.IOException]::new('Fixed cache write was denied after proved native and daemon cleanup'),'PSCAN_CACHE_EXPECTED_DENIAL',[Management.Automation.ErrorCategory]::WriteError,$outcome)
            throw $record
        }
        throw 'Module-cache canary failed unexpectedly at the Docker boundary'
    }
}

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
    if ($null -ne $HostUID -or $null -ne $HostGID) {
        if ($null -eq $HostUID -or $null -eq $HostGID -or $HostUID -lt 0 -or $HostGID -lt 0) { throw 'Cache proof requires a complete non-negative UID/GID pair' }
        $parameters.HostUID=$HostUID; $parameters.HostGID=$HostGID
    } elseif ($IsLinux) {
        throw 'Linux cache proof requires the admitted host UID/GID pair'
    }
    $json = & (Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1') @parameters
    $result = ($json -join "`n") | ConvertFrom-Json
    Assert-ReleaseCacheResult $result $ExpectedDockerSHA256
}
