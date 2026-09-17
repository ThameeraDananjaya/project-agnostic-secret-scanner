$ErrorActionPreference='Stop'
. (Join-Path $PSScriptRoot 'execution-profile.ps1')
$cases=0
function Check([bool]$Value,[string]$Message) { $script:cases++;if(!$Value){throw $Message} }
function Reject([scriptblock]$Body) { $failed=$false;try{&$Body}catch{$failed=$true};Check $failed 'Untrusted profile or receipt accepted' }
function Parse([string]$Name) {
    $tokens=$null;$errors=$null
    $tree=[Management.Automation.Language.Parser]::ParseFile((Join-Path $PSScriptRoot $Name),[ref]$tokens,[ref]$errors)
    if($errors.Count){throw 'Profile caller does not parse'}
    return $tree
}
$boundary=Parse 'docker-execution.ps1';$child=Parse 'invoke-docker-boundary.ps1'
$operationAssignment=$boundary.Find({param($n)$n-is[Management.Automation.Language.AssignmentStatementAst]-and$n.Left.Extent.Text-ceq'$dockerBudgetMilliseconds'},$true)
$transportAssignment=$child.Find({param($n)$n-is[Management.Automation.Language.AssignmentStatementAst]-and$n.Left.Extent.Text-ceq'$childTransportBudgetMilliseconds'},$true)
Check ($null-ne$operationAssignment-and$null-ne$transportAssignment) 'Actual production budget selection missing'
$expected=[ordered]@{EngineInspection=15000;ExactImageInventory=15000;ApprovedImagePull=180000;RepositoryDigestInspection=15000;ContainerCacheProof=15000;ContainerCrlfParse=15000;DependencyAcquisition=300000;ReleaseBuild=900000;ReleasePackage=120000}
foreach($Operation in $expected.Keys) {
    . ([scriptblock]::Create($operationAssignment.Extent.Text))
    . ([scriptblock]::Create($transportAssignment.Extent.Text))
    Check ($dockerBudgetMilliseconds-eq$expected[$Operation]) 'Actual Docker operation selection differs'
    Check ($childTransportBudgetMilliseconds-eq$expected[$Operation]+20000) 'Parent transport and child operation ceilings conflict'
}
foreach($Operation in @('','releasebuild','ReleaseBuild ','NativeFixture','Arbitrary')) { Reject { Get-ReleaseOperationBudget $Operation } }
foreach($tree in @($boundary,$child)) {
    Check (@($tree.ParamBlock.Parameters|Where-Object{$_.Name.VariablePath.UserPath-match'(?i)budget|timeout|duration|profile'}).Count-eq 0) 'Caller-selected timing override exists'
}
$profile=Get-ReleaseExecutionProfile;$profile.operationBudgetMilliseconds.ReleaseBuild=1
Check ((Get-ReleaseOperationBudget 'ReleaseBuild')-eq900000) 'Returned profile mutation leaked into production table'
$saved=$env:PSCAN_DOCKER_BUDGET_MILLISECONDS
try{$env:PSCAN_DOCKER_BUDGET_MILLISECONDS='2147483647';Check ((Get-ReleaseOperationBudget 'ReleaseBuild')-eq900000) 'Environment duration override admitted'}finally{$env:PSCAN_DOCKER_BUDGET_MILLISECONDS=$saved}

