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
. (Join-Path $PSScriptRoot 'execution-profile.ps1')
$dockerOutputLimit = 131072
$dockerBudgetMilliseconds = Get-ReleaseOperationBudget $Operation
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

public sealed class PscanOperationBudget {
    readonly Stopwatch watch = Stopwatch.StartNew();
    readonly int operationMs, cleanupMs, streamLimit;
    long cleanupStart = -1;
    int stdoutUsed, stderrUsed;
    public PscanOperationBudget(int operation, int cleanup, int limit) {
        if (operation < 1 || operation > 900000 || cleanup < 0 || cleanup > 2000 || limit < 1 || limit > 131072) throw new ArgumentOutOfRangeException();
        operationMs=operation; cleanupMs=cleanup; streamLimit=limit;
    }
    public void BeginCleanup() { Interlocked.CompareExchange(ref cleanupStart, watch.ElapsedMilliseconds, -1); }
    public bool CleanupStarted { get { return Interlocked.Read(ref cleanupStart)>=0; } }
    public int Remaining {
        get {
            long start=Interlocked.Read(ref cleanupStart);
            long left=start<0 ? operationMs-watch.ElapsedMilliseconds : cleanupMs-(watch.ElapsedMilliseconds-start);
            return (int)Math.Max(0,left);
        }
    }
    public int StdOutUsed { get { return Volatile.Read(ref stdoutUsed); } }
    public int StdErrUsed { get { return Volatile.Read(ref stderrUsed); } }
    public bool Consume(int bytes, bool stderr) {
        if(bytes<0)throw new ArgumentOutOfRangeException();
        int used=stderr ? Interlocked.Add(ref stderrUsed,bytes) : Interlocked.Add(ref stdoutUsed,bytes);
        return used<=streamLimit;
    }
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

    [StructLayout(LayoutKind.Sequential)] struct LinuxPollFd { public int fd; public short events; public short revents; }
    [DllImport("libc", EntryPoint="syscall", SetLastError=true)] static extern long PidOpen(long number, int pid, uint flags);
    [DllImport("libc", EntryPoint="syscall", SetLastError=true)] static extern long PidSignal(long number, SafeFileHandle fd, int signal, IntPtr info, uint flags);
    [DllImport("libc", EntryPoint="poll", SetLastError=true)] static extern int PidPoll(ref LinuxPollFd fd, UIntPtr count, int timeout);

