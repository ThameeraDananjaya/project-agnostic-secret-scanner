# Fixed release-tooling profile; no environment or caller-selected durations.
function Get-ReleaseExecutionProfile {
    return [ordered]@{
        identifier='pscan-release-execution-v2'
        operationBudgetMilliseconds=[ordered]@{
            EngineInspection=15000; ExactImageInventory=15000
            ApprovedImagePull=180000; RepositoryDigestInspection=15000
            ContainerCacheProof=15000; ContainerCrlfParse=15000
            DependencyAcquisition=300000; ReleaseBuild=900000; ReleasePackage=120000
        }
        streamLimitBytes=131072; cleanupGraceMilliseconds=2000
        childTransportAllowanceMilliseconds=20000; maximumProtocolCalls=9
    }
}

function Get-ReleaseOperationBudget([string]$Operation) {
    $profile=Get-ReleaseExecutionProfile
    if ($profile.operationBudgetMilliseconds.Keys -cnotcontains $Operation) { throw 'Unknown exact release operation' }
    return [int]$profile.operationBudgetMilliseconds[$Operation]
}

function Assert-ReleaseProfileShape($Actual,$Expected) {
    if ($Actual -isnot [Collections.IDictionary] -or $Actual.Count -ne $Expected.Count) { throw 'Release profile fields differ' }
    foreach ($key in $Expected.Keys) {
        if ($Actual.Keys -cnotcontains $key) { throw 'Release profile field missing or case differs' }
        $value=$Actual[$key];$wanted=$Expected[$key]
        if ($wanted -is [Collections.IDictionary]) { Assert-ReleaseProfileShape $value $wanted }
        elseif ($wanted -is [int]) {
            if (($value -isnot [int] -and $value -isnot [long]) -or $value -ne $wanted) { throw 'Release profile numeric value differs' }
        } elseif ($value -isnot [string] -or $value -cne $wanted) { throw 'Release profile value differs' }
    }
}

function Read-ReleaseImageAdmission([string]$Path) {
    # The local receipt is strict and versioned. Reject duplicates before a
    # PowerShell map could silently collapse conflicting JSON evidence.
    $stream=[IO.File]::OpenRead($Path)
    try {
        $buffer=[byte[]]::new(16385);$length=0
        while ($length -lt $buffer.Length) { $read=$stream.Read($buffer,$length,$buffer.Length-$length);if($read -eq 0){break};$length+=$read }
        if ($length -gt 16384) { throw 'Image admission receipt exceeds bound' }
        $bytes=[byte[]]::new($length);[Array]::Copy($buffer,$bytes,$length)
    } finally { $stream.Dispose() }
    $json=[Text.UTF8Encoding]::new($false,$true).GetString($bytes)
    $options=[Text.Json.JsonDocumentOptions]::new();$options.MaxDepth=8
    $document=[Text.Json.JsonDocument]::Parse($json,$options)
    try {
        $pending=[Collections.Generic.Stack[Text.Json.JsonElement]]::new();$pending.Push($document.RootElement);$count=0
        while ($pending.Count) {
            if (++$count -gt 128) { throw 'Image admission receipt structure exceeds bound' }
            $value=$pending.Pop()
            if ($value.ValueKind -eq [Text.Json.JsonValueKind]::Object) {
                $names=[Collections.Generic.HashSet[string]]::new([StringComparer]::OrdinalIgnoreCase)
                foreach ($property in $value.EnumerateObject()) {
                    if (!$names.Add($property.Name)) { throw 'Duplicate image admission receipt field' }
                    $pending.Push($property.Value)
                }
            } elseif ($value.ValueKind -eq [Text.Json.JsonValueKind]::Array) { throw 'Image admission receipt contains an array' }
        }
    } finally { $document.Dispose() }
    $receipt=ConvertFrom-Json -InputObject $json -AsHashtable -Depth 8
    $keys=@('schemaVersion','sourceRevision','image','dockerExecutableSHA256','pulledDuringAdmission','hostIdentityMode','hostUID','hostGID','containment','executionProfile')
    if ($receipt -isnot [Collections.IDictionary] -or $receipt.Count -ne $keys.Count) { throw 'Image admission receipt fields differ' }
    foreach ($key in $keys) { if ($receipt.Keys -cnotcontains $key) { throw 'Image admission receipt field missing or case differs' } }
    foreach ($key in @('schemaVersion','sourceRevision','image','dockerExecutableSHA256','hostIdentityMode','containment')) {
        if ($receipt[$key] -isnot [string]) { throw 'Image admission receipt string field malformed' }
    }
    if ($receipt.schemaVersion -isnot [string] -or $receipt.schemaVersion -cne '2.0' -or
        $receipt.sourceRevision -isnot [string] -or $receipt.sourceRevision -cnotmatch '^[0-9a-f]{40}$' -or
        $receipt.dockerExecutableSHA256 -isnot [string] -or $receipt.dockerExecutableSHA256 -cnotmatch '^[0-9a-f]{64}$' -or
        $receipt.image -cne 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452' -or
        $receipt.pulledDuringAdmission -isnot [bool] -or $receipt.containment -cne 'empty-after-every-operation') { throw 'Image admission receipt identity malformed' }
    if ($receipt.hostIdentityMode -ceq 'linux-host-numeric-uid-gid') {
        foreach ($key in @('hostUID','hostGID')) {
            if (($receipt[$key] -isnot [int] -and $receipt[$key] -isnot [long]) -or $receipt[$key] -lt 0 -or $receipt[$key] -gt [int]::MaxValue) { throw 'Image admission numeric identity malformed' }
        }
    } elseif ($receipt.hostIdentityMode -cne 'windows-invoking-host' -or $null -ne $receipt.hostUID -or $null -ne $receipt.hostGID) { throw 'Image admission host identity malformed' }
    Assert-ReleaseProfileShape $receipt.executionProfile (Get-ReleaseExecutionProfile)
    return $receipt
}
