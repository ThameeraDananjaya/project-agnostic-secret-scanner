$ErrorActionPreference = 'Stop'
$source=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'native-profile-driver.ps1'))
# Replace only external invocation sites in a memory-only copy. Preserve the
# actual driver's try/catch/finally and error propagation; never run host tools.
$python='& /usr/bin/python3 -I -B'
$fixture='& "$PSScriptRoot/test-docker-execution.ps1" -InternalMode CleanNative'
if (($source.Split($python).Count-1) -ne 4 -or ($source.Split($fixture).Count-1) -ne 1) {
    throw 'Driver invocation sites changed; inert test needs review'
}
$body=$source.Replace($python,'Invoke-InertContext').Replace($fixture,'Invoke-InertFixture')
$cases=@(
    @{Fail='none';Expected='';Calls=5},
    @{Fail='fixture';Expected='original fixture failure';Calls=5},
    @{Fail='before';Expected='Named profile context admission failed';Calls=3},
    @{Fail='after';Expected='Named profile after-context mismatch';Calls=5},
    @{Fail='fixture-and-after';Expected='original fixture failure';Calls=5},
    @{Fail='before-facts';Expected='Before-fixture host facts unavailable';Calls=4},
    @{Fail='after-facts';Expected='After-fixture host facts unavailable';Calls=5}
)
$savedLastExitCode=$global:LASTEXITCODE
try {
foreach($case in $cases) {
    $script:observed=[Collections.Generic.List[string]]::new()
    $script:failure=$case.Fail
    function Invoke-InertContext {
        $kind=if($args[0].EndsWith('native-profile-context.py')){'context'}else{'facts'}
        $phase=if($kind-eq'context'){$args[4]}else{$args[2]}
        $script:observed.Add("$kind-$phase")
        $global:LASTEXITCODE=if(($kind-eq'context'-and$phase-eq'before'-and$script:failure-eq'before')-or
            ($kind-eq'context'-and$phase-eq'after'-and$script:failure-in@('after','fixture-and-after'))-or
            ($kind-eq'facts'-and$script:failure-eq"$phase-facts")){1}else{0}
    }
    function Invoke-InertFixture {
        $script:observed.Add('fixture')
        if($script:failure-in@('fixture','fixture-and-after')){throw 'original fixture failure'}
    }
    $caught=''
    try{& ([scriptblock]::Create($body)) -BaselinePath '/inert/baseline' -BaselineSHA256 ('a'*64)}
    catch{$caught=$_.Exception.Message}
    if($caught-cne$case.Expected-or$script:observed.Count-ne$case.Calls){
        throw "Inert driver case $($case.Fail) failed: $caught; $($script:observed -join ',')"
    }
    if($script:observed[0]-cne'context-before'-or$script:observed-notcontains'context-after'){
        throw 'Driver context ordering changed'
    }
}
} finally { $global:LASTEXITCODE=$savedLastExitCode }
Write-Output 'Named-profile driver inert failure-preservation PASS cases=7 host-calls=0 fixture-calls=0'
