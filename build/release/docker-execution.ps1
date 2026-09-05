[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [ValidateSet('EngineInspection','ExactImageInventory','ApprovedImagePull','RepositoryDigestInspection','ContainerCacheProof','ContainerCrlfParse','DependencyAcquisition','ReleaseBuild','ReleasePackage')]
    [string]$Operation,
    [ValidatePattern('^[0-9a-f]{64}$')][string]$ExpectedDockerSHA256,
    [string]$CacheDirectory,
    [Nullable[int]]$HostUID,
    [Nullable[int]]$HostGID,
    [switch]$ReadOnlyCache,
    [ValidateSet('CacheCanary','Acquisition','Build','Package')][string]$PayloadKind,
    [string]$SourceRoot,
    [string]$RunnerGoArchive,
    [string]$EngineGoArchive,
    [string]$GitleaksArchive,
    [string]$ProductArchive,
    [string]$ProductBlobManifest,
    [string]$ProductPathManifest,
    [string]$ProductModeManifest,
    [string]$ToolingArchive,
    [string]$ToolingBlobManifest,
    [string]$ToolingPathManifest,
    [string]$ToolingModeManifest,
    [string]$RawOutputDirectory,
    [string]$LinuxStageDirectory,
    [string]$WindowsStageDirectory,
    [string]$DistributionDirectory,
    [ValidatePattern('^[0-9]+$')][string]$SourceDateEpoch,
    [ValidatePattern('^[0-9a-f]{40}$')][string]$ProductRevision,
    [ValidatePattern('^[0-9a-f]{40}$')][string]$ToolingRevision,
    [ValidatePattern('^[0-9a-f]{40}$')][string]$ToolingTree,
    [ValidatePattern('^[0-9a-f]{64}$')][string]$ProductArchiveSHA256,
    [ValidatePattern('^[0-9a-f]{64}$')][string]$ToolingArchiveSHA256,
    [ValidateRange(1, 1000000)][int]$ProductFileCount,
    [ValidateRange(1, 1000000)][int]$ToolingFileCount,
    [ValidatePattern('^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$')][string]$Created
)

if ($MyInvocation.InvocationName -eq '.' -or $MyInvocation.Line -match '^\s*\.\s') {
    throw 'The closed Docker execution entrypoint cannot be dot-sourced'
}

$ErrorActionPreference = 'Stop'
$dockerOutputLimit = 131072
$dockerBudgetMilliseconds = 15000
$cleanupGraceMilliseconds = 2000
$image = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
if ($Operation -cne 'EngineInspection' -and [string]::IsNullOrEmpty($ExpectedDockerSHA256)) { throw "$Operation requires the exact Docker digest admitted by EngineInspection" }

$nativeBoundarySource = @'
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Win32.SafeHandles;

public sealed class PscanBoundaryResult {
    public int ExitCode;
    public byte[] StdOut;
    public byte[] StdErr;
    public string Terminal;
    public long[] ObservedMembers;
    public bool ContainmentEmpty;
}

public static class PscanNativeBoundary {
    const uint CREATE_SUSPENDED = 0x00000004;
    const uint CREATE_NO_WINDOW = 0x08000000;
    const uint CREATE_UNICODE_ENVIRONMENT = 0x00000400;
    const uint STARTF_USESTDHANDLES = 0x00000100;
    const uint HANDLE_FLAG_INHERIT = 0x00000001;
    const uint WAIT_OBJECT_0 = 0;
    const uint WAIT_TIMEOUT = 258;
    const int JobObjectBasicProcessIdList = 3;
    const int JobObjectExtendedLimitInformation = 9;
    const uint JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x00002000;

    [StructLayout(LayoutKind.Sequential)] struct SECURITY_ATTRIBUTES { public int nLength; public IntPtr lpSecurityDescriptor; public int bInheritHandle; }
    [StructLayout(LayoutKind.Sequential, CharSet=CharSet.Unicode)] struct STARTUPINFO { public int cb; public string lpReserved; public string lpDesktop; public string lpTitle; public int dwX; public int dwY; public int dwXSize; public int dwYSize; public int dwXCountChars; public int dwYCountChars; public int dwFillAttribute; public int dwFlags; public short wShowWindow; public short cbReserved2; public IntPtr lpReserved2; public IntPtr hStdInput; public IntPtr hStdOutput; public IntPtr hStdError; }
    [StructLayout(LayoutKind.Sequential)] struct PROCESS_INFORMATION { public IntPtr hProcess; public IntPtr hThread; public int dwProcessId; public int dwThreadId; }
    [StructLayout(LayoutKind.Sequential)] struct IO_COUNTERS { public ulong ReadOperationCount, WriteOperationCount, OtherOperationCount, ReadTransferCount, WriteTransferCount, OtherTransferCount; }
    [StructLayout(LayoutKind.Sequential)] struct JOBOBJECT_BASIC_LIMIT_INFORMATION { public long PerProcessUserTimeLimit, PerJobUserTimeLimit; public uint LimitFlags; public UIntPtr MinimumWorkingSetSize, MaximumWorkingSetSize; public uint ActiveProcessLimit; public UIntPtr Affinity; public uint PriorityClass, SchedulingClass; }
    [StructLayout(LayoutKind.Sequential)] struct JOBOBJECT_EXTENDED_LIMIT_INFORMATION { public JOBOBJECT_BASIC_LIMIT_INFORMATION BasicLimitInformation; public IO_COUNTERS IoInfo; public UIntPtr ProcessMemoryLimit, JobMemoryLimit, PeakProcessMemoryUsed, PeakJobMemoryUsed; }

