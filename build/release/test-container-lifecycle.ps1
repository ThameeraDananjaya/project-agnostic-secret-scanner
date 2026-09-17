$ErrorActionPreference='Stop'
. (Join-Path $PSScriptRoot 'docker-container-lifecycle.ps1')
$tokens=$null;$errors=$null
$tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot 'docker-execution.ps1'),[ref]$tokens,[ref]$errors)
if($errors.Count){throw 'Boundary does not parse'}
$native=$tree.Find({param($n)$n-is[Management.Automation.Language.StringConstantExpressionAst]-and$n.Value.Contains('public static class PscanNativeBoundary')},$true).Value
[void](Add-Type -TypeDefinition $native)
$image='docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$imageID='sha256:'+('a'*64);$ownedID='b'*64
$vector=@('run','--rm','--pull=never','--network','none','--read-only','--cap-drop','ALL','--security-opt','no-new-privileges','--pids-limit','32','--memory','128m','--memory-swap','128m','--cpus','1','--user','1001:1002','--tmpfs','/work:rw,noexec,nosuid,nodev,size=16m,mode=0700,uid=1001,gid=1002','--mount','type=bind,src=/owned/cache,dst=/gomodcache','--workdir','/work',$image,'/bin/sh','-ceu','fixed inert payload')
function Record([string]$State) {
    @{
        Id=$ownedID;Name='/'+$script:ownedName;Image=$imageID
        Config=@{Image=$image;Labels=@{'org.pscan.invocation'=$script:ownedToken};User='1001:1002';WorkingDir='/work';Tty=$false;OpenStdin=$false;Entrypoint=$null;Volumes=$null;Cmd=@('/bin/sh','-ceu','fixed inert payload')}
        HostConfig=@{
            Privileged=$false;ReadonlyRootfs=$true;AutoRemove=$false;PublishAllPorts=$false;RestartPolicy=@{Name='no';MaximumRetryCount=0}
            PidMode='';UsernsMode='';IpcMode='private';CapAdd=$null;CapDrop=@('ALL');SecurityOpt=@('no-new-privileges');NetworkMode='none'
            Memory=134217728;MemorySwap=134217728;PidsLimit=32;NanoCpus=1000000000
            Tmpfs=@{'/work'='rw,noexec,nosuid,nodev,size=16m,mode=0700,uid=1001,gid=1002'}
            Mounts=@(@{Type='bind';Source='/owned/cache';Target='/gomodcache';BindOptions=@{}})
        }
        State=@{Status=$State;Running=$false;Restarting=$false;OOMKilled=$false;ExitCode=$(if($script:scenario-eq'nonzero'){7}else{0})}
    }
}
function Invoke-HeldDockerCall([string[]]$Arguments,$Context) {
    # Only this test scope replaces I/O. Production state-machine functions above
    # execute unchanged, including identity/config admission, ordering and cleanup.
    if($Context.Budget.Remaining-le 0){throw 'inert shared deadline exhausted'}
    $verb=$Arguments[1];$script:events.Add($Arguments -join ' ')
    $stage=if($verb-eq'inspect'-and$Arguments[0]-eq'container'){if($script:started){'exit-inspect'}else{'created-inspect'}}elseif($verb-eq'ls'){'absence'}else{$verb}
    if($script:scenario-eq('fail-'+$stage)){
        $Context.Budget.BeginCleanup()
        if($stage-eq'create'){$script:lateCreated=$true}
        throw "inert $stage failure"
    }
    $stdout='';$stderr='';$exitCode=0
    switch($verb){
        'inspect' {
            if($Arguments[0]-eq'image'){$stdout=$imageID+"`n";break}
            if($Arguments[-1]-cne$ownedID){throw 'Inspection escaped exact owned identity'}
            $record=Record $(if($script:started){'exited'}else{'created'})
            switch($script:scenario){
                'foreign-id' {$record.Id='c'*64}
                'foreign-name' {$record.Name='/unrelated'}
                'foreign-label' {$record.Config.Labels['org.pscan.invocation']='foreign'}
                'wrong-image' {$record.Image='sha256:'+('c'*64)}
                'privileged' {$record.HostConfig.Privileged=$true}
                'writable-root' {$record.HostConfig.ReadonlyRootfs=$false}
                'extra-cap' {$record.HostConfig.CapAdd=@('SYS_ADMIN')}
                'missing-cap-drop' {$record.HostConfig.CapDrop=@()}
                'missing-nnp' {$record.HostConfig.SecurityOpt=@()}
                'host-network' {$record.HostConfig.NetworkMode='host'}
                'host-pid' {$record.HostConfig.PidMode='host'}
                'memory' {$record.HostConfig.Memory=268435456}
                'cpu' {$record.HostConfig.NanoCpus=2000000000}
                'process-limit' {$record.HostConfig.PidsLimit=128}
                'wrong-user' {$record.Config.User='0:0'}
                'workdir' {$record.Config.WorkingDir='/'}
                'extra-command' {$record.Config.Cmd+=@('extra')}
                'wrong-bind' {$record.HostConfig.Mounts[0].Source='/unrelated'}
                'bind-override' {$record.HostConfig.Mounts[0].BindOptions.CreateMountpoint=$true}
                'wrong-tmpfs' {$record.HostConfig.Tmpfs['/work']='rw'}
                'restart' {$record.HostConfig.RestartPolicy.Name='always'}
                'running-before-start' {if(!$script:started){$record.State.Running=$true}}
                'oom' {if($script:started){$record.State.OOMKilled=$true}}
                'still-running' {if($script:started){$record.State.Running=$true;$record.State.Status='running'}}
                'exit-mismatch' {if($script:started){$record.State.ExitCode=7}}
                'malformed-exit' {if($script:started){$record.State.ExitCode='zero'}}
            }
            $stdout=($record|ConvertTo-Json -Depth 12 -Compress)+"`n"
            if($script:scenario-eq'duplicate-json'){$stdout=$stdout.Replace('"Id":',('"Id":"'+$ownedID+'","Id":'))}
        }
        'create' {
            if($Arguments-contains'--rm'-or$Arguments-notcontains'--restart=no'){throw 'Create could auto-remove or restart'}
            $script:ownedName=$Arguments[$Arguments.IndexOf('--name')+1]
            $script:ownedToken=$Arguments[$Arguments.IndexOf('--label')+1].Substring('org.pscan.invocation='.Length)
            if($script:ownedName-cne('pscan-'+$script:ownedToken)){throw 'Internal ownership binding differs'}
            $stdout=if($script:scenario-eq'malformed-create'){'not-an-id'}else{$ownedID+"`n"}
        }
        'start' {
            if($Arguments.Count-ne 4-or$Arguments[-1]-cne$ownedID-or$script:started){throw 'Start escaped exact single owned target'}
            $script:started=$true;$stdout='workload-output';$exitCode=if($script:scenario-eq'nonzero'){7}else{0}
            if($script:scenario-eq'stdout-total'){$stdout='x'*131072}
            if($script:scenario-eq'deadline'){Start-Sleep -Milliseconds 180}
            if($script:scenario-eq'exhaust-cleanup'){$Context.Budget.BeginCleanup();Start-Sleep -Milliseconds 90;throw 'native cleanup consumed the one grace'}
        }
        'rm' {if($Arguments[-1]-cne$ownedID){throw 'Removal escaped exact owned ID'};$stdout=$ownedID+"`n";$script:removed=$true}
        'ls' {
            if($Arguments-notcontains"id=$ownedID"-or$Arguments-notcontains'--all'-or$Arguments-notcontains'--no-trunc'){throw 'Absence query broadened'}
            if($script:scenario-eq'absence-nonempty'){$stdout='"'+$ownedID+'"'}
        }
        default{throw 'Unexpected inert protocol command'}
    }
    if(!$Context.Budget.Consume([Text.Encoding]::UTF8.GetByteCount($stdout),$false)-or!$Context.Budget.Consume([Text.Encoding]::UTF8.GetByteCount($stderr),$true)){$Context.Budget.BeginCleanup();throw 'inert aggregate stream overflow'}
    return [pscustomobject]@{ExitCode=$exitCode;StdOut=$stdout;StdErr=$stderr}
}
$cases=0
foreach($scenario in @('success','nonzero','fail-inspect','fail-create','fail-created-inspect','fail-start','fail-exit-inspect','fail-rm','fail-absence',
    'foreign-id','foreign-name','foreign-label','wrong-image','privileged','writable-root','extra-cap','missing-cap-drop','missing-nnp','host-network','host-pid','memory','cpu','process-limit','wrong-user','workdir','extra-command','wrong-bind','bind-override','wrong-tmpfs','restart','running-before-start','oom','still-running','exit-mismatch','malformed-exit','duplicate-json','malformed-create','stdout-total','deadline','exhaust-cleanup','absence-nonempty')) {
    $script:scenario=$scenario;$script:events=[Collections.Generic.List[string]]::new();$script:started=$false;$script:removed=$false;$script:lateCreated=$false
    $script:ownedName='';$script:ownedToken=''
    $operationMs=if($scenario-eq'deadline'){150}else{15000};$cleanupMs=if($scenario-eq'exhaust-cleanup'){50}else{2000}
    $context=@{Image=$image;Budget=[PscanOperationBudget]::new($operationMs,$cleanupMs,131072)}
    $passed=$false;$message='';try{$result=Invoke-ContainerProtocol $vector $context;$passed=$result.DaemonEmpty}catch{$message=$_.Exception.Message}
    if($passed-ne($scenario-in@('success','nonzero'))){throw "Protocol case mismatch: $scenario $message"}
    if($passed-and(!$script:removed-or$result.ExitCode-ne$(if($scenario-eq'nonzero'){7}else{0}))){throw 'Workload outcome or removal was not preserved'}
    if($scenario-in@('malformed-create','foreign-id','foreign-name','foreign-label','fail-create','fail-created-inspect','duplicate-json')-and($script:started-or$script:removed)){throw 'Unadmitted object was started or removed'}
    if($scenario-eq'exhaust-cleanup'-and$script:removed){throw 'Cleanup obtained a fresh grace after native exhaustion'}
    $cases++
}
# Exercise the actual shared capture counters, independent streams and immutable
# cleanup start; these are not modeled by the protocol I/O fixture.
$budget=[PscanOperationBudget]::new(15000,100,8)
$first=[PscanNativeBoundary]::CaptureShared([IO.MemoryStream]::new([byte[]]@(1,2,3,4,5)),131072,$budget,$false).GetAwaiter().GetResult()
$second=[PscanNativeBoundary]::CaptureShared([IO.MemoryStream]::new([byte[]]@(1,2,3,4,5,6,7,8)),131072,$budget,$true).GetAwaiter().GetResult()
$rejected=$false;try{[void][PscanNativeBoundary]::CaptureShared([IO.MemoryStream]::new([byte[]]@(1,2,3,4)),131072,$budget,$false).GetAwaiter().GetResult()}catch{$rejected=$true}
if(!$rejected-or$first.Length-ne 5-or$second.Length-ne 8){throw 'Actual per-stream aggregate capture failed'}
$budget.BeginCleanup();Start-Sleep -Milliseconds 120;$budget.BeginCleanup();if($budget.Remaining-ne 0){throw 'Cleanup deadline was renewed'}
Write-Output "Container lifecycle inert PASS scenarios=$cases shared-capture=PASS shared-cleanup-clock=PASS Docker-calls=0 daemon-runtime-proof=UNAVAILABLE"
