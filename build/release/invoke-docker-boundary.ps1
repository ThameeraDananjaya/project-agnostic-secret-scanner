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

# A fixed fresh PowerShell transport; Docker remains created only by the closed child.
if ($MyInvocation.InvocationName -eq '.') { throw 'Boundary child launcher cannot be dot-sourced' }
$ErrorActionPreference='Stop'
function Invoke-BoundaryChild([Diagnostics.ProcessStartInfo]$StartInfo) {
    $process=[Diagnostics.Process]::new();$process.StartInfo=$StartInfo
    $watch=[Diagnostics.Stopwatch]::StartNew()
    $started=$false;$completed=$false
    $streams=@();$captures=@();$buffers=@();$reads=@();$closed=@($false,$false)
    try {
        if(!$process.Start()){throw 'Boundary child did not start'};$started=$true
        $streams=@($process.StandardOutput.BaseStream,$process.StandardError.BaseStream)
        $captures=@([IO.MemoryStream]::new(),[IO.MemoryStream]::new())
        $buffers=@([byte[]]::new(8192),[byte[]]::new(8192))
        $reads=@($streams[0].ReadAsync($buffers[0],0,8192),$streams[1].ReadAsync($buffers[1],0,8192))
        while($true){
            for($i=0;$i-lt 2;$i++){
                if(!$closed[$i]-and$reads[$i].IsCompleted){
                    $n=$reads[$i].GetAwaiter().GetResult()
                    if($n-eq 0){$closed[$i]=$true}
                    else{
                        # JSON escaping can expand each admitted Docker byte sixfold.
                        # 2MiB admits both 131072-byte streams plus fixed metadata.
                        if($captures[$i].Length+$n-gt 2097152){throw 'Boundary child transport overflow'}
                        $captures[$i].Write($buffers[$i],0,$n)
                        $reads[$i]=$streams[$i].ReadAsync($buffers[$i],0,8192)
                    }
                }
            }
            if($process.HasExited-and$closed[0]-and$closed[1]){break}
            # Independent startup/compilation/JSON transport allowance. The child
            # still enforces the unchanged 15000ms Docker + 2000ms cleanup bounds.
            if($watch.ElapsedMilliseconds-ge 35000){throw 'Boundary child transport timeout'}
            Start-Sleep -Milliseconds 2
        }
        $utf8=[Text.UTF8Encoding]::new($false,$true)
        $stdout=$utf8.GetString($captures[0].ToArray());$stderr=$utf8.GetString($captures[1].ToArray())
        if($process.ExitCode-ne 0-or$stderr.Length-ne 0){throw "Boundary child failed exit=$($process.ExitCode): $stderr"}
        $completed=$true
        return $stdout
    } finally {
        if($started-and!$completed){
            try{if(!$process.HasExited){$process.Kill($true)}}catch{}
            if(!$process.WaitForExit(2000)){[Console]::Error.WriteLine('Boundary child cleanup uncertain')}
            # A transport failure never becomes trusted containment evidence.
        }
        foreach($stream in $captures){$stream.Dispose()}
        $process.Dispose()
    }
}
$start=[Diagnostics.ProcessStartInfo]::new()
$start.FileName=if($IsLinux){'/opt/microsoft/powershell/7/pwsh'}elseif($IsWindows){Join-Path $PSHOME 'pwsh.exe'}else{throw 'Unsupported boundary child host'}
$start.UseShellExecute=$false;$start.CreateNoWindow=$true
$start.RedirectStandardOutput=$true;$start.RedirectStandardError=$true
foreach($arg in @('-NoLogo','-NoProfile','-NonInteractive','-File',(Join-Path $PSScriptRoot 'docker-execution.ps1'))){$start.ArgumentList.Add($arg)}
foreach($name in @($PSBoundParameters.Keys|Sort-Object)){
    $value=$PSBoundParameters[$name]
    if($value-is[Management.Automation.SwitchParameter]){if($value.IsPresent){$start.ArgumentList.Add('-'+$name)};continue}
    $start.ArgumentList.Add('-'+$name)
    $start.ArgumentList.Add([Convert]::ToString($value,[Globalization.CultureInfo]::InvariantCulture))
}
$json=Invoke-BoundaryChild $start
$document=[Text.Json.JsonDocument]::Parse($json)
try {
    if($document.RootElement.ValueKind-ne[Text.Json.JsonValueKind]::Object){throw 'Boundary child JSON must be one object'}
    $names=[Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
    foreach($property in $document.RootElement.EnumerateObject()){if(!$names.Add($property.Name)){throw 'Duplicate boundary child JSON field'}}
} finally {$document.Dispose()}
$result=$json|ConvertFrom-Json -AsHashtable
$expected=@('Operation','DockerSHA256','ExitCode','StdOut','StdErr','StdOutByteCount','StdErrByteCount','ContainmentMembersObserved','ContainmentEmpty','DaemonContainerID','DaemonContainerRemoved','ProtocolCallCount','AggregateStdOutByteCount','AggregateStdErrByteCount')
if($result-isnot[Collections.IDictionary]-or$result.Count-ne$expected.Count-or@($result.Keys|Where-Object{$_-cnotin$expected}).Count-ne 0-or
    $result.Operation-cne$Operation-or$result.DockerSHA256-notmatch'^[0-9a-f]{64}$'-or$result.ContainmentEmpty-isnot[bool]-or!$result.ContainmentEmpty-or
    $result.ExitCode-isnot[long]-and$result.ExitCode-isnot[int]-or$result.StdOut-isnot[string]-or$result.StdErr-isnot[string]){throw 'Boundary child returned malformed evidence'}
foreach($pair in @(@('StdOut','StdOutByteCount'),@('StdErr','StdErrByteCount'))){
    $count=$result[$pair[1]]
    if(($count-isnot[long]-and$count-isnot[int])-or$count-lt 0-or$count-gt 131072-or[Text.Encoding]::UTF8.GetByteCount($result[$pair[0]])-ne$count){throw 'Boundary child stream binding mismatch'}
}
if(($result.ContainmentMembersObserved-isnot[long]-and$result.ContainmentMembersObserved-isnot[int])-or$result.ContainmentMembersObserved-lt 0){throw 'Boundary member count invalid'}
$containerOperation=$Operation-in@('ContainerCacheProof','ContainerCrlfParse','DependencyAcquisition','ReleaseBuild','ReleasePackage')
if($result.DaemonContainerRemoved-isnot[bool]-or$result.DaemonContainerRemoved-ne$containerOperation-or
    ($containerOperation-and($result.DaemonContainerID-isnot[string]-or$result.DaemonContainerID-cnotmatch'^[0-9a-f]{64}$'))-or
    (!$containerOperation-and$null-ne$result.DaemonContainerID)-or
    ($result.ProtocolCallCount-isnot[long]-and$result.ProtocolCallCount-isnot[int])-or$result.ProtocolCallCount-ne$(if($containerOperation){7}else{1})){throw 'Daemon lifecycle binding mismatch'}
foreach($stream in @('StdOut','StdErr')){
    $aggregate=$result['Aggregate'+$stream+'ByteCount']
    if(($aggregate-isnot[long]-and$aggregate-isnot[int])-or$aggregate-lt$result[$stream+'ByteCount']-or$aggregate-gt 131072){throw 'Aggregate protocol stream evidence invalid'}
}
$global:LASTEXITCODE=0
Write-Host ('PSCAN_DOCKER_LIFECYCLE '+(@{
    operation=$Operation;container=$result.DaemonContainerID;daemonRemoved=$result.DaemonContainerRemoved
    nativeEmpty=$result.ContainmentEmpty;calls=$result.ProtocolCallCount
    stdoutBytes=$result.AggregateStdOutByteCount;stderrBytes=$result.AggregateStdErrByteCount
}|ConvertTo-Json -Compress))
Write-Output $json
