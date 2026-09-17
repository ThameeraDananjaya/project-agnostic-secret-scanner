$ErrorActionPreference='Stop'
$tokens=$null;$errors=$null
$tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'docker-execution.ps1'),[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Production cleanup does not parse'}
$native=$tree.Find({param($n)$n-is[Management.Automation.Language.StringConstantExpressionAst]-and$n.Value.Contains('public static class PscanNativeBoundary')},$true).Value
[void](Add-Type -TypeDefinition $native)
Add-Type -TypeDefinition @'
using System;
using System.IO;
using System.Collections.Generic;
public static class PscanCleanupTestIO {
    public static bool InitExited,RootExited,InitAdmitted,DelayRoot,HoldInit,FailSignal,FailPoll;
    public static string RecordPath;
    public static List<string> Signals=new List<string>();
    public static bool LinuxExited(string handle) {
        if(FailPoll){FailPoll=false;throw new IOException("inert poll failure");}
        if(handle=="init")return InitExited;
        if(handle=="root")return RootExited;
        throw new Exception("unadmitted handle");
    }
    public static bool LinuxSignal(string handle,int signal) {
        if(signal!=9||(handle!="init"&&handle!="root"))throw new Exception("unadmitted signal");
        Signals.Add(handle);
        if(FailSignal){FailSignal=false;throw new IOException("inert signal failure");}
        if(handle=="init"){if(!HoldInit)InitExited=true;}
        else {
            if(InitAdmitted && File.Exists(RecordPath))throw new Exception("supervisor killed before reaping");
            if(!DelayRoot)RootExited=true;
        }
        return true;
    }
}
'@
$controller=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Invoke-LinuxPinnedCleanup'},$true)
$reaping=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Test-LinuxInitReaped'},$true)
if($null-eq$controller-or$null-eq$reaping){throw 'Actual cleanup functions missing'}
# Only OS I/O is substituted in this test scope. Execute actual controller,
# parser and monotonic budget; production accepts no callback or I/O override.
$source=$controller.Extent.Text.Replace('[PscanNativeBoundary]','[PscanCleanupTestIO]')
. ([scriptblock]::Create($source))
$source=$reaping.Extent.Text.Replace('[PscanNativeBoundary]','[PscanCleanupTestIO]').Replace('"/proc/$($Init.PID)/stat"','$script:initRecordPath')
if($source.Contains('"/proc/')){throw 'Reaping test did not bind owned file I/O'}
. ([scriptblock]::Create($source))
$cases=0
function Check([bool]$Value,[string]$Message){$script:cases++;if(!$Value){throw $Message}}
Check (!$controller.Extent.Text.Contains('Pin-Member')-and!$controller.Extent.Text.Contains('LinuxPin(')) 'Cleanup reacquires process identity'
$root=Join-Path ([IO.Path]::GetTempPath()) ('pscan-cleanup-'+[guid]::NewGuid().ToString('N'))
[void][IO.Directory]::CreateDirectory($root)
$script:initRecordPath=Join-Path $root 'init.stat'
[PscanCleanupTestIO]::RecordPath=$script:initRecordPath
function Stat([string]$State='Z',[string]$Start='7') {
    $fields=@($State,'111','111','111')+@('0')*15+@($Start)
    [IO.File]::WriteAllText($script:initRecordPath,'222 (inert init) '+($fields-join' '))
}
function Get-LinuxSessionMembers([int]$Session,[hashtable]$Ledger) {
    if($Session-ne111){throw 'Session scope changed'}
    $script:tick++
    if($script:scenario-ceq'member-error'-and$script:tick-eq1){throw 'inert member inspection failure'}
    if($script:tick-eq1-and$script:scenario-ceq'live-then-reaped'){Stat 'Z'}
    if($script:tick-ge2-and$script:scenario-cnotin@('never-reaped','absent-before-exit')){[IO.File]::Delete($script:initRecordPath)}
    if($script:scenario-ceq'delayed-supervisor'-and$script:tick-ge4){[PscanCleanupTestIO]::RootExited=$true}
    if($script:scenario-ceq'pipe-delay'-and$script:tick-ge6){$script:pipesClosed=$true}
    if(![PscanCleanupTestIO]::RootExited){[pscustomobject]@{PID=111;StartTime=8}}
    if([IO.File]::Exists($script:initRecordPath)){[pscustomobject]@{PID=222;StartTime=7}}
}
function Test-Path([string]$LiteralPath,[object]$ErrorAction) {
    if($LiteralPath-cne'/proc/222/ns/pid'){throw 'Namespace inspection scope changed'}
    if($script:scenario-ceq'namespace-error'){throw 'inert namespace inspection failure'}
    return ($script:scenario-ceq'namespace-remains')
}
try {
    foreach($scenario in @('live-then-reaped','zombie-then-reaped','already-reaped','delayed-supervisor','never-reaped','absent-before-exit',
        'missing-init','missing-root','reused-record','malformed-record','member-error','poll-error','signal-error','pipe-delay','pipe-never','namespace-error','namespace-remains','initial-error','expired')) {
        $script:scenario=$scenario;$script:tick=0;$script:pipesClosed=$scenario-cnotin@('pipe-delay','pipe-never')
        [PscanCleanupTestIO]::Signals.Clear();[PscanCleanupTestIO]::InitExited=$scenario-cne'live-then-reaped'
        [PscanCleanupTestIO]::RootExited=$false;[PscanCleanupTestIO]::InitAdmitted=$scenario-cne'missing-init'
        [PscanCleanupTestIO]::DelayRoot=$scenario-ceq'delayed-supervisor';[PscanCleanupTestIO]::HoldInit=$scenario-ceq'absent-before-exit'
        [PscanCleanupTestIO]::FailSignal=$scenario-ceq'signal-error';[PscanCleanupTestIO]::FailPoll=$scenario-ceq'poll-error'
        Stat $(if($scenario-ceq'live-then-reaped'){'S'}else{'Z'})
        if($scenario-in@('already-reaped','absent-before-exit','missing-init')){[IO.File]::Delete($script:initRecordPath)}
        if($scenario-ceq'absent-before-exit'){[PscanCleanupTestIO]::InitExited=$false}
        if($scenario-ceq'reused-record'){Stat 'S' '99'}
        if($scenario-ceq'malformed-record'){[IO.File]::WriteAllText($script:initRecordPath,'unreadable identity')}
        $init=[pscustomobject]@{PID=222;StartTime=[uint64]7};$initHandle='init';$rootHandle='root'
        if($scenario-ceq'missing-init'){$init=$null;$initHandle=$null}
        if($scenario-ceq'missing-root'){$rootHandle=$null}
        $process=[pscustomobject]@{};$process|Add-Member -MemberType ScriptProperty -Name HasExited -Value {[PscanCleanupTestIO]::RootExited}
        $task=[pscustomobject]@{};$task|Add-Member -MemberType ScriptProperty -Name IsCompleted -Value {$script:pipesClosed}
        $budget=[PscanOperationBudget]::new(15000,150,131072);$budget.BeginCleanup()
        if($scenario-ceq'expired'){Start-Sleep -Milliseconds 170}
        $clock=[Diagnostics.Stopwatch]::StartNew()
        $result=Invoke-LinuxPinnedCleanup $init $initHandle $rootHandle $process $task $task 111 @{} $budget ($scenario-ceq'initial-error')
        $expected=$scenario-in@('live-then-reaped','zombie-then-reaped','already-reaped','delayed-supervisor','pipe-delay')
        Check ($result.Proved-eq$expected) "Cleanup outcome differs: $scenario $($result|ConvertTo-Json -Compress)"
        Check ($clock.ElapsedMilliseconds-lt1000) "Cleanup obtained unbounded time: $scenario"
        Check (@([PscanCleanupTestIO]::Signals|Where-Object{$_-ceq'root'}).Count-le1) 'Repeated supervisor signal'
        if($scenario-in@('never-reaped','absent-before-exit','expired','missing-root')){Check (![PscanCleanupTestIO]::Signals.Contains('root')) 'Root signalled without init reaping or budget/handle'}
        if($scenario-in@('reused-record','malformed-record','poll-error','member-error','signal-error','initial-error','namespace-error')){Check ($null-ne$result.FailureStage) 'Actual error latch lost'}
    }
} finally {
    $resolved=[IO.Path]::GetFullPath($root);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)-or![IO.Path]::GetFileName($resolved).StartsWith('pscan-cleanup-')){throw 'Unsafe test cleanup'}
    [IO.Directory]::Delete($resolved,$true)
}
# Execute the actual Linux result-construction tail with real Task<byte[]>.
# Only process/capture I/O is inert. The complete held-call function below keeps
# production containment, deadline, type validation and strict UTF-8 decoding.
$boundary=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Invoke-LinuxSessionBoundary'},$true)
$captureStart=$boundary.Find({param($n)$n-is[Management.Automation.Language.AssignmentStatementAst]-and$n.Left.Extent.Text-ceq'$stdout'-and$n.Right.Extent.Text-ceq'[byte[]]::new(0)'},$true)
$captureTry=$boundary.Body.EndBlock.Statements|Where-Object{$_-is[Management.Automation.Language.TryStatementAst]-and$null-ne$_.Finally-and$_.Body.Extent.StartOffset-lt$captureStart.Extent.StartOffset-and$_.Body.Extent.EndOffset-gt$captureStart.Extent.EndOffset}
if($null-eq$captureStart-or@($captureTry).Count-ne1){throw 'Actual capture construction missing or ambiguous'}
$tail=$tree.Extent.Text.Substring($captureStart.Extent.StartOffset,$captureTry.Body.Extent.EndOffset-1-$captureStart.Extent.StartOffset)
$assemble=[scriptblock]::Create(@'
param($stdoutTask,$stderrTask,$terminal=$null)
$started=$true;$process=[pscustomobject]@{HasExited=$true;ExitCode=7}
$ledger=@{};$namespaceIdentity='pid:[222]';$cleanupDetails=$null
'@ + "`n" + $tail)
$held=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Invoke-HeldDockerCall'},$true)
$nativeDispatch=$held.Body.EndBlock.Statements|Where-Object{$_-is[Management.Automation.Language.IfStatementAst]-and$_.Clauses[0].Item1.Extent.Text-ceq'$IsWindows'}
if(@($nativeDispatch).Count-ne1){throw 'Actual native dispatch missing or ambiguous'}
. ([scriptblock]::Create($held.Extent.Text.Replace($nativeDispatch.Extent.Text,'$result=$script:captureResult')))
function Decode($Result){$script:captureResult=$Result;Invoke-HeldDockerCall @('inert') @{Budget=[PscanOperationBudget]::new(15000,2000,131072);Calls=0;Observed=0}}
function Completed($Bytes){[Threading.Tasks.Task]::FromResult[byte[]]($Bytes)}
$captureChecksBefore=$cases
foreach($pair in @(
    @{Out=[byte[]]::new(0);Err=[byte[]]::new(0);TextOut='';TextErr=''},
    @{Out=[byte[]]::new(0);Err=[byte[]]@(65);TextOut='';TextErr='A'},
    @{Out=[byte[]]@(65);Err=[byte[]]::new(0);TextOut='A';TextErr=''},
    @{Out=[byte[]]@(65);Err=[byte[]]@(226,130,172);TextOut='A';TextErr=[char]0x20ac},
    @{Out=[byte[]]@(226,130,172);Err=[byte[]]@(65,66);TextOut=[char]0x20ac;TextErr='AB'}
)){
    $result=& $assemble (Completed $pair.Out) (Completed $pair.Err)
    Check ($result.StdOut-is[byte[]]-and$result.StdErr-is[byte[]]) 'Actual construction lost buffer type'
    Check ([object]::ReferenceEquals($result.StdOut,$pair.Out)-and[object]::ReferenceEquals($result.StdErr,$pair.Err)) 'Actual construction changed captured buffers'
    $decoded=Decode $result
    Check ($decoded.StdOut-ceq$pair.TextOut-and$decoded.StdErr-ceq$pair.TextErr-and$decoded.ExitCode-eq7) 'Actual decoder changed stream or exit evidence'
}
foreach($bytes in @([byte[]]@(195,40),[byte[]]@(226,130))){
    foreach($side in @('stdout','stderr')){
        $out=[byte[]]::new(0);$err=[byte[]]::new(0)
        if($side-ceq'stdout'){$out=$bytes}else{$err=$bytes}
        $result=& $assemble (Completed $out) (Completed $err)
        $rejected=$false;try{[void](Decode $result)}catch{$rejected=$_.Exception.ToString().Contains('DecoderFallbackException')}
        Check $rejected 'Actual decoder accepted malformed/incomplete UTF-8'
    }
}
foreach($invalid in @(@{Value=$null},@{Value=[byte]65},@{Value=[object[]]@(65)},@{Value=[int[]]@(65)},@{Value='A'},@{Value=[object[]]::new(0)},@{Value=@{}})){
    foreach($side in @('stdout','stderr')){
        $good=Completed ([byte[]]::new(0));$bad=[pscustomobject]@{IsCompletedSuccessfully=$true;Result=$invalid.Value}
        $outTask=$good;$errTask=$good
        if($side-ceq'stdout'){$outTask=$bad}else{$errTask=$bad}
        $result=& $assemble $outTask $errTask
        Check ($result.Terminal-ceq'capture evidence unavailable'-and!$result.ContainmentEmpty) 'Invalid capture construction was accepted as empty'
        # Independently exercise the production decoder type gate even if a
        # malformed result claims successful containment.
        $result.Terminal=$null;$result.ContainmentEmpty=$true
        $rejected=$false;try{[void](Decode $result)}catch{$rejected=$_.Exception.Message-ceq'Docker protocol capture buffers are not byte arrays'}
        Check $rejected 'Actual decoder coerced invalid capture type'
    }
}
$pending=[Threading.Tasks.TaskCompletionSource[byte[]]]::new()
$faulted=[Threading.Tasks.TaskCompletionSource[byte[]]]::new();$faulted.SetException([IO.IOException]::new('inert capture failure'))
$cancelled=[Threading.Tasks.TaskCompletionSource[byte[]]]::new();$cancelled.SetCanceled()
foreach($taskCase in @(@{Task=$null},@{Task=$pending.Task},@{Task=$faulted.Task},@{Task=$cancelled.Task})){
    foreach($side in @('stdout','stderr')){
        $outTask=Completed ([byte[]]::new(0));$errTask=$outTask
        if($side-ceq'stdout'){$outTask=$taskCase.Task}else{$errTask=$taskCase.Task}
        $result=& $assemble $outTask $errTask
        Check ($result.StdOut-is[byte[]]-and$result.StdErr-is[byte[]]-and$result.Terminal-ceq'capture evidence unavailable'-and!$result.ContainmentEmpty) 'Incomplete capture did not remain terminal'
        $rejected=$false;try{[void](Decode $result)}catch{$rejected=$_.Exception.Message.StartsWith('Docker protocol native containment failed:')}
        Check $rejected 'Actual held call decoded uncompleted capture'
        $prior=& $assemble $outTask $errTask 'running: stream overflow or read failure; cleanup uncertainty'
        Check ($prior.Terminal-ceq'running: stream overflow or read failure; cleanup uncertainty') 'Capture construction replaced existing terminal evidence'
    }
}
Write-Output "Linux pinned cleanup inert PASS cases=$cases scenarios=19 capture-checks=$($cases-$captureChecksBefore) actual-controller=PASS actual-stat-parser=PASS actual-capture-construction=PASS actual-strict-decoder=PASS shared-clock=PASS Linux-runtime-proof=UNAVAILABLE"
