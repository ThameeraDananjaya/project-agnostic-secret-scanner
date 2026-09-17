$ErrorActionPreference='Stop'
$source=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'unsigned-build-profile-driver.ps1'))
$body=$source.Replace('& /usr/bin/python3 -I -B','Invoke-InertContext')
$body=[regex]::Replace($body,'& "\$PSScriptRoot/([a-z-]+\.ps1)"','Invoke-InertOperation ''$1''')
if($body.Contains('& "$PSScriptRoot/')-or$body.Contains('& /usr/bin/python3')){throw 'Unreplaced external driver operation'}
$saved=@{GITHUB_SHA=$env:GITHUB_SHA;GITHUB_REF=$env:GITHUB_REF;GITHUB_WORKSPACE=$env:GITHUB_WORKSPACE;RUNNER_TEMP=$env:RUNNER_TEMP}
$savedExit=$global:LASTEXITCODE;$cases=0
try {
    $env:GITHUB_SHA='1'*40;$env:GITHUB_REF='refs/tags/release-tooling-v1.0.0-c2-linux-boundary'
    $env:GITHUB_WORKSPACE=Join-Path ([IO.Path]::GetTempPath()) 'inert-repository';$env:RUNNER_TEMP=Join-Path ([IO.Path]::GetTempPath()) 'inert-runner'
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
            $caught='';try{& ([scriptblock]::Create($body)) -Stage $stage -BaselinePath '/inert/baseline' -BaselineSHA256 ('a'*64)}catch{$caught=$_.Exception.Message}
            $expected=switch($failure){'none'{''}'before'{'Named profile context admission failed'}'after'{'Named profile after-context mismatch'}default{'original stage failure'}}
            if($caught-cne$expected-or$script:events[0]-cne'context-before'-or$script:events-notcontains'context-after'-or$script:events-notcontains'facts-after'){throw "Stage failure propagation mismatch: $stage/$failure"}
            $operations=@($script:events|Where-Object{$_-like'operation:*'})
            if($operations.Count-ne$(if($failure-eq'before'){0}else{1})){throw 'Stage dispatch did not remain bounded to one operation'}
            $cases++
        }
    }
} finally {
    foreach($key in $saved.Keys){[Environment]::SetEnvironmentVariable($key,$saved[$key],'Process')}
    $global:LASTEXITCODE=$savedExit
}
Write-Output "Unsigned build driver inert PASS cases=$cases host-calls=0 Docker-calls=0"
