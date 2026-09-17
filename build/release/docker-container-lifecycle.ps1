# Internal fixed Docker protocol. This file exposes no command-line entrypoint
# or callback; it is materialized and admitted with the closed boundary source.
function Read-ContainerObject([string]$Json) {
    if ([Text.Encoding]::UTF8.GetByteCount($Json)-gt 131072) { throw 'Container inspection exceeds stream bound' }
    $options=[Text.Json.JsonDocumentOptions]::new(); $options.MaxDepth=32
    $document=[Text.Json.JsonDocument]::Parse($Json,$options)
    try {
        $pending=[Collections.Generic.Stack[Text.Json.JsonElement]]::new(); $pending.Push($document.RootElement); $nodes=0
        while ($pending.Count) {
            if (++$nodes-gt 16384) { throw 'Container inspection structure exceeds bound' }
            $value=$pending.Pop()
            if ($value.ValueKind-eq[Text.Json.JsonValueKind]::Object) {
                $names=[Collections.Generic.HashSet[string]]::new([StringComparer]::Ordinal)
                foreach ($property in $value.EnumerateObject()) { if (!$names.Add($property.Name)) { throw 'Duplicate container inspection field' }; $pending.Push($property.Value) }
            } elseif ($value.ValueKind-eq[Text.Json.JsonValueKind]::Array) {
                foreach ($item in $value.EnumerateArray()) { $pending.Push($item) }
            }
        }
        if ($document.RootElement.ValueKind-ne[Text.Json.JsonValueKind]::Object) { throw 'Container inspection is not one object' }
    } finally { $document.Dispose() }
    return ConvertFrom-Json -InputObject $Json -AsHashtable -Depth 32
}

function Assert-ContainerIdentity($Object,[string]$Identifier,[string]$Name,[string]$Token) {
    if ($Identifier-cnotmatch'^[0-9a-f]{64}$'-or$Object.Id-cne$Identifier-or$Object.Name-cne('/'+$Name)-or
        $Object.Config.Labels['org.pscan.invocation']-cne$Token) { throw 'Exact container ownership is unproved' }
}

