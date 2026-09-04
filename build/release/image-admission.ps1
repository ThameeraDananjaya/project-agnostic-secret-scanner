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

function Get-ExactDockerExecutableForVerification {
    if ($IsWindows) {
        $root = [Environment]::GetFolderPath([Environment+SpecialFolder]::ProgramFiles)
        $candidate = [IO.Path]::GetFullPath([IO.Path]::Combine($root, 'Docker', 'Docker', 'resources', 'bin', 'docker.exe'))
        if (![IO.File]::Exists($candidate)) { throw 'Docker executable is absent from the fixed Docker Desktop installation path' }
        $signature = Microsoft.PowerShell.Security\Get-AuthenticodeSignature -LiteralPath $candidate
        if ($signature.Status -ne [Management.Automation.SignatureStatus]::Valid -or $null -eq $signature.SignerCertificate -or $signature.SignerCertificate.Subject -notmatch '(?i)\bDocker\b') { throw 'Docker executable signature or publisher identity is untrusted' }
    } elseif ($IsLinux) {
        $candidate = '/usr/bin/docker'
        if (![IO.File]::Exists($candidate)) { throw 'Docker executable is absent from the fixed Linux system path' }
        $item = [IO.FileInfo]::new($candidate)
        if ($null -ne $item.LinkTarget) { throw 'Docker executable identity is untrusted because the fixed Linux path is a symbolic link' }
        $mode = [IO.File]::GetUnixFileMode($candidate)
        if (($mode -band ([IO.UnixFileMode]::GroupWrite -bor [IO.UnixFileMode]::OtherWrite)) -ne 0) { throw 'Docker executable identity is mutable' }
    } else { throw 'Docker executable identity is unsupported on this operating system' }
    $candidate
}

function Invoke-ExactReleaseImageInspectProcess {
    $limit = 131072
    $budget = 15000
    $cleanupGrace = 2000
    $process = [Diagnostics.Process]::new()
    $stdout = [IO.MemoryStream]::new(); $stderr = [IO.MemoryStream]::new()
    $stdoutBuffer = [byte[]]::new(4096); $stderrBuffer = [byte[]]::new(4096)
    $clock = [Diagnostics.Stopwatch]::StartNew(); $reason = $null
    try {
        $start = [Diagnostics.ProcessStartInfo]::new()
        $start.FileName = Get-ExactDockerExecutableForVerification
        $start.UseShellExecute = $false; $start.CreateNoWindow = $true
        $start.RedirectStandardOutput = $true; $start.RedirectStandardError = $true
        $start.WorkingDirectory = [IO.Path]::GetTempPath(); $start.Environment.Clear()
        if ($IsWindows) {
            $systemRoot = [Environment]::GetFolderPath([Environment+SpecialFolder]::Windows)
            if ([string]::IsNullOrWhiteSpace($systemRoot)) { throw 'Docker inspection cannot establish a fixed system root' }
            $start.Environment['SystemRoot']=$systemRoot; $start.Environment['WINDIR']=$systemRoot
            $start.Environment['TEMP']=[IO.Path]::GetTempPath(); $start.Environment['TMP']=[IO.Path]::GetTempPath()
        } else { $start.Environment['HOME']='/tmp'; $start.Environment['LANG']='C.UTF-8'; $start.Environment['LC_ALL']='C.UTF-8' }
        $start.Environment['DOCKER_CLI_HINTS']='false'; $start.Environment['DOCKER_CONFIG']=[IO.Path]::Combine([IO.Path]::GetTempPath(), 'pscan-empty-docker-config')
        foreach ($argument in @('image','inspect','--format','{{json .RepoDigests}}','docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452')) { [void]$start.ArgumentList.Add($argument) }
        $process.StartInfo=$start
        try { if (!$process.Start()) { throw 'start returned false' } } catch { throw 'Docker repository-digest inspection start failed and is terminal untrusted evidence' }
        $outTask=$process.StandardOutput.BaseStream.ReadAsync($stdoutBuffer,0,$stdoutBuffer.Length); $errTask=$process.StandardError.BaseStream.ReadAsync($stderrBuffer,0,$stderrBuffer.Length)
        $outClosed=$false; $errClosed=$false; $complete=$false
        while (!$complete -and $null -eq $reason) {
            if ($clock.ElapsedMilliseconds -ge $budget) { $reason='Docker repository-digest inspection exceeded the complete monotonic wall-clock budget'; break }
            if (!$outClosed -and $outTask.IsCompleted) { try{$n=$outTask.GetAwaiter().GetResult()}catch{$reason='Docker repository-digest inspection stdout read failed';break}; if($n -eq 0){$outClosed=$true}else{if(($stdout.Length+$n)-gt$limit){$reason='Docker repository-digest inspection stdout exceeded the byte limit';break};$stdout.Write($stdoutBuffer,0,$n);$outTask=$process.StandardOutput.BaseStream.ReadAsync($stdoutBuffer,0,$stdoutBuffer.Length)} }
            if (!$errClosed -and $errTask.IsCompleted) { try{$n=$errTask.GetAwaiter().GetResult()}catch{$reason='Docker repository-digest inspection stderr read failed';break}; if($n -eq 0){$errClosed=$true}else{if(($stderr.Length+$n)-gt$limit){$reason='Docker repository-digest inspection stderr exceeded the byte limit';break};$stderr.Write($stderrBuffer,0,$n);$errTask=$process.StandardError.BaseStream.ReadAsync($stderrBuffer,0,$stderrBuffer.Length)} }
            try{$exited=$process.HasExited}catch{$reason='Docker repository-digest inspection exit state is untrusted';break}
            $complete=$exited -and $outClosed -and $errClosed; if(!$complete){[Threading.Thread]::Sleep(1)}
        }
        if($null -ne $reason){$cleanup=[Diagnostics.Stopwatch]::StartNew();try{if(!$process.HasExited){$process.Kill($true)};$left=$cleanupGrace-[int]$cleanup.ElapsedMilliseconds;$killed=if($left-gt 0){$process.WaitForExit($left)}else{$process.HasExited}}catch{throw "$reason; process-tree termination is uncertain and terminal"};if(!$killed){throw "$reason; process-tree termination exceeded the fixed cleanup grace and is terminal"};while($cleanup.ElapsedMilliseconds-lt$cleanupGrace-and(!$outTask.IsCompleted-or!$errTask.IsCompleted)){[Threading.Thread]::Sleep(1)};if(!$outTask.IsCompleted-or!$errTask.IsCompleted){throw "$reason; redirected pipe closure is uncertain and terminal"};throw "$reason; command evidence is terminal and cannot be retried or reclassified"}
        $utf8=[Text.UTF8Encoding]::new($false,$true);try{$outText=$utf8.GetString($stdout.ToArray());$errText=$utf8.GetString($stderr.ToArray())}catch{throw 'Docker repository-digest inspection returned invalid or incomplete UTF-8 and is terminal untrusted evidence'}
        [pscustomobject]@{ExitCode=$process.ExitCode;StdOut=$outText;StdErr=$errText}
    } finally {$stdout.Dispose();$stderr.Dispose();$process.Dispose();$clock.Stop()}
}

