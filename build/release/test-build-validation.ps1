$ErrorActionPreference='Stop'
. (Join-Path $PSScriptRoot 'build-validation.ps1')
$valid=@{
    GITHUB_ACTIONS='true';GITHUB_EVENT_NAME='workflow_dispatch';GITHUB_REPOSITORY='ThameeraDananjaya/project-agnostic-secret-scanner'
    GITHUB_REPOSITORY_OWNER_ID='50274860';GITHUB_REF='refs/heads/main';GITHUB_SHA=('a'*40)
    PSCAN_VALIDATION_WORKFLOW_SHA=('a'*40);PSCAN_VALIDATION_WORKFLOW_REF='ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release-validate-linux.yml@refs/heads/main'
    GITHUB_RUN_ID='123';GITHUB_RUN_ATTEMPT='1'
}
$saved=@{};foreach($key in $valid.Keys){$saved[$key]=[Environment]::GetEnvironmentVariable($key,'Process')}
foreach($key in @('PSCAN_TRUSTED_LAUNCHER_REVISION','PSCAN_TRUSTED_LAUNCHER_TREE')){$saved[$key]=[Environment]::GetEnvironmentVariable($key,'Process')}
$cases=0
function Reject([scriptblock]$Body){$failed=$false;try{&$Body}catch{$failed=$true};if(!$failed){throw 'Validation gate unexpectedly admitted invalid evidence'};$script:cases++}
$base=Join-Path ([IO.Path]::GetTempPath()) ('pscan-validation-mode-'+[guid]::NewGuid().ToString('N'))
[void][IO.Directory]::CreateDirectory($base)
try {
    foreach($key in $valid.Keys){[Environment]::SetEnvironmentVariable($key,$valid[$key],'Process')}
    $identity=Assert-BuildValidationInvocation ('a'*40);$cases++
    foreach($key in $valid.Keys){
        [Environment]::SetEnvironmentVariable($key,'wrong','Process')
        Reject {Assert-BuildValidationInvocation ('a'*40)}
        [Environment]::SetEnvironmentVariable($key,$valid[$key],'Process')
    }
    Reject {Assert-BuildValidationInvocation ('b'*40)}
    $builder=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'build.ps1'))
    $begin=$builder.IndexOf('$validationInvocation=$null');$end=$builder.IndexOf('$resolvedProductTag =',$begin)
    if($begin-lt 0-or$end-le$begin){throw 'Actual mode admission site missing'}
    $modeGate=[scriptblock]::Create($builder.Substring($begin,$end-$begin))
    $Mode='Candidate';$toolingRevision='a'*40
    Reject {&$modeGate}
    $Mode='Validation';. $modeGate
    if($null-ne$toolingTag-or$workflow-cne'.github/workflows/release-validate-linux.yml'-or$workflowRef-cne'refs/heads/main'){throw 'Validation retained a candidate identity'};$cases++
    # Drive the actual whole builder with invalid exact source in both modes.
    # Neither may reach output creation, acquisition, Docker or mode admission.
    $repository=[IO.Path]::GetFullPath((Join-Path $PSScriptRoot '../..'))
    $env:PSCAN_TRUSTED_LAUNCHER_REVISION='0'*40;$env:PSCAN_TRUSTED_LAUNCHER_TREE='1'*40
    foreach($modeName in @('Candidate','Validation')){
        $message='';try{& (Join-Path $PSScriptRoot 'build.ps1') -RepositoryRoot $repository -ExpectedToolingRevision ('0'*40) -AcquisitionDirectory (Join-Path $base 'absent-cache') -OutputDirectory (Join-Path $base 'must-not-exist') -Mode $modeName}catch{$message=$_.Exception.Message}
        if(!$message-or(Test-Path -LiteralPath (Join-Path $base 'must-not-exist'))){throw 'Mode bypassed exact source gate'};$cases++
    }
    $distribution=Join-Path $base 'inert-writer-fixture';[void][IO.Directory]::CreateDirectory($distribution)
    foreach($name in @('project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz','project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip','BUILD-PROVENANCE.json','CHECKSUMS.sha256')){[IO.File]::WriteAllText((Join-Path $distribution $name),'inert writer fixture, not a built package')}
    $sourceTrust=[pscustomobject]@{Commit=('a'*40);Tree=('b'*40)}
    [void](Complete-BuildValidation $distribution $sourceTrust $identity '2026-09-17T00:00:00Z')
    $record=Get-Content -Raw -LiteralPath (Join-Path $distribution 'BUILD-VALIDATION.json')|ConvertFrom-Json -AsHashtable
    if($record.candidate-ne$false-or$record.schema-cne'pscan-build-validation-v1'-or$record.invocation.workflowSha-cne('a'*40)-or$record.invocation.runId-cne'123'-or$record.Contains('releaseIdentity')-or$record.Contains('releaseTooling')-or(Test-Path -LiteralPath (Join-Path $distribution 'release-manifest.json'))){throw 'Validation writer manufactured candidate claims'};$cases++
    Reject {Complete-BuildValidation $distribution $sourceTrust $identity '2026-09-17T00:00:00Z'}
    [IO.File]::WriteAllText((Join-Path $distribution 'release-manifest.json'),'untrusted candidate mixture')
    Reject {Complete-BuildValidation $distribution $sourceTrust $identity '2026-09-17T00:00:00Z'}
    $returnGate=$builder.IndexOf("if(`$Mode-ceq'Validation'){`n    Complete-BuildValidation")
    $trustWriter=$builder.IndexOf("Write-Utf8 (Join-Path `$dist 'release-manifest.json')")
    if($returnGate-lt 0-or$trustWriter-le$returnGate-or!$builder.Substring($returnGate,$trustWriter-$returnGate).Contains("`n    return`n}")){throw 'Validation does not terminate before trust-manifest emission'};$cases++
} finally {
    foreach($key in $saved.Keys){[Environment]::SetEnvironmentVariable($key,$saved[$key],'Process')}
    $resolved=[IO.Path]::GetFullPath($base);$temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if(!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)-or![IO.Path]::GetFileName($resolved).StartsWith('pscan-validation-mode-')){throw 'Unsafe fixture cleanup'}
    Remove-Item -LiteralPath $resolved -Recurse -Force
}
Write-Output "Build validation separation PASS cases=$cases writer-fixtures=INERT candidate-manifests=0 actual-builds=0"
