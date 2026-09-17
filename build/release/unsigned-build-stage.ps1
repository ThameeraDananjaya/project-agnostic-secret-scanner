param([Parameter(Mandatory)][ValidateSet('NativePrerequisites','ImageAdmission','ContainerCache','ContainerCRLF','Acquire','BuildA','BuildB','Compare')][string]$Stage)
$ErrorActionPreference='Stop'
$baseline=Join-Path $env:RUNNER_TEMP ("pscan-unsigned-profile-"+$Stage+".json")
& /usr/bin/python3 -I -B "$PSScriptRoot/unsigned-build-profile-context.py" snapshot $baseline
if($LASTEXITCODE-ne 0){throw 'Unsigned build ordinary caller baseline unavailable'}
$baselineHash=(Get-FileHash -LiteralPath $baseline -Algorithm SHA256).Hash.ToLowerInvariant()
& /usr/bin/aa-exec --profile pscan-native-diagnostic -- /opt/microsoft/powershell/7/pwsh -NoLogo -NoProfile -NonInteractive -File "$PSScriptRoot/unsigned-build-profile-driver.ps1" -Stage $Stage -BaselinePath $baseline -BaselineSHA256 $baselineHash
if($LASTEXITCODE-ne 0){throw "Unsigned build stage failed: $Stage"}