function Assert-ContainerConfiguration($Object,[string[]]$RunArguments,[string]$Image,[string]$ImageID) {
    function Value([string]$Key) {
        $locations=@(for($i=0;$i-lt$RunArguments.Count;$i++){if($RunArguments[$i]-ceq$Key){$i}})
        if ($locations.Count-ne 1-or$locations[0]+1-ge$RunArguments.Count) { throw "Fixed Docker option is missing or duplicated: $Key" }
        return $RunArguments[$locations[0]+1]
    }
    function SameList($Actual,$Expected,[string]$Label) {
        $a=@($Actual); $e=@($Expected)
        if ($a.Count-ne$e.Count) { throw "$Label count differs" }
        for($i=0;$i-lt$a.Count;$i++){if($a[$i]-cne$e[$i]){throw "$Label differs"}}
    }
    function Bytes([string]$Text) {
        if ($Text-cnotmatch'^(\d+)([mg])$') { throw 'Unsupported fixed memory option' }
        return [long]$Matches[1]*$(if($Matches[2]-ceq'm'){1048576}else{1073741824})
    }
    $config=$Object.Config; $hostConfig=$Object.HostConfig
    if ($Object.Image-cne$ImageID-or$config.Image-cne$Image-or$hostConfig.Privileged-isnot[bool]-or$hostConfig.Privileged-or
        $hostConfig.ReadonlyRootfs-isnot[bool]-or!$hostConfig.ReadonlyRootfs-or$hostConfig.AutoRemove-isnot[bool]-or$hostConfig.AutoRemove-or
        $hostConfig.PublishAllPorts-isnot[bool]-or$hostConfig.PublishAllPorts-or
        $hostConfig.RestartPolicy.Name-cne'no'-or$hostConfig.RestartPolicy.MaximumRetryCount-ne 0-or
        $hostConfig.PidMode-cne''-or$hostConfig.UsernsMode-cne''-or$hostConfig.IpcMode-cne'private') { throw 'Created container identity or isolation configuration differs' }
    foreach($field in @('CapAdd','Devices','DeviceRequests','Binds','VolumesFrom','PortBindings')) {
        if ($null-ne$hostConfig[$field]-and$hostConfig[$field].Count-ne 0) { throw "Unexpected container privilege or mount field: $field" }
    }
    SameList $hostConfig.CapDrop @('ALL') 'Dropped capabilities'
    SameList $hostConfig.SecurityOpt @('no-new-privileges') 'Security options'
    $network=if($RunArguments-contains'--network'){Value '--network'}else{'default'}
    if ($hostConfig.NetworkMode-cne$network-or$hostConfig.Memory-ne(Bytes (Value '--memory'))-or$hostConfig.MemorySwap-ne(Bytes (Value '--memory-swap'))-or
        $hostConfig.PidsLimit-ne[int](Value '--pids-limit')-or$hostConfig.NanoCpus-ne([long](Value '--cpus')*1000000000)) { throw 'Created container network or resource limits differ' }
    $user=if($RunArguments-contains'--user'){Value '--user'}else{''}
    $workdir=if($RunArguments-contains'--workdir'){Value '--workdir'}else{''}
    if ($config.User-cne$user-or$config.WorkingDir-cne$workdir-or$config.Tty-isnot[bool]-or$config.Tty-or$config.OpenStdin-isnot[bool]-or$config.OpenStdin-or
        ($null-ne$config.Entrypoint-and$config.Entrypoint.Count-ne 0)-or($null-ne$config.Volumes-and$config.Volumes.Count-ne 0)) { throw 'Created container user, process or implicit volume configuration differs' }
    $imageIndex=$RunArguments.IndexOf($Image)
    if ($imageIndex-lt 0-or$imageIndex+1-ge$RunArguments.Count) { throw 'Fixed container command missing' }
    SameList $config.Cmd $RunArguments[($imageIndex+1)..($RunArguments.Count-1)] 'Container command'
    $tmpfs=(Value '--tmpfs').Split(':',2)
    if ($hostConfig.Tmpfs.Count-ne 1-or$hostConfig.Tmpfs[$tmpfs[0]]-cne$tmpfs[1]) { throw 'Created container tmpfs differs' }
    $mounts=@(for($i=0;$i-lt$imageIndex;$i++){if($RunArguments[$i]-ceq'--mount'){$RunArguments[$i+1]}})
    if (@($hostConfig.Mounts).Count-ne$mounts.Count) { throw 'Created container bind mount count differs' }
    for($i=0;$i-lt$mounts.Count;$i++) {
        if ($mounts[$i]-cnotmatch'^type=bind,src=(.+),dst=([^,]+)(,readonly)?$') { throw 'Unsupported fixed bind syntax' }
        $source=$Matches[1];$target=$Matches[2];$readOnly=$Matches[3]-ceq',readonly';$mount=$hostConfig.Mounts[$i]
        # Docker's mount API omits false ReadOnly (json omitempty); omission is
        # only acceptable for the explicitly writable expected bind.
        $actualReadOnly=if($mount.Contains('ReadOnly')) { if($mount.ReadOnly-isnot[bool]){throw 'Malformed bind read-only value'}; $mount.ReadOnly } else { $false }
        if ($mount.Type-cne'bind'-or$mount.Source-cne$source-or$mount.Target-cne$target-or$actualReadOnly-ne$readOnly-or
            ($null-ne$mount.BindOptions-and$mount.BindOptions.Propagation-notin@($null,'','rprivate'))-or$null-ne$mount.VolumeOptions-or$null-ne$mount.TmpfsOptions-or$null-ne$mount.ImageOptions-or$null-ne$mount.ClusterOptions) { throw 'Created container bind mount differs' }
        if ($null-ne$mount.BindOptions) {
            foreach($option in @('NonRecursive','CreateMountpoint','ReadOnlyNonRecursive','ReadOnlyForceRecursive')) {
                if($mount.BindOptions.Contains($option)-and($mount.BindOptions[$option]-isnot[bool]-or$mount.BindOptions[$option])){throw 'Unexpected bind behavior override'}
            }
        }
    }
}

