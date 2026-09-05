$ErrorActionPreference = 'Stop'

$productionPath = Join-Path $PSScriptRoot 'docker-execution.ps1'
$releaseFiles = @(
    'admit-image.ps1','acquire.ps1','build.ps1','cache-canary.ps1',
    'image-admission.ps1','test-cache-boundary.ps1','test-crlf-shell-payloads.ps1'
)
$allSource = foreach ($name in $releaseFiles) { Get-Content -Raw -LiteralPath (Join-Path $PSScriptRoot $name) }
$production = Get-Content -Raw -LiteralPath $productionPath

if (($allSource -join "`n") -match '(?m)&\s+(docker|docker\.exe)\b') { throw 'A workflow-reachable release script retains an ambient Docker invocation' }
if (($allSource -join "`n").Contains('Invoke-ExactReleaseImageInspectProcess')) { throw 'A duplicated Docker runner remains outside the closed entrypoint' }
foreach ($forbidden in @('ScriptBlock','Callback','Invoker','ExecutablePath','ArgumentListInput','DOCKER_CONTEXT','ReleaseDockerInvoker')) {
    if ($production.Substring(0,$production.IndexOf("if (`$MyInvocation")).Contains($forbidden)) { throw "Closed entrypoint exposes forbidden parameter $forbidden" }
}
foreach ($required in @(
    "ValidateSet('EngineInspection','ExactImageInventory','ApprovedImagePull','RepositoryDigestInspection','ContainerCacheProof','ContainerCrlfParse','DependencyAcquisition','ReleaseBuild','ReleasePackage')",
    'CREATE_SUSPENDED','AssignProcessToJobObject','JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE',
    '/usr/bin/setsid','Get-LinuxSessionMembers','Test-LinuxMemberAlive',
    '131072','15000','2000','Environment.Clear()','DOCKER_CONFIG','FileShare]::Read'
)) { if (!$production.Contains($required)) { throw "Closed entrypoint omits required marker: $required" } }

$tokens=$null;$errors=$null
$ast=[Management.Automation.Language.Parser]::ParseFile($productionPath,[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Closed Docker entrypoint does not parse'}
$nativeSources=@($ast.FindAll({param($node)$node-is[Management.Automation.Language.StringConstantExpressionAst]-and$node.Value.Contains('public static class PscanNativeBoundary')},$true))
if($nativeSources.Count-ne 1){throw 'Expected one immutable native boundary source'}
if(!('PscanNativeBoundary'-as[type])){Add-Type -TypeDefinition $nativeSources[0].Value -Language CSharp}

$temporaryRoot=Join-Path ([IO.Path]::GetTempPath()) ('pscan-c2-iteration-004-'+[guid]::NewGuid().ToString('N'))
[void](New-Item -ItemType Directory -Path $temporaryRoot)
try {
    $probe=Join-Path $temporaryRoot 'dot-source-probe.ps1'
    $escaped=$productionPath.Replace("'","''")
    Set-Content -LiteralPath $probe -Encoding utf8NoBOM -Value "try { . '$escaped' -Operation EngineInspection } catch { if (`$_.Exception.Message -notmatch 'cannot be dot-sourced') { throw }; exit 0 }; exit 1"
    & (Get-Process -Id $PID).Path -NoProfile -NonInteractive -File $probe
    if($LASTEXITCODE-ne 0){throw 'Dot-source rejection failed'}

    if($IsLinux){
        foreach($functionName in @('Get-LinuxSessionMembers','Test-LinuxMemberAlive','Invoke-LinuxSessionBoundary')){
            $definition=$ast.Find({param($node)$node-is[Management.Automation.Language.FunctionDefinitionAst]-and$node.Name-ceq$functionName},$true)
            if($null-eq$definition){throw "Missing Linux boundary function $functionName"};. ([scriptblock]::Create($definition.Extent.Text))
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
'detach'{$p=Start-Process -FilePath /usr/bin/setsid -ArgumentList @((Get-Process -Id $PID).Path,'-NoProfile','-NonInteractive','-File',$PSCommandPath,'-Mode','hang') -PassThru -NoNewWindow;Set-Content $PidFile $p.Id;Start-Sleep -Milliseconds 500}
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
        Write-Output 'Docker execution iteration-004 PASS single-boundary=PASS private-session=PASS identity-ledger=PASS streams=PASS utf8=PASS process-tree=EMPTY'
        return
    }
    if(!$IsWindows){throw 'Native containment fixtures require Windows or Linux'}

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
    Write-Output 'Docker execution iteration-004 PASS single-boundary=PASS job-assignment-before-resume=PASS kill-on-close=PASS streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT'
} finally {
    $resolved=[IO.Path]::GetFullPath($temporaryRoot);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)){throw "Unsafe cleanup path: $resolved"}
    if(Test-Path -LiteralPath $resolved){Remove-Item -LiteralPath $resolved -Recurse -Force}
}
