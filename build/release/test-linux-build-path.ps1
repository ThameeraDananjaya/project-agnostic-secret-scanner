$ErrorActionPreference = 'Stop'
$assertions = 0
function Check([bool]$Value, [string]$Message) { $script:assertions++; if (!$Value) { throw $Message } }
function Reject([scriptblock]$Body) { $failed=$false; try { & $Body } catch { $failed=$true }; Check $failed 'Expected rejection' }
$tokens=$null; $errors=$null
$tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'docker-execution.ps1'),[ref]$tokens,[ref]$errors)
if ($errors.Count) { throw 'Production boundary does not parse' }
foreach ($name in @('Resolve-ContainerUser','Require-Path','ConvertTo-LF','Get-CachePayload','Get-AcquisitionPayload','Get-BuildPayload','Get-PackagePayload')) {
    $node=$tree.Find({param($n) $n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq$name},$true)
    if ($null-eq$node) { throw "Missing actual production function: $name" }
    . ([scriptblock]::Create($node.Extent.Text))
}
$assignment=$tree.Find({param($n) $n-is[Management.Automation.Language.AssignmentStatementAst]-and$n.Left.Extent.Text-ceq'$arguments'-and$n.Right.Extent.Text.StartsWith('switch')},$true)
if ($null-eq$assignment) { throw 'Missing actual fixed argument table' }
foreach ($pair in @(@(1001,1001),@(0,0),@(1001,0))) {
    $resolved=Resolve-ContainerUser $pair[0] $pair[1] $true
    Check ($resolved.Arguments.Count-eq 2-and$resolved.Arguments[1]-ceq"$($pair[0]):$($pair[1])") 'Boxed nullable values lost'
    Check ($resolved.TmpfsSuffix-ceq",mode=0700,uid=$($pair[0]),gid=$($pair[1])") 'Tmpfs identity lost'
}
foreach ($pair in @(@($null,1),@(1,$null),@(-1,1),@(1,-1))) { Reject { Resolve-ContainerUser $pair[0] $pair[1] $true } }
Reject { Resolve-ContainerUser $null $null $true }
$absent=Resolve-ContainerUser $null $null $false
Check ($absent.Arguments.Count-eq 0-and$absent.TmpfsSuffix-ceq'') 'Optional non-Linux identity changed'

