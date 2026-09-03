function ConvertTo-LFPosixShellPayload {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Payload
    )

    $normalized = $Payload.Replace("`r`n", "`n").Replace("`r", "`n")
    Assert-LFPosixShellPayload -Payload $normalized
    return $normalized
}

function Assert-LFPosixShellPayload {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][AllowEmptyString()][string]$Payload
    )

    if ($Payload.IndexOf([char]13) -ge 0) {
        throw 'POSIX shell payload contains a carriage return and cannot be passed to Docker'
    }
    if ($Payload.IndexOf([char]0) -ge 0) {
        throw 'POSIX shell payload contains a NUL and cannot be passed to Docker'
    }
}