    [DllImport("kernel32.dll", SetLastError=true)] static extern bool CreatePipe(out IntPtr read, out IntPtr write, ref SECURITY_ATTRIBUTES sa, int size);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool SetHandleInformation(IntPtr handle, uint mask, uint flags);
    [DllImport("kernel32.dll", SetLastError=true, CharSet=CharSet.Unicode)] static extern IntPtr CreateJobObject(IntPtr attrs, string name);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool SetInformationJobObject(IntPtr job, int infoClass, IntPtr info, uint length);
    [DllImport("kernel32.dll", SetLastError=true, CharSet=CharSet.Unicode)] static extern bool CreateProcess(string app, StringBuilder command, IntPtr pa, IntPtr ta, bool inherit, uint flags, IntPtr environment, string cwd, ref STARTUPINFO si, out PROCESS_INFORMATION pi);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool AssignProcessToJobObject(IntPtr job, IntPtr process);
    [DllImport("kernel32.dll", SetLastError=true)] static extern uint ResumeThread(IntPtr thread);
    [DllImport("kernel32.dll", SetLastError=true)] static extern uint WaitForSingleObject(IntPtr handle, uint ms);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool GetExitCodeProcess(IntPtr process, out uint code);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool TerminateJobObject(IntPtr job, uint code);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool QueryInformationJobObject(IntPtr job, int infoClass, IntPtr info, uint length, out uint returned);
    [DllImport("kernel32.dll", SetLastError=true)] static extern bool CloseHandle(IntPtr handle);

    static Exception Win32(string operation) { return new System.ComponentModel.Win32Exception(Marshal.GetLastWin32Error(), operation); }
    static string Quote(string value) {
        if (value.Length > 0 && value.IndexOfAny(new[]{' ', '\t', '\n', '\v', '"'}) < 0) return value;
        var b = new StringBuilder("\""); int slashes = 0;
        foreach (char c in value) {
            if (c == '\\') { slashes++; continue; }
            if (c == '"') { b.Append('\\', slashes * 2 + 1); b.Append(c); slashes = 0; continue; }
            b.Append('\\', slashes); slashes = 0; b.Append(c);
        }
        b.Append('\\', slashes * 2); b.Append('"'); return b.ToString();
    }
    static IntPtr EnvironmentBlock(IDictionary<string,string> environment) {
        var keys = new List<string>(environment.Keys); keys.Sort(StringComparer.OrdinalIgnoreCase);
        var b = new StringBuilder(); foreach (string key in keys) b.Append(key).Append('=').Append(environment[key]).Append('\0'); b.Append('\0');
        return Marshal.StringToHGlobalUni(b.ToString());
    }
    static long[] Members(IntPtr job) {
        int capacity = 64;
        while (capacity <= 4096) {
            int bytes = 8 + IntPtr.Size * capacity; IntPtr p = Marshal.AllocHGlobal(bytes);
            try {
                for (int i=0;i<bytes;i++) Marshal.WriteByte(p,i,0);
                uint returned; if (!QueryInformationJobObject(job, JobObjectBasicProcessIdList, p, (uint)bytes, out returned)) {
                    int error=Marshal.GetLastWin32Error(); if (error==122) { capacity*=2; continue; } throw Win32("QueryInformationJobObject");
                }
                int count=Marshal.ReadInt32(p,4); if (count>capacity) { capacity*=2; continue; }
                var ids=new long[count]; for(int i=0;i<count;i++) ids[i]=IntPtr.Size==8?Marshal.ReadInt64(p,8+i*IntPtr.Size):(uint)Marshal.ReadInt32(p,8+i*IntPtr.Size); return ids;
            } finally { Marshal.FreeHGlobal(p); }
        }
        throw new InvalidOperationException("containment member limit exceeded");
    }
    static byte[] ReadBounded(IntPtr handle, int limit, Action<string> fail, string stream) {
        using (var fs=new FileStream(new SafeFileHandle(handle,true),FileAccess.Read,4096,false)) using (var ms=new MemoryStream()) {
            var buffer=new byte[4096];
            try { for (;;) { int n=fs.Read(buffer,0,buffer.Length); if(n==0) break; if(ms.Length+n>limit) { fail(stream+" overflow"); break; } ms.Write(buffer,0,n); } }
            catch(Exception e) { fail(stream+" read failure: "+e.GetType().Name); }
            return ms.ToArray();
        }
    }
    public static Task<byte[]> Capture(Stream stream, int limit) {
        return Task.Run(() => {
            using (stream) using (var ms=new MemoryStream()) {
                var buffer=new byte[4096]; for(;;){int n=stream.Read(buffer,0,buffer.Length);if(n==0)break;if(ms.Length+n>limit)throw new InvalidDataException("stream overflow");ms.Write(buffer,0,n);} return ms.ToArray();
            }
        });
    }
    public static PscanBoundaryResult RunWindows(string executable, string[] arguments, IDictionary<string,string> environment, string cwd, int streamLimit, int budgetMs, int cleanupMs) {
        IntPtr job=IntPtr.Zero, outRead=IntPtr.Zero, outWrite=IntPtr.Zero, errRead=IntPtr.Zero, errWrite=IntPtr.Zero, env=IntPtr.Zero; PROCESS_INFORMATION pi=new PROCESS_INFORMATION(); bool resumed=false;
        var observed=new HashSet<long>(); string terminal=null; object gate=new object(); var watch=Stopwatch.StartNew();
        Action<string> fail=(reason)=>{ lock(gate){ if(terminal==null) terminal=reason; } if(job!=IntPtr.Zero) TerminateJobObject(job,197); };
        try {
            job=CreateJobObject(IntPtr.Zero,null); if(job==IntPtr.Zero) throw Win32("CreateJobObject");
            var info=new JOBOBJECT_EXTENDED_LIMIT_INFORMATION(); info.BasicLimitInformation.LimitFlags=JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE; int infoSize=Marshal.SizeOf(info); IntPtr infoPtr=Marshal.AllocHGlobal(infoSize);
            try { Marshal.StructureToPtr(info,infoPtr,false); if(!SetInformationJobObject(job,JobObjectExtendedLimitInformation,infoPtr,(uint)infoSize)) throw Win32("SetInformationJobObject"); } finally { Marshal.FreeHGlobal(infoPtr); }
            var sa=new SECURITY_ATTRIBUTES{nLength=Marshal.SizeOf(typeof(SECURITY_ATTRIBUTES)),bInheritHandle=1};
            if(!CreatePipe(out outRead,out outWrite,ref sa,0)||!SetHandleInformation(outRead,HANDLE_FLAG_INHERIT,0)) throw Win32("stdout pipe");
            if(!CreatePipe(out errRead,out errWrite,ref sa,0)||!SetHandleInformation(errRead,HANDLE_FLAG_INHERIT,0)) throw Win32("stderr pipe");
            var si=new STARTUPINFO{cb=Marshal.SizeOf(typeof(STARTUPINFO)),dwFlags=(int)STARTF_USESTDHANDLES,hStdInput=IntPtr.Zero,hStdOutput=outWrite,hStdError=errWrite};
            var command=new StringBuilder(Quote(executable)); foreach(string arg in arguments) command.Append(' ').Append(Quote(arg)); env=EnvironmentBlock(environment);
            if(!CreateProcess(executable,command,IntPtr.Zero,IntPtr.Zero,true,CREATE_SUSPENDED|CREATE_NO_WINDOW|CREATE_UNICODE_ENVIRONMENT,env,cwd,ref si,out pi)) throw Win32("CreateProcess");
            if(!AssignProcessToJobObject(job,pi.hProcess)) throw Win32("AssignProcessToJobObject"); observed.Add(pi.dwProcessId);
            CloseHandle(outWrite);outWrite=IntPtr.Zero;CloseHandle(errWrite);errWrite=IntPtr.Zero;
            IntPtr stdoutHandle=outRead;outRead=IntPtr.Zero;IntPtr stderrHandle=errRead;errRead=IntPtr.Zero;
            var outTask=Task.Run(()=>ReadBounded(stdoutHandle,streamLimit,fail,"stdout"));
            var errTask=Task.Run(()=>ReadBounded(stderrHandle,streamLimit,fail,"stderr"));
            if(ResumeThread(pi.hThread)==0xffffffff) throw Win32("ResumeThread"); resumed=true;
            bool rootExited=false;
            while(watch.ElapsedMilliseconds<budgetMs){
                foreach(long id in Members(job)) observed.Add(id);
                if(terminal!=null) break;
                rootExited=WaitForSingleObject(pi.hProcess,0)==WAIT_OBJECT_0;
                if(rootExited&&outTask.IsCompleted&&errTask.IsCompleted&&Members(job).Length==0) break;
                Thread.Sleep(5);
            }
            if(!(rootExited&&outTask.IsCompleted&&errTask.IsCompleted&&Members(job).Length==0)) { if(terminal==null) terminal="timeout or incomplete lifecycle"; TerminateJobObject(job,198); }
            if(terminal!=null){
                var cleanup=Stopwatch.StartNew(); while(cleanup.ElapsedMilliseconds<cleanupMs){foreach(long id in Members(job))observed.Add(id);if(Members(job).Length==0&&outTask.IsCompleted&&errTask.IsCompleted)break;Thread.Sleep(5);} if(Members(job).Length!=0||!outTask.IsCompleted||!errTask.IsCompleted) terminal+="; cleanup uncertainty";
            }
            uint exit; if(!GetExitCodeProcess(pi.hProcess,out exit)) { fail("exit read failure"); exit=199; }
            var stdout=outTask.IsCompleted?outTask.Result:new byte[0];var stderr=errTask.IsCompleted?errTask.Result:new byte[0];bool empty=Members(job).Length==0;
            return new PscanBoundaryResult{ExitCode=(int)exit,StdOut=stdout,StdErr=stderr,Terminal=terminal,ObservedMembers=new List<long>(observed).ToArray(),ContainmentEmpty=empty};
        } finally {
            if(!resumed&&pi.hProcess!=IntPtr.Zero&&job!=IntPtr.Zero)TerminateJobObject(job,200);
            if(pi.hThread!=IntPtr.Zero)CloseHandle(pi.hThread);if(pi.hProcess!=IntPtr.Zero)CloseHandle(pi.hProcess);
            if(outRead!=IntPtr.Zero)CloseHandle(outRead);if(outWrite!=IntPtr.Zero)CloseHandle(outWrite);if(errRead!=IntPtr.Zero)CloseHandle(errRead);if(errWrite!=IntPtr.Zero)CloseHandle(errWrite);
            if(env!=IntPtr.Zero)Marshal.FreeHGlobal(env);if(job!=IntPtr.Zero)CloseHandle(job);
        }
    }
}
'@

