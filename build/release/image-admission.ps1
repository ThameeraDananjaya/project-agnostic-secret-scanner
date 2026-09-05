function ConvertFrom-ReleaseJsonDocument {
    param([Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json, [Parameter(Mandatory = $true)][string]$Operation)
    if ([string]::IsNullOrWhiteSpace($Json)) { throw "$Operation returned empty structured data" }
    try { [Text.Json.JsonDocument]::Parse($Json) } catch { throw "$Operation returned malformed structured data" }
}

function Assert-ReleaseEngineEvidence {
    param([Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json)
    $document = ConvertFrom-ReleaseJsonDocument -Json $Json -Operation 'Docker engine reachability inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Object) { throw 'Docker engine reachability inspection returned a non-object result' }
        $properties = @{}
        foreach ($property in $document.RootElement.EnumerateObject()) { if ($properties.ContainsKey($property.Name)) { throw 'Docker engine reachability inspection returned duplicate fields' }; $properties[$property.Name] = $property.Value }
        foreach ($name in @('Version','ApiVersion','Os','Arch')) {
            if (!$properties.ContainsKey($name) -or $properties[$name].ValueKind -ne [Text.Json.JsonValueKind]::String -or [string]::IsNullOrWhiteSpace($properties[$name].GetString())) {
                throw 'Docker engine reachability inspection returned incomplete structured data'
            }
        }
        'Responsive'
    } finally { $document.Dispose() }
}

function Resolve-ReleaseImageListEvidence {
    param([Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json)
    if ([string]::IsNullOrWhiteSpace($Json)) { return 'ConclusiveAbsent' }
    $lines = @([regex]::Split($Json.Trim(), '\r?\n'))
    if ($lines.Count -ne 1 -or [string]::IsNullOrWhiteSpace($lines[0])) { throw 'Docker exact-image inventory inspection returned ambiguous structured data' }
    $document = ConvertFrom-ReleaseJsonDocument -Json $lines[0] -Operation 'Docker exact-image inventory inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Object) { throw 'Docker exact-image inventory inspection returned a non-object result' }
        $properties = @{}
        foreach ($property in $document.RootElement.EnumerateObject()) { if ($properties.ContainsKey($property.Name)) { throw 'Docker exact-image inventory inspection returned duplicate fields' }; $properties[$property.Name] = $property.Value }
        foreach ($name in @('Repository','Digest','ID')) { if (!$properties.ContainsKey($name) -or $properties[$name].ValueKind -ne [Text.Json.JsonValueKind]::String) { throw 'Docker exact-image inventory inspection returned incomplete structured data' } }
        if ($properties.Repository.GetString() -cne 'golang' -or $properties.Digest.GetString() -cne 'sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452' -or $properties.ID.GetString() -cnotmatch '^sha256:[0-9a-f]{64}$') {
            throw 'Docker exact-image inventory inspection returned unexpected identity data'
        }
        'CandidatePresent'
    } finally { $document.Dispose() }
}

function Read-ReleaseRepoDigestsEvidence {
    param([Parameter(Mandatory = $true)][AllowEmptyString()][string]$Json)
    $document = ConvertFrom-ReleaseJsonDocument -Json $Json -Operation 'Docker repository-digest inspection'
    try {
        if ($document.RootElement.ValueKind -ne [Text.Json.JsonValueKind]::Array) { throw 'Docker repository-digest inspection returned a non-array result' }
        $values = [Collections.Generic.List[string]]::new()
        foreach ($element in $document.RootElement.EnumerateArray()) { if ($element.ValueKind -ne [Text.Json.JsonValueKind]::String) { throw 'Docker repository-digest inspection returned a non-string item' }; $values.Add($element.GetString()) }
        ,$values.ToArray()
    } finally { $document.Dispose() }
}

function Assert-ReleaseImageIdentityEvidence {
    param([Parameter(Mandatory = $true)][string]$Image, [Parameter(Mandatory = $true)][AllowNull()]$RepoDigests)
    $exact = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
    if ($Image -cne $exact) { throw 'Build image must use the exact canonical Correction C2 repository and digest reference' }
    if ($null -eq $RepoDigests -or $RepoDigests -isnot [array]) { throw 'Docker repository-digest evidence must be one structured array' }
    $values = @($RepoDigests)
    if ($values.Count -ne 1 -or $values[0] -isnot [string] -or $values[0] -cne 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452') {
        throw 'Docker repository-digest evidence must contain exactly one canonical engine identity'
    }
    $Image
}
