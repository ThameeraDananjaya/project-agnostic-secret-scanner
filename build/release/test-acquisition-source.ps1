$ErrorActionPreference='Stop'
$source=[IO.File]::ReadAllText((Join-Path $PSScriptRoot 'acquire.ps1'))
$end=$source.IndexOf('$cache = [IO.Path]::GetFullPath($CacheDirectory)')
if ($end-lt 0) { throw 'Actual pre-acquisition source gate is missing' }
# Execute the exact parameter/source-admission prefix only. No cache, Docker or
# network instruction is included; test success is source proof only.
$prefix=$source.Substring(0,$end)
$sourceTrustPath=(Join-Path $PSScriptRoot 'source-trust.ps1').Replace("'","''")
# A dynamically created scriptblock has no PSScriptRoot. Bind only that literal
# include location to the same actual file; all source checks remain unchanged.
$prefix=$prefix.Replace("(Join-Path `$PSScriptRoot 'source-trust.ps1')", "'$sourceTrustPath'")
$gate=[scriptblock]::Create($prefix+"`n`$sourceTrust")
$base=Join-Path ([IO.Path]::GetTempPath()) ('pscan-acquisition-source-'+[guid]::NewGuid().ToString('N'))
$repository=Join-Path $base 'repository'; [void][IO.Directory]::CreateDirectory($repository)
$savedRevision=$env:PSCAN_TRUSTED_LAUNCHER_REVISION
$savedRoot=$env:PSCAN_VERIFIED_REPOSITORY_ROOT
$cases=0
function Invoke-FixtureGit([string[]]$Arguments) {
    $result=& git -C $repository @Arguments
    if ($LASTEXITCODE-ne 0) { throw 'Owned source fixture Git operation failed' }
    return $result
}
function Prove([string]$Root,[string]$Revision,[bool]$Expected) {
    $passed=$false
    try { $proof=& $gate -CacheDirectory (Join-Path $base 'must-not-exist') -SourceRepository $Root -SourceRevision $Revision; $passed=$proof.RawEqualCount-eq 2-and$proof.WorkingTreeInputsTrusted }
    catch { if ($Expected) { throw } }
    if ($passed-ne$Expected-or(Test-Path -LiteralPath (Join-Path $base 'must-not-exist'))) { throw 'Pre-acquisition admission result or side-effect mismatch' }
    $script:cases++
}
try {
    [void](Invoke-FixtureGit @('init','--quiet'))
    [void](Invoke-FixtureGit @('config','core.autocrlf','false'))
    [IO.File]::WriteAllText((Join-Path $repository 'go.mod'),"module example.invalid/fixture`n`ngo 1.27.0`n",[Text.UTF8Encoding]::new($false))
    [IO.File]::WriteAllText((Join-Path $repository 'go.sum'),"fixture source binding only`n",[Text.UTF8Encoding]::new($false))
    [void](Invoke-FixtureGit @('add','--','go.mod','go.sum'))
    [void](Invoke-FixtureGit @('-c','user.name=PSCAN local fixture','-c','user.email=fixture@example.invalid','commit','--quiet','-m','owned source admission fixture'))
    $revision=(Invoke-FixtureGit @('rev-parse','HEAD')).Trim()
    $env:PSCAN_TRUSTED_LAUNCHER_REVISION=$revision
    $env:PSCAN_VERIFIED_REPOSITORY_ROOT=Join-Path $base 'unverified-environment-path'
    Prove $repository $revision $true
    Prove $repository ('0'*40) $false
    $entrypointOnly=Join-Path $base 'entrypoint-only/build/release'; [void][IO.Directory]::CreateDirectory($entrypointOnly)
    Prove $entrypointOnly $revision $false
    foreach ($name in @('go.mod','go.sum')) {
        $path=Join-Path $repository $name; $original=[IO.File]::ReadAllBytes($path)
        [IO.File]::AppendAllText($path,"uncommitted drift`n")
        Prove $repository $revision $false
        [IO.File]::WriteAllBytes($path,$original)
    }
    $sum=Join-Path $repository 'go.sum'; $original=[IO.File]::ReadAllBytes($sum)
    [IO.File]::Delete($sum); Prove $repository $revision $false; [IO.File]::WriteAllBytes($sum,$original)
    $untracked=Join-Path $repository 'untracked'; [IO.File]::WriteAllText($untracked,'not admitted')
    Prove $repository $revision $false; [IO.File]::Delete($untracked)
    Prove $repository $revision $true
} finally {
    $env:PSCAN_TRUSTED_LAUNCHER_REVISION=$savedRevision; $env:PSCAN_VERIFIED_REPOSITORY_ROOT=$savedRoot
    $resolved=[IO.Path]::GetFullPath($base); $temp=[IO.Path]::GetFullPath([IO.Path]::GetTempPath())
    if (!$resolved.StartsWith($temp,[StringComparison]::OrdinalIgnoreCase)-or![IO.Path]::GetFileName($resolved).StartsWith('pscan-acquisition-source-')) { throw 'Unsafe fixture cleanup' }
    Remove-Item -LiteralPath $resolved -Recurse -Force
}
Write-Output "Acquisition source gate PASS cases=$cases cache-mutations=0 Docker-calls=0 network-calls=0"
