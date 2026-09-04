param(
    [Parameter(Mandatory = $true)][string]$CacheDirectory,
    [Parameter(Mandatory = $true)][string]$SourceRepository,
    [Parameter(Mandatory = $true)][string]$SourceRevision,
    [Parameter(Mandatory = $true)][string]$WorkingDirectory
)

$ErrorActionPreference = 'Stop'
if ($MyInvocation.InvocationName -eq '.') { throw 'The release image admission entrypoint cannot be dot-sourced' }

$releaseImage = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$dockerOutputLimit = 131072
$dockerBudgetMilliseconds = 15000
$cleanupGraceMilliseconds = 2000
. ([IO.Path]::Combine($PSScriptRoot, 'image-admission.ps1'))

function Get-ClosedDockerExecutable {
    if ($IsWindows) {
        $programFiles = [Environment]::GetFolderPath([Environment+SpecialFolder]::ProgramFiles)
        if ([string]::IsNullOrWhiteSpace($programFiles)) { throw 'Docker executable identity is untrusted because Program Files is unavailable' }
        $candidate = [IO.Path]::GetFullPath([IO.Path]::Combine($programFiles, 'Docker', 'Docker', 'resources', 'bin', 'docker.exe'))
        if (![IO.File]::Exists($candidate)) { throw 'Docker executable is absent from the fixed Docker Desktop installation path' }
        $signature = Microsoft.PowerShell.Security\Get-AuthenticodeSignature -LiteralPath $candidate
        if ($signature.Status -ne [Management.Automation.SignatureStatus]::Valid -or $null -eq $signature.SignerCertificate -or $signature.SignerCertificate.Subject -notmatch '(?i)\bDocker\b') {
            throw 'Docker executable signature or publisher identity is untrusted'
        }
    } elseif ($IsLinux) {
        $candidate = '/usr/bin/docker'
        if (![IO.File]::Exists($candidate)) { throw 'Docker executable is absent from the fixed Linux system path' }
        $item = [IO.FileInfo]::new($candidate)
        if ($null -ne $item.LinkTarget) { throw 'Docker executable identity is untrusted because the fixed Linux path is a symbolic link' }
        $mode = [IO.File]::GetUnixFileMode($candidate)
        $mutable = [IO.UnixFileMode]::GroupWrite -bor [IO.UnixFileMode]::OtherWrite
        if (($mode -band $mutable) -ne 0) { throw 'Docker executable identity is untrusted because it is group- or world-writable' }
    } else { throw 'Docker executable identity is unsupported on this operating system' }
    $digest = (Microsoft.PowerShell.Utility\Get-FileHash -LiteralPath $candidate -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($digest -cnotmatch '^[0-9a-f]{64}$') { throw 'Docker executable digest evidence is malformed' }
    [pscustomobject]@{ Path = $candidate; Sha256 = $digest }
}

function Invoke-ClosedNativeProcess {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$ExecutablePath,
        [Parameter(Mandatory = $true)][string[]]$Arguments,
        [Parameter(Mandatory = $true)][string]$Operation
    )
    $clock = [Diagnostics.Stopwatch]::StartNew()
    $process = [Diagnostics.Process]::new()
    $stdout = [IO.MemoryStream]::new()
    $stderr = [IO.MemoryStream]::new()
    $stdoutBuffer = [byte[]]::new(4096)
    $stderrBuffer = [byte[]]::new(4096)
    $terminalReason = $null
    try {
        $startInfo = [Diagnostics.ProcessStartInfo]::new()
        $startInfo.FileName = $ExecutablePath
        $startInfo.UseShellExecute = $false
        $startInfo.CreateNoWindow = $true
        $startInfo.RedirectStandardOutput = $true
        $startInfo.RedirectStandardError = $true
        $startInfo.WorkingDirectory = [IO.Path]::GetTempPath()
        $startInfo.Environment.Clear()
        if ($IsWindows) {
            $systemRoot = [Environment]::GetFolderPath([Environment+SpecialFolder]::Windows)
            if ([string]::IsNullOrWhiteSpace($systemRoot)) { throw "$Operation cannot establish a fixed system root" }
            $startInfo.Environment['SystemRoot'] = $systemRoot
            $startInfo.Environment['WINDIR'] = $systemRoot
            $startInfo.Environment['TEMP'] = [IO.Path]::GetTempPath()
            $startInfo.Environment['TMP'] = [IO.Path]::GetTempPath()
        } else {
            $startInfo.Environment['HOME'] = '/tmp'
            $startInfo.Environment['LANG'] = 'C.UTF-8'
            $startInfo.Environment['LC_ALL'] = 'C.UTF-8'
        }
        $startInfo.Environment['DOCKER_CLI_HINTS'] = 'false'
        $startInfo.Environment['DOCKER_CONFIG'] = [IO.Path]::Combine([IO.Path]::GetTempPath(), 'pscan-empty-docker-config')
        foreach ($argument in $Arguments) { [void]$startInfo.ArgumentList.Add($argument) }
        $process.StartInfo = $startInfo
        try { if (!$process.Start()) { throw "$Operation could not be started" } }
        catch { throw "$Operation start failed and is terminal untrusted evidence" }

        $stdoutTask = $process.StandardOutput.BaseStream.ReadAsync($stdoutBuffer, 0, $stdoutBuffer.Length)
        $stderrTask = $process.StandardError.BaseStream.ReadAsync($stderrBuffer, 0, $stderrBuffer.Length)
        $stdoutClosed = $false
        $stderrClosed = $false
        $trustedComplete = $false
        while (!$trustedComplete -and $null -eq $terminalReason) {
            if ($clock.ElapsedMilliseconds -ge $dockerBudgetMilliseconds) { $terminalReason = "$Operation exceeded the complete monotonic wall-clock budget"; break }
            if (!$stdoutClosed -and $stdoutTask.IsCompleted) {
                try { $count = $stdoutTask.GetAwaiter().GetResult() } catch { $terminalReason = "$Operation stdout read failed"; break }
                if ($count -eq 0) { $stdoutClosed = $true } else {
                    if (($stdout.Length + $count) -gt $dockerOutputLimit) { $terminalReason = "$Operation stdout exceeded the byte limit"; break }
                    $stdout.Write($stdoutBuffer, 0, $count)
                    $stdoutTask = $process.StandardOutput.BaseStream.ReadAsync($stdoutBuffer, 0, $stdoutBuffer.Length)
                }
            }
            if (!$stderrClosed -and $stderrTask.IsCompleted) {
                try { $count = $stderrTask.GetAwaiter().GetResult() } catch { $terminalReason = "$Operation stderr read failed"; break }
                if ($count -eq 0) { $stderrClosed = $true } else {
                    if (($stderr.Length + $count) -gt $dockerOutputLimit) { $terminalReason = "$Operation stderr exceeded the byte limit"; break }
                    $stderr.Write($stderrBuffer, 0, $count)
                    $stderrTask = $process.StandardError.BaseStream.ReadAsync($stderrBuffer, 0, $stderrBuffer.Length)
                }
            }
            try { $exited = $process.HasExited } catch { $terminalReason = "$Operation exit state is untrusted"; break }
            $trustedComplete = $exited -and $stdoutClosed -and $stderrClosed
            if (!$trustedComplete) { [Threading.Thread]::Sleep(1) }
        }
        if ($null -ne $terminalReason) {
            $cleanupClock = [Diagnostics.Stopwatch]::StartNew()
            try {
                if (!$process.HasExited) { $process.Kill($true) }
                $remaining = $cleanupGraceMilliseconds - [int]$cleanupClock.ElapsedMilliseconds
                $killConfirmed = if ($remaining -gt 0) { $process.WaitForExit($remaining) } else { $process.HasExited }
            } catch { throw "$terminalReason; process-tree termination is uncertain and terminal" }
            if (!$killConfirmed) { throw "$terminalReason; process-tree termination exceeded the fixed cleanup grace and is terminal" }
            while ($cleanupClock.ElapsedMilliseconds -lt $cleanupGraceMilliseconds -and (!$stdoutTask.IsCompleted -or !$stderrTask.IsCompleted)) { [Threading.Thread]::Sleep(1) }
            if (!$stdoutTask.IsCompleted -or !$stderrTask.IsCompleted) { throw "$terminalReason; redirected pipe closure is uncertain and terminal" }
            throw "$terminalReason; command evidence is terminal and cannot be retried or reclassified"
        }
        $utf8 = [Text.UTF8Encoding]::new($false, $true)
        try { $stdoutText = $utf8.GetString($stdout.ToArray()); $stderrText = $utf8.GetString($stderr.ToArray()) }
        catch { throw "$Operation returned invalid or incomplete UTF-8 and is terminal untrusted evidence" }
        [pscustomobject]@{
            ExitCode = $process.ExitCode; StdOut = $stdoutText; StdErr = $stderrText
            StdOutByteCount = [int]$stdout.Length; StdErrByteCount = [int]$stderr.Length
            ElapsedMilliseconds = [int64]$clock.ElapsedMilliseconds
        }
    } finally { $stdout.Dispose(); $stderr.Dispose(); $process.Dispose(); $clock.Stop() }
}

