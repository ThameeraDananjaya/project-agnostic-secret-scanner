param(
    [ValidateSet('Orchestrate','CleanNative','HostileCompatible','HostileStale')]
    [string]$InternalMode = 'Orchestrate',
    [string]$ControlPath
)

$ErrorActionPreference = 'Stop'
$productionPath = Join-Path $PSScriptRoot 'docker-execution.ps1'
$hostExecutable = (Get-Process -Id $PID).Path

function Get-NativeBoundarySource {
    $tokens = $null
    $errors = $null
    $ast = [Management.Automation.Language.Parser]::ParseFile($productionPath, [ref]$tokens, [ref]$errors)
    if ($errors.Count -ne 0) { throw 'Closed Docker entrypoint does not parse' }
    $sources = @($ast.FindAll({
        param($node)
        $node -is [Management.Automation.Language.StringConstantExpressionAst] -and
            $node.Value.Contains('public static class PscanNativeBoundary')
    }, $true))
    if ($sources.Count -ne 1) { throw 'Expected one immutable native boundary source' }
    [pscustomobject]@{ Ast=$ast; Source=$sources[0].Value }
}

function Assert-CompiledBoundaryIdentity([type[]]$Types) {
    $boundary = @($Types | Where-Object { $_.IsPublic -and $_.FullName -ceq 'PscanNativeBoundary' })
    $result = @($Types | Where-Object { $_.IsPublic -and $_.FullName -ceq 'PscanBoundaryResult' })
    if ($boundary.Count -ne 1 -or $result.Count -ne 1 -or
        ![object]::ReferenceEquals($boundary[0].Assembly, $result[0].Assembly) -or
        ![object]::ReferenceEquals($boundary[0], ('PscanNativeBoundary' -as [type])) -or
        ![object]::ReferenceEquals($result[0], ('PscanBoundaryResult' -as [type]))) {
        throw 'Clean fixture compiled an unexpected native boundary identity'
    }
}

function Get-BoundaryDirectories {
    @(Get-ChildItem -LiteralPath ([IO.Path]::GetTempPath()) -Directory -Filter 'pscan-docker-boundary-*' -ErrorAction Stop |
        ForEach-Object { $_.FullName } | Sort-Object)
}

function Assert-NoNewBoundaryDirectory([string[]]$Before, [string]$Name) {
    $after = @(Get-BoundaryDirectories)
    $new = @($after | Where-Object { $Before -cnotcontains $_ })
    if ($new.Count -ne 0) { throw "$Name reached Docker-boundary directory creation: $($new[0])" }
}