    static void LinuxHost() {
        // The release contract's Linux native host and assets are amd64. Do not
        // guess syscall numbers on an unadmitted architecture or emulate pidfds.
        if (!OperatingSystem.IsLinux() || RuntimeInformation.ProcessArchitecture != Architecture.X64)
            throw new PlatformNotSupportedException("Linux pidfd boundary requires the supported Linux amd64 host");
    }
    public static SafeFileHandle LinuxPin(int pid) {
        LinuxHost();
        if (pid <= 1) throw new ArgumentOutOfRangeException(nameof(pid));
        long fd = PidOpen(434, pid, 0); // Linux x86-64 __NR_pidfd_open
        if (fd < 0) throw Win32("pidfd_open: unsupported or unavailable process identity");
        return new SafeFileHandle(new IntPtr(fd), true);
    }
    public static bool LinuxSignal(SafeFileHandle fd, int signal) {
        LinuxHost();
        if (fd == null || fd.IsInvalid || fd.IsClosed) throw new ArgumentException("Invalid pidfd");
        if (signal != 9 && signal != 18 && signal != 19) throw new ArgumentOutOfRangeException(nameof(signal));
        long result = PidSignal(424, fd, signal, IntPtr.Zero, 0);
        return LinuxSignalResult(result, result < 0 ? Marshal.GetLastWin32Error() : 0);
    }
    static bool LinuxSignalResult(long result, int error) {
        if (result == 0) return true;
        if (result == -1 && error == 3) return false; // ESRCH: never signal a replacement PID.
        throw new System.ComponentModel.Win32Exception(error, "pidfd_send_signal failed");
    }
    static bool LinuxPollResult(int count, short events, int error) {
        if (count < 0) throw new System.ComponentModel.Win32Exception(error, "pidfd poll failed");
        if (count > 1 || (count == 0 && events != 0) || (count == 1 && events == 0) || (events & ~(1 | 16)) != 0)
            throw new IOException("Untrusted pidfd poll status");
        return (events & (1 | 16)) != 0;
    }
    public static bool LinuxExited(SafeFileHandle fd) {
        LinuxHost();
        if (fd == null || fd.IsInvalid || fd.IsClosed) throw new ArgumentException("Invalid pidfd");
        bool held = false;
        try {
            fd.DangerousAddRef(ref held);
            var value = new LinuxPollFd { fd = fd.DangerousGetHandle().ToInt32(), events = 1 };
            int count = PidPoll(ref value, new UIntPtr(1), 0);
            return LinuxPollResult(count, value.revents, count < 0 ? Marshal.GetLastWin32Error() : 0);
        } finally { if (held) fd.DangerousRelease(); }
    }

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
    static byte[] ReadBounded(IntPtr handle, int limit, Action<string> fail, string stream, PscanOperationBudget budget) {
        using (var fs=new FileStream(new SafeFileHandle(handle,true),FileAccess.Read,4096,false)) using (var ms=new MemoryStream()) {
            var buffer=new byte[4096];
            try { for (;;) { int n=fs.Read(buffer,0,buffer.Length); if(n==0) break; if(ms.Length+n>limit || !budget.Consume(n,stream=="stderr")) { fail(stream+" overflow"); break; } ms.Write(buffer,0,n); } }
            catch(Exception e) { fail(stream+" read failure: "+e.GetType().Name); }
            return ms.ToArray();
        }
    }
    public static Task<byte[]> Capture(Stream stream, int limit) { return CaptureShared(stream,limit,null,false); }
    public static Task<byte[]> CaptureShared(Stream stream, int limit, PscanOperationBudget budget, bool stderr) {
        return Task.Run(() => {
            using (stream) using (var ms=new MemoryStream()) {
                var buffer=new byte[4096]; for(;;){int n=stream.Read(buffer,0,buffer.Length);if(n==0)break;if(ms.Length+n>limit || (budget!=null && !budget.Consume(n,stderr)))throw new InvalidDataException("stream overflow");ms.Write(buffer,0,n);} return ms.ToArray();
            }
        });
    }
    public static PscanBoundaryResult RunWindows(string executable, string[] arguments, IDictionary<string,string> environment, string cwd, int streamLimit, int budgetMs, int cleanupMs) {
        return RunWindowsShared(executable,arguments,environment,cwd,streamLimit,new PscanOperationBudget(budgetMs,cleanupMs,streamLimit));
    }
    public static PscanBoundaryResult RunWindowsShared(string executable, string[] arguments, IDictionary<string,string> environment, string cwd, int streamLimit, PscanOperationBudget budget) {
        if(budget.Remaining<=0)throw new TimeoutException("Shared Docker operation deadline exhausted before process creation");
        IntPtr job=IntPtr.Zero, outRead=IntPtr.Zero, outWrite=IntPtr.Zero, errRead=IntPtr.Zero, errWrite=IntPtr.Zero, env=IntPtr.Zero; PROCESS_INFORMATION pi=new PROCESS_INFORMATION(); bool resumed=false;
        var observed=new HashSet<long>(); string terminal=null; object gate=new object(); var watch=Stopwatch.StartNew();
        Action<string> fail=(reason)=>{ lock(gate){ if(terminal==null) { terminal=reason; budget.BeginCleanup(); } } if(job!=IntPtr.Zero) TerminateJobObject(job,197); };
        try {
            job=CreateJobObject(IntPtr.Zero,null); if(job==IntPtr.Zero) throw Win32("CreateJobObject");
            var info=new JOBOBJECT_EXTENDED_LIMIT_INFORMATION(); info.BasicLimitInformation.LimitFlags=JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE; int infoSize=Marshal.SizeOf(info); IntPtr infoPtr=Marshal.AllocHGlobal(infoSize);
            try { Marshal.StructureToPtr(info,infoPtr,false); if(!SetInformationJobObject(job,JobObjectExtendedLimitInformation,infoPtr,(uint)infoSize)) throw Win32("SetInformationJobObject"); } finally { Marshal.FreeHGlobal(infoPtr); }
            var sa=new SECURITY_ATTRIBUTES{nLength=Marshal.SizeOf(typeof(SECURITY_ATTRIBUTES)),bInheritHandle=1};
            if(!CreatePipe(out outRead,out outWrite,ref sa,0)||!SetHandleInformation(outRead,HANDLE_FLAG_INHERIT,0)) throw Win32("stdout pipe");
            if(!CreatePipe(out errRead,out errWrite,ref sa,0)||!SetHandleInformation(errRead,HANDLE_FLAG_INHERIT,0)) throw Win32("stderr pipe");
            var si=new STARTUPINFO{cb=Marshal.SizeOf(typeof(STARTUPINFO)),dwFlags=(int)STARTF_USESTDHANDLES,hStdInput=IntPtr.Zero,hStdOutput=outWrite,hStdError=errWrite};
            var command=new StringBuilder(Quote(executable)); foreach(string arg in arguments) command.Append(' ').Append(Quote(arg)); env=EnvironmentBlock(environment);
            if(budget.Remaining<=0)throw new TimeoutException("Shared deadline exhausted before native creation");
            if(!CreateProcess(executable,command,IntPtr.Zero,IntPtr.Zero,true,CREATE_SUSPENDED|CREATE_NO_WINDOW|CREATE_UNICODE_ENVIRONMENT,env,cwd,ref si,out pi)) throw Win32("CreateProcess");
            if(!AssignProcessToJobObject(job,pi.hProcess)) throw Win32("AssignProcessToJobObject"); observed.Add(pi.dwProcessId);
            CloseHandle(outWrite);outWrite=IntPtr.Zero;CloseHandle(errWrite);errWrite=IntPtr.Zero;
            IntPtr stdoutHandle=outRead;outRead=IntPtr.Zero;IntPtr stderrHandle=errRead;errRead=IntPtr.Zero;
            var outTask=Task.Run(()=>ReadBounded(stdoutHandle,streamLimit,fail,"stdout",budget));
            var errTask=Task.Run(()=>ReadBounded(stderrHandle,streamLimit,fail,"stderr",budget));
            if(ResumeThread(pi.hThread)==0xffffffff) throw Win32("ResumeThread"); resumed=true;
            bool rootExited=false;
            while(budget.Remaining>0){
                foreach(long id in Members(job)) observed.Add(id);
                if(terminal!=null) break;
                rootExited=WaitForSingleObject(pi.hProcess,0)==WAIT_OBJECT_0;
                if(rootExited&&outTask.IsCompleted&&errTask.IsCompleted&&Members(job).Length==0) break;
                Thread.Sleep(5);
            }
            if(!(rootExited&&outTask.IsCompleted&&errTask.IsCompleted&&Members(job).Length==0)) { if(terminal==null) terminal="timeout or incomplete lifecycle"; budget.BeginCleanup(); TerminateJobObject(job,198); }
            if(terminal!=null){
                budget.BeginCleanup(); while(budget.Remaining>0){foreach(long id in Members(job))observed.Add(id);if(Members(job).Length==0&&outTask.IsCompleted&&errTask.IsCompleted)break;Thread.Sleep(5);} if(Members(job).Length!=0||!outTask.IsCompleted||!errTask.IsCompleted) terminal+="; cleanup uncertainty";
            }
            uint exit; if(!GetExitCodeProcess(pi.hProcess,out exit)) { fail("exit read failure"); exit=199; }
            var stdout=outTask.IsCompleted?outTask.Result:new byte[0];var stderr=errTask.IsCompleted?errTask.Result:new byte[0];bool empty=Members(job).Length==0;
            return new PscanBoundaryResult{ExitCode=(int)exit,StdOut=stdout,StdErr=stderr,Terminal=terminal,ObservedMembers=new List<long>(observed).ToArray(),ContainmentEmpty=empty};
        } catch { budget.BeginCleanup(); throw; } finally {
            if(!resumed&&pi.hProcess!=IntPtr.Zero&&job!=IntPtr.Zero)TerminateJobObject(job,200);
            if(pi.hThread!=IntPtr.Zero)CloseHandle(pi.hThread);if(pi.hProcess!=IntPtr.Zero)CloseHandle(pi.hProcess);
            if(outRead!=IntPtr.Zero)CloseHandle(outRead);if(outWrite!=IntPtr.Zero)CloseHandle(outWrite);if(errRead!=IntPtr.Zero)CloseHandle(errRead);if(errWrite!=IntPtr.Zero)CloseHandle(errWrite);
            if(env!=IntPtr.Zero)Marshal.FreeHGlobal(env);if(job!=IntPtr.Zero)CloseHandle(job);
        }
    }
}
'@