function Invoke-ClosedDockerOperation {
    param([Parameter(Mandatory = $true)][ValidateSet('Engine','Inventory','Inspect','Pull')][string]$Operation, [Parameter(Mandatory = $true)]$DockerIdentity)
    $arguments = switch ($Operation) {
        'Engine' { @('version','--format','{{json .Server}}') }
        'Inventory' { @('image','ls','--all','--no-trunc','--digests','--filter',"reference=$releaseImage",'--format','{{json .}}') }
        'Inspect' { @('image','inspect','--format','{{json .RepoDigests}}',$releaseImage) }
        'Pull' { @('pull',$releaseImage) }
    }
    Invoke-ClosedNativeProcess -ExecutablePath $DockerIdentity.Path -Arguments $arguments -Operation "Docker $Operation"
}

function Assert-ClosedCommandSuccess($Result, [string]$Operation) {
    if ($Result.ExitCode -ne 0 -or ![string]::IsNullOrEmpty($Result.StdErr)) { throw "$Operation failed and is terminal untrusted evidence" }
}

function Invoke-GenuineHostPrerequisites {
    $cacheRoot = [IO.Path]::GetFullPath((Microsoft.PowerShell.Management\Resolve-Path -LiteralPath $CacheDirectory).Path)
    $downloads = [IO.Path]::Combine($cacheRoot, 'downloads')
    $moduleCache = [IO.Path]::Combine($cacheRoot, 'gomodcache')
    foreach ($path in @($downloads, $moduleCache)) { if (!(Microsoft.PowerShell.Management\Test-Path -LiteralPath $path -PathType Container)) { throw "Exact acquisition cache path is absent before image admission: $path" } }
    $hostProof = & {
        . ([IO.Path]::Combine($PSScriptRoot, 'host-cache-canary.ps1'))
        $proof = $null
        foreach ($path in @($cacheRoot, $downloads, $moduleCache)) { $proof = Invoke-HostCacheCanary -CacheDirectory $path; if ($null -eq $proof) { throw 'Host cache prerequisite returned no proof' } }
        $proof
    }
    & ([IO.Path]::Combine($PSScriptRoot, 'test-crlf-shell-payloads.ps1')) -SourceRepository $SourceRepository -SourceRevision $SourceRevision -WorkingDirectory $WorkingDirectory -Phase HostOnly
    [pscustomobject]@{ CacheRoot=$cacheRoot; Downloads=$downloads; ModuleCache=$moduleCache; HostProof=$hostProof }
}

