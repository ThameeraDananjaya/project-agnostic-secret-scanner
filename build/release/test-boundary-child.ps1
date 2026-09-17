$ErrorActionPreference='Stop'
$tokens=$null;$errors=$null
$tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1'),[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Boundary child launcher does not parse'}
$function=$tree.Find({param($n)$n-is[Management.Automation.Language.FunctionDefinitionAst]-and$n.Name-ceq'Invoke-BoundaryChild'},$true)
if($null-eq$function){throw 'Actual child transport missing'}
. ([scriptblock]::Create($function.Extent.Text))
$hostPath=(Get-Process -Id $PID).Path
function StartInfo([string]$Code){
    $s=[Diagnostics.ProcessStartInfo]::new();$s.FileName=$hostPath;$s.UseShellExecute=$false;$s.CreateNoWindow=$true
    $s.RedirectStandardOutput=$true;$s.RedirectStandardError=$true
    foreach($arg in @('-NoLogo','-NoProfile','-NonInteractive','-Command',$Code)){$s.ArgumentList.Add($arg)}
    return $s
}
$cases=0
$first=(Invoke-BoundaryChild (StartInfo '[Console]::Out.Write($PID)')).Trim()
$second=(Invoke-BoundaryChild (StartInfo '[Console]::Out.Write($PID)')).Trim()
if($first-notmatch'^\d+$'-or$second-notmatch'^\d+$'-or$first-eq$second-or$first-eq[string]$PID-or$second-eq[string]$PID){throw 'Sequential operations did not use fresh children'};$cases++
$large=Invoke-BoundaryChild (StartInfo '$v=[string]::new([char]1,131072);[Console]::Out.Write((@{StdOut=$v;StdErr=$v}|ConvertTo-Json -Compress))')
$value=$large|ConvertFrom-Json
if($value.StdOut.Length-ne 131072-or$value.StdErr.Length-ne 131072-or[Text.Encoding]::UTF8.GetByteCount($large)-le 131072){throw 'Transport narrowed escaped Docker stream evidence'};$cases++
foreach($code in @('exit 7','[Console]::Error.Write("failure")','[Console]::OpenStandardOutput().Write([byte[]]@(0xc3,0x28))','[Console]::Out.Write([string]::new([char]65,2097153))')){
    $rejected=$false;try{[void](Invoke-BoundaryChild (StartInfo $code))}catch{$rejected=$true}
    if(!$rejected){throw 'Child failure or malformed/oversized transport was accepted'};$cases++
}
# Exercise the actual deadline/cleanup path with a shorter test-only deadline.
$short=$function.Extent.Text.Replace('-ge 35000','-ge 2000')
if($short-ceq$function.Extent.Text){throw 'Transport deadline site changed'}
. ([scriptblock]::Create($short))
$watch=[Diagnostics.Stopwatch]::StartNew();$rejected=$false
try{[void](Invoke-BoundaryChild (StartInfo 'Start-Sleep -Seconds 30'))}catch{$rejected=$true}
if(!$rejected-or$watch.ElapsedMilliseconds-gt 6000){throw 'Timeout/cleanup did not fail within the test bound'};$cases++
$launcher=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'invoke-docker-boundary.ps1'))
$begin=$launcher.IndexOf('$document=[Text.Json.JsonDocument]::Parse($json)')
$end=$launcher.IndexOf('$global:LASTEXITCODE=0',$begin)
if($begin-lt 0-or$end-le$begin){throw 'Actual evidence validator sites changed'}
$validator=[scriptblock]::Create($launcher.Substring($begin,$end-$begin))
$Operation='EngineInspection'
$valid=@{Operation=$Operation;DockerSHA256=('a'*64);ExitCode=0;StdOut='ok';StdErr='';StdOutByteCount=2;StdErrByteCount=0;ContainmentMembersObserved=0;ContainmentEmpty=$true}|ConvertTo-Json -Compress
$json=$valid;&$validator;$cases++
foreach($bad in @($valid.Replace('"StdOutByteCount":2','"StdOutByteCount":3'),$valid.Replace('"ContainmentEmpty":true','"ContainmentEmpty":false'),$valid.Replace('"Operation":','"Extra":0,"Operation":'),$valid.Replace('"Operation":','"Operation":"duplicate","Operation":'),($valid+'{}'),'[]')){
    $json=$bad;$rejected=$false;try{&$validator}catch{$rejected=$true};if(!$rejected){throw 'Malformed child evidence accepted'};$cases++
}
Write-Output "Fresh boundary child transport PASS cases=$cases synthetic-native-children=YES Docker-execution=NONE"
