$ErrorActionPreference = 'Stop'

$productionPath = Join-Path $PSScriptRoot 'admit-image.ps1'
$sharedPath = Join-Path $PSScriptRoot 'image-admission.ps1'
$workflowPath = Join-Path $PSScriptRoot '..\..\.github\workflows\release-recovery-v1.0.0.yml'
foreach ($path in @($productionPath,$sharedPath,$PSCommandPath)) { [void][scriptblock]::Create((Get-Content -Raw -LiteralPath $path)) }
. $sharedPath

$exact = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineDigest = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$wrongDigest = 'golang@sha256:1111111111111111111111111111111111111111111111111111111111111111'
$wrongRepository = 'example.invalid/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineJson = '{"Version":"29.7.2","ApiVersion":"1.52","Os":"linux","Arch":"amd64"}'
$listJson = '{"Repository":"golang","Digest":"sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452","ID":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}'
$commandVersion = 'docker|version|--format|{{json .Server}}'
$commandList = "docker|image|ls|--all|--no-trunc|--digests|--filter|reference=$exact|--format|{{json .}}"
$commandInspect = "docker|image|inspect|--format|{{json .RepoDigests}}|$exact"
$commandPull = "docker|pull|$exact"

function Assert-Equal($Actual,$Expected,[string]$Message) { if ($Actual -cne $Expected) { throw "$Message expected=[$Expected] actual=[$Actual]" } }
function Assert-Sequence([string[]]$Actual,[string[]]$Expected,[string]$Name) {
    if ($Actual.Count -ne $Expected.Count) { throw "$Name count mismatch expected=$($Expected.Count) actual=$($Actual.Count): $($Actual -join '; ')" }
    for ($index=0;$index-lt$Expected.Count;$index++) { if ($Actual[$index] -cne $Expected[$index]) { throw "$Name order mismatch at $index expected=[$($Expected[$index])] actual=[$($Actual[$index])]" } }
}
function New-TestResult([int]$ExitCode=0,[string]$StdOut='',[string]$StdErr='',[string]$Terminal='') {
    [pscustomobject]@{ ExitCode=$ExitCode; StdOut=$StdOut; StdErr=$StdErr; Terminal=$Terminal }
}
function Assert-TestCommandSuccess($Result,[string]$Operation) {
    if ($null -eq $Result -or $Result.PSObject.Properties.Name -notcontains 'Terminal' -or $Result.PSObject.Properties.Name -notcontains 'ExitCode' -or $Result.PSObject.Properties.Name -notcontains 'StdOut' -or $Result.PSObject.Properties.Name -notcontains 'StdErr') { throw "$Operation returned invalid lifecycle evidence" }
    if (![string]::IsNullOrEmpty($Result.Terminal)) { throw "$Operation returned terminal untrusted lifecycle evidence" }
    if ($Result.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($Result.StdErr)) { throw "$Operation failed and is untrusted" }
    if ([Text.UTF8Encoding]::new($false,$true).GetByteCount($Result.StdOut) -gt 131072 -or [Text.UTF8Encoding]::new($false,$true).GetByteCount($Result.StdErr) -gt 131072) { throw "$Operation exceeded its byte limit" }
}

function Invoke-TestAdmissionModel {
    param([string]$Name,[object[]]$Results,[string[]]$ExpectedEvents,[bool]$ExpectSuccess=$false,[bool]$ExpectPulled=$false)
    $queue=[Collections.Generic.Queue[object]]::new();foreach($result in $Results){$queue.Enqueue($result)}
    $events=[Collections.Generic.List[string]]::new();foreach($event in @('host-cache','host-cache','host-cache','host-crlf')){[void]$events.Add($event)}
    $laterActions=0;$success=$false;$pulled=$false
    function Take([string]$command){[void]$events.Add($command);if($queue.Count-eq 0){throw 'Unexpected command'};$queue.Dequeue()}
    try {
        $result=Take $commandVersion;Assert-TestCommandSuccess $result 'engine';[void](Assert-ReleaseEngineEvidence -Json $result.StdOut)
        $result=Take $commandList;Assert-TestCommandSuccess $result 'inventory';$state=Resolve-ReleaseImageListEvidence -Json $result.StdOut
        if($state -eq 'ConclusiveAbsent'){$result=Take $commandPull;Assert-TestCommandSuccess $result 'pull';$pulled=$true}
        $result=Take $commandInspect;Assert-TestCommandSuccess $result 'inspect';$digests=Read-ReleaseRepoDigestsEvidence -Json $result.StdOut;[void](Assert-ReleaseImageIdentityEvidence -Image $exact -RepoDigests $digests)
        $success=$true;$laterActions++
    } catch { if($ExpectSuccess){throw "$Name unexpectedly rejected: $($_.Exception.Message)"} }
    if(!$ExpectSuccess-and$success){throw "$Name unexpectedly admitted"}
    Assert-Equal $laterActions ($(if($ExpectSuccess){1}else{0})) "$Name later action"
    if($ExpectSuccess){Assert-Equal $pulled $ExpectPulled "$Name pull state"}
    Assert-Sequence $events.ToArray() $ExpectedEvents $Name
}

