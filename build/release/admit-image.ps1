param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [switch]$AllowImagePull
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'image-admission.ps1')
$result = Invoke-ReleaseImageAdmission -CacheDirectory $CacheDirectory -AllowImagePull:$AllowImagePull
Write-Output "Release image admission PASS image=$($result.Image) pulled=$($result.Pulled) host-identity=$($result.HostIdentityMode)"