if ($null -ne ('PscanNativeBoundary' -as [type])) {
    throw 'Pre-existing PscanNativeBoundary type is terminal ambient state'
}
foreach($nativeTypeName in @('PscanOperationBudget','PscanBoundaryResult')){
    if($null-ne($nativeTypeName-as[type])){throw "Pre-existing $nativeTypeName type is terminal ambient state"}
}
$compiledNativeBoundaryTypes = @(Add-Type -TypeDefinition $nativeBoundarySource -Language CSharp -PassThru)
$compiledBoundary = @($compiledNativeBoundaryTypes | Where-Object { $_.IsPublic -and $_.FullName -ceq 'PscanNativeBoundary' })
$compiledResult = @($compiledNativeBoundaryTypes | Where-Object { $_.IsPublic -and $_.FullName -ceq 'PscanBoundaryResult' })
$compiledBudget = @($compiledNativeBoundaryTypes | Where-Object { $_.IsPublic -and $_.FullName -ceq 'PscanOperationBudget' })
if($compiledBudget.Count-ne 1-or![object]::ReferenceEquals($compiledBudget[0],('PscanOperationBudget'-as[type]))-or
    ![object]::ReferenceEquals($compiledBudget[0].Assembly,$compiledBoundary[0].Assembly)-or![object]::ReferenceEquals($compiledResult[0],('PscanBoundaryResult'-as[type]))){throw 'Compiled native budget/result identity is unexpected'}