function Invoke-CleanNativeFixture {
    if ($null -ne ('PscanNativeBoundary' -as [type])) { throw 'Clean fixture process began with the native boundary already loaded' }
    $parsed = Get-NativeBoundarySource
    $compiled = @(Add-Type -TypeDefinition $parsed.Source -Language CSharp -PassThru)
    Assert-CompiledBoundaryIdentity -Types $compiled

    $temporaryRoot = Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-005-clean-' + [guid]::NewGuid().ToString('N'))
    [void](New-Item -ItemType Directory -Path $temporaryRoot)
    try {
        if ($IsLinux) {
            foreach ($functionName in @('Get-LinuxSessionMembers','Test-LinuxMemberAlive','Invoke-LinuxSessionBoundary')) {
                $definition = $parsed.Ast.Find({
                    param($node)
                    $node -is [Management.Automation.Language.FunctionDefinitionAst] -and $node.Name -ceq $functionName
                }, $true)
                if ($null -eq $definition) { throw "Missing Linux boundary function $functionName" }
                . ([scriptblock]::Create($definition.Extent.Text))
            }
            $script:dockerOutputLimit=131072;$script:dockerBudgetMilliseconds=15000;$script:cleanupGraceMilliseconds=2000;$script:Operation='NativeFixture'
            $native=(Get-Process -Id $PID).Path;$identity=[IO.File]::Open($native,[IO.FileMode]::Open,[IO.FileAccess]::Read,[IO.FileShare]::Read)
            $fixture=Join-Path $temporaryRoot 'native-fixture.ps1'
            Set-Content -LiteralPath $fixture -Encoding utf8NoBOM -Value @'
param([string]$Mode,[int]$Count=0,[string]$PidFile='')
$out=[Console]::OpenStandardOutput();$err=[Console]::OpenStandardError()
function Bytes($stream,[byte]$value,[int]$count){$chunk=[byte[]]::new(1024);[Array]::Fill($chunk,$value);for($left=$count;$left-gt 0;$left-=$n){$n=[Math]::Min($left,$chunk.Length);$stream.Write($chunk,0,$n);$stream.Flush()}}
switch($Mode){
'stdout'{Bytes $out 97 $Count}'stderr'{Bytes $err 98 $Count}'both'{for($left=$Count;$left-gt 0;$left-=$n){$n=[Math]::Min(1024,$left);Bytes $out 97 $n;Bytes $err 98 $n}}
'split-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush();Start-Sleep -Milliseconds 25;$out.WriteByte(0xAC);$out.Flush()}
'invalid-utf8'{$out.WriteByte(0xC3);$out.WriteByte(0x28);$out.Flush()}'incomplete-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush()}'immediate'{}'nonzero'{exit 7}'hang'{Start-Sleep -Seconds 30}
'child'{$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang') -PassThru -NoNewWindow;Set-Content $PidFile $p.Id}
'detach'{$p=Start-Process -FilePath /usr/bin/setsid -ArgumentList @((Get-Process -Id $PID).Path,'-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang') -PassThru -NoNewWindow;Set-Content $PidFile $p.Id;Start-Sleep -Seconds 30}
'grandchild'{$childFile=$PidFile+'.child';$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','child','-PidFile',$childFile) -PassThru -NoNewWindow;Set-Content $PidFile $p.Id;Start-Sleep -Milliseconds 500}
default{throw 'unknown fixture mode'}}
'@
            $environment=[Collections.Generic.Dictionary[string,string]]::new([StringComparer]::Ordinal);$environment.Add('HOME',$temporaryRoot);$environment.Add('TMPDIR',$temporaryRoot);$environment.Add('LANG','C.UTF-8')
            $reference="/proc/$PID/fd/$($identity.SafeFileHandle.DangerousGetHandle().ToInt64())"
            function RunLinux([string]$mode,[int]$count=0,[string]$pidFile=''){$a=@('-NoProfile','-NonInteractive','-File',$fixture,'-Mode',$mode,'-Count',[string]$count);if($pidFile){$a+=@('-PidFile',$pidFile)};Invoke-LinuxSessionBoundary -ExecutableReference $reference -Arguments $a -Environment $environment -WorkingDirectory $temporaryRoot}
            function LinuxSuccess($r,[string]$name){if($r.Terminal-or!$r.ContainmentEmpty){throw "$name returned terminal lifecycle evidence"}}
            function LinuxTerminal($r,[string]$name){if(!$r.Terminal){throw "$name did not return terminal evidence"};if($r.Terminal.Contains('cleanup uncertainty')){throw "$name left uncertain containment"}}
            foreach($n in @(131071,131072)){$r=RunLinux stdout $n;LinuxSuccess $r "stdout-$n";if($r.StdOut.Length-ne$n){throw 'stdout limit mismatch'};$r=RunLinux stderr $n;LinuxSuccess $r "stderr-$n";if($r.StdErr.Length-ne$n){throw 'stderr limit mismatch'}}
            $r=RunLinux immediate;LinuxSuccess $r immediate;$r=RunLinux both 131072;LinuxSuccess $r both;$utf8=[Text.UTF8Encoding]::new($false,$true);$r=RunLinux split-utf8;LinuxSuccess $r split;if($utf8.GetString($r.StdOut)-cne '€'){throw 'split UTF-8 mismatch'};$r=RunLinux nonzero;LinuxSuccess $r nonzero;if($r.ExitCode-ne 7){throw 'nonzero mismatch'}
            foreach($case in @(@{M='stdout';C=131073},@{M='stderr';C=131073},@{M='both';C=131073},@{M='hang';C=0})){LinuxTerminal (RunLinux $case.M $case.C) $case.M}
            foreach($mode in @('invalid-utf8','incomplete-utf8')){$r=RunLinux $mode;LinuxSuccess $r $mode;$failed=$false;try{[void]$utf8.GetString($r.StdOut)}catch{$failed=$true};if(!$failed){throw "$mode unexpectedly decoded"}}
            foreach($mode in @('child','detach','grandchild')){$pidFile=Join-Path $temporaryRoot "$mode.pid";$r=RunLinux $mode 0 $pidFile;LinuxTerminal $r $mode;foreach($path in @($pidFile,$pidFile+'.child')){if(Test-Path $path){$id=[int](Get-Content -Raw $path);if(Test-Path "/proc/$id"){throw "$mode left live member $id"}}}}
            $identity.Dispose()
            $containment = 'LINUX-PID-NS'; $replacementRace = 'NOT-APPLICABLE'
        } elseif ($IsWindows) {
            $native=(Get-Process -Id $PID).Path
            $signature=Get-AuthenticodeSignature -LiteralPath $native
            if($signature.Status-ne[Management.Automation.SignatureStatus]::Valid){throw 'Native fixture host is not signed system tooling'}
            $fixture=Join-Path $temporaryRoot 'native-fixture.ps1'
            Set-Content -LiteralPath $fixture -Encoding utf8NoBOM -Value @'
param([string]$Mode,[int]$Count=0,[string]$PidFile='')
$out=[Console]::OpenStandardOutput();$err=[Console]::OpenStandardError()
function Bytes($stream,[byte]$value,[int]$count){$chunk=[byte[]]::new(1024);[Array]::Fill($chunk,$value);for($left=$count;$left-gt 0;$left-=$n){$n=[Math]::Min($left,$chunk.Length);$stream.Write($chunk,0,$n);$stream.Flush()}}
switch($Mode){
'stdout'{Bytes $out 97 $Count}
'stderr'{Bytes $err 98 $Count}
'both'{for($left=$Count;$left-gt 0;$left-=$n){$n=[Math]::Min(1024,$left);Bytes $out 97 $n;Bytes $err 98 $n}}
'split-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush();Start-Sleep -Milliseconds 25;$out.WriteByte(0xAC);$out.Flush()}
'invalid-utf8'{$out.WriteByte(0xC3);$out.WriteByte(0x28);$out.Flush()}
'incomplete-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush()}
'immediate'{}
'nonzero'{exit 7}
'hang'{Start-Sleep -Seconds 30}
'child'{$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang') -PassThru -NoNewWindow;Set-Content -LiteralPath $PidFile -Value $p.Id}
'grandchild'{$childFile=$PidFile+'.child';$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','child','-PidFile',$childFile) -PassThru -NoNewWindow;Set-Content -LiteralPath $PidFile -Value $p.Id;Start-Sleep -Milliseconds 500}
default{throw 'unknown fixture mode'}
}
'@
            $environment=[Collections.Generic.Dictionary[string,string]]::new([StringComparer]::Ordinal)
            $environment.Add('SystemRoot',$env:SystemRoot);$environment.Add('WINDIR',$env:WINDIR);$environment.Add('TEMP',$temporaryRoot);$environment.Add('TMP',$temporaryRoot)
            function Run([string]$mode,[int]$count=0,[string]$pidFile=''){
                $a=@('-NoProfile','-NonInteractive','-File',$fixture,'-Mode',$mode,'-Count',[string]$count)
                if($pidFile){$a+=@('-PidFile',$pidFile)}
                [PscanNativeBoundary]::RunWindows($native,$a,$environment,$temporaryRoot,131072,15000,2000)
            }
            function Success($r,[string]$name){if($r.Terminal-or!$r.ContainmentEmpty){throw "$name returned terminal lifecycle evidence"}}
            function Terminal($r,[string]$name){if(!$r.Terminal-or!$r.ContainmentEmpty){throw "$name did not terminate an empty job"}}
            $r=Run immediate;Success $r immediate
            foreach($n in @(131071,131072)){$r=Run stdout $n;Success $r "stdout-$n";if($r.StdOut.Length-ne$n){throw 'stdout boundary mismatch'};$r=Run stderr $n;Success $r "stderr-$n";if($r.StdErr.Length-ne$n){throw 'stderr boundary mismatch'}}
            $r=Run both 131072;Success $r 'both';if($r.StdOut.Length-ne 131072-or$r.StdErr.Length-ne 131072){throw 'simultaneous boundary mismatch'}
            $utf8=[Text.UTF8Encoding]::new($false,$true);$r=Run split-utf8;Success $r 'split-utf8';if($utf8.GetString($r.StdOut)-cne '€'){throw 'split UTF-8 mismatch'}
            $r=Run nonzero;Success $r 'nonzero';if($r.ExitCode-ne 7){throw 'nonzero exit mismatch'}
            foreach($case in @(@{M='stdout';C=131073},@{M='stderr';C=131073},@{M='both';C=131073},@{M='hang';C=0})){Terminal (Run $case.M $case.C) $case.M}
            foreach($mode in @('invalid-utf8','incomplete-utf8')){$r=Run $mode;Success $r $mode;$failed=$false;try{[void]$utf8.GetString($r.StdOut)}catch{$failed=$true};if(!$failed){throw "$mode unexpectedly decoded"}}
            $startFailed=$false;try{[void][PscanNativeBoundary]::RunWindows((Join-Path $temporaryRoot 'missing.exe'),@(),$environment,$temporaryRoot,131072,15000,2000)}catch{$startFailed=$true};if(!$startFailed){throw 'Start failure unexpectedly returned trusted evidence'}
            foreach($mode in @('child','grandchild')){$pidFile=Join-Path $temporaryRoot "$mode.pid";$r=Run $mode 0 $pidFile;Terminal $r $mode;foreach($path in @($pidFile,$pidFile+'.child')){if(Test-Path $path){$id=[int](Get-Content -Raw $path);if(Get-Process -Id $id -ErrorAction SilentlyContinue){throw "$mode left live member $id"}}}}
            $copy=Join-Path $temporaryRoot 'identity.exe';Copy-Item -LiteralPath $native -Destination $copy
            $replacement=Join-Path $temporaryRoot 'replacement.exe';Copy-Item -LiteralPath $native -Destination $replacement
            $handle=[IO.File]::Open($copy,[IO.FileMode]::Open,[IO.FileAccess]::Read,[IO.FileShare]::Read)
            try{$replaced=$true;try{[IO.File]::Move($replacement,$copy,$true)}catch{$replaced=$false};if($replaced){throw 'Validated executable path was replaceable while the identity handle was held'}}finally{$handle.Dispose()}
            $containment = 'WINDOWS-JOB'; $replacementRace = 'REJECT'
        } else {
            throw 'Native containment fixtures require Windows or Linux'
        }

        $before = @(Get-BoundaryDirectories)
        $secondOutput = @()
        $secondError = $null
        try { $secondOutput = @(& $productionPath -Operation EngineInspection) }
        catch { $secondError = $_.Exception.Message }
        if ($secondError -cne 'Pre-existing PscanNativeBoundary type is terminal ambient state') { throw "Second invocation did not reject ambient state: $secondError" }
        if ($secondOutput.Count -ne 0) { throw 'Second invocation emitted a result despite terminal ambient state' }
        Assert-NoNewBoundaryDirectory -Before $before -Name 'Second invocation'
        Write-Output "Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=$containment streams=PASS utf8=PASS process-tree=EMPTY replacement-race=$replacementRace second-invocation=REJECT"
    } finally {
        $resolved=[IO.Path]::GetFullPath($temporaryRoot);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
        if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)){throw "Unsafe cleanup path: $resolved"}
        if(Test-Path -LiteralPath $resolved){Remove-Item -LiteralPath $resolved -Recurse -Force}
    }
}

function Invoke-HostileFixture([bool]$Stale) {
    if ([string]::IsNullOrWhiteSpace($ControlPath)) { throw 'Hostile fixture requires a control path' }
    $startState = if ($null -eq ('PscanNativeBoundary' -as [type])) { 'ABSENT' } else { 'PRESENT' }
    if ($startState -cne 'ABSENT') { throw 'Hostile fixture process did not start clean' }
    $compatibleSource = @'
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Threading.Tasks;
public sealed class PscanBoundaryResult { public int ExitCode; public byte[] StdOut; public byte[] StdErr; public string Terminal; public long[] ObservedMembers; public bool ContainmentEmpty; }
public static class PscanNativeBoundary {
    public static int FakeBoundaryCalls = 0;
    public static int DockerCalls = 0;
    public static Task<byte[]> Capture(Stream stream, int limit) { FakeBoundaryCalls++; DockerCalls++; return Task.FromResult(Encoding.UTF8.GetBytes("FABRICATED-CAPTURE")); }
    public static PscanBoundaryResult RunWindows(string executable, string[] arguments, IDictionary<string,string> environment, string cwd, int streamLimit, int budgetMs, int cleanupMs) {
        FakeBoundaryCalls++; DockerCalls++;
        return new PscanBoundaryResult { ExitCode=0, StdOut=Encoding.UTF8.GetBytes("{\"Version\":\"FABRICATED\",\"ApiVersion\":\"FABRICATED\",\"Os\":\"linux\",\"Arch\":\"amd64\"}"), StdErr=new byte[0], Terminal=null, ObservedMembers=new long[0], ContainmentEmpty=true };
    }
}
'@
    $staleSource = @'
using System;
using System.Collections.Generic;
using System.IO;
using System.Text;
using System.Threading.Tasks;
public sealed class PscanBoundaryResult { public int ExitCode; public byte[] StdOut; public byte[] StdErr; public string Terminal; public long[] ObservedMembers; public bool ContainmentEmpty; }
public static class PscanNativeBoundary {
    public const int HistoricalGeneration = 3;
    public static int FakeBoundaryCalls = 0;
    public static int DockerCalls = 0;
    public static Task<byte[]> Capture(Stream stream, int limit) { FakeBoundaryCalls++; DockerCalls++; return Task.FromResult(Encoding.UTF8.GetBytes("STALE-FABRICATED-CAPTURE")); }
    public static PscanBoundaryResult RunWindows(string executable, string[] arguments, IDictionary<string,string> environment, string cwd, int streamLimit, int budgetMs, int cleanupMs) {
        FakeBoundaryCalls++; DockerCalls++;
        return new PscanBoundaryResult { ExitCode=0, StdOut=Encoding.UTF8.GetBytes("{\"Version\":\"STALE-FABRICATED\",\"ApiVersion\":\"STALE\",\"Os\":\"linux\",\"Arch\":\"amd64\"}"), StdErr=new byte[0], Terminal=null, ObservedMembers=new long[0], ContainmentEmpty=true };
    }
}
'@
    $fixtureSource = if ($Stale) { $staleSource } else { $compatibleSource }
    Add-Type -TypeDefinition $fixtureSource -Language CSharp
    $before = @(Get-BoundaryDirectories)
    $output = @()
    $errorMessage = $null
    try { $output = @(& $productionPath -Operation EngineInspection) }
    catch { $errorMessage = $_.Exception.Message }
    Assert-NoNewBoundaryDirectory -Before $before -Name "Hostile $(if($Stale){'stale'}else{'compatible'}) fixture"
    $fabricated = (($output -join "`n") -match 'FABRICATED|ApiVersion|ContainmentEmpty|DockerSHA256')
    $record = [ordered]@{
        Case=$(if($Stale){'stale'}else{'compatible'}); StartState=$startState; Preloaded=('PscanNativeBoundary' -as [type]).FullName
        Rejected=($errorMessage -ceq 'Pre-existing PscanNativeBoundary type is terminal ambient state')
        Error=$errorMessage; FakeBoundaryCalls=[PscanNativeBoundary]::FakeBoundaryCalls; DockerCalls=[PscanNativeBoundary]::DockerCalls
        OutputCount=$output.Count; TrustedResult=($output.Count -ne 0); FabricatedFieldsObserved=[bool]$fabricated
    }
    [IO.File]::WriteAllText([IO.Path]::GetFullPath($ControlPath), ($record | ConvertTo-Json -Compress), [Text.UTF8Encoding]::new($false))
    if (!$record.Rejected -or $record.FakeBoundaryCalls -ne 0 -or $record.DockerCalls -ne 0 -or $record.TrustedResult -or $record.FabricatedFieldsObserved) { exit 0 }
    exit 197
}

if ($InternalMode -ceq 'CleanNative') { Invoke-CleanNativeFixture; return }
if ($InternalMode -ceq 'HostileCompatible') { Invoke-HostileFixture -Stale $false }
if ($InternalMode -ceq 'HostileStale') { Invoke-HostileFixture -Stale $true }

if ($null -ne ('PscanNativeBoundary' -as [type])) { throw 'Orchestrator process began with an ambient native boundary type' }
$releaseFiles = @(
    'admit-image.ps1','acquire.ps1','build.ps1','cache-canary.ps1',
    'image-admission.ps1','test-cache-boundary.ps1','test-crlf-shell-payloads.ps1'
)
$allSource = foreach ($name in $releaseFiles) { Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot $name) }
$production = Get-Content -Raw -LiteralPath $productionPath
$parsedProduction = Get-NativeBoundarySource
$linuxSessionFunction = $parsedProduction.Ast.Find({
    param($node)
    $node -is [Management.Automation.Language.FunctionDefinitionAst] -and
        $node.Name -ceq 'Get-LinuxSessionMembers'
}, $true)
if ($null -eq $linuxSessionFunction) { throw 'Production omits Get-LinuxSessionMembers' }
$linuxSessionAssignments = @($linuxSessionFunction.FindAll({
    param($node)
    $node -is [Management.Automation.Language.AssignmentStatementAst]
}, $true))
$getAssignmentTargetVariableName = {
    param([Management.Automation.Language.AssignmentStatementAst]$Assignment)
    $target = $Assignment.Left
    while ($target -is [Management.Automation.Language.AttributedExpressionAst]) {
        $target = $target.Child
    }
    if ($target -isnot [Management.Automation.Language.VariableExpressionAst]) { return $null }
    ($target.VariablePath.UserPath -split ':')[-1]
}
$isReservedPidAssignment = {
    param([Management.Automation.Language.AssignmentStatementAst]$Assignment)
    (& $getAssignmentTargetVariableName $Assignment) -ieq 'PID'
}
$assignmentPredicateCases = @(
    @{ Name='untyped-mixed-case'; Source='$pId = 1'; Expected=1 },
    @{ Name='typed-mixed-case'; Source='[int]$PiD = 1'; Expected=1 },
    @{ Name='rhs-and-ledger-property'; Source='$value = $PID; $ledger = [pscustomobject]@{ PID=$value }; [int]$other = 1'; Expected=0 }
)
foreach ($case in $assignmentPredicateCases) {
    $probeTokens = $null
    $probeErrors = $null
    $probeAst = [Management.Automation.Language.Parser]::ParseInput($case.Source, [ref]$probeTokens, [ref]$probeErrors)
    if ($probeErrors.Count -ne 0) { throw "PID assignment predicate self-test did not parse: $($case.Name)" }
    $probeAssignments = @($probeAst.FindAll({
        param($node)
        $node -is [Management.Automation.Language.AssignmentStatementAst]
    }, $true))
    $probeReserved = @($probeAssignments | Where-Object { & $isReservedPidAssignment $_ })
    if ($probeReserved.Count -ne $case.Expected) {
        throw "PID assignment predicate self-test failed: $($case.Name)"
    }
}
$reservedPidAssignments = @($linuxSessionAssignments | Where-Object {
    & $isReservedPidAssignment $_
})
if ($reservedPidAssignments.Count -ne 0) { throw 'Get-LinuxSessionMembers assigns to reserved automatic variable PID' }
$identifierAssignments = @($linuxSessionAssignments | Where-Object {
    $_.Left -is [Management.Automation.Language.VariableExpressionAst] -and
        $_.Left.VariablePath.UserPath -ceq 'linuxProcessIdentifier'
})
if ($identifierAssignments.Count -ne 1 -or
    $identifierAssignments[0].Right.Extent.Text -cne "[int]`$stat.Substring(0, `$stat.IndexOf(' '))") {
    throw 'Linux process identifier is not bound once from the parsed stat PID as int'
}
$identityAssignments = @($linuxSessionAssignments | Where-Object {
    $_.Left -is [Management.Automation.Language.VariableExpressionAst] -and
        $_.Left.VariablePath.UserPath -ceq 'identity'
})
if ($identityAssignments.Count -ne 1 -or
    $identityAssignments[0].Right.Extent.Text -cne '"${linuxProcessIdentifier}:$startTime"') {
    throw 'Linux process identity does not preserve <process-id>:<start-time> from the renamed identifier'
}
$ledgerAssignments = @($linuxSessionAssignments | Where-Object {
    $_.Left -is [Management.Automation.Language.IndexExpressionAst] -and
        $_.Left.Extent.Text -ceq '$Ledger[$identity]'
})
if ($ledgerAssignments.Count -ne 1) { throw 'Linux session ledger assignment is not exact' }
$ledgerTables = @($ledgerAssignments[0].Right.FindAll({
    param($node)
    $node -is [Management.Automation.Language.HashtableAst]
}, $true))
if ($ledgerTables.Count -ne 1 -or $ledgerTables[0].KeyValuePairs.Count -ne 2) {
    throw 'Linux session ledger does not preserve its two-field shape'
}
$ledgerPid = @($ledgerTables[0].KeyValuePairs | Where-Object { $_.Item1.Value -ceq 'PID' })
$ledgerStartTime = @($ledgerTables[0].KeyValuePairs | Where-Object { $_.Item1.Value -ceq 'StartTime' })
if ($ledgerPid.Count -ne 1 -or $ledgerStartTime.Count -ne 1) {
    throw 'Linux session ledger does not preserve PID and StartTime fields'
}
$ledgerPidVariables = @($ledgerPid[0].Item2.FindAll({
    param($node)
    $node -is [Management.Automation.Language.VariableExpressionAst]
}, $true))
$ledgerStartTimeVariables = @($ledgerStartTime[0].Item2.FindAll({
    param($node)
    $node -is [Management.Automation.Language.VariableExpressionAst]
}, $true))
if ($ledgerPidVariables.Count -ne 1 -or $ledgerPidVariables[0].VariablePath.UserPath -cne 'linuxProcessIdentifier' -or
    $ledgerStartTimeVariables.Count -ne 1 -or $ledgerStartTimeVariables[0].VariablePath.UserPath -cne 'startTime') {
    throw 'Linux session ledger does not preserve PID and StartTime bindings'
}
Write-Output 'Docker execution iteration-006 PID source regression PASS untyped-mixed-case=REJECT typed-mixed-case=REJECT rhs-and-ledger-property=ALLOW data-flow=PASS'
if (($allSource -join "`n") -match '(?m)&\s+(docker|docker\.exe)\b') { throw 'A workflow-reachable release script retains an ambient Docker invocation' }
if (($allSource -join "`n").Contains('Invoke-ExactReleaseImageInspectProcess')) { throw 'A duplicated Docker runner remains outside the closed entrypoint' }
foreach ($forbidden in @('ScriptBlock','Callback','Invoker','ExecutablePath','ArgumentListInput','DOCKER_CONTEXT','ReleaseDockerInvoker')) {
    if ($production.Substring(0,$production.IndexOf("if (`$MyInvocation")).Contains($forbidden)) { throw "Closed entrypoint exposes forbidden parameter $forbidden" }
}
foreach ($required in @(
    "ValidateSet('EngineInspection','ExactImageInventory','ApprovedImagePull','RepositoryDigestInspection','ContainerCacheProof','ContainerCrlfParse','DependencyAcquisition','ReleaseBuild','ReleasePackage')",
    'CREATE_SUSPENDED','AssignProcessToJobObject','JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE',
    '/usr/bin/setsid','/usr/bin/unshare','--kill-child=SIGKILL','kill -STOP $$','Get-LinuxSessionMembers','Test-LinuxMemberAlive',
    '131072','15000','2000','Environment.Clear()','DOCKER_CONFIG','FileShare]::Read'
)) { if (!$production.Contains($required)) { throw "Closed entrypoint omits required marker: $required" } }
$gate = $production.IndexOf("if (`$null -ne ('PscanNativeBoundary' -as [type]))")
$compile = $production.IndexOf('Add-Type -TypeDefinition $nativeBoundarySource')
$dispatch = $production.IndexOf('[PscanNativeBoundary]::')
$dockerCall = $production.IndexOf('$result = Invoke-BoundDocker -Arguments $arguments')
if ($gate -lt 0 -or $compile -lt 0 -or $dispatch -lt 0 -or $dockerCall -lt 0 -or !($gate -lt $compile -and $compile -lt $dispatch -and $dispatch -lt $dockerCall)) {
    throw 'Static native-type isolation ordering is not fail-closed'
}
$legacyConditional = 'if (!(' + "'PscanNativeBoundary' -as [type]))"
if ($production.Contains($legacyConditional)) { throw 'Production retains conditional native-type reuse' }
$harnessSource = Get-Content -Raw -LiteralPath $PSCommandPath
if ($harnessSource.Contains($legacyConditional)) { throw 'Harness retains conditional ambient-type reuse' }

function Invoke-IsolatedScript([string]$ScriptPath, [string[]]$Arguments, [int]$TimeoutMilliseconds) {
    $start = [Diagnostics.ProcessStartInfo]::new()
    $start.FileName = $hostExecutable
    $start.UseShellExecute = $false
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    foreach ($argument in @('-NoProfile','-NonInteractive','-File',$ScriptPath) + $Arguments) { [void]$start.ArgumentList.Add($argument) }
    $process = [Diagnostics.Process]::new(); $process.StartInfo = $start
    if (!$process.Start()) { throw "Failed to start isolated script: $ScriptPath" }
    $stdoutRead = $process.StandardOutput.ReadToEndAsync(); $stderrRead = $process.StandardError.ReadToEndAsync()
    if (!$process.WaitForExit($TimeoutMilliseconds)) { $process.Kill($true); $process.WaitForExit(); throw "Isolated script exceeded its fixed local bound: $ScriptPath" }
    [pscustomobject]@{ ExitCode=$process.ExitCode; StdOut=$stdoutRead.GetAwaiter().GetResult(); StdErr=$stderrRead.GetAwaiter().GetResult() }
}

$temporaryRoot = Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-005-orchestrator-' + [guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $temporaryRoot)
try {
    $probe=Join-Path $temporaryRoot 'dot-source-probe.ps1';$escaped=$productionPath.Replace("'","''")
    Set-Content -LiteralPath $probe -Encoding utf8NoBOM -Value "try { . '$escaped' -Operation EngineInspection } catch { if (`$_.Exception.Message -notmatch 'cannot be dot-sourced') { throw }; exit 0 }; exit 1"
    $dot = Invoke-IsolatedScript -ScriptPath $probe -Arguments @() -TimeoutMilliseconds 15000
    if ($dot.ExitCode -ne 0) { throw "Dot-source rejection failed: $($dot.StdErr)" }

    foreach ($case in @(@{Mode='HostileCompatible';Name='compatible'},@{Mode='HostileStale';Name='stale'})) {
        $control = Join-Path $temporaryRoot "$($case.Name).json"
        $run = Invoke-IsolatedScript -ScriptPath $PSCommandPath -Arguments @('-InternalMode',$case.Mode,'-ControlPath',$control) -TimeoutMilliseconds 15000
        if ($run.ExitCode -ne 197 -or !(Test-Path -LiteralPath $control -PathType Leaf)) { throw "Hostile $($case.Name) process did not reject nonzero: $($run.StdErr)" }
        $record = Get-Content -Raw -LiteralPath $control | ConvertFrom-Json
        if ($record.StartState -cne 'ABSENT' -or $record.Preloaded -cne 'PscanNativeBoundary' -or !$record.Rejected -or
            $record.FakeBoundaryCalls -ne 0 -or $record.DockerCalls -ne 0 -or $record.OutputCount -ne 0 -or
            $record.TrustedResult -or $record.FabricatedFieldsObserved) {
            throw "Hostile $($case.Name) isolation record is untrusted"
        }
    }

    $clean = Invoke-IsolatedScript -ScriptPath $PSCommandPath -Arguments @('-InternalMode','CleanNative') -TimeoutMilliseconds 180000
    if ($clean.ExitCode -ne 0 -or $clean.StdOut -notmatch 'Docker execution iteration-005 PASS') { throw "Clean isolated native matrix failed: $($clean.StdErr)" }
    Write-Output $clean.StdOut.Trim()
    Write-Output 'Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0'
} finally {
    $resolved=[IO.Path]::GetFullPath($temporaryRoot);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)){throw "Unsafe cleanup path: $resolved"}
    if(Test-Path -LiteralPath $resolved){Remove-Item -LiteralPath $resolved -Recurse -Force}
}
