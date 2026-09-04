param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [Parameter(Mandatory = $true)][string]$SourceRepository,
    [Parameter(Mandatory = $true)][string]$SourceRevision,
    [Parameter(Mandatory = $true)][string]$WorkingDirectory
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'image-admission.ps1')
$result = Invoke-ReleaseImageBootstrap @PSBoundParameters
Write-Output "Release image admission PASS image=$($result.Image) pulled=$($result.Pulled) host-identity=$($result.HostIdentityMode)"
