$ErrorActionPreference = 'Stop'

function Invoke-SourceTrustGit {
    param(
        [Parameter(Mandatory = $true)][string]$Repository,
        [Parameter(Mandatory = $true)][string[]]$Arguments,
        [int[]]$AllowedExitCodes = @(0),
        [switch]$InspectConfiguration
    )

    $startInfo = [Diagnostics.ProcessStartInfo]::new()
    $startInfo.FileName = 'git'
    $startInfo.UseShellExecute = $false
    $startInfo.RedirectStandardOutput = $true
    $startInfo.RedirectStandardError = $true
    $startInfo.Environment['GIT_CONFIG_NOSYSTEM'] = '1'
    $startInfo.Environment['GIT_CONFIG_GLOBAL'] = $(if ($IsWindows) { 'NUL' } else { '/dev/null' })
    $startInfo.Environment['GIT_ATTR_NOSYSTEM'] = '1'
    $startInfo.Environment['GIT_NO_REPLACE_OBJECTS'] = '1'
    $startInfo.Environment['GIT_NO_LAZY_FETCH'] = '1'
    $startInfo.Environment['GIT_OPTIONAL_LOCKS'] = '0'
    foreach ($name in @('GIT_DIR','GIT_WORK_TREE','GIT_INDEX_FILE','GIT_OBJECT_DIRECTORY','GIT_ALTERNATE_OBJECT_DIRECTORIES','GIT_COMMON_DIR','GIT_CEILING_DIRECTORIES','GIT_CONFIG_COUNT')) {
        [void]$startInfo.Environment.Remove($name)
    }
    foreach ($name in @($startInfo.Environment.Keys)) {
        if ($name -match '^GIT_CONFIG_(KEY|VALUE)_\d+$') { [void]$startInfo.Environment.Remove($name) }
    }
    [void]$startInfo.ArgumentList.Add('-C')
    [void]$startInfo.ArgumentList.Add($Repository)
    if (!$InspectConfiguration) {
        foreach ($override in @('core.fsmonitor=false','core.untrackedCache=false','core.ignoreStat=false','core.autocrlf=false','core.eol=lf')) {
            [void]$startInfo.ArgumentList.Add('-c')
            [void]$startInfo.ArgumentList.Add($override)
        }
    }
    foreach ($argument in $Arguments) { [void]$startInfo.ArgumentList.Add($argument) }

    $process = [Diagnostics.Process]::new()
    $process.StartInfo = $startInfo
    [void]$process.Start()
    $memory = [IO.MemoryStream]::new()
    $copy = $process.StandardOutput.BaseStream.CopyToAsync($memory)
    $errorRead = $process.StandardError.ReadToEndAsync()
    $process.WaitForExit()
    $copy.GetAwaiter().GetResult()
    $errorText = $errorRead.GetAwaiter().GetResult()
    $bytes = $memory.ToArray()
    $memory.Dispose()
    if ($AllowedExitCodes -notcontains $process.ExitCode) {
        throw "Git command failed closed (exit $($process.ExitCode)): git $($Arguments -join ' ')`n$errorText"
    }
    return [pscustomobject]@{ ExitCode=$process.ExitCode; Bytes=$bytes; Error=$errorText }
}

function Convert-StrictUtf8([AllowEmptyCollection()][byte[]]$Bytes, [string]$Label) {
    if ($null -eq $Bytes) { $Bytes = [byte[]]::new(0) }
    try { return [Text.UTF8Encoding]::new($false, $true).GetString($Bytes) }
    catch { throw "$Label is not strict UTF-8 and cannot be trusted" }
}

