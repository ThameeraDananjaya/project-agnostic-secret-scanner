$ErrorActionPreference = 'Stop'

$admissionPath = Join-Path $PSScriptRoot 'admit-image.ps1'
$parsersPath = Join-Path $PSScriptRoot 'image-admission.ps1'
$boundaryPath = Join-Path $PSScriptRoot 'docker-execution.ps1'
foreach($path in @($admissionPath,$parsersPath,$boundaryPath,$PSCommandPath)){[void][scriptblock]::Create((Get-Content -Raw -LiteralPath $path))}
. $parsersPath

$exact='docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineDigest='golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$wrongDigest='golang@sha256:1111111111111111111111111111111111111111111111111111111111111111'
$wrongRepository='example.invalid/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineJson='{"Version":"29.7.2","ApiVersion":"1.52","Os":"linux","Arch":"amd64"}'
$listJson='{"Repository":"golang","Digest":"sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452","ID":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}'
$version='EngineInspection';$inventory='ExactImageInventory';$inspect='RepositoryDigestInspection';$pull='ApprovedImagePull'

function Result([int]$ExitCode=0,[string]$StdOut='',[string]$StdErr='',[string]$Terminal=''){[pscustomobject]@{ExitCode=$ExitCode;StdOut=$StdOut;StdErr=$StdErr;Terminal=$Terminal;ContainmentEmpty=([string]::IsNullOrEmpty($Terminal));DockerSHA256=('a'*64)}}
function Equal($actual,$expected,[string]$message){if($actual-cne$expected){throw "$message expected=[$expected] actual=[$actual]"}}
function Assert-Sequence([string[]]$actual,[string[]]$expected,[string]$name){if($actual.Count-ne$expected.Count){throw "$name command count mismatch"};for($i=0;$i-lt$expected.Count;$i++){if($actual[$i]-cne$expected[$i]){throw "$name command order mismatch at $i"}}}
function Success($result,[string]$operation){if($result.Terminal-or!$result.ContainmentEmpty-or$result.ExitCode-ne 0-or![string]::IsNullOrEmpty($result.StdErr)){throw "$operation failed closed"}}
function Model([string]$name,[object[]]$results,[string[]]$expected,[bool]$expectSuccess=$false,[bool]$expectPull=$false){
    $queue=[Collections.Generic.Queue[object]]::new();foreach($r in $results){$queue.Enqueue($r)}
    $events=[Collections.Generic.List[string]]::new();foreach($e in @('host-cache-root','host-cache-downloads','host-cache-modules','host-crlf')){[void]$events.Add($e)}
    $success=$false;$pulled=$false;$later=0
    function Take([string]$operation){[void]$events.Add($operation);if(!$queue.Count){throw 'unexpected operation'};$queue.Dequeue()}
    try{
        $r=Take $version;Success $r $version;[void](Assert-ReleaseEngineEvidence -Json $r.StdOut)
        $r=Take $inventory;Success $r $inventory;$state=Resolve-ReleaseImageListEvidence -Json $r.StdOut
        if($state-eq'ConclusiveAbsent'){$r=Take $pull;Success $r $pull;$pulled=$true}
        $r=Take $inspect;Success $r $inspect;$digests=Read-ReleaseRepoDigestsEvidence -Json $r.StdOut;[void](Assert-ReleaseImageIdentityEvidence -Image $exact -RepoDigests $digests)
        $success=$true;$later++
    }catch{if($expectSuccess){throw "$name unexpectedly rejected: $($_.Exception.Message)"}}
    if(!$expectSuccess-and$success){throw "$name unexpectedly admitted"};Equal $later $(if($expectSuccess){1}else{0}) "$name later action";if($expectSuccess){Equal $pulled $expectPull "$name pull"};Assert-Sequence $events.ToArray() $expected $name
}

$hostEvents=@('host-cache-root','host-cache-downloads','host-cache-modules','host-crlf')
$present=@($hostEvents+@($version,$inventory,$inspect));$absent=@($hostEvents+@($version,$inventory,$pull,$inspect));$versionOnly=@($hostEvents+@($version));$inventoryReached=@($hostEvents+@($version,$inventory))
Model 'pre-existing-exact' @((Result -StdOut $engineJson),(Result -StdOut $listJson),(Result -StdOut ('["'+$engineDigest+'"]'))) $present $true $false
Model 'conclusive-absence-single-pull' @((Result -StdOut $engineJson),(Result),(Result -StdOut 'pulled'),(Result -StdOut ('["'+$engineDigest+'"]'))) $absent $true $true
foreach($case in @(
    @{N='daemon';R=@((Result -ExitCode 1 -StdErr daemon));E=$versionOnly},@{N='permission';R=@((Result -StdOut $engineJson),(Result -ExitCode 1 -StdErr denied));E=$inventoryReached},
    @{N='timeout';R=@((Result -Terminal timeout));E=$versionOnly},@{N='overflow';R=@((Result -Terminal overflow));E=$versionOnly},
    @{N='malformed-engine';R=@((Result -StdOut '{'));E=$versionOnly},@{N='malformed-list';R=@((Result -StdOut $engineJson),(Result -StdOut '{'));E=$inventoryReached},
    @{N='ambiguous-list';R=@((Result -StdOut $engineJson),(Result -StdOut ($listJson+"`n"+$listJson)));E=$inventoryReached},
    @{N='pull-failure';R=@((Result -StdOut $engineJson),(Result),(Result -ExitCode 1 -StdErr failed));E=@($hostEvents+@($version,$inventory,$pull))}
)){Model $case.N $case.R $case.E}
foreach($case in @(
    @{N='empty';J=''},@{N='null';J='null'},@{N='scalar';J=('"'+$engineDigest+'"')},@{N='empty-array';J='[]'},@{N='malformed';J='['},
    @{N='mixed';J=('["'+$engineDigest+'","'+$wrongDigest+'"]')},@{N='duplicate';J=('["'+$engineDigest+'","'+$engineDigest+'"]')},
    @{N='alias';J=('["'+$exact+'"]')},@{N='wrong-repository';J=('["'+$wrongRepository+'"]')},@{N='wrong-digest';J=('["'+$wrongDigest+'"]')}
)){Model $case.N @((Result -StdOut $engineJson),(Result -StdOut $listJson),(Result -StdOut $case.J)) $present}

$admission=Get-Content -Raw $admissionPath;$parsers=Get-Content -Raw $parsersPath;$boundary=Get-Content -Raw $boundaryPath
if($parsers-match'Diagnostics\.Process|Get-AuthenticodeSignature|/usr/bin/docker'){throw 'Pure admission parser retains an execution boundary'}
if(($boundary.Split("'ApprovedImagePull'",[StringSplitOptions]::None).Count-1)-lt 2){throw 'Closed operation table omits the one approved pull operation'}
foreach($forbidden in @('ScriptBlock','Callback','Invoker','ExecutablePath')){if($admission.Substring(0,$admission.IndexOf('$ErrorActionPreference')).Contains($forbidden)){throw "Admission exposes $forbidden"}}

& (Join-Path $PSScriptRoot 'test-docker-execution.ps1')
Write-Output 'Image admission iteration-004 PASS state-matrix=PASS parser-only=PASS closed-docker-boundary=PASS'
