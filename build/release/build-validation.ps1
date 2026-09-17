function Assert-BuildValidationInvocation([string]$Revision) {
    $repository='ThameeraDananjaya/project-agnostic-secret-scanner'
    $workflow='.github/workflows/release-validate-linux.yml'
    $ref='refs/heads/main'
    if($Revision-cnotmatch'^[0-9a-f]{40}$'-or$env:GITHUB_ACTIONS-cne'true'-or$env:GITHUB_EVENT_NAME-cne'workflow_dispatch'-or
        $env:GITHUB_REPOSITORY-cne$repository-or$env:GITHUB_REPOSITORY_OWNER_ID-cne'50274860'-or$env:GITHUB_REF-cne$ref-or$env:GITHUB_SHA-cne$Revision-or
        $env:PSCAN_VALIDATION_WORKFLOW_SHA-cne$Revision-or$env:PSCAN_VALIDATION_WORKFLOW_REF-cne"$repository/$workflow@$ref"-or
        $env:GITHUB_RUN_ID-cnotmatch'^[1-9][0-9]{0,19}$'-or$env:GITHUB_RUN_ATTEMPT-cnotmatch'^[1-9][0-9]{0,5}$') { throw 'Exact validation workflow invocation is unproved' }
    return [ordered]@{repository=$repository;repositoryOwnerId='50274860';workflow=$workflow;ref=$ref;workflowSha=$Revision;trigger='workflow_dispatch';runId=$env:GITHUB_RUN_ID;runAttempt=$env:GITHUB_RUN_ATTEMPT}
}

function Complete-BuildValidation([string]$Distribution,$SourceTrust,$Invocation,[string]$Created) {
    if(Test-Path -LiteralPath (Join-Path $Distribution 'release-manifest.json')){throw 'Validation directory contains a candidate trust manifest'}
    $expected=Assert-BuildValidationInvocation $SourceTrust.Commit
    foreach($key in $expected.Keys){if($Invocation[$key]-cne$expected[$key]){throw 'Validation invocation changed before output'}}
    $files=@(Get-ChildItem -LiteralPath $Distribution -Force|Sort-Object Name)
    $assets=@(foreach($file in $files){
        if($file.PSIsContainer-or($file.Attributes-band[IO.FileAttributes]::ReparsePoint)-ne 0){throw 'Validation output is not a regular file set'}
        [ordered]@{name=$file.Name;bytes=$file.Length;sha256=(Get-FileHash -LiteralPath $file.FullName -Algorithm SHA256).Hash.ToLowerInvariant()}
    })
    foreach($required in @('project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz','project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip','BUILD-PROVENANCE.json','CHECKSUMS.sha256')){
        if($assets.name-cnotcontains$required){throw 'Validation omitted required actual build output'}
    }
    $record=[ordered]@{
        schema='pscan-build-validation-v1';purpose='validation-only';candidate=$false;createdAt=$Created
        source=[ordered]@{commit=$SourceTrust.Commit;tree=$SourceTrust.Tree}
        productSource=[ordered]@{tag='v1.0.0';commit='a13c28fe7273bc8dc6545f97966a02889524eb4c';tree='217b711ddea51fd0ea7e808edd2e27fdecef8427'}
        invocation=$Invocation;outputs=$assets
    }
    $path=Join-Path $Distribution 'BUILD-VALIDATION.json'
    $stream=[IO.File]::Open($path,[IO.FileMode]::CreateNew,[IO.FileAccess]::Write,[IO.FileShare]::None)
    try{$bytes=[Text.UTF8Encoding]::new($false).GetBytes(($record|ConvertTo-Json -Depth 12)+"`n");$stream.Write($bytes,0,$bytes.Length)}finally{$stream.Dispose()}
    Write-Output "Actual build validation complete; candidate=false source=$($SourceTrust.Commit) outputs=$($assets.Count)"
}