if (!('PscanNativeBoundary' -as [type])) {
    Add-Type -TypeDefinition $nativeBoundarySource -Language CSharp
}

function Require-Path([string]$Name, [string]$Value, [bool]$Container = $false) {
    if ([string]::IsNullOrWhiteSpace($Value)) { throw "$Operation requires $Name" }
    $resolved = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $Value).Path)
    if ($Container -and !(Test-Path -LiteralPath $resolved -PathType Container)) { throw "$Name must be a directory" }
    if (!$Container -and !(Test-Path -LiteralPath $resolved -PathType Leaf)) { throw "$Name must be a regular file" }
    return $resolved
}

function New-PrivateDirectory([string]$Path) {
    [void](New-Item -ItemType Directory -Path $Path)
    if ($IsWindows) {
        $identity = [Security.Principal.WindowsIdentity]::GetCurrent().User
        $acl = [Security.AccessControl.DirectorySecurity]::new()
        $acl.SetOwner($identity)
        $acl.SetAccessRuleProtection($true, $false)
        $rule = [Security.AccessControl.FileSystemAccessRule]::new($identity, [Security.AccessControl.FileSystemRights]::FullControl, [Security.AccessControl.InheritanceFlags]'ContainerInherit,ObjectInherit', [Security.AccessControl.PropagationFlags]::None, [Security.AccessControl.AccessControlType]::Allow)
        $acl.AddAccessRule($rule)
        Microsoft.PowerShell.Security\Set-Acl -LiteralPath $Path -AclObject $acl
    } else {
        [IO.File]::SetUnixFileMode($Path, [IO.UnixFileMode]'UserRead,UserWrite,UserExecute')
    }
    if (Get-ChildItem -LiteralPath $Path -Force | Select-Object -First 1) { throw 'Private directory was not created empty' }
}

