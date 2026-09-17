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

function Get-LinuxFixtureMarkerPaths([string]$Mode,[string]$PidFile) {
    if($Mode-notin@('child','grandchild','hold-child','hold-grandchild','detach')){throw 'Unknown descendant fixture mode'}
    if($Mode.EndsWith('grandchild')){return @($PidFile,($PidFile+'.child'),($PidFile+'.child.ready'))}
    return @($PidFile,($PidFile+'.ready'))
}

function Get-WindowsFixtureMarkerPaths([string]$Mode,[string]$PidFile) {
    if($Mode-eq'child'){return @($PidFile)}
    if($Mode-eq'grandchild'){return @($PidFile,($PidFile+'.child'))}
    throw 'Unknown Windows descendant fixture mode'
}

function Invoke-CleanNativeFixture {
    . (Join-Path $PSScriptRoot 'native-fixture-diagnostics.ps1')
    if ($null -ne ('PscanNativeBoundary' -as [type])) { throw 'Clean fixture process began with the native boundary already loaded' }
    $parsed = Get-NativeBoundarySource
    $compiled = @(Add-Type -TypeDefinition $parsed.Source -Language CSharp -PassThru)
    Assert-CompiledBoundaryIdentity -Types $compiled

    $temporaryRoot = Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-005-clean-' + [guid]::NewGuid().ToString('N'))
    [void](New-Item -ItemType Directory -Path $temporaryRoot)
    try {
        if ($IsLinux) {
            foreach ($functionName in @('Get-LinuxSessionMembers','Test-LinuxMemberAlive','Get-LinuxBoundarySnapshot','New-LinuxBoundaryStartInfo','Get-LinuxGateAction','Test-LinuxInitReaped','Invoke-LinuxPinnedCleanup','Invoke-LinuxSessionBoundary')) {
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
function Marker([string]$Path,[int]$Identifier,[string]$Role){
 $text=[IO.File]::ReadAllText("/proc/$Identifier/stat");$close=$text.LastIndexOf(')');$fields=$text.Substring($close+2).Split(' ',[StringSplitOptions]::RemoveEmptyEntries)
 $ns=[string]((Get-Item -LiteralPath "/proc/$Identifier/ns/pid" -Force).Target)
 $value=@{PID=$Identifier;StartTime=[uint64]$fields[19];Session=[int]$fields[3];Namespace=$ns;Role=$Role}|ConvertTo-Json -Compress
 [IO.File]::WriteAllText($Path+'.tmp',$value);[IO.File]::Move($Path+'.tmp',$Path)
}
function WaitMarker([string]$Path){$clock=[Diagnostics.Stopwatch]::StartNew();while(!(Test-Path -LiteralPath $Path)){if($clock.ElapsedMilliseconds-ge 10000){throw 'Synthetic descendant readiness missing'};Start-Sleep -Milliseconds 5}}
$hold=$Mode.StartsWith('hold-');$fixtureMode=$Mode.Replace('hold-','')
switch($fixtureMode){
'stdout'{Bytes $out 97 $Count}'stderr'{Bytes $err 98 $Count}'both'{for($left=$Count;$left-gt 0;$left-=$n){$n=[Math]::Min(1024,$left);Bytes $out 97 $n;Bytes $err 98 $n}}
'split-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush();Start-Sleep -Milliseconds 25;$out.WriteByte(0xAC);$out.Flush()}
'invalid-utf8'{$out.WriteByte(0xC3);$out.WriteByte(0x28);$out.Flush()}'incomplete-utf8'{$out.WriteByte(0xE2);$out.WriteByte(0x82);$out.Flush()}'immediate'{}'nonzero'{exit 7}'hang'{if($PidFile){Marker $PidFile $PID 'ready'};Start-Sleep -Seconds 30}
 'child'{$ready=$PidFile+'.ready';$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang','-PidFile',$ready) -PassThru -NoNewWindow;Marker $PidFile $p.Id 'child';WaitMarker $ready;if($hold){Start-Sleep -Seconds 30}}
'detach'{$ready=$PidFile+'.ready';$p=Start-Process -FilePath /usr/bin/setsid -ArgumentList @((Get-Process -Id $PID).Path,'-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang','-PidFile',$ready) -PassThru -NoNewWindow;Marker $PidFile $p.Id 'detach';WaitMarker $ready;Start-Sleep -Seconds 30}
'grandchild'{$childFile=$PidFile+'.child';$p=Start-Process -FilePath (Get-Process -Id $PID).Path -ArgumentList @('-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hold-child','-PidFile',$childFile) -PassThru -NoNewWindow;Marker $PidFile $p.Id 'grandchild';WaitMarker ($childFile+'.ready');if($hold){Start-Sleep -Seconds 30}}
default{throw 'unknown fixture mode'}}
'@
            $environment=[Collections.Generic.Dictionary[string,string]]::new([StringComparer]::Ordinal);$environment.Add('HOME',$temporaryRoot);$environment.Add('TMPDIR',$temporaryRoot);$environment.Add('LANG','C.UTF-8')
            $reference="/proc/$PID/fd/$($identity.SafeFileHandle.DangerousGetHandle().ToInt64())"
            function RunLinux([string]$mode,[int]$count=0,[string]$pidFile=''){$a=@('-NoProfile','-NonInteractive','-File',$fixture,'-Mode',$mode,'-Count',[string]$count);if($pidFile){$a+=@('-PidFile',$pidFile)};Invoke-LinuxSessionBoundary -ExecutableReference $reference -Arguments $a -Environment $environment -WorkingDirectory $temporaryRoot}
            function LinuxSuccess($r,[string]$name){if($r.Terminal-or!$r.ContainmentEmpty){[Console]::Error.WriteLine((ConvertTo-PscanNativeFixtureDiagnostic -Result $r -CaseName $name));throw "$name returned terminal lifecycle evidence"}}
            function LinuxTerminal($r,[string]$name){if(!$r.Terminal-or$r.Terminal.Contains('cleanup uncertainty')){[Console]::Error.WriteLine((ConvertTo-PscanNativeFixtureDiagnostic -Result $r -CaseName $name));throw "$name did not prove expected terminal cleanup"}}
            foreach($n in @(131071,131072)){$r=RunLinux stdout $n;LinuxSuccess $r "stdout-$n";if($r.StdOut.Length-ne$n){throw 'stdout limit mismatch'};$r=RunLinux stderr $n;LinuxSuccess $r "stderr-$n";if($r.StdErr.Length-ne$n){throw 'stderr limit mismatch'}}
            $r=RunLinux immediate;LinuxSuccess $r immediate;$r=RunLinux both 131072;LinuxSuccess $r both;$utf8=[Text.UTF8Encoding]::new($false,$true);$r=RunLinux split-utf8;LinuxSuccess $r split;if($utf8.GetString($r.StdOut)-cne '€'){throw 'split UTF-8 mismatch'};$r=RunLinux nonzero;LinuxSuccess $r nonzero;if($r.ExitCode-ne 7){throw 'nonzero mismatch'}
            foreach($case in @(@{M='stdout';C=131073},@{M='stderr';C=131073},@{M='both';C=131073},@{M='hang';C=0})){LinuxTerminal (RunLinux $case.M $case.C) $case.M}
            foreach($mode in @('invalid-utf8','incomplete-utf8')){$r=RunLinux $mode;LinuxSuccess $r $mode;$failed=$false;try{[void]$utf8.GetString($r.StdOut)}catch{$failed=$true};if(!$failed){throw "$mode unexpectedly decoded"}}
            foreach($mode in @('child','grandchild','hold-child','hold-grandchild','detach')){
                $pidFile=Join-Path $temporaryRoot "$mode.pid";$r=RunLinux $mode 0 $pidFile
                if($mode-in@('child','grandchild')){LinuxSuccess $r $mode;if($r.ExitCode-ne 0){throw "$mode root exit failed"}}
                else{LinuxTerminal $r $mode}
                $namespaces=@($r.ObservedMembers|Where-Object{$_-match'^pid:\[[0-9]+\]$'})
                if($namespaces.Count-ne 1){throw "$mode namespace identity was not proved"}
                $paths=@(Get-LinuxFixtureMarkerPaths $mode $pidFile)
                $markers=@(foreach($path in $paths){
                    if(!(Test-Path -LiteralPath $path -PathType Leaf)-or(Get-Item -LiteralPath $path).Length-gt 2048){throw "$mode descendant marker missing or oversized"}
                    $record=[IO.File]::ReadAllText($path)|ConvertFrom-Json
                    if($record.PID-le 1-or$record.StartTime-le 0-or$record.Namespace-cne$namespaces[0]){throw "$mode descendant identity is not namespace-bound"}
                    $record
                })
                $ready=$markers[-1];$spawned=$markers[-2]
                if($ready.Role-cne'ready'-or$ready.PID-ne$spawned.PID-or$ready.StartTime-ne$spawned.StartTime){throw "$mode intended descendant did not become ready"}
                if($mode-eq'detach'-and$ready.Session-ne$ready.PID){throw 'Detached descendant did not establish a new session'}
                # Namespace-local IDs are never interpreted as host PIDs. The
                # boundary pins init and proves its death, which kernel-terminates
                # every namespace descendant, including detached sessions. Also
                # reject any surviving exact host identity in the observed ledger.
                foreach($entry in $r.ObservedMembers){
                    if($entry-match'^([0-9]+):([0-9]+)$'){
                        $member=[pscustomobject]@{PID=[int]$Matches[1];StartTime=[uint64]$Matches[2]}
                        if(Test-LinuxMemberAlive $member){throw "$mode left an observed host process alive"}
                    }
                }
            }
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
            foreach($mode in @('child','grandchild')){
                $pidFile=Join-Path $temporaryRoot "$mode.pid";$r=Run $mode 0 $pidFile;Terminal $r $mode
                foreach($path in @(Get-WindowsFixtureMarkerPaths $mode $pidFile)){
                    if(!(Test-Path -LiteralPath $path -PathType Leaf)-or(Get-Item -LiteralPath $path).Length-gt 64){throw "$mode required Windows descendant marker missing or oversized"}
                    $id=[int](Get-Content -Raw -LiteralPath $path)
                    if($id-le 0){throw "$mode invalid Windows descendant identifier"}
                    if(Get-Process -Id $id -ErrorAction SilentlyContinue){throw "$mode left live member $id"}
                }
            }
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
    'image-admission.ps1','test-cache-boundary.ps1','test-crlf-shell-payloads.ps1',
    'docker-container-lifecycle.ps1','execution-profile.ps1','build-validation.ps1','invoke-docker-boundary.ps1'
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
$isAutomaticPidVariablePath = {
    param([Management.Automation.VariablePath]$Path)

    if (!$Path.IsVariable -or $Path.IsDriveQualified) { return $false }
    if ($Path.IsUnqualified) { return $Path.UserPath -ieq 'PID' }
    if ($Path.IsUnscopedVariable) { return $Path.UserPath -ieq 'variable:PID' }
    if ($Path.IsGlobal) { return $Path.UserPath -ieq 'global:PID' }
    if ($Path.IsScript) { return $Path.UserPath -ieq 'script:PID' }
    if ($Path.IsLocal) { return $Path.UserPath -ieq 'local:PID' }
    if ($Path.IsPrivate) { return $Path.UserPath -ieq 'private:PID' }
    return $false
}
$getAssignmentTargetVariableNames = {
    param([Management.Automation.Language.Ast]$Target)

    $targetVariableNames = [Collections.Generic.List[string]]::new()
    $walkAssignmentTarget = $null
    $walkAssignmentTarget = {
        param([Management.Automation.Language.Ast]$Node)

        if ($Node -is [Management.Automation.Language.AttributedExpressionAst]) {
            & $walkAssignmentTarget -Node $Node.Child
            return
        }
        if ($Node -is [Management.Automation.Language.ParenExpressionAst]) {
            & $walkAssignmentTarget -Node $Node.Pipeline
            return
        }
        if ($Node -is [Management.Automation.Language.PipelineAst]) {
            if ($Node.PipelineElements.Count -ne 1) {
                throw "PID assignment target walker found a non-assignable pipeline: $($Node.Extent.Text)"
            }
            & $walkAssignmentTarget -Node $Node.PipelineElements[0]
            return
        }
        if ($Node -is [Management.Automation.Language.CommandExpressionAst]) {
            if ($Node.Redirections.Count -ne 0) {
                throw "PID assignment target walker found a redirected command expression: $($Node.Extent.Text)"
            }
            & $walkAssignmentTarget -Node $Node.Expression
            return
        }
        if ($Node -is [Management.Automation.Language.ArrayLiteralAst]) {
            foreach ($element in $Node.Elements) {
                & $walkAssignmentTarget -Node $element
            }
            return
        }
        if ($Node -is [Management.Automation.Language.VariableExpressionAst]) {
            if (!$Node.Splatted -and (& $isAutomaticPidVariablePath -Path $Node.VariablePath)) {
                [void]$targetVariableNames.Add('PID')
            }
            return
        }
        if ($Node -is [Management.Automation.Language.MemberExpressionAst] -or
            $Node -is [Management.Automation.Language.IndexExpressionAst]) {
            return
        }
        throw "PID assignment target walker found an unsupported AST shape: $($Node.GetType().FullName)"
    }

    & $walkAssignmentTarget -Node $Target
    $targetVariableNames.ToArray()
}
$isReservedPidAssignment = {
    param([Management.Automation.Language.AssignmentStatementAst]$Assignment)
    $targetVariableNames = @(& $getAssignmentTargetVariableNames -Target $Assignment.Left)
    @($targetVariableNames | Where-Object { $_ -ieq 'PID' }).Count -gt 0
}
$nonTargetPidSource = @'
$value = $PID
$text = '$PID = 1'
# $PiD = 1
$ledger = [pscustomobject]@{ PID=$value }
[int]$other = 1
'@
$assignmentPredicateCases = @(
    @{ Name='untyped-mixed-case'; Source='$pId = 1'; Expected=1 },
    @{ Name='typed-mixed-case'; Source='[int]$PiD = 1'; Expected=1 },
    @{ Name='parenthesized-mixed-case'; Source='($PiD) = 1'; Expected=1 },
    @{ Name='multi-target-mixed-case'; Source='$PiD, $other = 1, 2'; Expected=1 },
    @{ Name='parenthesized-multi-target-mixed-case'; Source='($PiD, $other) = 1, 2'; Expected=1 },
    @{ Name='nested-typed-parenthesized-scoped'; Source='[int](($script:PiD)) = 1'; Expected=1 },
    @{ Name='nested-multi-target-elements'; Source='(($other)), (([int]$PiD)) = 1, 2'; Expected=1 },
    @{ Name='background-parenthesized'; Source='($PiD &) = 1'; Expected=1 },
    @{ Name='automatic-variable-scope-targets'; Source='$PID = 1; ${PID} = 1; $variable:PID = 1; ${variable:PID} = 1; $global:PID = 1; ${global:PID} = 1; $script:PID = 1; ${script:PID} = 1; $local:PID = 1; ${local:PID} = 1; $private:PID = 1; ${private:PID} = 1'; Expected=12 },
    @{ Name='member-and-index-targets'; Source='$array[$PID] = 1; $object.PID = 1; ($array[$PID]) = 2; ($object.PID) = 2'; Expected=0 },
    @{ Name='typed-member-and-index-targets'; Source='[int]$array[$PID] = 1; [int]$object.PID = 1'; Expected=0 },
    @{ Name='drive-qualified-non-automatic-targets'; Source='$env:PID = 1; ${env:PID} = 1; ${function:PID} = { 1 }; ${foo:bar:PID} = 1'; Expected=0 },
    @{ Name='multi-colon-true-variable-targets'; Source='${variable:env:PID} = 1; ${global:env:PID} = 1; ${local:foo:PID} = 1; ${script:foo:PID} = 1; ${private:foo:PID} = 1'; Expected=0 },
    @{ Name='rhs-string-comment-and-ledger-property'; Source=$nonTargetPidSource; Expected=0 }
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
$assignmentRuntimeCases = @(
    @{ Name='direct'; Source='$PID = 1'; Collision=$true },
    @{ Name='braced-direct'; Source='${PID} = 1'; Collision=$true },
    @{ Name='variable-scope'; Source='${variable:PID} = 1'; Collision=$true },
    @{ Name='global-scope'; Source='$global:PID = 1'; Collision=$true },
    @{ Name='braced-global-scope'; Source='${global:PID} = 1'; Collision=$true },
    @{ Name='script-scope'; Source='$script:PID = 1'; Collision=$true },
    @{ Name='braced-script-scope'; Source='${script:PID} = 1'; Collision=$true },
    @{ Name='local-scope'; Source='$local:PID = 1'; Collision=$true },
    @{ Name='braced-local-scope'; Source='${local:PID} = 1'; Collision=$true },
    @{ Name='private-scope'; Source='$private:PID = 1'; Collision=$true },
    @{ Name='braced-private-scope'; Source='${private:PID} = 1'; Collision=$true },
    @{ Name='variable-multi-colon'; Source='${variable:env:PID} = 1'; Collision=$false },
    @{ Name='global-multi-colon'; Source='${global:env:PID} = 1'; Collision=$false },
    @{ Name='local-multi-colon'; Source='${local:foo:PID} = 1'; Collision=$false },
    @{ Name='script-multi-colon'; Source='${script:foo:PID} = 1'; Collision=$false },
    @{ Name='private-multi-colon'; Source='${private:foo:PID} = 1'; Collision=$false },
    @{ Name='environment-provider'; Source='$env:PID = ''control'''; Collision=$false },
    @{ Name='braced-environment-provider'; Source='${env:PID} = ''control'''; Collision=$false },
    @{ Name='function-provider'; Source='${function:PID} = { ''control'' }'; Collision=$false }
)
foreach ($case in $assignmentRuntimeCases) {
    $runtimeBody = "`$ErrorActionPreference='Stop'; try { function Invoke-PidAssignmentProbe { $($case.Source) }; Invoke-PidAssignmentProbe; [Console]::Out.WriteLine('EXECUTED_WITHOUT_PID_COLLISION'); exit 0 } catch { [Console]::Error.WriteLine(`$_.Exception.Message); exit 73 }"
    $runtimeStart = [Diagnostics.ProcessStartInfo]::new()
    $runtimeStart.FileName = $hostExecutable
    $runtimeStart.UseShellExecute = $false
    $runtimeStart.CreateNoWindow = $true
    $runtimeStart.RedirectStandardOutput = $true
    $runtimeStart.RedirectStandardError = $true
    foreach ($argument in @('-NoProfile','-NonInteractive','-Command',$runtimeBody)) {
        [void]$runtimeStart.ArgumentList.Add($argument)
    }
    $runtimeProcess = [Diagnostics.Process]::new()
    $runtimeProcess.StartInfo = $runtimeStart
    try {
        [void]$runtimeProcess.Start()
        $runtimeOutputRead = $runtimeProcess.StandardOutput.ReadToEndAsync()
        $runtimeErrorRead = $runtimeProcess.StandardError.ReadToEndAsync()
        if (!$runtimeProcess.WaitForExit(10000)) {
            $runtimeProcess.Kill($true)
            throw "PID assignment runtime self-test timed out: $($case.Name)"
        }
        $runtimeOutput = $runtimeOutputRead.GetAwaiter().GetResult().Trim()
        $runtimeError = $runtimeErrorRead.GetAwaiter().GetResult().Trim()
        $runtimeExitCode = $runtimeProcess.ExitCode
    } finally {
        $runtimeProcess.Dispose()
    }
    $runtimeCollision = $runtimeExitCode -eq 73 -and
        $runtimeError -match '(?i)Cannot overwrite variable PID because it is read-only or constant'
    $runtimeAllowed = $runtimeExitCode -eq 0 -and $runtimeOutput -ceq 'EXECUTED_WITHOUT_PID_COLLISION'
    if (($case.Collision -and !$runtimeCollision) -or (!$case.Collision -and !$runtimeAllowed)) {
        throw "PID assignment runtime self-test failed: $($case.Name) exit=$runtimeExitCode output=$runtimeOutput error=$runtimeError"
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
Write-Output 'Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT runtime-scopes=REJECT multi-colon-provider-controls=ALLOW member-index=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS'
if (($allSource -join "`n") -match '(?m)&\s+(docker|docker\.exe)\b') { throw 'A workflow-reachable release script retains an ambient Docker invocation' }
if (($allSource -join "`n").Contains('Invoke-ExactReleaseImageInspectProcess')) { throw 'A duplicated Docker runner remains outside the closed entrypoint' }
foreach ($forbidden in @('ScriptBlock','Callback','Invoker','ExecutablePath','ArgumentListInput','DOCKER_CONTEXT','ReleaseDockerInvoker')) {
    if ($production.Substring(0,$production.IndexOf("if (`$MyInvocation")).Contains($forbidden)) { throw "Closed entrypoint exposes forbidden parameter $forbidden" }
}
foreach ($required in @(
    "ValidateSet('EngineInspection','ExactImageInventory','ApprovedImagePull','RepositoryDigestInspection','ContainerCacheProof','ContainerCrlfParse','DependencyAcquisition','ReleaseBuild','ReleasePackage')",
    'CREATE_SUSPENDED','AssignProcessToJobObject','JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE',
    '/usr/bin/setsid','/usr/bin/unshare','--kill-child=SIGKILL','Get-LinuxGateAction','Get-LinuxSessionMembers','Test-LinuxMemberAlive',
    '131072','Get-ReleaseOperationBudget $Operation','2000','Environment.Clear()','DOCKER_CONFIG','FileShare]::Read'
)) { if (!$production.Contains($required)) { throw "Closed entrypoint omits required marker: $required" } }
# Exercise the actual ancestor-stop, pinned identity, stopped-state and inherited
# inode gate before the isolated native matrix. A retired shell self-stop string
# cannot prove this protocol (namespace PID 1 does not accept that self-stop).
& (Join-Path $PSScriptRoot 'test-linux-handoff.ps1')
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