$hostEvents=@('host-cache','host-cache','host-cache','host-crlf')
$present=@($hostEvents+@($commandVersion,$commandList,$commandInspect));$absent=@($hostEvents+@($commandVersion,$commandList,$commandPull,$commandInspect))
$versionOnly=@($hostEvents+@($commandVersion));$listReached=@($hostEvents+@($commandVersion,$commandList))
Invoke-TestAdmissionModel 'pre-existing-exact' @((New-TestResult -StdOut $engineJson),(New-TestResult -StdOut $listJson),(New-TestResult -StdOut ('["'+$engineDigest+'"]'))) $present $true $false
Invoke-TestAdmissionModel 'conclusive-absence-single-pull' @((New-TestResult -StdOut $engineJson),(New-TestResult),(New-TestResult -StdOut 'pulled'),(New-TestResult -StdOut ('["'+$engineDigest+'"]'))) $absent $true $true
Invoke-TestAdmissionModel 'daemon-unavailable' @((New-TestResult -ExitCode 1 -StdErr 'daemon')) $versionOnly
Invoke-TestAdmissionModel 'permission-denial' @((New-TestResult -StdOut $engineJson),(New-TestResult -ExitCode 1 -StdErr 'permission')) $listReached
Invoke-TestAdmissionModel 'timeout-terminal' @((New-TestResult -Terminal 'timeout')) $versionOnly
Invoke-TestAdmissionModel 'overflow-terminal' @((New-TestResult -Terminal 'overflow')) $versionOnly
Invoke-TestAdmissionModel 'read-terminal' @((New-TestResult -Terminal 'read')) $versionOnly
Invoke-TestAdmissionModel 'cleanup-terminal' @((New-TestResult -Terminal 'cleanup uncertainty')) $versionOnly
Invoke-TestAdmissionModel 'invalid-invocation' @((New-TestResult -ExitCode 125 -StdErr 'invalid')) $versionOnly
Invoke-TestAdmissionModel 'unexpected-success-stderr' @((New-TestResult -StdOut $engineJson -StdErr 'warning')) $versionOnly
Invoke-TestAdmissionModel 'malformed-engine' @((New-TestResult -StdOut '{')) $versionOnly
Invoke-TestAdmissionModel 'malformed-list' @((New-TestResult -StdOut $engineJson),(New-TestResult -StdOut '{')) $listReached
Invoke-TestAdmissionModel 'ambiguous-list' @((New-TestResult -StdOut $engineJson),(New-TestResult -StdOut ($listJson+"`n"+$listJson))) $listReached
Invoke-TestAdmissionModel 'deceptive-absence' @((New-TestResult -StdOut $engineJson),(New-TestResult -StdOut $listJson),(New-TestResult -ExitCode 1 -StdErr 'No such image')) $present

foreach($case in @(
    @{N='empty';J=''},@{N='null';J='null'},@{N='scalar';J=('"'+$engineDigest+'"')},@{N='empty-array';J='[]'},@{N='malformed';J='['},
    @{N='exact-plus-wrong';J=('["'+$engineDigest+'","'+$wrongDigest+'"]')},@{N='duplicate';J=('["'+$engineDigest+'","'+$engineDigest+'"]')},
    @{N='alias';J=('["'+$exact+'"]')},@{N='wrong-repository';J=('["'+$wrongRepository+'"]')},@{N='wrong-digest';J=('["'+$wrongDigest+'"]')}
)) { Invoke-TestAdmissionModel $case.N @((New-TestResult -StdOut $engineJson),(New-TestResult -StdOut $listJson),(New-TestResult -StdOut $case.J)) $present }
Invoke-TestAdmissionModel 'pull-failure-no-retry' @((New-TestResult -StdOut $engineJson),(New-TestResult),(New-TestResult -ExitCode 1 -StdErr 'pull failed')) @($hostEvents+@($commandVersion,$commandList,$commandPull))
Invoke-TestAdmissionModel 'post-pull-inspect-failure' @((New-TestResult -StdOut $engineJson),(New-TestResult),(New-TestResult -StdOut 'pulled'),(New-TestResult -ExitCode 1 -StdErr 'inspect failed')) $absent
Invoke-TestAdmissionModel 'post-pull-mixed' @((New-TestResult -StdOut $engineJson),(New-TestResult),(New-TestResult -StdOut 'pulled'),(New-TestResult -StdOut ('["'+$engineDigest+'","'+$wrongDigest+'"]'))) $absent