$prerequisites = Invoke-GenuineHostPrerequisites
$dockerIdentity = Get-ClosedDockerExecutable
$engine = Invoke-ClosedDockerOperation -Operation Engine -DockerIdentity $dockerIdentity
Assert-ClosedCommandSuccess $engine 'Docker engine reachability inspection'
[void](Assert-ReleaseEngineEvidence -Json $engine.StdOut)
$inventory = Invoke-ClosedDockerOperation -Operation Inventory -DockerIdentity $dockerIdentity
Assert-ClosedCommandSuccess $inventory 'Docker exact-image inventory inspection'
$state = Resolve-ReleaseImageListEvidence -Json $inventory.StdOut
$pulled = $false
if ($state -eq 'ConclusiveAbsent') {
    $pull = Invoke-ClosedDockerOperation -Operation Pull -DockerIdentity $dockerIdentity
    Assert-ClosedCommandSuccess $pull 'Pinned build image pull'
    $pulled = $true
}
$inspect = Invoke-ClosedDockerOperation -Operation Inspect -DockerIdentity $dockerIdentity
Assert-ClosedCommandSuccess $inspect 'Docker repository-digest inspection'
$repoDigests = Read-ReleaseRepoDigestsEvidence -Json $inspect.StdOut
[void](Assert-ReleaseImageIdentityEvidence -Image $releaseImage -RepoDigests $repoDigests)
$result = [pscustomobject]@{
    Image=$releaseImage; Pulled=$pulled; DockerExecutableSha256=$dockerIdentity.Sha256
    HostIdentityMode=$prerequisites.HostProof.IdentityMode; HostUID=$prerequisites.HostProof.HostUID; HostGID=$prerequisites.HostProof.HostGID
    ProvedCachePaths=@($prerequisites.CacheRoot,$prerequisites.Downloads,$prerequisites.ModuleCache)
}
Write-Output "Release image admission PASS image=$($result.Image) pulled=$($result.Pulled) host-identity=$($result.HostIdentityMode) docker-sha256=$($result.DockerExecutableSha256)"
