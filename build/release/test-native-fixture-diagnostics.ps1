$ErrorActionPreference='Stop'
. (Join-Path $PSScriptRoot 'native-fixture-diagnostics.ps1')
$script:checks=0
function Assert($Condition,[string]$Name) {
    if (!$Condition) { throw "Diagnostic formatter assertion failed: $Name" }
    $script:checks++
}
function Fixture {
    [pscustomobject]@{Terminal='timeout or incomplete lifecycle';ContainmentEmpty=$false;ExitCode=1;
        StdOut=[byte[]]::new(0);StdErr=[Text.Encoding]::UTF8.GetBytes("unshare: unshare failed: Operation not permitted`n")}
}
function Format($Value,$Name='stdout-131071') {
    $json=ConvertTo-PscanNativeFixtureDiagnostic -Result $Value -CaseName $Name
    Assert ($json -is [string]) 'one string result'
    Assert ([Text.Encoding]::UTF8.GetByteCount($json) -le 2048) 'bounded serialized result'
    Assert (!$json.Contains('DO-NOT-LOG-THIS-SECRET')) 'no supplied stream or unknown reason content'
    $json | ConvertFrom-Json
}
$r=Format (Fixture)
Assert ($r.terminalReason -ceq 'timeout or incomplete lifecycle') 'known reason retained'
Assert ($r.ContainmentEmpty.state -eq 'present' -and $r.ContainmentEmpty.value -ceq $false) 'false preserved'
Assert ($r.ExitCode.state -eq 'present' -and $r.ExitCode.value -ceq 1) 'nonzero preserved'
Assert ($r.stdout.bytes -eq 0 -and $r.stdout.state -eq 'present') 'actual empty byte array'
Assert ($r.stderr.classification -eq 'unshare-operation-not-permitted') 'exact fixed launcher classification'

foreach ($reason in @('ambiguous PID-namespace init membership','malformed PID-namespace identity',
    'PID-namespace init did not stop at the pre-execution gate','PID-namespace init resume failed',
    'stream overflow or read failure','timeout or incomplete lifecycle')) {
    foreach ($suffix in @('','; cleanup uncertainty')) {
        $value=Fixture; $value.Terminal=$reason+$suffix; $r=Format $value
        Assert ($r.terminalReason -ceq $value.Terminal) 'known reason and cleanup suffix retained'
    }
}
foreach ($name in @('Terminal','ContainmentEmpty','ExitCode','StdOut','StdErr')) {
    $value=Fixture; $value.PSObject.Properties.Remove($name); $r=Format $value
    $state=switch ($name) {'Terminal' {$r.terminalState} 'StdOut' {$r.stdout.state} 'StdErr' {$r.stderr.state} default {$r.$name.state}}
    Assert ($state -eq 'missing') 'missing explicitly marked'
}
$r=Format $null
Assert ($r.ExitCode.value -eq $null -and $r.ExitCode.state -eq 'invalid-result') 'no default zero for absent result'
Assert ($r.ContainmentEmpty.value -eq $null -and $r.ContainmentEmpty.state -eq 'invalid-result') 'no default success for absent result'
$value=Fixture; $value.ExitCode='1'; $value.ContainmentEmpty='false'; $value.StdOut='abc'; $value.Terminal=1
$r=Format $value
Assert ($r.ExitCode.state -eq 'invalid-type' -and $r.ContainmentEmpty.state -eq 'invalid-type') 'no scalar coercion'
Assert ($r.stdout.state -eq 'invalid-type' -and $r.terminalState -eq 'invalid-type') 'invalid stream and reason types'
$value=Fixture; $value.StdOut=$null; $value.StdErr=$null; $value.Terminal=$null; $r=Format $value
Assert ($r.stdout.state -eq 'unavailable' -and $r.stdout.bytes -eq $null) 'null stream not presumed empty'
Assert ($r.stderr.classification -eq 'unavailable' -and $r.terminalState -eq 'null') 'null values explicit'
$value=Fixture; $value.StdOut=[object[]]@([byte]1,[byte]2); $r=Format $value
Assert ($r.stdout.state -eq 'present' -and $r.stdout.bytes -eq 2) 'PowerShell enumerated bytes'
$value.StdOut=[object[]]@(1,'x'); $r=Format $value
Assert ($r.stdout.state -eq 'invalid-elements' -and $r.stdout.bytes -eq $null) 'nonbyte elements rejected'
$value.StdOut=[byte[,]]::new(2,2); $r=Format $value
Assert ($r.stdout.state -in @('invalid-type','invalid-rank')) 'multidimensional stream rejected'
$value=Fixture; $value.StdOut=[byte[]]::new(131073); $r=Format $value
Assert ($r.stdout.state -eq 'oversize' -and $r.stdout.bytes -eq 131073) 'oversize metadata only'
$value=Fixture; $value.StdErr=[byte[]]::new(4097); $r=Format $value
Assert ($r.stderr.classification -eq 'truncated-without-content' -and $r.stderr.bytes -eq 4097) 'stderr inspection cap'
$value=Fixture; $value.Terminal='DO-NOT-LOG-THIS-SECRET'*100; $r=Format $value
Assert ($r.terminalState -eq 'truncated-without-content' -and $r.terminalReason -eq $null) 'terminal truncation without raw prefix'
$value.Terminal='DO-NOT-LOG-THIS-SECRET'; $value.StdOut=[Text.Encoding]::UTF8.GetBytes('DO-NOT-LOG-THIS-SECRET')
$value.StdErr=[Text.Encoding]::UTF8.GetBytes('DO-NOT-LOG-THIS-SECRET'); $r=Format $value
Assert ($r.terminalState -eq 'unrecognized-without-content' -and $r.stderr.classification -eq 'unrecognized-without-content') 'unknown content omitted'
$value=Fixture; $value.StdErr=[byte[]]@(0xc3,0x28); $r=Format $value
Assert ($r.stderr.classification -eq 'invalid-utf8-without-content') 'invalid UTF8 metadata'
$value=Fixture; $value.StdErr=[byte[]]::new(0); $r=Format $value
Assert ($r.stderr.classification -eq 'empty') 'empty stderr distinct from unavailable'
$value=Fixture; $value.PSObject.Properties.Remove('ExitCode')
$value | Add-Member -MemberType ScriptProperty -Name ExitCode -Value { throw 'must not invoke property code' }
$r=Format $value
Assert ($r.ExitCode.state -eq 'invalid-property' -and $r.ExitCode.value -eq $null) 'no property getter execution'
$r=Format (Fixture) 'DO-NOT-LOG-THIS-SECRET'
Assert ($r.caseState -eq 'invalid' -and $r.case -eq $null) 'closed case-name set'

# Validate integration syntax without executing the native fixture or production boundary.
foreach ($name in @('native-fixture-diagnostics.ps1','test-native-fixture-diagnostics.ps1','test-docker-execution.ps1')) {
    $tokens=$null; $parseErrors=$null
    [void][Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $name),[ref]$tokens,[ref]$parseErrors)
    Assert ($parseErrors.Count -eq 0) "parse $name"
}
Write-Output "Synthetic diagnostic formatter PASS assertions=$script:checks native-execution=NONE docker-execution=NONE"