if ($compiledBoundary.Count -ne 1 -or $compiledResult.Count -ne 1 -or
    ![object]::ReferenceEquals($compiledBoundary[0].Assembly, $compiledResult[0].Assembly) -or
    ![object]::ReferenceEquals($compiledBoundary[0], ('PscanNativeBoundary' -as [type])) -or
    ![object]::ReferenceEquals($compiledResult[0], ('PscanBoundaryResult' -as [type]))) {
    throw 'Compiled native boundary type identity is unexpected'
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
        $linuxProcessIdentifier = [int]$stat.Substring(0, $stat.IndexOf(' ')); $fields = $stat.Substring($close + 2).Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
        if ($fields.Count -lt 20 -or [int]$fields[3] -ne $Session) { continue }
        $startTime = [uint64]$fields[19]; $identity = "${linuxProcessIdentifier}:$startTime"
        $Ledger[$identity] = [pscustomobject]@{ PID=$linuxProcessIdentifier; StartTime=$startTime }
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

function Get-LinuxBoundarySnapshot([int]$Identifier) {
    try {
        $text=[IO.File]::ReadAllText("/proc/$Identifier/stat")
        $close=$text.LastIndexOf(')')
        if($close-lt 1){throw 'Malformed Linux boundary process record'}
        $fields=$text.Substring($close+2).Split(' ',[StringSplitOptions]::RemoveEmptyEntries)
        if($fields.Count-lt 20){throw 'Incomplete Linux boundary process record'}
        $pidNamespace=[string]((Get-Item -LiteralPath "/proc/$Identifier/ns/pid" -Force -ErrorAction Stop).Target)
        $userNamespace=[string]((Get-Item -LiteralPath "/proc/$Identifier/ns/user" -Force -ErrorAction Stop).Target)
        $mountNamespace=[string]((Get-Item -LiteralPath "/proc/$Identifier/ns/mnt" -Force -ErrorAction Stop).Target)
        if($mountNamespace-notmatch '^mnt:\[[0-9]+\]$'){throw 'Malformed Linux mount namespace'}
        if($pidNamespace-notmatch '^pid:\[[0-9]+\]$'-or$userNamespace-notmatch '^user:\[[0-9]+\]$'){throw 'Malformed Linux boundary namespace'}
        [pscustomobject]@{PID=$Identifier;Parent=[int]$fields[1];Session=[int]$fields[3];State=$fields[0];StartTime=[uint64]$fields[19];PidNamespace=$pidNamespace;UserNamespace=$userNamespace;MountNamespace=$mountNamespace}
    } catch [IO.FileNotFoundException] { return $null }
      catch [IO.DirectoryNotFoundException] { return $null }
      catch [Management.Automation.ItemNotFoundException] { return $null }
}

function New-LinuxBoundaryStartInfo([string]$ExecutableReference,[string[]]$Arguments,[Collections.Generic.Dictionary[string,string]]$Environment,[string]$WorkingDirectory) {
    if($ExecutableReference-cnotmatch "^/proc/$PID/fd/[0-9]+$"){throw 'Linux executable must be the current caller held descriptor'}
    $start=[Diagnostics.ProcessStartInfo]::new()
    $start.FileName='/usr/bin/setsid';$start.UseShellExecute=$false
    $start.RedirectStandardInput=$true;$start.RedirectStandardOutput=$true;$start.RedirectStandardError=$true
    $start.WorkingDirectory=$WorkingDirectory;$start.Environment.Clear()
    foreach($pair in $Environment.GetEnumerator()){$start.Environment.Add($pair.Key,$pair.Value)}
    foreach($value in @('--wait','/bin/sh','-c',
        'IFS= read -r gate || exit 125; [ "$gate" = PSCAN_OPEN_V1 ] || exit 125; exec 3<"$1" || exit 125; shift; exec "$@"',
        'pscan-open',$ExecutableReference,'/usr/bin/unshare','--user','--map-current-user','--pid','--mount-proc','--propagation','private','--fork','--kill-child=SIGKILL','--','/bin/sh','-c',
        'IFS= read -r gate || exit 125; [ "$gate" = PSCAN_EXEC_V1 ] || exit 125; exec /proc/self/fd/3 "$@"',
        'pscan-docker')){$start.ArgumentList.Add($value)}
    foreach($argument in $Arguments){$start.ArgumentList.Add($argument)}
    return $start
}

function Get-LinuxGateAction($Snapshot,$Original,[int]$Supervisor,[string]$OuterPidNamespace,[string]$OuterUserNamespace,[string]$OuterMountNamespace,[string]$HeldIdentity,[string]$InheritedIdentity) {
    if($null-eq$Snapshot){throw 'Namespace init disappeared before release'}
    if($Snapshot.Parent-ne$Supervisor-or$Snapshot.Session-ne$Supervisor){throw 'Namespace init supervisor or session mismatch'}
    if($Snapshot.PidNamespace-eq$OuterPidNamespace-or$Snapshot.UserNamespace-eq$OuterUserNamespace-or$Snapshot.MountNamespace-eq$OuterMountNamespace){throw 'Required new namespaces unavailable'}
    if($null-eq$Original){return 'Stop'}
    foreach($key in @('PID','StartTime','Parent','Session','PidNamespace','UserNamespace','MountNamespace')){
        if($Snapshot.$key-cne$Original.$key){throw 'Namespace init identity changed before release'}
    }
    if($Snapshot.State-cne'T'){return 'Wait'}
    if($HeldIdentity-cnotmatch '^[0-9]+:[0-9]+$'-or$InheritedIdentity-cne$HeldIdentity){throw 'Inherited executable descriptor identity mismatch'}
    return 'Release'
}

function Test-LinuxInitReaped($Init,$InitHandle) {
    if ($null-eq$Init-or$null-eq$InitHandle) { return $false }
    if (![PscanNativeBoundary]::LinuxExited($InitHandle)) { return $false }
    # Exited pidfd / missing namespace link alone can still describe a zombie.
    # Preserve its supervisor until the original process record is absent.
    try { $stat=[IO.File]::ReadAllText("/proc/$($Init.PID)/stat") }
    catch [IO.FileNotFoundException] { return $true }
    catch [IO.DirectoryNotFoundException] { return $true }
    $close=$stat.LastIndexOf(')');$space=$stat.IndexOf(' ')
    if($close-lt 1-or$space-lt 1){throw 'Malformed init reaping identity'}
    $fields=$stat.Substring($close+2).Split(' ',[StringSplitOptions]::RemoveEmptyEntries)
    if($fields.Count-lt 20-or$stat.Substring(0,$space)-cne[string]$Init.PID-or
        $fields[19]-notmatch'^[0-9]+$'-or[uint64]$fields[19]-ne$Init.StartTime){throw 'Init identity changed before reaping observation'}
    return $false
}

function Invoke-LinuxPinnedCleanup($Init,$InitHandle,$RootHandle,$Process,$StdOutTask,$StdErrTask,[int]$Session,[hashtable]$Ledger,$Budget,[bool]$InitialUncertainty) {
    $Budget.BeginCleanup() # CAS preserves the original terminal deadline.
    $uncertain=$InitialUncertainty;$failure=if($uncertain){'close-gate'}else{$null}
    $initReaped=$null;$rootExited=$null;$membersEmpty=$null;$streamsClosed=$null;$namespaceGone=$null
    $proved=$false;$rootSignalled=$false;$stage='signal-init'
    try {
        if($null-ne$InitHandle-and$Budget.Remaining-gt 0){[void][PscanNativeBoundary]::LinuxSignal($InitHandle,9)}
        elseif($null-eq$InitHandle-and$null-ne$RootHandle-and$Budget.Remaining-gt 0){
            # Before init admission there is no invented namespace-closure proof.
            $stage='signal-root';[void][PscanNativeBoundary]::LinuxSignal($RootHandle,9);$rootSignalled=$true
        }
    } catch {$uncertain=$true;if($null-eq$failure){$failure=$stage}}
    while($Budget.Remaining-gt 0){
        try {
            $stage='inspect-init';$initReaped=Test-LinuxInitReaped $Init $InitHandle
            $stage='inspect-root';$rootExited=$Process.HasExited-and$null-ne$RootHandle-and[PscanNativeBoundary]::LinuxExited($RootHandle)
            if($initReaped-and!$rootExited-and!$rootSignalled-and$null-ne$RootHandle-and$Budget.Remaining-gt 0){
                $stage='signal-root';[void][PscanNativeBoundary]::LinuxSignal($RootHandle,9);$rootSignalled=$true
            }
            $stage='inspect-members';$membersEmpty=@(Get-LinuxSessionMembers -Session $Session -Ledger $Ledger).Count-eq 0
            $stage='inspect-namespace';$namespaceGone=$null-ne$Init-and$initReaped-and!(Test-Path -LiteralPath "/proc/$($Init.PID)/ns/pid" -ErrorAction Stop)
            $stage='inspect-streams';$streamsClosed=$null-ne$StdOutTask-and$null-ne$StdErrTask-and$StdOutTask.IsCompleted-and$StdErrTask.IsCompleted
            if($initReaped-and$rootExited-and$membersEmpty-and$namespaceGone-and$streamsClosed-and$Budget.Remaining-gt 0){$proved=!$uncertain;break}
        } catch {$uncertain=$true;if($null-eq$failure){$failure=$stage}}
        Start-Sleep -Milliseconds 2
    }
    [pscustomobject]@{Proved=$proved;FailureStage=$failure;DeadlineExpired=($Budget.Remaining-le 0);
        InitReaped=$initReaped;RootExited=$rootExited;MembersEmpty=$membersEmpty;NamespaceGone=$namespaceGone;StreamsClosed=$streamsClosed}
}

function Invoke-LinuxSessionBoundary([string]$ExecutableReference,[string[]]$Arguments,[Collections.Generic.Dictionary[string,string]]$Environment,[string]$WorkingDirectory, $SharedBudget=$null) {
    if ($null-eq$SharedBudget) { $SharedBudget=[PscanOperationBudget]::new($dockerBudgetMilliseconds,$cleanupGraceMilliseconds,$dockerOutputLimit) }
    if ($SharedBudget.Remaining-le 0) { throw 'Shared Docker deadline exhausted before Linux process creation' }
    # Reject an unsupported ABI/kernel before creating any child process.
    $probe=[PscanNativeBoundary]::LinuxPin($PID)
    try {if([PscanNativeBoundary]::LinuxExited($probe)){throw 'Calling Linux process identity unavailable'}}
    finally {$probe.Dispose()}
    foreach($path in @('/usr/bin/setsid','/usr/bin/unshare')){
        $item=Get-Item -LiteralPath $path -Force
        $owner=(& /usr/bin/stat -c '%u' -- $path).Trim()
        if($LASTEXITCODE-ne 0-or!$item-or$item.LinkType-or($item.UnixFileMode-band[IO.UnixFileMode]'GroupWrite,OtherWrite')-ne 0-or$owner-ne'0'){throw 'Fixed Linux launcher is untrusted'}
    }
    $shellIdentity=(& /usr/bin/stat -Lc '%u:%a' -- /bin/sh).Trim()
    if($LASTEXITCODE-ne 0-or$shellIdentity-notmatch '^0:[1357][0145][0145]$'){throw 'Fixed Linux gate shell is untrusted'}
    $start=New-LinuxBoundaryStartInfo $ExecutableReference $Arguments $Environment $WorkingDirectory
    $heldIdentity=(& /usr/bin/stat -Lc '%d:%i' -- $ExecutableReference).Trim()
    if($LASTEXITCODE-ne 0-or$heldIdentity-notmatch '^[0-9]+:[0-9]+$'){throw 'Held Linux executable identity unavailable'}
    $outer=Get-LinuxBoundarySnapshot $PID
    if($null-eq$outer){throw 'Calling Linux context unavailable'}
    $process=[Diagnostics.Process]::new();$process.StartInfo=$start
    $handles=@{};$ledger=@{};$init=$null;$initHandle=$null;$rootHandle=$null
    $namespaceIdentity=$null;$resumed=$false;$stdinClosed=$false;$started=$false;$terminal=$null
    $stdoutTask=$null;$stderrTask=$null;$session=0;$phase='root-pin';$cleanupUncertain=$false
    $cleanupDetails=$null
    $watch=[Diagnostics.Stopwatch]::StartNew()

    function Pin-Member($member) {
        $key="$($member.PID):$($member.StartTime)"
        if($handles.ContainsKey($key)){return $handles[$key]}
        $handle=[PscanNativeBoundary]::LinuxPin($member.PID)
        try {
            $again=Get-LinuxBoundarySnapshot $member.PID
            if($null-eq$again-or$again.StartTime-ne$member.StartTime-or$again.Parent-ne$member.Parent-or$again.Session-ne$member.Session-or[PscanNativeBoundary]::LinuxExited($handle)){throw 'Linux process identity changed during pin acquisition'}
            $handles.Add($key,$handle)
            return $handle
        } catch {$handle.Dispose();throw}
    }
    function Close-Gate {
        if(!$stdinClosed-and$started){$process.StandardInput.Close()}
    }
    try {
        try {
            if($SharedBudget.Remaining-le 0){throw 'Shared deadline exhausted before Linux creation'}
            if(!$process.Start()){throw 'Linux session process did not start'}
            $started=$true;$session=$process.Id
            $stdoutTask=[PscanNativeBoundary]::CaptureShared($process.StandardOutput.BaseStream,$dockerOutputLimit,$SharedBudget,$false)
            $stderrTask=[PscanNativeBoundary]::CaptureShared($process.StandardError.BaseStream,$dockerOutputLimit,$SharedBudget,$true)
            while($SharedBudget.Remaining-gt 0){
                if($stdoutTask.IsFaulted-or$stderrTask.IsFaulted){throw 'stream overflow or read failure'}
                if($phase-eq'root-pin'){
                    if($process.HasExited){throw 'Launcher exited before supervisor admission'}
                    $root=Get-LinuxBoundarySnapshot $session
                    if($null-ne$root){
                        if($root.Parent-ne$PID){throw 'Launcher parent identity mismatch'}
                        if($root.Session-eq$session){
                            $rootHandle=Pin-Member $root
                            if($process.HasExited){throw 'Launcher exited during supervisor admission'}
                            $process.StandardInput.WriteLine('PSCAN_OPEN_V1');$process.StandardInput.Flush()
                            $phase='init-discovery'
                        }
                    }
                } elseif(!$resumed){
                    if($process.HasExited-or[PscanNativeBoundary]::LinuxExited($rootHandle)){throw 'Launcher exited before namespace admission'}
                    $childrenText=[IO.File]::ReadAllText("/proc/$session/task/$session/children").Trim()
                    $children=@(($childrenText-split'\s+')|Where-Object{$_-ne''})
                    if($children.Count-gt 1-or@($children|Where-Object{$_-notmatch'^[0-9]+$'}).Count){throw 'Ambiguous namespace init membership'}
                    if($children.Count-eq 1){
                        $snapshot=Get-LinuxBoundarySnapshot ([int]$children[0])
                        $inherited=''
                        if($null-ne$init-and$null-ne$snapshot-and$snapshot.State-ceq'T'){
                            $inherited=(& /usr/bin/stat -Lc '%d:%i' -- "/proc/$($snapshot.PID)/fd/3").Trim()
                            if($LASTEXITCODE-ne 0){throw 'Inherited held descriptor unavailable'}
                        }
                        $action=Get-LinuxGateAction $snapshot $init $session $outer.PidNamespace $outer.UserNamespace $outer.MountNamespace $heldIdentity $inherited
                        if($action-eq'Stop'){
                            $initHandle=Pin-Member $snapshot;$init=$snapshot;$namespaceIdentity=$snapshot.PidNamespace
                            if(![PscanNativeBoundary]::LinuxSignal($initHandle,19)){throw 'Namespace init exited before ancestor stop'}
                            $phase='init-stopping'
                        } elseif($action-eq'Release'){
                            if([PscanNativeBoundary]::LinuxExited($initHandle)){throw 'Pinned namespace init exited at stopped gate'}
                            # Init cannot consume this token while stopped. Closing the owned
                            # pipe leaves EOF for the eventual Docker command, never host stdin.
                            $process.StandardInput.WriteLine('PSCAN_EXEC_V1');$process.StandardInput.Flush()
                            Close-Gate;$stdinClosed=$true
                            if(![PscanNativeBoundary]::LinuxSignal($initHandle,18)){throw 'Namespace init resume failed'}
                            $resumed=$true;$phase='running'
                        }
                    } elseif($null-ne$init){throw 'Namespace init membership disappeared before release'}
                }
                $members=@(Get-LinuxSessionMembers -Session $session -Ledger $ledger)
                if($resumed-and$process.HasExited-and$stdoutTask.IsCompleted-and$stderrTask.IsCompleted-and$members.Count-eq 0-and[PscanNativeBoundary]::LinuxExited($initHandle)-and!(Test-Path -LiteralPath "/proc/$($init.PID)/ns/pid")){break}
                Start-Sleep -Milliseconds 2
            }
            if(!$resumed-or!$process.HasExited-or!$stdoutTask.IsCompleted-or!$stderrTask.IsCompleted-or@(Get-LinuxSessionMembers -Session $session -Ledger $ledger).Count-ne 0-or![PscanNativeBoundary]::LinuxExited($initHandle)-or(Test-Path -LiteralPath "/proc/$($init.PID)/ns/pid")){throw 'timeout or incomplete lifecycle'}
        } catch {$terminal="${phase}: $($_.Exception.Message)";$SharedBudget.BeginCleanup()}

        if($terminal-and$started){
            $SharedBudget.BeginCleanup()
            try {Close-Gate;$stdinClosed=$true} catch {$cleanupUncertain=$true}
            # Namespace-init death kills all namespace descendants, including
            # detached sessions. Preserve the exact supervisor to reap it; never
            # reacquire/signal numeric descendants racing with their own exit.
            $cleanupDetails=Invoke-LinuxPinnedCleanup $init $initHandle $rootHandle $process $stdoutTask $stderrTask $session $ledger $SharedBudget $cleanupUncertain
            if(!$cleanupDetails.Proved){$terminal+='; cleanup uncertainty'}
        }
        # Assign array members directly: a PowerShell statement pipeline would
        # enumerate empty/single-byte buffers into null/scalar values.
        $stdout=[byte[]]::new(0);$stderr=[byte[]]::new(0)
        $captureComplete=$null-ne$stdoutTask-and$null-ne$stderrTask-and$stdoutTask.IsCompletedSuccessfully-and$stderrTask.IsCompletedSuccessfully
        if($null-ne$stdoutTask-and$stdoutTask.IsCompletedSuccessfully){$stdout=$stdoutTask.Result}
        if($null-ne$stderrTask-and$stderrTask.IsCompletedSuccessfully){$stderr=$stderrTask.Result}
        if(!$captureComplete-or$stdout-isnot[byte[]]-or$stderr-isnot[byte[]]){
            if(!$terminal){$terminal='capture evidence unavailable'}
        }
        [pscustomobject]@{ExitCode=$(if($started-and$process.HasExited){$process.ExitCode}else{199});StdOut=$stdout;StdErr=$stderr;Terminal=$terminal;ObservedMembers=@($ledger.Keys)+@($namespaceIdentity);ContainmentEmpty=(!$terminal);CleanupDetails=$cleanupDetails}
    } finally {
        try{Close-Gate}catch{}
        foreach($handle in $handles.Values){$handle.Dispose()}
        $process.Dispose()
    }
}

function Invoke-HeldDockerCall([string[]]$Arguments,$Context) {
    if ($Context.Budget.Remaining-le 0) { throw 'Shared Docker deadline exhausted before protocol call' }
    $Context.Calls++
    if ($Context.Calls-gt 9) { throw 'Fixed Docker protocol call count exceeded' }
    if ($IsWindows) {
        $result=[PscanNativeBoundary]::RunWindowsShared($Context.Path,$Arguments,$Context.Environment,$Context.Working,$dockerOutputLimit,$Context.Budget)
    } else {
        $result=Invoke-LinuxSessionBoundary -ExecutableReference $Context.Reference -Arguments $Arguments -Environment $Context.Environment -WorkingDirectory $Context.Working -SharedBudget $Context.Budget
    }
    $Context.Observed += $result.ObservedMembers.Count
    if ($result.Terminal-or!$result.ContainmentEmpty) { throw "Docker protocol native containment failed: $($result.Terminal)" }
    if ($Context.Budget.Remaining-le 0) { throw 'Shared Docker deadline exhausted after protocol call' }
    if($result.StdOut-isnot[byte[]]-or$result.StdErr-isnot[byte[]]){throw 'Docker protocol capture buffers are not byte arrays'}
    $utf8=[Text.UTF8Encoding]::new($false,$true)
    return [pscustomobject]@{ExitCode=$result.ExitCode;StdOut=$utf8.GetString($result.StdOut);StdErr=$utf8.GetString($result.StdErr)}
}

. (Join-Path $PSScriptRoot 'docker-container-lifecycle.ps1')

function Invoke-BoundDocker([string[]]$Arguments) {
    $budget=[PscanOperationBudget]::new($dockerBudgetMilliseconds,$cleanupGraceMilliseconds,$dockerOutputLimit)
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
        } else {
            $environment.Add('TMPDIR', $temp); $environment.Add('DOCKER_HOST', 'unix:///var/run/docker.sock'); $environment.Add('LANG', 'C.UTF-8')
        }
        $context=@{Budget=$budget;Path=$dockerPath;Reference=$reference;Environment=$environment;Working=$working;Image=$image;Observed=0;Calls=0}
        $result=if ($Arguments[0]-ceq'run') { Invoke-ContainerProtocol $Arguments $context } else { Invoke-HeldDockerCall $Arguments $context }
        if ($budget.CleanupStarted-or$budget.Remaining-le 0) { throw 'Shared Docker lifecycle is terminal or expired' }
        $identityStream.Position = 0; $after = (Get-FileHash -InputStream $identityStream -Algorithm SHA256).Hash.ToLowerInvariant()
        if ($after -cne $digest) { throw 'Docker executable identity changed during execution' }
        $stdout=$result.StdOut; $stderr=$result.StdErr
        if (Get-ChildItem -LiteralPath $config -Force | Select-Object -First 1) { throw 'Docker wrote unexpected configuration or credential state' }
        if ($budget.Remaining-le 0) { throw 'Shared Docker deadline exhausted before final evidence' }
        return [pscustomobject]@{
            Operation=$Operation; DockerSHA256=$digest; ExitCode=$result.ExitCode; StdOut=$stdout; StdErr=$stderr
            StdOutByteCount=[Text.Encoding]::UTF8.GetByteCount($stdout); StdErrByteCount=[Text.Encoding]::UTF8.GetByteCount($stderr)
            ContainmentMembersObserved=$context.Observed; ContainmentEmpty=$true
            DaemonContainerID=$result.ContainerID; DaemonContainerRemoved=($Arguments[0]-ceq'run'-and$result.DaemonEmpty)
            ProtocolCallCount=$context.Calls; AggregateStdOutByteCount=$budget.StdOutUsed; AggregateStdErrByteCount=$budget.StdErrUsed
        }
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
if ! printf '%s\n' 'PSCAN-06-C2-CONTAINER-CACHE-CANARY' > "$canary"; then
    printf '%s\n' '{"schema":"pscan-cache-write-v1","phase":"create-write","outcome":"failed"}'
    exit 74
fi
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
cp /work/tooling/contracts/release-manifest/schema-2.2.json /out/tooling-materialized/contracts/release-manifest/schema-2.2.json
cp /work/tooling/contracts/release-manifest/schema-2.3.json /out/tooling-materialized/contracts/release-manifest/schema-2.3.json
cp /work/tooling/docs/release/OFFLINE-VERIFICATION-RUNBOOK.md /out/tooling-materialized/docs/release/OFFLINE-VERIFICATION-RUNBOOK.md
cp /work/tooling/docs/release/SCANNER-IO-REFERENCE.md /out/tooling-materialized/docs/release/SCANNER-IO-REFERENCE.md
'@ }

function Get-PackagePayload([string]$Epoch) {
    ConvertTo-LF "cp /tools/packager-linux-amd64 /tmp/packager; chmod 0755 /tmp/packager; /tmp/packager -root /input/linux -output /dist/project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz -format tar.gz -epoch $Epoch; /tmp/packager -root /input/windows -output /dist/project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip -format zip -epoch $Epoch"
}

function Resolve-ContainerUser([Nullable[int]]$UIDValue, [Nullable[int]]$GIDValue, [bool]$Required) {
    if ($null -eq $UIDValue -and $null -eq $GIDValue) {
        if ($Required) { throw 'Linux container operation requires the admitted host UID/GID pair' }
        return [pscustomobject]@{ Arguments=@(); TmpfsSuffix='' }
    }
    if ($null -eq $UIDValue -or $null -eq $GIDValue -or $UIDValue -lt 0 -or $GIDValue -lt 0) {
        throw 'Container operation requires a complete non-negative UID/GID pair'
    }
    return [pscustomobject]@{
        Arguments=@('--user', "${UIDValue}:${GIDValue}")
        TmpfsSuffix=",mode=0700,uid=$UIDValue,gid=$GIDValue"
    }
}

$containerUser = $null
if ($Operation -in @('ContainerCacheProof','DependencyAcquisition','ReleaseBuild','ReleasePackage')) {
    $containerUser = Resolve-ContainerUser $HostUID $HostGID $IsLinux
    if ($IsLinux) {
        $currentUID = (& /usr/bin/id -u).Trim()
        if ($LASTEXITCODE -ne 0 -or $currentUID -notmatch '^\d+$' -or [int]$currentUID -ne $HostUID) { throw 'Container UID differs from the invoking host identity' }
        $currentGID = (& /usr/bin/id -g).Trim()
        if ($LASTEXITCODE -ne 0 -or $currentGID -notmatch '^\d+$' -or [int]$currentGID -ne $HostGID) { throw 'Container GID differs from the invoking host identity' }
    }
}
$arguments = switch ($Operation) {
    'EngineInspection' { @('version','--format','{{json .Server}}') }
    'ExactImageInventory' { @('image','ls','--all','--no-trunc','--digests','--filter',"reference=$image",'--format','{{json .}}') }
    'ApprovedImagePull' { @('pull',$image) }
    'RepositoryDigestInspection' { @('image','inspect','--format','{{json .RepoDigests}}',$image) }
    'ContainerCacheProof' {
        $cache = Require-Path 'CacheDirectory' $CacheDirectory $true
        $tmpfs='/work:rw,noexec,nosuid,nodev,size=16m'; $a=@('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1')
        $a += $containerUser.Arguments
        $tmpfs += $(if ($containerUser.TmpfsSuffix) { $containerUser.TmpfsSuffix } else { ',mode=0700' })
        $mount="type=bind,src=$cache,dst=/gomodcache";if($ReadOnlyCache){$mount+=',readonly'}
        $a+@('--tmpfs',$tmpfs,'--mount',$mount,'--workdir','/work',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work','/bin/sh','-ceu',(Get-CachePayload))
    }
    'ContainerCrlfParse' {
        if([string]::IsNullOrEmpty($PayloadKind)){throw 'ContainerCrlfParse requires PayloadKind'}
        $payload=switch($PayloadKind){'CacheCanary'{Get-CachePayload};'Acquisition'{Get-AcquisitionPayload};'Build'{Get-BuildPayload};'Package'{Get-PackagePayload '0'}}
        @('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1','--tmpfs','/work:rw,noexec,nosuid,nodev,size=16m,mode=0700','--workdir','/work',$image,'/bin/sh','-n','-c',$payload)
    }
    'DependencyAcquisition' {
        $root=Require-Path 'SourceRoot' $SourceRoot $true;$runner=Require-Path 'RunnerGoArchive' $RunnerGoArchive;$engine=Require-Path 'EngineGoArchive' $EngineGoArchive;$gitleaks=Require-Path 'GitleaksArchive' $GitleaksArchive;$cache=Require-Path 'CacheDirectory' $CacheDirectory $true
        $a=@('run','--rm','--pull=never','--network','bridge','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','256','--memory','4g','--memory-swap','4g','--cpus','2','--tmpfs','/work:rw,exec,nosuid,nodev,size=1g','--mount',"type=bind,src=$root,dst=/src,readonly",'--mount',"type=bind,src=$runner,dst=/input/runner-go.tar.gz,readonly",'--mount',"type=bind,src=$engine,dst=/input/engine-go.tar.gz,readonly",'--mount',"type=bind,src=$gitleaks,dst=/input/gitleaks.tar.gz,readonly",'--mount',"type=bind,src=$cache,dst=/gomodcache",'--workdir','/src',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work','/bin/sh','-ceu',(Get-AcquisitionPayload))
        $index=$a.IndexOf('--tmpfs');$a=$a[0..($index-1)]+$containerUser.Arguments+$a[$index..($a.Count-1)]
        $a[$a.IndexOf('/work:rw,exec,nosuid,nodev,size=1g')] += $containerUser.TmpfsSuffix
        $a
    }
    'ReleaseBuild' {
        foreach($required in @('SourceDateEpoch','ProductRevision','ToolingRevision','ToolingTree','ProductArchiveSHA256','ToolingArchiveSHA256','Created')){if([string]::IsNullOrEmpty((Get-Variable -Name $required -ValueOnly))){throw "ReleaseBuild requires $required"}}
        $product=Require-Path 'ProductArchive' $ProductArchive;$productBlobs=Require-Path 'ProductBlobManifest' $ProductBlobManifest;$productPaths=Require-Path 'ProductPathManifest' $ProductPathManifest;$productModes=Require-Path 'ProductModeManifest' $ProductModeManifest
        $tooling=Require-Path 'ToolingArchive' $ToolingArchive;$toolingBlobs=Require-Path 'ToolingBlobManifest' $ToolingBlobManifest;$toolingPaths=Require-Path 'ToolingPathManifest' $ToolingPathManifest;$toolingModes=Require-Path 'ToolingModeManifest' $ToolingModeManifest
        $runner=Require-Path 'RunnerGoArchive' $RunnerGoArchive;$engine=Require-Path 'EngineGoArchive' $EngineGoArchive;$gitleaks=Require-Path 'GitleaksArchive' $GitleaksArchive;$cache=Require-Path 'CacheDirectory' $CacheDirectory $true;$raw=Require-Path 'RawOutputDirectory' $RawOutputDirectory $true
        $a=@('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','512','--memory','8g','--memory-swap','8g','--cpus','2','--tmpfs','/work:rw,exec,nosuid,nodev,size=4g','--mount',"type=bind,src=$product,dst=/input/product-source.tar,readonly",'--mount',"type=bind,src=$productBlobs,dst=/input/product-source.blobs.sha256,readonly",'--mount',"type=bind,src=$productPaths,dst=/input/product-source.paths,readonly",'--mount',"type=bind,src=$productModes,dst=/input/product-source.modes,readonly",'--mount',"type=bind,src=$tooling,dst=/input/release-tooling-source.tar,readonly",'--mount',"type=bind,src=$toolingBlobs,dst=/input/release-tooling-source.blobs.sha256,readonly",'--mount',"type=bind,src=$toolingPaths,dst=/input/release-tooling-source.paths,readonly",'--mount',"type=bind,src=$toolingModes,dst=/input/release-tooling-source.modes,readonly",'--mount',"type=bind,src=$runner,dst=/input/runner-go.tar.gz,readonly",'--mount',"type=bind,src=$engine,dst=/input/engine-go.tar.gz,readonly",'--mount',"type=bind,src=$gitleaks,dst=/input/gitleaks.tar.gz,readonly",'--mount',"type=bind,src=$cache,dst=/gomodcache,readonly",'--mount',"type=bind,src=$raw,dst=/out",'--workdir','/work',$image,'/usr/bin/env','-i','PATH=/usr/bin:/bin','HOME=/work',"SOURCE_DATE_EPOCH=$SourceDateEpoch","PSCAN_PRODUCT_REVISION=$ProductRevision","PSCAN_TOOLING_REVISION=$ToolingRevision","PSCAN_TOOLING_TREE=$ToolingTree","PSCAN_CREATED=$Created","PSCAN_PRODUCT_ARCHIVE_SHA256=$ProductArchiveSHA256","PSCAN_TOOLING_ARCHIVE_SHA256=$ToolingArchiveSHA256","PSCAN_PRODUCT_FILE_COUNT=$ProductFileCount","PSCAN_TOOLING_FILE_COUNT=$ToolingFileCount",'/bin/sh','-ceu',(Get-BuildPayload))
        $index=$a.IndexOf('--tmpfs');$a=$a[0..($index-1)]+$containerUser.Arguments+$a[$index..($a.Count-1)]
        $a[$a.IndexOf('--tmpfs')+1] += $containerUser.TmpfsSuffix
        $a
    }
    'ReleasePackage' {
        if([string]::IsNullOrEmpty($SourceDateEpoch)){throw 'ReleasePackage requires SourceDateEpoch'}
        $raw=Require-Path 'RawOutputDirectory' $RawOutputDirectory $true;$linux=Require-Path 'LinuxStageDirectory' $LinuxStageDirectory $true;$windows=Require-Path 'WindowsStageDirectory' $WindowsStageDirectory $true;$dist=Require-Path 'DistributionDirectory' $DistributionDirectory $true
        $a=@('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','64','--memory','1g','--memory-swap','1g','--cpus','1','--tmpfs','/tmp:rw,exec,nosuid,nodev,size=64m','--mount',"type=bind,src=$raw,dst=/tools,readonly",'--mount',"type=bind,src=$linux,dst=/input/linux,readonly",'--mount',"type=bind,src=$windows,dst=/input/windows,readonly",'--mount',"type=bind,src=$dist,dst=/dist",'--workdir','/tmp',$image,'/bin/sh','-ceu',(Get-PackagePayload $SourceDateEpoch))
        $index=$a.IndexOf('--tmpfs');$a=$a[0..($index-1)]+$containerUser.Arguments+$a[$index..($a.Count-1)]
        $a[$a.IndexOf('--tmpfs')+1] += $containerUser.TmpfsSuffix
        $a
    }
}

$result = Invoke-BoundDocker -Arguments $arguments
$result | ConvertTo-Json -Compress
