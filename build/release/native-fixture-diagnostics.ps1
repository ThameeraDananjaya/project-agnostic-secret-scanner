# Pure formatter for the fixed synthetic fixture. No stdout or environment dump.
function ConvertTo-PscanNativeFixtureDiagnostic($Result, $CaseName) {
    function Field([string]$Name) {
        if ($Result -isnot [pscustomobject]) { return @{ state='invalid-result'; value=$null } }
        $property = $Result.PSObject.Properties[$Name]
        if ($null -eq $property) { return @{ state='missing'; value=$null } }
        if ($property.MemberType -ne [Management.Automation.PSMemberTypes]::NoteProperty) { return @{ state='invalid-property'; value=$null } }
        return @{ state='present'; value=$property.Value }
    }
    function StreamMetadata([string]$Name) {
        $field = Field $Name
        if ($field.state -ne 'present') { return @{ state=$field.state; bytes=$null; data=$null } }
        $value = $field.value
        if ($null -eq $value) { return @{ state='unavailable'; bytes=$null; data=$null } }
        if ($value -isnot [byte[]] -and $value -isnot [object[]]) { return @{ state='invalid-type'; bytes=$null; data=$null } }
        if ($value.Rank -ne 1) { return @{ state='invalid-rank'; bytes=$null; data=$null } }
        if ($value.Length -gt 131072) { return @{ state='oversize'; bytes=$value.Length; data=$null } }
        if ($value -is [object[]]) {
            foreach ($item in $value) { if ($item -isnot [byte]) { return @{ state='invalid-elements'; bytes=$null; data=$null } } }
        }
        return @{ state='present'; bytes=$value.Length; data=([byte[]]$value) }
    }
    $record = [ordered]@{ schema='pscan-native-fixture-diagnostic-v2'; case=$null; caseState='invalid' }
    $knownCases=@('stdout-131071','stdout-131072','stderr-131071','stderr-131072',
        'immediate','both','split','nonzero','invalid-utf8','incomplete-utf8')
    if ($CaseName -is [string] -and $CaseName.Length -le 64 -and $knownCases -ccontains $CaseName) {
        $record.case=$CaseName; $record.caseState='present'
    }
    $terminal = Field 'Terminal'
    $record.terminalState=$terminal.state; $record.terminalReason=$null
    if ($terminal.state -eq 'present') {
        if ($null -eq $terminal.value) { $record.terminalState='null' }
        elseif ($terminal.value -isnot [string]) { $record.terminalState='invalid-type' }
        elseif ($terminal.value.Length -gt 256) { $record.terminalState='truncated-without-content' }
        else {
            $known = @('ambiguous PID-namespace init membership','malformed PID-namespace identity',
                'PID-namespace init did not stop at the pre-execution gate','PID-namespace init resume failed',
                'stream overflow or read failure','timeout or incomplete lifecycle')
            $reason = $terminal.value -creplace '; cleanup uncertainty$', ''
            if ($known -ccontains $reason) { $record.terminalReason=$terminal.value }
            else { $record.terminalState='unrecognized-without-content' }
        }
    }
    foreach ($name in @('ExitCode','ContainmentEmpty')) {
        $field = Field $name; $value=$null; $state=$field.state
        if ($state -eq 'present') {
            if (($name -eq 'ExitCode' -and $field.value -is [int]) -or
                ($name -eq 'ContainmentEmpty' -and $field.value -is [bool])) { $value=$field.value }
            else { $state='invalid-type' }
        }
        $record[$name]=[ordered]@{state=$state;value=$value}
    }
    $stdout=StreamMetadata 'StdOut'; $stderr=StreamMetadata 'StdErr'
    $record.stdout=[ordered]@{state=$stdout.state;bytes=$stdout.bytes}
    $record.stderr=[ordered]@{state=$stderr.state;bytes=$stderr.bytes;
        excerptEncoding='escaped-bytes';excerptSourceBytes=$null;excerptTruncated=$null;excerpt=$null}
    if ($stderr.state -eq 'present') {
        $count=[Math]::Min(256,$stderr.bytes)
        $escaped=[Text.StringBuilder]::new()
        for ($index=0;$index-lt$count;$index++) {
            $b=$stderr.data[$index]
            if ($b -eq 92) { [void]$escaped.Append('\\') }
            elseif ($b -ge 32 -and $b -le 126) { [void]$escaped.Append([char]$b) }
            else { [void]$escaped.Append('\x').Append($b.ToString('x2')) }
        }
        $record.stderr.excerptSourceBytes=$count
        $record.stderr.excerptTruncated=($stderr.bytes -gt $count)
        $record.stderr.excerpt=$escaped.ToString()
    }
    $json=$record | ConvertTo-Json -Depth 5 -Compress
    if ([Text.Encoding]::UTF8.GetByteCount($json) -gt 2048) { throw 'Synthetic diagnostic exceeded its fixed encoding bound' }
    return $json
}