$root=Join-Path ([IO.Path]::GetTempPath()) ('pscan-profile-'+[guid]::NewGuid().ToString('N'))
[void][IO.Directory]::CreateDirectory($root);$path=Join-Path $root 'receipt.json'
function GoodReceipt {
    return [ordered]@{schemaVersion='2.0';sourceRevision=('a'*40);image='docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452';dockerExecutableSHA256=('b'*64);pulledDuringAdmission=$false;hostIdentityMode='linux-host-numeric-uid-gid';hostUID=1001;hostGID=1002;containment='empty-after-every-operation';executionProfile=(Get-ReleaseExecutionProfile)}
}
function Save($Receipt) { [IO.File]::WriteAllText($path,($Receipt|ConvertTo-Json -Depth 8),[Text.UTF8Encoding]::new($false)) }
try {
    Save (GoodReceipt);$read=Read-ReleaseImageAdmission $path;Check ($read.executionProfile.operationBudgetMilliseconds.ReleaseBuild-eq900000) 'Exact receipt rejected'
    $receipt=GoodReceipt;$receipt.hostIdentityMode='windows-invoking-host';$receipt.hostUID=$null;$receipt.hostGID=$null;Save $receipt
    Check ((Read-ReleaseImageAdmission $path).hostIdentityMode-ceq'windows-invoking-host') 'Exact Windows receipt rejected'
    foreach($field in @('identifier','streamLimitBytes','cleanupGraceMilliseconds','childTransportAllowanceMilliseconds','maximumProtocolCalls')) {
        $receipt=GoodReceipt;$receipt.executionProfile[$field]='untrusted';Save $receipt;Reject { Read-ReleaseImageAdmission $path }
    }
    foreach($name in $expected.Keys) {
        foreach($bad in @([string]$expected[$name],($expected[$name]+1),($expected[$name]-1))) {
            $receipt=GoodReceipt;$receipt.executionProfile.operationBudgetMilliseconds[$name]=$bad;Save $receipt;Reject { Read-ReleaseImageAdmission $path }
        }
    }
    foreach($change in @('old-version','old-budget','extra-profile','extra-operation','missing-operation','missing-profile','bad-owner','partial-owner','string-pull','wrong-image','wrong-containment')) {
        $receipt=GoodReceipt
        switch($change) {
            'old-version'{$receipt.schemaVersion='1.0'}
            'old-budget'{$receipt.commandBudgetMilliseconds=15000}
            'extra-profile'{$receipt.executionProfile.extra=1}
            'extra-operation'{$receipt.executionProfile.operationBudgetMilliseconds.Unapproved=900000}
            'missing-operation'{$receipt.executionProfile.operationBudgetMilliseconds.Remove('ReleaseBuild')}
            'missing-profile'{$receipt.Remove('executionProfile')}
            'bad-owner'{$receipt.hostUID='1001'}
            'partial-owner'{$receipt.hostGID=$null}
            'string-pull'{$receipt.pulledDuringAdmission='false'}
            'wrong-image'{$receipt.image='golang:latest'}
            'wrong-containment'{$receipt.containment='unproved'}
        }
        Save $receipt;Reject { Read-ReleaseImageAdmission $path }
    }
    $good=GoodReceipt|ConvertTo-Json -Depth 8 -Compress
    foreach($bad in @($good.Replace('"ReleaseBuild":900000','"ReleaseBuild":900000,"ReleaseBuild":15000'),$good.Replace('"schemaVersion":','"SchemaVersion":"2.0","schemaVersion":'),$good.Replace('"schemaVersion":','"schemaversion":'),($good+'{}'),'[]','null',(' '*16385))) {
        [IO.File]::WriteAllText($path,$bad);Reject { Read-ReleaseImageAdmission $path }
    }
    [IO.File]::WriteAllBytes($path,[byte[]]@(0xc3,0x28));Reject { Read-ReleaseImageAdmission $path }
} finally {
    $resolved=[IO.Path]::GetFullPath($root);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)-or![IO.Path]::GetFileName($resolved).StartsWith('pscan-profile-')){throw 'Unsafe owned fixture cleanup'}
    [IO.Directory]::Delete($resolved,$true)
}
$native=$boundary.Find({param($n)$n-is[Management.Automation.Language.StringConstantExpressionAst]-and$n.Value.Contains('public static class PscanNativeBoundary')},$true).Value
[void](Add-Type -TypeDefinition $native)
foreach($milliseconds in $expected.Values) {
    $budget=[PscanOperationBudget]::new($milliseconds,2000,131072)
    Check ($budget.Remaining-gt$milliseconds-1000-and$budget.Remaining-le$milliseconds) 'Actual monotonic budget narrowed selected operation'
}
foreach($args in @(@(900001,2000,131072),@(15000,2001,131072),@(15000,2000,131073),@(0,2000,131072))) { Reject { [PscanOperationBudget]::new($args[0],$args[1],$args[2]) } }
$budget=[PscanOperationBudget]::new(50,50,131072);Start-Sleep -Milliseconds 70
Check ($budget.Remaining-eq0-and!$budget.CleanupStarted) 'Operation clock renewed or granted cleanup as success'
$budget.BeginCleanup();Start-Sleep -Milliseconds 70;$budget.BeginCleanup()
Check ($budget.Remaining-eq0-and$budget.CleanupStarted) 'Cleanup clock renewed'
Write-Output "Fixed execution profile PASS cases=$cases actual-selection=PASS strict-receipt=PASS shared-clock=PASS Docker-calls=0"