$temporaryRoot=Join-Path ([IO.Path]::GetTempPath()) ('pscan-linux-build-path-'+[guid]::NewGuid().ToString('N'))
[void][IO.Directory]::CreateDirectory($temporaryRoot)
try {
    $file=Join-Path $temporaryRoot 'input'; [IO.File]::WriteAllText($file,'owned test input')
    foreach ($name in @('SourceRoot','CacheDirectory','RawOutputDirectory','LinuxStageDirectory','WindowsStageDirectory','DistributionDirectory')) { Set-Variable -Name $name -Value $temporaryRoot }
    foreach ($name in @('RunnerGoArchive','EngineGoArchive','GitleaksArchive','ProductArchive','ProductBlobManifest','ProductPathManifest','ProductModeManifest','ToolingArchive','ToolingBlobManifest','ToolingPathManifest','ToolingModeManifest')) { Set-Variable -Name $name -Value $file }
    foreach ($name in @('ProductRevision','ToolingRevision','ToolingTree')) { Set-Variable -Name $name -Value ('a'*40) }
    $ProductArchiveSHA256='a'*64; $ToolingArchiveSHA256='b'*64; $Created='2026-09-17T00:00:00Z'; $SourceDateEpoch='0'; $ProductFileCount=1; $ToolingFileCount=1
    $image='docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
    foreach ($Operation in @('ContainerCacheProof','DependencyAcquisition','ReleaseBuild','ReleasePackage')) {
        $containerUser=Resolve-ContainerUser 1001 1002 $true; $ReadOnlyCache=$false
        . ([scriptblock]::Create($assignment.Extent.Text))
        Check (@($arguments|Where-Object{$_-ceq'--user'}).Count-eq 1-and$arguments[$arguments.IndexOf('--user')+1]-ceq'1001:1002') "Missing actual $Operation user mapping"
        Check ($arguments[$arguments.IndexOf('--tmpfs')+1].EndsWith(',mode=0700,uid=1001,gid=1002')) "Missing actual $Operation tmpfs mapping"
        Check ($arguments[$arguments.IndexOf('--cap-drop')+1]-ceq'ALL'-and$arguments-contains'no-new-privileges'-and$arguments-contains'--read-only'-and$arguments-contains'--pull=never') 'Container restrictions changed'
        if ($Operation-ne'DependencyAcquisition') { Check ($arguments[$arguments.IndexOf('--network')+1]-ceq'none') 'Offline operation gained network' }
        if ($Operation-eq'ReleaseBuild') {
            Check ($arguments-contains"type=bind,src=$temporaryRoot,dst=/gomodcache,readonly") 'Build cache is writable'
            Check ($arguments[$arguments.IndexOf('--workdir')+1]-ceq'/work') 'Build requests a daemon-created subdirectory instead of owned tmpfs root'
        }
        if ($Operation-eq'ContainerCacheProof') {
            $ReadOnlyCache=$true; . ([scriptblock]::Create($assignment.Extent.Text))
            Check ($arguments-contains"type=bind,src=$temporaryRoot,dst=/gomodcache,readonly") 'Negative canary lost read-only mount'
        }
    }
    . (Join-Path $PSScriptRoot 'cache-canary.ps1')
    foreach($reason in @('Permission denied','Read-only file system','timeout',('Permission denied'+"`n"+'unexpected second error'))){
        $r=[pscustomobject]@{ContainmentEmpty=$true;DaemonContainerRemoved=$true;DaemonContainerID=('b'*64);DockerSHA256=('a'*64);Operation='ContainerCacheProof';ExitCode=74;StdOut=('{"schema":"pscan-cache-write-v1","phase":"create-write","outcome":"failed"}'+"`n");StdErr="/bin/sh: 7: cannot create /gomodcache/.pscan-cache-canary-1: $reason`n"}
        $errorID='';try{Assert-ReleaseCacheResult $r ('a'*64)}catch{$errorID=$_.FullyQualifiedErrorId}
        Check ($errorID-ne''-and($errorID-like'PSCAN_CACHE_EXPECTED_DENIAL*')-eq($reason-in@('Permission denied','Read-only file system'))) 'Unexpected failure was reclassified as expected cache denial'
    }
    foreach($field in @('ContainmentEmpty','DaemonContainerRemoved','DockerSHA256','Operation','ExitCode','StdOut')){
        $r=[pscustomobject]@{ContainmentEmpty=$true;DaemonContainerRemoved=$true;DaemonContainerID=('b'*64);DockerSHA256=('a'*64);Operation='ContainerCacheProof';ExitCode=74;StdOut=('{"schema":"pscan-cache-write-v1","phase":"create-write","outcome":"failed"}'+"`n");StdErr="/bin/sh: 7: cannot create /gomodcache/.pscan-cache-canary-1: Permission denied`n"}
        $r.$field=switch($field){'DockerSHA256'{'c'*64}'Operation'{'wrong'}'ExitCode'{199}'StdOut'{'{malformed'}default{$false}}
        $errorID='';try{Assert-ReleaseCacheResult $r ('a'*64)}catch{$errorID=$_.FullyQualifiedErrorId}
        Check ($errorID-ne''-and$errorID-notlike'PSCAN_CACHE_EXPECTED_DENIAL*') 'Untrusted lifecycle became an expected negative control'
    }
    . (Join-Path $PSScriptRoot 'host-cache-canary.ps1')
    function New-CanaryFailure([Exception]$Inner) {
        $failure=[IO.IOException]::new('inert canary record',$Inner)
        $failure.Data['PSCANCachePhase']='create-write'
        $failure.Data['PSCANCacheCleanup']='not-needed'
        $failure.Data['PSCANCacheOwnedFile']=$false
        return $failure
    }
    $errno13=[IO.IOException]::new('arbitrary non-permission message',13)
    foreach ($inner in @(
        $errno13,
        [UnauthorizedAccessException]::new('inert unauthorized wrapper',$errno13),
        [Management.Automation.MethodInvocationException]::new('inert method wrapper',$errno13),
        [Management.Automation.MethodInvocationException]::new('inert method wrapper',[UnauthorizedAccessException]::new('inert unauthorized wrapper',$errno13))
    )) {
        $failure=New-CanaryFailure $inner
        Check (Test-LinuxHostCacheWriteDenial $failure) 'Verified Linux errno 13 chain rejected'
        Check (Test-LinuxHostCacheWriteDenial ([Management.Automation.MethodInvocationException]::new('outer wrapper',$failure))) 'Outer method wrapper lost exact canary record'
    }
    foreach ($inner in @(
        [IO.IOException]::new('Permission denied'),
        [UnauthorizedAccessException]::new('Permission denied'),
        [IO.IOException]::new('Permission denied',30), # EROFS belongs to a mounted-filesystem control, not chmod.
        [IO.IOException]::new('Permission denied',1),
        [IO.IOException]::new('Permission denied',5),
        [IO.IOException]::new('Permission denied',9),
        [IO.IOException]::new('Permission denied',17),
        [IO.IOException]::new('Permission denied',28),
        [IO.IOException]::new('Permission denied',-2147024891), # Windows access denied HRESULT.
        [IO.IOException]::new('Permission denied',-2147024877), # Windows write-protected HRESULT.
        [Exception]::new('unknown wrapper',$errno13),
        [IO.IOException]::new('unstructured intermediate I/O failure',$errno13),
        [UnauthorizedAccessException]::new('wrong errno',[IO.IOException]::new('Permission denied',9))
    )) { Check (!(Test-LinuxHostCacheWriteDenial (New-CanaryFailure $inner))) 'Unproved platform denial admitted' }
    Check (!(Test-LinuxHostCacheWriteDenial $errno13)) 'Raw errno without canary phase/ownership record admitted'
    Check (!(Test-LinuxHostCacheWriteDenial $null)) 'Missing failure admitted'
    foreach ($phase in @('atomic-rename','read','delete','verify-cleanup','')) {
        $failure=New-CanaryFailure $errno13; $failure.Data['PSCANCachePhase']=$phase
        Check (!(Test-LinuxHostCacheWriteDenial $failure)) 'Later or missing phase admitted'
    }
    foreach ($cleanup in @('completed','failed','unproved','')) {
        $failure=New-CanaryFailure $errno13; $failure.Data['PSCANCacheCleanup']=$cleanup
        Check (!(Test-LinuxHostCacheWriteDenial $failure)) 'Uncertain or required cleanup admitted'
    }
    foreach ($owned in @($true,'false',0,$null)) {
        $failure=New-CanaryFailure $errno13; $failure.Data['PSCANCacheOwnedFile']=$owned
        Check (!(Test-LinuxHostCacheWriteDenial $failure)) 'Owned file or malformed ownership proof admitted'
    }
    $failure=New-CanaryFailure (New-CanaryFailure $errno13)
    Check (!(Test-LinuxHostCacheWriteDenial $failure)) 'Duplicate canary metadata admitted'
    $failure=New-CanaryFailure $errno13
    for($n=0;$n-lt 6;$n++) { $failure=[Management.Automation.MethodInvocationException]::new('bounded wrapper',$failure) }
    Check (Test-LinuxHostCacheWriteDenial $failure) 'Eight-exception bounded chain rejected'
    $failure=[Management.Automation.MethodInvocationException]::new('excess wrapper',$failure)
    Check (!(Test-LinuxHostCacheWriteDenial $failure)) 'Oversize exception chain admitted'
    $positive=Join-Path $temporaryRoot 'positive'; [void][IO.Directory]::CreateDirectory($positive)
    $identity=Invoke-HostCacheCanary $positive
    Check ($null-ne$identity-and[IO.Directory]::GetFileSystemEntries($positive).Count-eq 0) 'Actual host canary failed or left files'
    $failure=$null; try { Invoke-HostCacheCanary $file } catch { $failure=$_.Exception }
    Check ($null-ne$failure-and$null-ne$failure.InnerException) 'Host canary lost original exception'
    Check ($failure.Message.Contains('"phase":"create-write"')-and$failure.Message.Contains('"cleanup":"not-needed"')) 'Host failure did not identify operation and cleanup'
    Check ([IO.File]::ReadAllText($file)-ceq'owned test input') 'Failed canary modified unrelated input'
} finally {
    $resolved=[IO.Path]::GetFullPath($temporaryRoot)
    $temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if (!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)-or![IO.Path]::GetFileName($resolved).StartsWith('pscan-linux-build-path-')) { throw 'Unsafe test cleanup' }
    [IO.Directory]::Delete($resolved,$true)
}
Write-Output "Linux build path inert PASS assertions=$script:assertions actual-host-canary=PASS Docker-calls=0 Linux-runtime-proof=UNAVAILABLE"