function Get-BoundDockerExecutable {
    if ($IsWindows) {
        $path = 'C:\Program Files\Docker\Docker\resources\bin\docker.exe'
        if (!(Test-Path -LiteralPath $path -PathType Leaf)) { throw 'Fixed Docker Desktop executable is unavailable' }
        return $path
    }
    if (!$IsLinux) { throw 'Docker execution is unsupported on this platform' }
    $path = '/usr/bin/docker'
    $item = Get-Item -LiteralPath $path -Force
    if (!$item -or $item.LinkType) { throw 'Fixed Linux Docker executable path is absent or linked' }
    return $path
}

function Get-LinuxSessionMembers([int]$Session, [hashtable]$Ledger) {
    $current = [Collections.Generic.List[object]]::new()
    foreach ($directory in Get-ChildItem -LiteralPath /proc -Directory -ErrorAction Stop) {
        if ($directory.Name -notmatch '^\d+$') { continue }
        try { $stat = [IO.File]::ReadAllText((Join-Path $directory.FullName 'stat')) }
        catch [IO.FileNotFoundException] { continue }
        catch [IO.DirectoryNotFoundException] { continue }
        $close = $stat.LastIndexOf(')'); if ($close -lt 1) { throw 'Malformed Linux process identity record' }
        $pid = [int]$stat.Substring(0, $stat.IndexOf(' ')); $fields = $stat.Substring($close + 2).Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
        if ($fields.Count -lt 20 -or [int]$fields[3] -ne $Session) { continue }
        $startTime = [uint64]$fields[19]; $identity = "${pid}:$startTime"
        $Ledger[$identity] = [pscustomobject]@{ PID=$pid; StartTime=$startTime }
        $current.Add($Ledger[$identity])
    }
    $current.ToArray()
}

function Test-LinuxMemberAlive($Member) {
    try { $stat = [IO.File]::ReadAllText("/proc/$($Member.PID)/stat") }
    catch [IO.FileNotFoundException] { return $false }
    catch [IO.DirectoryNotFoundException] { return $false }
    $close = $stat.LastIndexOf(')'); if ($close -lt 1) { return $false }
    $fields = $stat.Substring($close + 2).Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
    $fields.Count -ge 20 -and [uint64]$fields[19] -eq $Member.StartTime
}

function Invoke-LinuxSessionBoundary([string]$ExecutableReference, [string[]]$Arguments, [Collections.Generic.Dictionary[string,string]]$Environment, [string]$WorkingDirectory) {
    $setsid = '/usr/bin/setsid'
    $setsidItem = Get-Item -LiteralPath $setsid -Force
    if (!$setsidItem -or $setsidItem.LinkType -or ($setsidItem.UnixFileMode -band [IO.UnixFileMode]'GroupWrite,OtherWrite') -ne 0 -or (& /usr/bin/stat -c '%u' -- $setsid).Trim() -ne '0') { throw 'Fixed Linux session launcher is untrusted' }
    $start = [Diagnostics.ProcessStartInfo]::new()
    $start.FileName = $setsid
    $start.UseShellExecute = $false
    $start.RedirectStandardOutput = $true
    $start.RedirectStandardError = $true
    $start.WorkingDirectory = $WorkingDirectory
    $start.Environment.Clear()
    foreach ($pair in $Environment.GetEnumerator()) { $start.Environment.Add($pair.Key, $pair.Value) }
    $start.ArgumentList.Add('--wait')
    $start.ArgumentList.Add($ExecutableReference)
    foreach ($argument in $Arguments) { $start.ArgumentList.Add($argument) }
    $process = [Diagnostics.Process]::new(); $process.StartInfo = $start
    if (!$process.Start()) { throw 'Linux session process did not start' }
    $session = $process.Id; $ledger = @{}; $ledger["$session:root"]=[pscustomobject]@{PID=$session;StartTime=0}
    $stdoutTask = [PscanNativeBoundary]::Capture($process.StandardOutput.BaseStream, $dockerOutputLimit)
    $stderrTask = [PscanNativeBoundary]::Capture($process.StandardError.BaseStream, $dockerOutputLimit)
    $watch = [Diagnostics.Stopwatch]::StartNew(); $terminal = $null
    while ($watch.ElapsedMilliseconds -lt $dockerBudgetMilliseconds) {
        [void](Get-LinuxSessionMembers -Session $session -Ledger $ledger)
        if ($stdoutTask.IsFaulted -or $stderrTask.IsFaulted) { $terminal='stream overflow or read failure'; break }
        $aliveLedger = @($ledger.Values | Where-Object { $_.StartTime -ne 0 -and (Test-LinuxMemberAlive $_) })
        if ($process.HasExited -and $stdoutTask.IsCompleted -and $stderrTask.IsCompleted -and @(Get-LinuxSessionMembers -Session $session -Ledger $ledger).Count -eq 0 -and $aliveLedger.Count -eq 0) { break }
        Start-Sleep -Milliseconds 2
    }
    $aliveLedger = @($ledger.Values | Where-Object { $_.StartTime -ne 0 -and (Test-LinuxMemberAlive $_) })
    if (!$process.HasExited -or !$stdoutTask.IsCompleted -or !$stderrTask.IsCompleted -or @(Get-LinuxSessionMembers -Session $session -Ledger $ledger).Count -ne 0 -or $aliveLedger.Count -ne 0) { if (!$terminal) {$terminal='timeout or incomplete lifecycle'} }
    if ($terminal) {
        & /usr/bin/kill -KILL -- "-$session" 2>$null
        foreach ($member in $ledger.Values) { if ($member.StartTime -ne 0 -and (Test-LinuxMemberAlive $member)) { & /usr/bin/kill -KILL -- "$($member.PID)" 2>$null } }
        $cleanup = [Diagnostics.Stopwatch]::StartNew()
        while ($cleanup.ElapsedMilliseconds -lt $cleanupGraceMilliseconds) {
            $alive=@($ledger.Values | Where-Object { $_.StartTime -ne 0 -and (Test-LinuxMemberAlive $_) }); $current=@(Get-LinuxSessionMembers -Session $session -Ledger $ledger)
            if ($alive.Count -eq 0 -and $current.Count -eq 0 -and $stdoutTask.IsCompleted -and $stderrTask.IsCompleted) { break }
            Start-Sleep -Milliseconds 2
        }
        $remaining=@($ledger.Values | Where-Object { $_.StartTime -ne 0 -and (Test-LinuxMemberAlive $_) })
        if ($remaining.Count -ne 0 -or @(Get-LinuxSessionMembers -Session $session -Ledger $ledger).Count -ne 0 -or !$stdoutTask.IsCompleted -or !$stderrTask.IsCompleted) { $terminal += '; cleanup uncertainty' }
    }
    $stdout = if($stdoutTask.IsCompletedSuccessfully){$stdoutTask.Result}else{[byte[]]::new(0)}
    $stderr = if($stderrTask.IsCompletedSuccessfully){$stderrTask.Result}else{[byte[]]::new(0)}
    [pscustomobject]@{ExitCode=$(if($process.HasExited){$process.ExitCode}else{199});StdOut=$stdout;StdErr=$stderr;Terminal=$terminal;ObservedMembers=@($ledger.Keys);ContainmentEmpty=(!$terminal)}
}