function Assert-AdmittedReleaseImage {
    param([string]$Image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452')
    $exact='docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
    if($Image -cne $exact){throw 'Build image must use the exact canonical Correction C2 repository and digest reference'}
    $result=Invoke-ExactReleaseImageInspectProcess
    if($result.ExitCode-ne 0-or![string]::IsNullOrEmpty($result.StdErr)){throw 'Docker repository-digest inspection failed and is terminal untrusted evidence'}
    $digests=Read-ReleaseRepoDigestsEvidence -Json $result.StdOut
    Assert-ReleaseImageIdentityEvidence -Image $exact -RepoDigests $digests
}

function Invoke-ReleaseImageAdmission {
    param([Parameter(Mandatory=$true)][string]$CacheDirectory,[Alias('AllowImagePull')][switch]$LegacyImagePullRequest)
    if($LegacyImagePullRequest){throw 'Image pull authority is restricted to the mandatory release-image bootstrap orchestrator'}
    $cacheRoot=[IO.Path]::GetFullPath((Microsoft.PowerShell.Management\Resolve-Path -LiteralPath $CacheDirectory).Path)
    $downloads=[IO.Path]::Combine($cacheRoot,'downloads');$moduleCache=[IO.Path]::Combine($cacheRoot,'gomodcache')
    foreach($path in @($downloads,$moduleCache)){if(!(Microsoft.PowerShell.Management\Test-Path -LiteralPath $path -PathType Container)){throw "Exact acquisition cache path is absent before image verification: $path"}}
    $hostProof=&{. ([IO.Path]::Combine($PSScriptRoot,'host-cache-canary.ps1'));$proof=$null;foreach($path in @($cacheRoot,$downloads,$moduleCache)){$proof=Invoke-HostCacheCanary -CacheDirectory $path};$proof}
    $canonical=Assert-AdmittedReleaseImage
    [pscustomobject]@{Image=$canonical;Pulled=$false;HostIdentityMode=$hostProof.IdentityMode;HostUID=$hostProof.HostUID;HostGID=$hostProof.HostGID;ProvedCachePaths=@($cacheRoot,$downloads,$moduleCache)}
}