$production=Get-Content -Raw -LiteralPath $productionPath;$shared=Get-Content -Raw -LiteralPath $sharedPath;$testSource=Get-Content -Raw -LiteralPath $PSCommandPath;$workflow=Get-Content -Raw -LiteralPath $workflowPath
foreach($forbidden in @('ReleaseDockerInvoker','ReleaseHostCacheCanaryInvoker','ReleaseHostOnlyCrlfProofInvoker','ReadToEnd','WaitForExit()','catch {}')) { if($production.Contains($forbidden)-or$shared.Contains($forbidden)){throw "Production retains forbidden boundary text: $forbidden"} }
if(($production.Split("@('pull',",[StringSplitOptions]::None).Count-1)-ne 1){throw 'Closed entrypoint must contain exactly one exact pull construction'}
if($shared.Contains("@('pull',")){throw 'Dot-sourceable helper retains pull authority'}
foreach($required in @('131072','15000','2000','UTF8Encoding','ReadAsync','Kill($true)','ArgumentList.Add','Environment.Clear','Get-AuthenticodeSignature','/usr/bin/docker')){if(!$production.Contains($required)){throw "Production omits required closed-boundary marker: $required"}}
foreach($forbiddenParam in @('ScriptBlock','Invoker','Callback','AllowImagePull','ExecutablePath')){if($production.Substring(0,$production.IndexOf('$ErrorActionPreference')).Contains($forbiddenParam)){throw "Production parameter block exposes $forbiddenParam"}}
foreach($required in @('-SourceRepository','-SourceRevision','-WorkingDirectory')){if(!$workflow.Contains($required)){throw "Workflow omits $required"}}