function Invoke-BoundDocker([string[]]$Arguments) {
    $privateRoot = Join-Path ([IO.Path]::GetTempPath()) ('pscan-docker-boundary-' + [guid]::NewGuid().ToString('N'))
    $working = Join-Path $privateRoot 'work'; $config = Join-Path $privateRoot 'docker-config'; $temp = Join-Path $privateRoot 'temp'
    New-PrivateDirectory $privateRoot; New-PrivateDirectory $working; New-PrivateDirectory $config; New-PrivateDirectory $temp
    $dockerPath = Get-BoundDockerExecutable
    $identityStream = [IO.File]::Open($dockerPath, [IO.FileMode]::Open, [IO.FileAccess]::Read, [IO.FileShare]::Read)
    try {
        if ($IsWindows) {
            $signature = Microsoft.PowerShell.Security\Get-AuthenticodeSignature -LiteralPath $dockerPath
            if ($signature.Status -ne [Management.Automation.SignatureStatus]::Valid -or $signature.SignerCertificate.Subject -notmatch '(^|, )O=Docker Inc(,|$)') { throw 'Held Docker executable publisher identity is untrusted' }
        } else {
            $fd = $identityStream.SafeFileHandle.DangerousGetHandle().ToInt64(); $reference = "/proc/$PID/fd/$fd"
            $identity = (& /usr/bin/stat -Lc '%u:%a' -- $reference).Trim()
            if ($LASTEXITCODE -ne 0 -or $identity -notmatch '^0:[1357][0145][0145]$') { throw 'Held Linux Docker executable owner or mode is untrusted' }
        }
        $digest = (Get-FileHash -InputStream $identityStream -Algorithm SHA256).Hash.ToLowerInvariant(); $identityStream.Position = 0
        if (![string]::IsNullOrEmpty($ExpectedDockerSHA256) -and $digest -cne $ExpectedDockerSHA256) { throw 'Docker executable digest changed after admission' }
        $environment = [Collections.Generic.Dictionary[string,string]]::new([StringComparer]::Ordinal)
        $environment.Add('DOCKER_CONFIG', $config); $environment.Add('HOME', $working)
        if ($IsWindows) {
            $environment.Add('SystemRoot', $env:SystemRoot); $environment.Add('WINDIR', $env:WINDIR); $environment.Add('TEMP', $temp); $environment.Add('TMP', $temp); $environment.Add('DOCKER_HOST', 'npipe:////./pipe/dockerDesktopLinuxEngine')
            $result = [PscanNativeBoundary]::RunWindows($dockerPath, $Arguments, $environment, $working, $dockerOutputLimit, $dockerBudgetMilliseconds, $cleanupGraceMilliseconds)
        } else {
            $environment.Add('TMPDIR', $temp); $environment.Add('DOCKER_HOST', 'unix:///var/run/docker.sock'); $environment.Add('LANG', 'C.UTF-8')
            $result = Invoke-LinuxSessionBoundary -ExecutableReference $reference -Arguments $Arguments -Environment $environment -WorkingDirectory $working
        }
        $identityStream.Position = 0; $after = (Get-FileHash -InputStream $identityStream -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($after -cne $digest) { throw 'Docker executable identity changed during execution' }
        if ($result.Terminal -or !$result.ContainmentEmpty) { throw "Docker $Operation returned terminal untrusted evidence: $($result.Terminal)" }
        $utf8 = [Text.UTF8Encoding]::new($false, $true)
        try { $stdout = $utf8.GetString($result.StdOut); $stderr = $utf8.GetString($result.StdErr) } catch { throw "Docker $Operation returned invalid UTF-8" }
        if (Get-ChildItem -LiteralPath $config -Force | Select-Object -First 1) { throw 'Docker wrote unexpected configuration or credential state' }
        return [pscustomobject]@{ Operation=$Operation; DockerSHA256=$digest; ExitCode=$result.ExitCode; StdOut=$stdout; StdErr=$stderr; StdOutByteCount=$result.StdOut.Length; StdErrByteCount=$result.StdErr.Length; ContainmentMembersObserved=$result.ObservedMembers.Count; ContainmentEmpty=$true }
    } finally {
        $identityStream.Dispose()
        if (Test-Path -LiteralPath $privateRoot) { Remove-Item -LiteralPath $privateRoot -Recurse -Force }
        if (Test-Path -LiteralPath $privateRoot) { throw 'Private Docker boundary cleanup failed' }
    }
}

function ConvertTo-LF([string]$Payload) {
    $value = $Payload.Replace("`r`n", "`n").Replace("`r", "`n")
    if ($value.IndexOf([char]0) -ge 0 -or $value.IndexOf([char]13) -ge 0) { throw 'POSIX payload is not exact LF text' }
    return $value
}

function Get-CachePayload { ConvertTo-LF @'
umask 077
canary=/gomodcache/.pscan-cache-canary-$$
renamed=/gomodcache/.pscan-cache-canary-ready-$$
cleanup() { rm -f "$canary" "$renamed"; }
trap cleanup EXIT HUP INT TERM
test ! -e "$canary"
test ! -e "$renamed"
printf '%s\n' 'PSCAN-06-C2-CONTAINER-CACHE-CANARY' > "$canary"
mv "$canary" "$renamed"
test "$(cat "$renamed")" = 'PSCAN-06-C2-CONTAINER-CACHE-CANARY'
rm "$renamed"
test ! -e "$canary"
test ! -e "$renamed"
trap - EXIT HUP INT TERM
'@ }

function Get-AcquisitionPayload { ConvertTo-LF @'
echo "63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445  /input/runner-go.tar.gz" | sha256sum -c -
echo "675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685  /input/engine-go.tar.gz" | sha256sum -c -
echo "6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115  /input/gitleaks.tar.gz" | sha256sum -c -
mkdir -p /work/runner /work/engine /work/gitleaks
tar -xzf /input/runner-go.tar.gz -C /work/runner
tar -xzf /input/engine-go.tar.gz -C /work/engine
tar -xzf /input/gitleaks.tar.gz -C /work/gitleaks --strip-components=1
test "$(/work/runner/go/bin/go version)" = "go version go1.27.1 linux/amd64"
test "$(/work/engine/go/bin/go version)" = "go version go1.27.0 linux/amd64"
cd /src
env GOTOOLCHAIN=local GOFLAGS=-mod=readonly GOMODCACHE=/gomodcache GOCACHE=/work/runner-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/runner/go/bin/go mod download
env GOTOOLCHAIN=local GOFLAGS=-mod=readonly GOMODCACHE=/gomodcache GOCACHE=/work/runner-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/runner/go/bin/go list -m all >/work/runner-modules.txt
cd /work/gitleaks
env GOTOOLCHAIN=local GOFLAGS=-mod=readonly GOMODCACHE=/gomodcache GOCACHE=/work/engine-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/engine/go/bin/go mod download all
env GOTOOLCHAIN=local GOFLAGS=-mod=readonly GOMODCACHE=/gomodcache GOCACHE=/work/engine-cache GOPROXY=https://proxy.golang.org GOSUMDB=sum.golang.org /work/engine/go/bin/go list -m all >/work/engine-modules.txt
'@ }

function Get-BuildPayload { ConvertTo-LF @'
echo "63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445  /input/runner-go.tar.gz" | sha256sum -c -
echo "675c26c449cbb18fc24b74650de1eabbae6e16f64326fd85a283fb3b58280685  /input/engine-go.tar.gz" | sha256sum -c -
echo "6b2638a733b85619dc80bdf28e84e4fed7e526a761ab5c148fbf67695aea2115  /input/gitleaks.tar.gz" | sha256sum -c -
test "$(sha256sum /input/product-source.tar | cut -d ' ' -f 1)" = "$PSCAN_PRODUCT_ARCHIVE_SHA256"
test "$(sha256sum /input/release-tooling-source.tar | cut -d ' ' -f 1)" = "$PSCAN_TOOLING_ARCHIVE_SHA256"
mkdir -p /work/runner /work/engine /work/gitleaks /work/product /work/tooling /work/cache /work/tmp /out/product /out/tooling-materialized/contracts/release-manifest /out/tooling-materialized/docs/release
tar -xzf /input/runner-go.tar.gz -C /work/runner
tar -xzf /input/engine-go.tar.gz -C /work/engine
tar -xzf /input/gitleaks.tar.gz -C /work/gitleaks --strip-components=1
tar -xf /input/product-source.tar -C /work/product
tar -xf /input/release-tooling-source.tar -C /work/tooling
(cd /work/product && sha256sum --quiet -c /input/product-source.blobs.sha256 && find . -type f -printf '%P\n' | LC_ALL=C sort > /work/product.paths && cmp /input/product-source.paths /work/product.paths)
(cd /work/tooling && sha256sum --quiet -c /input/release-tooling-source.blobs.sha256 && find . -type f -printf '%P\n' | LC_ALL=C sort > /work/tooling.paths && cmp /input/release-tooling-source.paths /work/tooling.paths)
while IFS="$(printf '\t')" read -r expected_mode path; do test "$(stat -c '%a' "/work/product/$path")" = "${expected_mode#100}"; done < /input/product-source.modes
while IFS="$(printf '\t')" read -r expected_mode path; do test "$(stat -c '%a' "/work/tooling/$path")" = "${expected_mode#100}"; done < /input/release-tooling-source.modes
test "$(wc -l < /input/product-source.paths)" -eq "$PSCAN_PRODUCT_FILE_COUNT"
test "$(wc -l < /input/release-tooling-source.paths)" -eq "$PSCAN_TOOLING_FILE_COUNT"
test "$(/work/runner/go/bin/go version)" = "go version go1.27.1 linux/amd64"
test "$(/work/engine/go/bin/go version)" = "go version go1.27.0 linux/amd64"
test "$(sha256sum /work/gitleaks/go.mod | cut -d ' ' -f 1)" = "607c140abf2a872e70423972d4dfc7fa658ebe10365d0ea995269ed292add7a3"
test "$(sha256sum /work/gitleaks/config/gitleaks.toml | cut -d ' ' -f 1)" = "e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf"
test "$(sha256sum /work/gitleaks/LICENSE | cut -d ' ' -f 1)" = "e3884b252b3bfc045e55be43a34d1e80da070bc6f804ac95bf4660e97d62ebc6"
test "$(sha256sum /work/gitleaks/.goreleaser.yml | cut -d ' ' -f 1)" = "1e6a76e13378b4ad215b423411ee6b1732f04c82224fbed4dcdb20b99eab9717"
export GOTOOLCHAIN=local GOMODCACHE=/gomodcache GOPROXY=off GOSUMDB=off CGO_ENABLED=0
export GOCACHE=/work/cache GOTMPDIR=/work/tmp TMPDIR=/work/tmp
cd /work/gitleaks
GOOS=linux GOARCH=amd64 /work/engine/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid= -X=github.com/zricethezav/gitleaks/v8/version.Version=8.30.1' -o /out/gitleaks-linux-amd64 .
GOOS=windows GOARCH=amd64 /work/engine/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid= -X=github.com/zricethezav/gitleaks/v8/version.Version=8.30.1' -o /out/gitleaks-windows-amd64.exe .
test "$(sha256sum /out/gitleaks-linux-amd64 | cut -d ' ' -f 1)" = "657ddddfb98e21052fb1a60d5d4e7d7534897347cb7df0031f13258a3f800586"
test "$(sha256sum /out/gitleaks-windows-amd64.exe | cut -d ' ' -f 1)" = "b2094b3534ce0abf9c74a4b251153f5a23ebb4e74d5ae4f6d6ceeb428aaf0178"
cd /work/product
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/scanner-runner-linux-amd64 ./cmd/scanner-runner
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/scanner-runner-windows-amd64.exe ./cmd/scanner-runner
cp -R rules contracts licenses /out/product/
cp LICENSE THIRD_PARTY_NOTICES.md /out/product/
cd /work/tooling
verifier_ldflags="-s -w -buildid= -X=main.releaseToolingCommit=$PSCAN_TOOLING_REVISION -X=main.releaseToolingTree=$PSCAN_TOOLING_TREE"
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags "$verifier_ldflags" -o /out/scanner-release-verifier-linux-amd64 ./build/release/cmd/release-verifier
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags "$verifier_ldflags" -o /out/scanner-release-verifier-windows-amd64.exe ./build/release/cmd/release-verifier
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/packager-linux-amd64 ./build/release/cmd/packager
GOOS=linux GOARCH=amd64 /work/runner/go/bin/go build -mod=readonly -trimpath -buildvcs=false -ldflags '-s -w -buildid=' -o /out/sbom-linux-amd64 ./build/release/cmd/sbom
/work/runner/go/bin/go list -mod=readonly -m -json all > /out/modules.json
PSCAN_GITLEAKS_BINARY=/out/gitleaks-linux-amd64 PSCAN_GITLEAKS_CONFIG=/work/product/rules/generic/gitleaks-v8.30.1.toml PSCAN_GITLEAKS_IGNORE=/work/product/rules/generic/gitleaks-ignore-empty-v1.txt /work/runner/go/bin/go test -p=1 -count=1 -run '^TestPinnedRuleAndCoverageIntegrityBindings$' ./tests/acceptance/gitleaks
PSCAN_GITLEAKS_BINARY=/out/gitleaks-linux-amd64 PSCAN_GITLEAKS_CONFIG=/work/product/rules/generic/gitleaks-v8.30.1.toml PSCAN_GITLEAKS_IGNORE=/work/product/rules/generic/gitleaks-ignore-empty-v1.txt /work/runner/go/bin/go test -p=1 -count=1 ./...
/work/runner/go/bin/go vet -p=1 ./...
GOOS=windows GOARCH=amd64 /work/runner/go/bin/go test -p=1 -exec /bin/true ./...
/out/sbom-linux-amd64 -input /out/modules.json -output /out/sbom.spdx.json -revision "$PSCAN_PRODUCT_REVISION" -created "$PSCAN_CREATED"
cp /work/tooling/contracts/release-manifest/schema-2.0.json /out/tooling-materialized/contracts/release-manifest/schema-2.0.json
cp /work/tooling/contracts/release-manifest/schema-2.1.json /out/tooling-materialized/contracts/release-manifest/schema-2.1.json
cp /work/tooling/docs/release/OFFLINE-VERIFICATION-RUNBOOK.md /out/tooling-materialized/docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
cp /work/tooling/docs/release/SCANNER-IO-REFERENCE.md /out/tooling-materialized/docs/release/SCANNER-IO-REFERENCE.md
'@ }

function Get-PackagePayload([string]$Epoch) {
    ConvertTo-LF "cp /tools/packager-linux-amd64 /tmp/packager; chmod 0755 /tmp/packager; /tmp/packager -root /input/linux -output /dist/project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz -format tar.gz -epoch $Epoch; /tmp/packager -root /input/windows -output /dist/project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip -format zip -epoch $Epoch"
}

$arguments = switch ($Operation) {
    'EngineInspection' { @('version','--format','{{json .Server}}') }
    'ExactImageInventory' { @('image','ls','--all','--no-trunc','--digests','--filter',"reference=$image",'--format','{{json .}}') }
    'ApprovedImagePull' { @('pull',$image) }
    'RepositoryDigestInspection' { @('image','inspect','--format','{{json .RepoDigests}}',$image) }
    'ContainerCacheProof' {
        $cache = Require-Path 'CacheDirectory' $CacheDirectory $true
        $tmpfs='/work:rw,noexec,nosuid,nodev,size=16m,mode=0700'; $a=@('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1')
        if($HostUID.HasValue -or $HostGID.HasValue){if(!$HostUID.HasValue -or !$HostGID.HasValue -or $HostUID.Value-lt 0 -or $HostGID.Value-lt 0){throw 'Cache proof requires a complete non-negative UID/GID pair'};$a+=@('--user',"$($HostUID.Value):$($HostGID.Value)");$tmpfs+=",uid=$($HostUID.Value),gid=$($HostGID.Value)"}
        $mount="type=bind,src=$cache,dst=/gomodcache";if($ReadOnlyCache){$mount+=',readonly'}
        $a+@('--tmpfs',$tmpfs,'--mount',$mount,'--workdir','/work',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work','/bin/sh','-ceu',(Get-CachePayload))
    }
    'ContainerCrlfParse' {
        if([string]::IsNullOrEmpty($PayloadKind)){throw 'ContainerCrlfParse requires PayloadKind'}
        $payload=switch($PayloadKind){'CacheCanary'{Get-CachePayload};'Acquisition'{Get-AcquisitionPayload};'Build'{Get-BuildPayload};'Package'{Get-PackagePayload '0'}}
        @('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1','--tmpfs','/work:rw,noexec,nosuid,nodev,size=16m,mode=0700',$image,'/bin/sh','-n','-c',$payload)
    }
    'DependencyAcquisition' {
        $root=Require-Path 'SourceRoot' $SourceRoot $true;$runner=Require-Path 'RunnerGoArchive' $RunnerGoArchive;$engine=Require-Path 'EngineGoArchive' $EngineGoArchive;$gitleaks=Require-Path 'GitleaksArchive' $GitleaksArchive;$cache=Require-Path 'CacheDirectory' $CacheDirectory $true
        $a=@('run','--rm','--pull=never','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','256','--memory','4g','--memory-swap','4g','--cpus','2','--tmpfs','/work:rw,exec,nosuid,nodev,size=1g','--mount',"type=bind,src=$root,dst=/src,readonly",'--mount',"type=bind,src=$runner,dst=/input/runner-go.tar.gz,readonly",'--mount',"type=bind,src=$engine,dst=/input/engine-go.tar.gz,readonly",'--mount',"type=bind,src=$gitleaks,dst=/input/gitleaks.tar.gz,readonly",'--mount',"type=bind,src=$cache,dst=/gomodcache",'--workdir','/src',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work','/bin/sh','-ceu',(Get-AcquisitionPayload))
        if($HostUID.HasValue -or $HostGID.HasValue){if(!$HostUID.HasValue-or!$HostGID.HasValue){throw 'Acquisition requires a complete UID/GID pair'};$index=$a.IndexOf('--tmpfs');$a=$a[0..($index-1)]+@('--user',"$($HostUID.Value):$($HostGID.Value)")+$a[$index..($a.Count-1)];$a[$a.IndexOf('/work:rw,exec,nosuid,nodev,size=1g')]="/work:rw,exec,nosuid,nodev,size=1g,mode=0700,uid=$($HostUID.Value),gid=$($HostGID.Value)"}
        $a
    }
    'ReleaseBuild' {
        foreach($required in @('SourceDateEpoch','ProductRevision','ToolingRevision','ToolingTree','ProductArchiveSHA256','ToolingArchiveSHA256','Created')){if([string]::IsNullOrEmpty((Get-Variable -Name $required -ValueOnly))){throw "ReleaseBuild requires $required"}}
        $product=Require-Path 'ProductArchive' $ProductArchive;$productBlobs=Require-Path 'ProductBlobManifest' $ProductBlobManifest;$productPaths=Require-Path 'ProductPathManifest' $ProductPathManifest;$productModes=Require-Path 'ProductModeManifest' $ProductModeManifest
        $tooling=Require-Path 'ToolingArchive' $ToolingArchive;$toolingBlobs=Require-Path 'ToolingBlobManifest' $ToolingBlobManifest;$toolingPaths=Require-Path 'ToolingPathManifest' $ToolingPathManifest;$toolingModes=Require-Path 'ToolingModeManifest' $ToolingModeManifest
        $runner=Require-Path 'RunnerGoArchive' $RunnerGoArchive;$engine=Require-Path 'EngineGoArchive' $EngineGoArchive;$gitleaks=Require-Path 'GitleaksArchive' $GitleaksArchive;$cache=Require-Path 'CacheDirectory' $CacheDirectory $true;$raw=Require-Path 'RawOutputDirectory' $RawOutputDirectory $true
        @('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','512','--memory','8g','--memory-swap','8g','--cpus','2','--tmpfs','/work:rw,exec,nosuid,nodev,size=4g','--mount',"type=bind,src=$product,dst=/input/product-source.tar,readonly",'--mount',"type=bind,src=$productBlobs,dst=/input/product-source.blobs.sha256,readonly",'--mount',"type=bind,src=$productPaths,dst=/input/product-source.paths,readonly",'--mount',"type=bind,src=$productModes,dst=/input/product-source.modes,readonly",'--mount',"type=bind,src=$tooling,dst=/input/release-tooling-source.tar,readonly",'--mount',"type=bind,src=$toolingBlobs,dst=/input/release-tooling-source.blobs.sha256,readonly",'--mount',"type=bind,src=$toolingPaths,dst=/input/release-tooling-source.paths,readonly",'--mount',"type=bind,src=$toolingModes,dst=/input/release-tooling-source.modes,readonly",'--mount',"type=bind,src=$runner,dst=/input/runner-go.tar.gz,readonly",'--mount',"type=bind,src=$engine,dst=/input/engine-go.tar.gz,readonly",'--mount',"type=bind,src=$gitleaks,dst=/input/gitleaks.tar.gz,readonly",'--mount',"type=bind,src=$cache,dst=/gomodcache,readonly",'--mount',"type=bind,src=$raw,dst=/out",'--workdir','/work/tooling',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work',"SOURCE_DATE_EPOCH=$SourceDateEpoch","PSCAN_PRODUCT_REVISION=$ProductRevision","PSCAN_TOOLING_REVISION=$ToolingRevision","PSCAN_TOOLING_TREE=$ToolingTree","PSCAN_CREATED=$Created","PSCAN_PRODUCT_ARCHIVE_SHA256=$ProductArchiveSHA256","PSCAN_TOOLING_ARCHIVE_SHA256=$ToolingArchiveSHA256","PSCAN_PRODUCT_FILE_COUNT=$ProductFileCount","PSCAN_TOOLING_FILE_COUNT=$ToolingFileCount",'/bin/sh','-ceu',(Get-BuildPayload))
    }
    'ReleasePackage' {
        if([string]::IsNullOrEmpty($SourceDateEpoch)){throw 'ReleasePackage requires SourceDateEpoch'}
        $raw=Require-Path 'RawOutputDirectory' $RawOutputDirectory $true;$linux=Require-Path 'LinuxStageDirectory' $LinuxStageDirectory $true;$windows=Require-Path 'WindowsStageDirectory' $WindowsStageDirectory $true;$dist=Require-Path 'DistributionDirectory' $DistributionDirectory $true
        @('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','64','--memory','1g','--memory-swap','1g','--cpus','1','--tmpfs','/tmp:rw,exec,nosuid,nodev,size=64m','--mount',"type=bind,src=$raw,dst=/tools,readonly",'--mount',"type=bind,src=$linux,dst=/input/linux,readonly",'--mount',"type=bind,src=$windows,dst=/input/windows,readonly",'--mount',"type=bind,src=$dist,dst=/dist",$image,'/bin/sh','-ceu',(Get-PackagePayload $SourceDateEpoch))
    }
}

$result = Invoke-BoundDocker -Arguments $arguments
$result | ConvertTo-Json -Compress