function Invoke-ContainerProtocol([string[]]$RunArguments,$Context) {
    if ($RunArguments.Count-lt 3-or$RunArguments[0]-cne'run'-or$RunArguments[1]-cne'--rm') { throw 'Expected closed container argument vector' }
    $token=[guid]::NewGuid().ToString('N');$name='pscan-'+$token
    $identifier=$null;$admitted=$false;$removed=$false;$failure=$null;$workload=$null;$exitCode=$null
    try {
        $imageResult=Invoke-HeldDockerCall @('image','inspect','--format','{{.Id}}',$Context.Image) $Context
        if ($imageResult.ExitCode-ne 0-or$imageResult.StdErr-ne''-or$imageResult.StdOut.Trim()-cnotmatch'^sha256:[0-9a-f]{64}$') { throw 'Exact created-image identity unavailable' }
        $imageID=$imageResult.StdOut.Trim()
        $create=@('container','create','--name',$name,'--label',"org.pscan.invocation=$token",'--restart=no')+$RunArguments[2..($RunArguments.Count-1)]
        $created=Invoke-HeldDockerCall $create $Context
        if ($created.ExitCode-ne 0-or$created.StdErr-ne''-or$created.StdOut.Trim()-cnotmatch'^[0-9a-f]{64}$') { throw 'Container create failed or returned uncertain identity; never start' }
        $identifier=$created.StdOut.Trim()
        $inspected=Invoke-HeldDockerCall @('container','inspect','--format','{{json .}}',$identifier) $Context
        if ($inspected.ExitCode-ne 0-or$inspected.StdErr-ne'') { throw 'Created container readback failed' }
        $record=Read-ContainerObject $inspected.StdOut
        Assert-ContainerIdentity $record $identifier $name $token
        $admitted=$true
        Assert-ContainerConfiguration $record $RunArguments $Context.Image $imageID
        if ($record.State.Status-cne'created'-or$record.State.Running-isnot[bool]-or$record.State.Running-or$record.State.Restarting-isnot[bool]-or$record.State.Restarting) { throw 'Container was not inert before its one admitted start' }
        if ($Context.Budget.Remaining-le 0) { throw 'Whole-operation deadline exhausted before container start' }
        $workload=Invoke-HeldDockerCall @('container','start','--attach',$identifier) $Context
        $inspected=Invoke-HeldDockerCall @('container','inspect','--format','{{json .}}',$identifier) $Context
        if ($inspected.ExitCode-ne 0-or$inspected.StdErr-ne'') { throw 'Container exit readback failed' }
        $record=Read-ContainerObject $inspected.StdOut
        Assert-ContainerIdentity $record $identifier $name $token
        if ($record.State.Status-cne'exited'-or$record.State.Running-isnot[bool]-or$record.State.Running-or$record.State.Restarting-isnot[bool]-or$record.State.Restarting-or
            $record.State.OOMKilled-isnot[bool]-or$record.State.OOMKilled-or$record.State.ExitCode-isnot[long]-and$record.State.ExitCode-isnot[int]-or$record.State.ExitCode-lt 0-or$record.State.ExitCode-gt 255) { throw 'Container workload exit is unproved or OOM killed' }
        $exitCode=[int]$record.State.ExitCode
        if ($workload.ExitCode-ne$exitCode) { throw 'Docker client and workload exit codes disagree' }
        $remove=Invoke-HeldDockerCall @('container','rm',$identifier) $Context
        if ($remove.ExitCode-ne 0-or$remove.StdErr-ne''-or$remove.StdOut.Trim()-cne$identifier) { throw 'Owned container removal failed' }
        $removed=$true
        $absence=Invoke-HeldDockerCall @('container','ls','--all','--no-trunc','--filter',"id=$identifier",'--format','{{json .ID}}') $Context
        if ($absence.ExitCode-ne 0-or$absence.StdErr-ne''-or$absence.StdOut.Trim()-cne'') { throw 'Owned container absence unproved' }
        if ($Context.Budget.Remaining-le 0-or$Context.Budget.CleanupStarted) { throw 'Container lifecycle exceeded shared operation deadline' }
    } catch { $failure=$_; $Context.Budget.BeginCleanup() }
    if ($failure) {
        $cleanup='unproved'
        if ($admitted-and!$removed-and$Context.Budget.Remaining-gt 0) {
            try {
                $remove=Invoke-HeldDockerCall @('container','rm','--force',$identifier) $Context
                if ($remove.ExitCode-ne 0-or$remove.StdErr-ne''-or$remove.StdOut.Trim()-cne$identifier) { throw 'Owned terminal removal failed' }
                $absence=Invoke-HeldDockerCall @('container','ls','--all','--no-trunc','--filter',"id=$identifier",'--format','{{json .ID}}') $Context
                if ($absence.ExitCode-ne 0-or$absence.StdErr-ne''-or$absence.StdOut.Trim()-cne''-or$Context.Budget.Remaining-le 0) { throw 'Owned terminal absence unproved' }
                $cleanup='proved'
            } catch { $cleanup='unproved' }
        }
        throw [IO.IOException]::new("Container protocol failed; container=$identifier admitted=$admitted removed=$removed daemon-cleanup=$cleanup; $($failure.Exception.Message)",$failure.Exception)
    }
    # A proved nonzero workload result is returned faithfully for negative tests.
    # Callers still reject it for acquisition/build/package success.
    return [pscustomobject]@{ExitCode=$exitCode;StdOut=$workload.StdOut;StdErr=$workload.StdErr;DaemonEmpty=$true;ContainerID=$identifier}
}