function Get-RawFileGitBlobID([string]$Path) {
    $stream = [IO.File]::Open($Path, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::Read)
    $sha1 = [Security.Cryptography.SHA1]::Create()
    try {
        $header = [Text.Encoding]::ASCII.GetBytes("blob $($stream.Length)`0")
        [void]$sha1.TransformBlock($header, 0, $header.Length, $null, 0)
        $buffer = [byte[]]::new(1024 * 1024)
        while (($read = $stream.Read($buffer, 0, $buffer.Length)) -gt 0) {
            [void]$sha1.TransformBlock($buffer, 0, $read, $null, 0)
        }
        [void]$sha1.TransformFinalBlock([byte[]]::new(0), 0, 0)
        return ([BitConverter]::ToString($sha1.Hash)).Replace('-','').ToLowerInvariant()
    } finally {
        $sha1.Dispose()
        $stream.Dispose()
    }
}

function Get-SourceTrustTreeEntries([string]$Repository, [string]$Tree) {
    $result = Invoke-SourceTrustGit -Repository $Repository -Arguments @('ls-tree','-r','-z','--full-tree',$Tree)
    $text = Convert-StrictUtf8 -Bytes $result.Bytes -Label 'Git tree enumeration'
    $records = @($text.Split([char]0, [StringSplitOptions]::RemoveEmptyEntries))
    if ($records.Count -eq 0) { throw 'Git tree is empty or cannot be enumerated' }
    $entries = [Collections.Generic.List[object]]::new()
    foreach ($record in $records) {
        if ($record -notmatch '^(?<mode>100644|100755) blob (?<object>[0-9a-f]{40})\t(?<path>[A-Za-z0-9._+@/-]+)$') {
            throw "Git tree contains an unsupported type, mode or path: $record"
        }
        $path = $Matches.path
        if ($path.StartsWith('/') -or $path.Contains('//') -or $path.Split('/') -contains '..') {
            throw "Git tree contains an unsafe path: $path"
        }
        $entries.Add([pscustomobject]@{ Mode=$Matches.mode; Object=$Matches.object; Path=$path })
    }
    return $entries.ToArray()
}

function Test-CanonicalLFToCRLFProjection([byte[]]$BlobBytes, [byte[]]$WorkingBytes) {
    $projected = [IO.MemoryStream]::new()
    try {
        for ($index = 0; $index -lt $BlobBytes.Length; $index++) {
            $value = $BlobBytes[$index]
            if ($value -eq 10 -and ($index -eq 0 -or $BlobBytes[$index - 1] -ne 13)) {
                $projected.WriteByte(13)
            }
            $projected.WriteByte($value)
        }
        $expected = $projected.ToArray()
        if ($expected.Length -ne $WorkingBytes.Length) { return $false }
        for ($index = 0; $index -lt $expected.Length; $index++) {
            if ($expected[$index] -ne $WorkingBytes[$index]) { return $false }
        }
        return $true
    } finally {
        $projected.Dispose()
    }
}

