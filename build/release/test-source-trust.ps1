param(
    [Parameter(Mandatory = $true)][string]$SourceRepository,
    [Parameter(Mandatory = $true)][string]$SourceRevision,
    [string]$WorkingDirectory
)

$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'source-trust.ps1')

$sourceRoot = [IO.Path]::GetFullPath((Resolve-Path -LiteralPath $SourceRepository).Path)
[void](Assert-ExactGitSourceTrust -Repository $sourceRoot -ExpectedRevision $SourceRevision)
if ([string]::IsNullOrWhiteSpace($WorkingDirectory)) {
    $WorkingDirectory = Join-Path ([IO.Path]::GetTempPath()) ("pscan-06-source-trust-" + [guid]::NewGuid().ToString('N'))
}
$testRoot = [IO.Path]::GetFullPath($WorkingDirectory)
if (Test-Path -LiteralPath $testRoot) {
    if (Get-ChildItem -LiteralPath $testRoot -Force | Select-Object -First 1) { throw 'Source-trust test directory must be new or empty' }
} else {
    New-Item -ItemType Directory -Path $testRoot | Out-Null
}

$launcher = Join-Path $PSScriptRoot 'invoke-exact-build.ps1'
$utf8 = [Text.UTF8Encoding]::new($false)

function New-Fixture([string]$Name) {
    $path = Join-Path $testRoot $Name
    & git clone --no-hardlinks --no-checkout --quiet $sourceRoot $path
    if ($LASTEXITCODE -ne 0) { throw "Fixture clone failed: $Name" }
    & git -C $path config core.autocrlf false
    & git -C $path config core.eol lf
    & git -C $path -c core.autocrlf=false -c core.eol=lf checkout --quiet --detach $SourceRevision
    if ($LASTEXITCODE -ne 0) { throw "Fixture checkout failed: $Name" }
    return $path
}

function Add-SyntheticLine([string]$Path, [string]$Line) {
    $value = [IO.File]::ReadAllText($Path)
    [IO.File]::WriteAllText($Path, $value + "`n$Line`n", $utf8)
}

function Assert-Rejected([string]$Name, [scriptblock]$Mutate, [switch]$ExpectDriverMarker) {
    $fixture = New-Fixture -Name $Name
    $output = Join-Path $fixture 'release-output'
    $marker = Join-Path $testRoot "$Name-untrusted-driver-ran.txt"
    & $Mutate $fixture $marker
    $log = Join-Path $testRoot "$Name.log"
    & pwsh -NoProfile -File $launcher -RepositoryRoot $fixture -ExpectedToolingRevision $SourceRevision -AcquisitionDirectory $testRoot -OutputDirectory $output -PreflightOnly *> $log
    $exitCode = $LASTEXITCODE
    if ($exitCode -eq 0) { throw "Source-trust adversarial case unexpectedly passed: $Name" }
    if (Test-Path -LiteralPath $output) {
        if (Get-ChildItem -LiteralPath $output -Force | Select-Object -First 1) { throw "Rejected case emitted build output: $Name" }
    }
    if (Test-Path -LiteralPath $marker) { throw "Rejected case executed the untrusted build driver: $Name" }
    Write-Output "SOURCE-TRUST REJECT case=$Name exit=$exitCode output_files=0 untrusted_driver_action=ABSENT"
}

Assert-Rejected 'assume-unchanged' {
    param($fixture,$marker)
    & git -C $fixture update-index --assume-unchanged -- 'build/release/build.ps1'
    Add-SyntheticLine (Join-Path $fixture 'build\release\build.ps1') "[IO.File]::WriteAllText('$($marker.Replace("'","''"))','ran')"
}
Assert-Rejected 'skip-worktree' {
    param($fixture,$marker)
    & git -C $fixture update-index --skip-worktree -- 'build/release/build.ps1'
    Add-SyntheticLine (Join-Path $fixture 'build\release\build.ps1') "[IO.File]::WriteAllText('$($marker.Replace("'","''"))','ran')"
}
Assert-Rejected 'staged-change' {
    param($fixture,$marker)
    Add-SyntheticLine (Join-Path $fixture 'README.md') 'synthetic staged source-trust mutation'
    & git -C $fixture add -- README.md
}
Assert-Rejected 'unstaged-change' {
    param($fixture,$marker)
    Add-SyntheticLine (Join-Path $fixture 'README.md') 'synthetic unstaged source-trust mutation'
}
Assert-Rejected 'untracked-file' {
    param($fixture,$marker)
    [IO.File]::WriteAllText((Join-Path $fixture 'synthetic-untracked.txt'), 'synthetic', $utf8)
}
Assert-Rejected 'fsmonitor-config' {
    param($fixture,$marker)
    & git -C $fixture config core.fsmonitor true
}
Assert-Rejected 'untracked-cache-config' {
    param($fixture,$marker)
    & git -C $fixture config core.untrackedCache true
}
Assert-Rejected 'driver-tampering' {
    param($fixture,$marker)
    $driver = Join-Path $fixture 'build\release\build.ps1'
    $value = [IO.File]::ReadAllText($driver)
    [IO.File]::WriteAllText($driver, "[IO.File]::WriteAllText('$($marker.Replace("'","''"))','ran')`n" + $value, $utf8)
}
Assert-Rejected 'path-mismatch' {
    param($fixture,$marker)
    Remove-Item -LiteralPath (Join-Path $fixture 'README.md') -Force
}
Assert-Rejected 'mode-mismatch' {
    param($fixture,$marker)
    & git -C $fixture update-index --chmod=+x -- README.md
}
Assert-Rejected 'blob-mismatch' {
    param($fixture,$marker)
    $synthetic = Join-Path $testRoot 'synthetic-index-blob.txt'
    [IO.File]::WriteAllText($synthetic, 'synthetic index blob', $utf8)
    $blob = (& git -C $fixture hash-object -w -- $synthetic).Trim()
    & git -C $fixture update-index --cacheinfo "100644,$blob,README.md"
}

$positive = New-Fixture -Name 'positive-exact'
& pwsh -NoProfile -File $launcher -RepositoryRoot $positive -ExpectedToolingRevision $SourceRevision -AcquisitionDirectory $testRoot -OutputDirectory (Join-Path $positive 'release-output') -PreflightOnly
if ($LASTEXITCODE -ne 0) { throw 'Exact clean source-trust positive case failed' }
Write-Output 'SOURCE-TRUST PASS adversarial_rejections=11 exact_positive=PASS zero_build_outputs=PASS untrusted_driver_action=ABSENT'
