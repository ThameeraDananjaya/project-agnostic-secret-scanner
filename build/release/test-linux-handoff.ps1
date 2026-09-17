$ErrorActionPreference='Stop'
$tokens=$null;$errors=$null
$tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'docker-execution.ps1'),[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Boundary source does not parse'}
foreach($name in @('New-LinuxBoundaryStartInfo','Get-LinuxGateAction')){
    $node=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq$name},$true)
    if($null-eq$node){throw 'Missing production function'}
    . ([scriptblock]::Create($node.Extent.Text))
}
$script:assertions=0
function Check([bool]$Condition,[string]$Message){$script:assertions++;if(!$Condition){throw $Message}}
function Reject([scriptblock]$Body){$rejected=$false;try{&$Body}catch{$rejected=$true};Check $rejected 'Expected rejected handoff'}
function Snapshot {
    [pscustomobject]@{PID=222;Parent=111;Session=111;State='S';StartTime=[uint64]9;PidNamespace='pid:[2]';UserNamespace='user:[2]';MountNamespace='mnt:[2]'}
}
function Action($current,$original,[string]$fd='7:9'){
    Get-LinuxGateAction $current $original 111 'pid:[1]' 'user:[1]' 'mnt:[1]' '7:9' $fd
}
$original=Snapshot
Check ((Action $original $null)-ceq'Stop') 'New init must be stopped by its ancestor'
Check ((Action (Snapshot) $original)-ceq'Wait') 'A merely present init must not be released'
$stopped=Snapshot;$stopped.State='T'
Check ((Action $stopped $original)-ceq'Release') 'Pinned stopped init with inherited inode should release'
Reject {Action $null $null}
Reject {Action $stopped $original '7:10'}
Reject {Action $stopped $original ''}
foreach($field in @('PID','StartTime','Parent','Session','PidNamespace','UserNamespace','MountNamespace')){
    $changed=Snapshot;$changed.State='T'
    $changed.$field=if($field.EndsWith('Namespace')){'other:[3]'}else{999}
    Reject {Action $changed $original}
}
foreach($field in @('PidNamespace','UserNamespace','MountNamespace')){
    $changed=Snapshot;$changed.$field=if($field-eq'PidNamespace'){'pid:[1]'}elseif($field-eq'UserNamespace'){'user:[1]'}else{'mnt:[1]'}
    Reject {Action $changed $null}
}
foreach($state in @('R','S','D','Z','X')){
    $changed=Snapshot;$changed.State=$state
    Check ((Action $changed $original)-ceq'Wait') 'Non-stopped state released the gate'
}
$environment=[Collections.Generic.Dictionary[string,string]]::new([StringComparer]::Ordinal)
$environment.Add('HOME','/inert home');$environment.Add('LANG','C.UTF-8')
$arguments=@('space value','single''quote','double"quote','; echo injected','$(`bad`)','',"line`nbreak",'unicode-☃')
$start=New-LinuxBoundaryStartInfo "/proc/$PID/fd/42" $arguments $environment '/inert work'
Check ($start.FileName-ceq'/usr/bin/setsid'-and!$start.UseShellExecute) 'Fixed session launcher required'
Check ($start.RedirectStandardInput-and$start.RedirectStandardOutput-and$start.RedirectStandardError) 'Owned gate and streams required'
Check ($start.Environment.Count-eq 2-and$start.Environment['HOME']-ceq'/inert home') 'Unexpected inherited environment'
Check ($start.WorkingDirectory-ceq'/inert work') 'Working-directory argument changed'
$offset=$start.ArgumentList.Count-$arguments.Count
for($i=0;$i-lt$arguments.Count;$i++){Check ($start.ArgumentList[$offset+$i]-ceq$arguments[$i]) 'Argument boundary or value changed'}
Check ($start.ArgumentList.Contains('--mount-proc')-and$start.ArgumentList.Contains('private')) 'Private proc namespace omitted'
Check (!$start.ArgumentList.Contains('--map-root-user')-and$start.ArgumentList.Contains('--map-current-user')) 'User mapping changed'
foreach($reference in @('/usr/bin/docker','/proc/self/fd/42','/proc/1/fd/42',"/proc/$PID/fd/42;bad")){
    Reject {New-LinuxBoundaryStartInfo $reference @() $environment '/inert'}
}
$native=$tree.Find({param($n)$n-is[Management.Automation.Language.StringConstantExpressionAst]-and$n.Value.Contains('public static class PscanNativeBoundary')},$true).Value
if($null-ne('PscanNativeBoundary'-as[type])){throw 'Inert compile needs a fresh process'}
$types=@(Add-Type -TypeDefinition $native -Language CSharp -PassThru)
Check ($null-ne('PscanNativeBoundary'-as[type])) 'Exact embedded native source did not compile'
if(!$IsLinux){Reject {[PscanNativeBoundary]::LinuxPin(222)}}
$flags=[Reflection.BindingFlags]'Static,NonPublic'
$signal=[PscanNativeBoundary].GetMethod('LinuxSignalResult',$flags)
$poll=[PscanNativeBoundary].GetMethod('LinuxPollResult',$flags)
Check ([bool]$signal.Invoke($null,@([long]0,[int]0))) 'Successful signal result rejected'
Check (![bool]$signal.Invoke($null,@([long]-1,[int]3))) 'ESRCH must mean the pinned process exited'
foreach($errorNumber in @(1,4,9,22,38)){
    $actualError=$null
    try{[void]$signal.Invoke($null,@([long]-1,[int]$errorNumber))}catch{
        $failure=$_.Exception
        while($failure-isnot[ComponentModel.Win32Exception]-and$null-ne$failure.InnerException){$failure=$failure.InnerException}
        if($failure-is[ComponentModel.Win32Exception]){$actualError=$failure.NativeErrorCode}
    }
    Check ($actualError-eq$errorNumber) 'Signal errno was not preserved as failure'
}
Check (![bool]$poll.Invoke($null,@([int]0,[short]0,[int]0))) 'Nonblocking no-event poll must remain live'
foreach($events in @(1,16,17)){Check ([bool]$poll.Invoke($null,@([int]1,[short]$events,[int]0))) 'Termination/hangup event rejected'}
foreach($pair in @(@(1,8),@(1,32),@(0,1),@(1,0),@(2,1))){
    Reject {$poll.Invoke($null,@([int]$pair[0],[short]$pair[1],[int]0))}
}
Reject {$poll.Invoke($null,@([int]-1,[short]0,[int]4))}
$fixtureTree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'test-docker-execution.ps1'),[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Native fixture source does not parse'}
$markerFunction=$fixtureTree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Get-LinuxFixtureMarkerPaths'},$true)
if($null-eq$markerFunction){throw 'Missing actual fixture path builder'}
. ([scriptblock]::Create($markerFunction.Extent.Text))
foreach($mode in @('child','grandchild','hold-child','hold-grandchild','detach')){
    $base='/inert directory/descendant.pid'
    $actual=@(Get-LinuxFixtureMarkerPaths $mode $base)
    $expected=if($mode.EndsWith('grandchild')){@('/inert directory/descendant.pid','/inert directory/descendant.pid.child','/inert directory/descendant.pid.child.ready')}else{@('/inert directory/descendant.pid','/inert directory/descendant.pid.ready')}
    Check ($actual.Count-eq$expected.Count) 'Marker paths contain an extra literal suffix'
    for($i=0;$i-lt$expected.Count;$i++){Check ($actual[$i]-ceq$expected[$i]) 'Marker path concatenation or order is wrong'}
}
Reject {Get-LinuxFixtureMarkerPaths 'unknown' '/inert/base'}
$windowsMarkerFunction=$fixtureTree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Get-WindowsFixtureMarkerPaths'},$true)
if($null-eq$windowsMarkerFunction){throw 'Missing actual Windows fixture path builder'}
. ([scriptblock]::Create($windowsMarkerFunction.Extent.Text))
foreach($mode in @('child','grandchild')){
    $base='C:/inert directory/descendant.pid'
    $actual=@(Get-WindowsFixtureMarkerPaths $mode $base)
    [string[]]$expected=if($mode-eq'grandchild'){@('C:/inert directory/descendant.pid','C:/inert directory/descendant.pid.child')}else{@('C:/inert directory/descendant.pid')}
    Check ($actual.Count-eq$expected.Count) 'Windows marker path count mismatch'
    for($i=0;$i-lt$expected.Count;$i++){Check ($actual[$i]-ceq$expected[$i]) 'Windows marker path concatenation or order is wrong'}
}
Reject {Get-WindowsFixtureMarkerPaths 'unknown' 'C:/inert/base'}
Write-Output "Linux handoff inert PASS assertions=$script:assertions native-process-calls=0 Linux-runtime-proof=UNAVAILABLE"
