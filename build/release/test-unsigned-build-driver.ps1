$ErrorActionPreference='Stop'
$source=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'unsigned-build-profile-driver.ps1'))
$body=$source.Replace('& /usr/bin/python3 -I -B','Invoke-InertContext')
$body=[regex]::Replace($body,'& "\$PSScriptRoot/([a-z-]+\.ps1)"','Invoke-InertOperation ''$1''')
if($body.Contains('& "$PSScriptRoot/')-or$body.Contains('& /usr/bin/python3')){throw 'Unreplaced external driver operation'}
$helperPath=(Join-Path $PSScriptRoot 'build-validation.ps1').Replace("'","''")
$body=$body.Replace("(Join-Path `$PSScriptRoot 'build-validation.ps1')", "'$helperPath'")
$saved=@{GITHUB_SHA=$env:GITHUB_SHA;GITHUB_REF=$env:GITHUB_REF;GITHUB_WORKSPACE=$env:GITHUB_WORKSPACE;RUNNER_TEMP=$env:RUNNER_TEMP}
foreach($key in @('PSCAN_VALIDATION_WORKFLOW_SHA','PSCAN_VALIDATION_WORKFLOW_REF','GITHUB_ACTIONS','GITHUB_EVENT_NAME','GITHUB_REPOSITORY','GITHUB_REPOSITORY_OWNER_ID','GITHUB_RUN_ID','GITHUB_RUN_ATTEMPT')){$saved[$key]=[Environment]::GetEnvironmentVariable($key,'Process')}
$savedExit=$global:LASTEXITCODE;$cases=0
try {
    $env:GITHUB_SHA='1'*40;$env:GITHUB_REF='refs/tags/release-tooling-v1.0.0-c2-sbom-v1'
    $env:GITHUB_WORKSPACE=Join-Path ([IO.Path]::GetTempPath()) 'inert-repository';$env:RUNNER_TEMP=Join-Path ([IO.Path]::GetTempPath()) 'inert-runner'
    $env:GITHUB_ACTIONS='true';$env:GITHUB_EVENT_NAME='workflow_dispatch';$env:GITHUB_REPOSITORY='ThameeraDananjaya/project-agnostic-secret-scanner';$env:GITHUB_REPOSITORY_OWNER_ID='50274860';$env:GITHUB_RUN_ID='123';$env:GITHUB_RUN_ATTEMPT='1'
    foreach($modeName in @('Candidate','Validation')){
    $env:GITHUB_REF=if($modeName-eq'Validation'){'refs/heads/main'}else{'refs/tags/release-tooling-v1.0.0-c2-sbom-v1'}
    $env:PSCAN_VALIDATION_WORKFLOW_SHA=if($modeName-eq'Validation'){'1'*40}else{$null}
    $env:PSCAN_VALIDATION_WORKFLOW_REF='ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-validate-linux.yml@refs/heads/main'
    foreach($stage in @('NativePrerequisites','ImageAdmission','ContainerCache','ContainerCRLF','Acquire','BuildA','BuildB','Compare')){
        foreach($failure in @('none','operation','before','after','operation-and-after')){
            $script:failure=$failure;$script:events=[Collections.Generic.List[string]]::new()
            function Invoke-InertContext {
                $kind=if($args[0].EndsWith('unsigned-build-profile-context.py')){'context'}else{'facts'}
                $phase=if($kind-eq'context'){$args[4]}else{$args[2]}
                $script:events.Add("$kind-$phase")
                $global:LASTEXITCODE=if($kind-eq'context'-and(($phase-eq'before'-and$script:failure-eq'before')-or($phase-eq'after'-and$script:failure-in@('after','operation-and-after')))){1}else{0}
            }
            function Invoke-InertOperation {
                $script:events.Add('operation:'+ $args[0])
                if($script:failure-in@('operation','operation-and-after')){throw 'original stage failure'}
            }
            $caught='';try{& ([scriptblock]::Create($body)) -Mode $modeName -Stage $stage -BaselinePath '/inert/baseline' -BaselineSHA256 ('a'*64)}catch{$caught=$_.Exception.Message}
            $expected=switch($failure){'none'{''}'before'{'Named profile context admission failed'}'after'{'Named profile after-context mismatch'}default{'original stage failure'}}
            if($caught-cne$expected-or$script:events[0]-cne'context-before'-or$script:events-notcontains'context-after'-or$script:events-notcontains'facts-after'){throw "Stage failure propagation mismatch: $stage/$failure"}
            $operations=@($script:events|Where-Object{$_-like'operation:*'})
            if($operations.Count-ne$(if($failure-eq'before'){0}else{1})){throw 'Stage dispatch did not remain bounded to one operation'}
            $cases++
        }
    }
    }
    foreach($badRef in @('refs/tags/release-tooling-v1.0.0-c2-linux-build-v2','refs/tags/release-tooling-v1.0.0-c2-linux-boundary','refs/tags/release-tooling-v1.0.0-c2-r6','refs/tags/release-tooling-v1.0.0-c2','refs/tags/release-tooling-v1.0.0-c1','refs/tags/v1.0.0','refs/heads/main')){
        $env:GITHUB_REF=$badRef;$env:PSCAN_VALIDATION_WORKFLOW_SHA=$null
        $script:failure='none';$script:events=[Collections.Generic.List[string]]::new()
        $caught='';try{& ([scriptblock]::Create($body)) -Mode Candidate -Stage BuildA -BaselinePath '/inert/baseline' -BaselineSHA256 ('a'*64)}catch{$caught=$_.Exception.Message}
        if($caught-cne'Unsigned candidate invocation identity is invalid'-or@($script:events|Where-Object{$_-like'operation:*'}).Count){throw 'Old or mutable ref reached candidate operation'}
        $cases++
    }
    $env:GITHUB_REF='refs/tags/release-tooling-v1.0.0-c2-sbom-v1';$env:PSCAN_VALIDATION_WORKFLOW_SHA='1'*40
    $script:events=[Collections.Generic.List[string]]::new()
    $caught='';try{& ([scriptblock]::Create($body)) -Mode Candidate -Stage BuildA -BaselinePath '/inert/baseline' -BaselineSHA256 ('a'*64)}catch{$caught=$_.Exception.Message}
    if($caught-cne'Unsigned candidate invocation identity is invalid'-or@($script:events|Where-Object{$_-like'operation:*'}).Count){throw 'Validation invocation reached candidate operation'}
    $cases++
} finally {
    foreach($key in $saved.Keys){[Environment]::SetEnvironmentVariable($key,$saved[$key],'Process')}
    $global:LASTEXITCODE=$savedExit
}
Write-Output "Unsigned build driver inert PASS cases=$cases host-calls=0 Docker-calls=0"