function Assert-ExactGitSourceTrust {
    param(
        [Parameter(Mandatory = $true)][string]$Repository,
        [Parameter(Mandatory = $true)][string]$ExpectedRevision,
        [switch]$AllowCanonicalEolProjection
    )

    if ($ExpectedRevision -notmatch '^[0-9a-f]{40}$') { throw 'Expected source revision must be one exact lowercase SHA-1 commit identity' }
    $requestedRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $Repository).Path)
    $rootResult = Invoke-SourceTrustGit -Repository $requestedRoot -Arguments @('rev-parse','--show-toplevel')
    $root = (Convert-StrictUtf8 -Bytes $rootResult.Bytes -Label 'Repository root identity').Trim()
    if ([IO.Path]::GetFullPath($root) -ne $requestedRoot) { throw 'Repository root identity is ambiguous' }

    foreach ($key in @('core.fsmonitor','core.untrackedCache','core.ignoreStat','core.sparseCheckout','extensions.partialClone')) {
        $config = Invoke-SourceTrustGit -Repository $root -Arguments @('config','--get-all',$key) -AllowedExitCodes @(0,1) -InspectConfiguration
        if ($config.ExitCode -eq 0 -and $config.Bytes.Length -gt 0) {
            throw "Unsupported Git shortcut or partial-clone configuration is present: $key"
        }
    }
    $promisor = Invoke-SourceTrustGit -Repository $root -Arguments @('config','--get-regexp','^remote\..*\.promisor$') -AllowedExitCodes @(0,1) -InspectConfiguration
    if ($promisor.ExitCode -eq 0 -and $promisor.Bytes.Length -gt 0) {
        throw 'Unsupported promisor-remote configuration is present'
    }

    $objectFormat = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse','--show-object-format')).Bytes -Label 'Git object format').Trim()
    if ($objectFormat -ne 'sha1') { throw "Unsupported Git object format: $objectFormat" }
    $head = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse','HEAD')).Bytes -Label 'HEAD identity').Trim()
    $commit = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse',"$ExpectedRevision`^{commit}")).Bytes -Label 'Commit identity').Trim()
    if ($head -ne $ExpectedRevision -or $commit -ne $ExpectedRevision) { throw 'HEAD does not equal the exact expected source commit' }
    $tree = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('rev-parse',"$ExpectedRevision`^{tree}")).Bytes -Label 'Tree identity').Trim()
    if ($tree -notmatch '^[0-9a-f]{40}$') { throw 'Expected source tree identity cannot be proven' }

    $entries = @(Get-SourceTrustTreeEntries -Repository $root -Tree $tree)
    $expectedIndex = [Collections.Generic.List[string]]::new()
    foreach ($entry in $entries) { $expectedIndex.Add("$($entry.Mode) $($entry.Object) 0`t$($entry.Path)") }
    $actualIndexText = Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('ls-files','--stage','-z')).Bytes -Label 'Git index'
    $actualIndex = @($actualIndexText.Split([char]0, [StringSplitOptions]::RemoveEmptyEntries))
    if ($actualIndex.Count -ne $expectedIndex.Count) { throw 'Complete index path count does not equal HEAD' }
    for ($index = 0; $index -lt $expectedIndex.Count; $index++) {
        if ($actualIndex[$index] -cne $expectedIndex[$index]) { throw "Index path/mode/object identity differs from HEAD: $($actualIndex[$index])" }
    }

    $flagsText = Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('ls-files','-v','-z')).Bytes -Label 'Git index flags'
    foreach ($record in @($flagsText.Split([char]0, [StringSplitOptions]::RemoveEmptyEntries))) {
        if ($record -cnotmatch '^H (?<path>.+)$') { throw "Assume-unchanged, skip-worktree or unsupported index state is present: $record" }
    }
    $fsmonitorText = Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('ls-files','-f','-z')).Bytes -Label 'Git fsmonitor flags'
    foreach ($record in @($fsmonitorText.Split([char]0, [StringSplitOptions]::RemoveEmptyEntries))) {
        if ($record -cnotmatch '^H (?<path>.+)$') { throw "Filesystem-monitor-valid or unsupported index state is present: $record" }
    }

    $untracked = Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $root -Arguments @('ls-files','--others','-z')).Bytes -Label 'Untracked-file enumeration'
    $untrackedPaths = @($untracked.Split([char]0, [StringSplitOptions]::RemoveEmptyEntries))
    if ($untrackedPaths.Count -ne 0) { throw "Untracked files are forbidden, including ignored files: $($untrackedPaths[0])" }

    $rawEqual = 0
    $canonicalEol = 0
    foreach ($entry in $entries) {
        $nativePath = Join-Path $root ($entry.Path.Replace('/', [IO.Path]::DirectorySeparatorChar))
        if (!(Test-Path -LiteralPath $nativePath -PathType Leaf)) { throw "Tracked working-tree path is missing or not a regular file: $($entry.Path)" }
        $item = Get-Item -LiteralPath $nativePath -Force
        if (($item.Attributes -band [IO.FileAttributes]::ReparsePoint) -ne 0) { throw "Tracked working-tree path is a reparse point: $($entry.Path)" }
        if (!$IsWindows) {
            $mode = [IO.File]::GetUnixFileMode($nativePath)
            $executeMask = [IO.UnixFileMode]::UserExecute -bor [IO.UnixFileMode]::GroupExecute -bor [IO.UnixFileMode]::OtherExecute
            $hasExecute = ($mode -band $executeMask) -ne 0
            if (($entry.Mode -eq '100755') -ne $hasExecute) { throw "Tracked working-tree executable mode differs from HEAD: $($entry.Path)" }
        } elseif ($entry.Mode -ne '100644') {
            throw "Executable Git mode cannot be proved on this Windows filesystem: $($entry.Path)"
        }

        $rawObject = Get-RawFileGitBlobID -Path $nativePath
        if ($rawObject -eq $entry.Object) {
            $rawEqual++
            continue
        }
        if (!$AllowCanonicalEolProjection) { throw "Tracked working-tree raw bytes differ from HEAD: $($entry.Path)" }
        $blobBytes = (Invoke-SourceTrustGit -Repository $root -Arguments @('cat-file','blob',$entry.Object)).Bytes
        $workingBytes = [IO.File]::ReadAllBytes($nativePath)
        if (!(Test-CanonicalLFToCRLFProjection -BlobBytes $blobBytes -WorkingBytes $workingBytes)) {
            throw "Tracked working-tree bytes are neither HEAD nor the allowed canonical CRLF projection: $($entry.Path)"
        }
        $canonicalEol++
    }

    return [pscustomobject]@{
        Repository=$root; Commit=$ExpectedRevision; Tree=$tree; FileCount=$entries.Count
        RawEqualCount=$rawEqual; CanonicalEolProjectionCount=$canonicalEol
        WorkingTreeInputsTrusted=($canonicalEol -eq 0)
    }
}