$temporaryRoot=Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-003-'+[guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $temporaryRoot)
try {
    $dotSourceProbe=Join-Path $temporaryRoot 'dot-source-probe.ps1'
    $escaped=$productionPath.Replace("'","''")
    Set-Content -LiteralPath $dotSourceProbe -Encoding utf8NoBOM -Value @"
function docker { throw 'ambient function reached' }
Set-Alias -Name Invoke-ClosedNativeProcess -Value Write-Output
`$script:ReleaseDockerInvoker={ throw 'ambient callback reached' }
`$env:PATH='$(Join-Path $temporaryRoot 'hostile-path')'
try { . '$escaped' -CacheDirectory '$escaped' -SourceRepository '$escaped' -SourceRevision '0000000000000000000000000000000000000000' -WorkingDirectory '$escaped' } catch { if (`$_.Exception.Message -notmatch 'cannot be dot-sourced') { throw }; exit 0 }
exit 1
"@
    & (Get-Process -Id $PID).Path -NoProfile -NonInteractive -File $dotSourceProbe
    if($LASTEXITCODE-ne 0){throw 'Dot-source and ambient-state bypass probe failed'}

    $tokens=$null;$errors=$null;$ast=[Management.Automation.Language.Parser]::ParseFile($productionPath,[ref]$tokens,[ref]$errors)
    if($errors.Count){throw 'Production AST parse failed'}
    $runner=$ast.Find({param($node)$node-is[Management.Automation.Language.FunctionDefinitionAst]-and$node.Name-ceq'Invoke-ClosedNativeProcess'},$true)
    if($null-eq$runner){throw 'Private bounded runner is absent'}
    . ([scriptblock]::Create($runner.Extent.Text))
    $dockerOutputLimit=131072;$dockerBudgetMilliseconds=15000;$cleanupGraceMilliseconds=2000
    $native=(Get-Process -Id $PID).Path
    if($IsWindows){$sig=Microsoft.PowerShell.Security\Get-AuthenticodeSignature -LiteralPath $native;if($sig.Status-ne[Management.Automation.SignatureStatus]::Valid){throw 'Native fixture host is not signed system tooling'}}
    $fixture=Join-Path $temporaryRoot 'native-fixture.ps1'
    Set-Content -LiteralPath $fixture -Encoding utf8NoBOM -Value @'
param([string]$Mode,[int]$Count=0,[string]$PidFile='')
$out=[Console]::OpenStandardOutput();$err=[Console]::OpenStandardError()
function Write-Bytes($stream,[byte[]]$bytes){$stream.Write($bytes,0,$bytes.Length);$stream.Flush()}
switch($Mode){
'stdout'{Write-Bytes $out ([byte[]](,[byte]97*$Count))}
'stderr'{Write-Bytes $err ([byte[]](,[byte]98*$Count))}
'both'{for($i=0;$i-lt$Count;$i+=1024){$n=[Math]::Min(1024,$Count-$i);Write-Bytes $out ([byte[]](,[byte]97*$n));Write-Bytes $err ([byte[]](,[byte]98*$n))}}
'split-utf8'{Write-Bytes $out ([byte[]](0xE2,0x82));Start-Sleep -Milliseconds 25;Write-Bytes $out ([byte[]](0xAC))}
'invalid-utf8'{Write-Bytes $out ([byte[]](0xC3,0x28))}
'incomplete-utf8'{Write-Bytes $out ([byte[]](0xE2,0x82))}
'nonzero'{exit 7}
'hang'{Start-Sleep -Seconds 30}
'child'{ $p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang') -PassThru -NoNewWindow;if($PidFile){Set-Content -LiteralPath $PidFile -Value $p.Id};Start-Sleep -Seconds 30 }
'grandchild'{ $childFile=$PidFile+'.child';$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','child','-PidFile',$childFile) -PassThru -NoNewWindow;Set-Content -LiteralPath $PidFile -Value $p.Id;Start-Sleep -Seconds 30 }
default{throw 'unknown fixture mode'}
}
'@
    function RunFixture([string]$mode,[int]$count=0,[string]$pidFile=''){
        $fixtureArguments=@('-NoProfile','-NonInteractive','-File',$fixture,'-Mode',$mode,'-Count',[string]$count)
        if(![string]::IsNullOrEmpty($pidFile)){$fixtureArguments+=@('-PidFile',$pidFile)}
        Invoke-ClosedNativeProcess -ExecutablePath $native -Arguments $fixtureArguments -Operation "fixture-$mode"
    }
    foreach($n in @(131071,131072)){ $r=RunFixture 'stdout' $n;Assert-Equal $r.StdOutByteCount $n "stdout-$n";$r=RunFixture 'stderr' $n;Assert-Equal $r.StdErrByteCount $n "stderr-$n" }
    $r=RunFixture 'both' 131072;Assert-Equal $r.StdOutByteCount 131072 'both stdout';Assert-Equal $r.StdErrByteCount 131072 'both stderr'
    $r=RunFixture 'split-utf8';Assert-Equal $r.StdOut '€' 'split UTF-8'
    $r=RunFixture 'nonzero';Assert-Equal $r.ExitCode 7 'nonzero exit'
    foreach($case in @(@{M='stdout';C=131073},@{M='stderr';C=131073},@{M='both';C=131073},@{M='invalid-utf8';C=0},@{M='incomplete-utf8';C=0},@{M='hang';C=0})){$failed=$false;$watch=[Diagnostics.Stopwatch]::StartNew();try{[void](RunFixture $case.M $case.C)}catch{$failed=$true};if(!$failed){throw "$($case.M) unexpectedly trusted"};if($watch.ElapsedMilliseconds-gt 17500){throw "$($case.M) exceeded bounded return"}}
    $missing=Join-Path $temporaryRoot 'missing-native.exe';$failed=$false;try{[void](Invoke-ClosedNativeProcess -ExecutablePath $missing -Arguments @() -Operation 'fixture-start-failure')}catch{$failed=$true};if(!$failed){throw 'Start failure unexpectedly trusted'}
    foreach($mode in @('child','grandchild')){$pidFile=Join-Path $temporaryRoot "$mode.pid";$failed=$false;try{[void](RunFixture $mode 0 $pidFile)}catch{$failed=$true};if(!$failed){throw "$mode unexpectedly trusted"};foreach($file in @($pidFile,$pidFile+'.child')){if(Test-Path $file){$fixturePid=[int](Get-Content -Raw $file);if(Get-Process -Id $fixturePid -ErrorAction SilentlyContinue){throw "$mode left a live descendant $fixturePid"}}}}
    Write-Output 'Image admission iteration-003 PASS private=BOUND state-matrix=PASS byte-caps=LIVE utf8=STRICT timeout=MONOTONIC process-tree=TERMINAL'
} finally {
    $resolved=[IO.Path]::GetFullPath($temporaryRoot);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath());if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)){throw "Unsafe cleanup path: $resolved"};if(Test-Path -LiteralPath $resolved){Remove-Item -LiteralPath $resolved -Recurse -Force}
}