function Export-ExactGitTreeMaterialization {
    param(
        [Parameter(Mandatory = $true)][string]$Repository,
        [Parameter(Mandatory = $true)][string]$Revision,
        [Parameter(Mandatory = $true)][string]$Destination
    )
    if (Test-Path -LiteralPath $Destination) {
        if (Get-ChildItem -LiteralPath $Destination -Force | Select-Object -First 1) { throw 'Exact materialization destination must be empty' }
    } else {
        New-Item -ItemType Directory -Path $Destination | Out-Null
    }
    $tree = (Convert-StrictUtf8 -Bytes (Invoke-SourceTrustGit -Repository $Repository -Arguments @('rev-parse',"$Revision`^{tree}")).Bytes -Label 'Materialization tree').Trim()
    $entries = @(Get-SourceTrustTreeEntries -Repository $Repository -Tree $tree)
    foreach ($entry in $entries) {
        $destinationPath = Join-Path $Destination ($entry.Path.Replace('/', [IO.Path]::DirectorySeparatorChar))
        $parent = Split-Path -Parent $destinationPath
        if (!(Test-Path -LiteralPath $parent -PathType Container)) { New-Item -ItemType Directory -Path $parent -Force | Out-Null }
        $blobBytes = (Invoke-SourceTrustGit -Repository $Repository -Arguments @('cat-file','blob',$entry.Object)).Bytes
        [IO.File]::WriteAllBytes($destinationPath, $blobBytes)
        if ((Get-RawFileGitBlobID -Path $destinationPath) -ne $entry.Object) { throw "Exact materialized blob mismatch: $($entry.Path)" }
        if (!$IsWindows -and $entry.Mode -eq '100755') {
            [IO.File]::SetUnixFileMode($destinationPath, [IO.UnixFileMode]::UserRead -bor [IO.UnixFileMode]::UserWrite -bor [IO.UnixFileMode]::UserExecute -bor [IO.UnixFileMode]::GroupRead -bor [IO.UnixFileMode]::GroupExecute -bor [IO.UnixFileMode]::OtherRead -bor [IO.UnixFileMode]::OtherExecute)
        }
    }
    return [pscustomobject]@{ Commit=$Revision; Tree=$tree; FileCount=$entries.Count; Destination=[IO.Path]::GetFullPath($Destination) }
}
