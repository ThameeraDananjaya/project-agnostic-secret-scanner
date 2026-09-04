$ErrorActionPreference = 'Stop'
. (Join-Path $PSScriptRoot 'image-admission.ps1')

$exact = 'docker.io/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
$engineDigest = 'golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452'
if ((Assert-ReleaseImageIdentityEvidence -Image $exact -RepoDigests @($engineDigest)) -ne $exact) {
    throw 'Exact canonical image evidence was not admitted'
}

$invalid = @(
    [ordered]@{Name='short-name';Image=$engineDigest;RepoDigests=@($engineDigest)},
    [ordered]@{Name='mutable-tag';Image='docker.io/library/golang:1.27.1';RepoDigests=@($engineDigest)},
    [ordered]@{Name='wrong-repository';Image=$exact;RepoDigests=@('example.invalid/library/golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452')},
    [ordered]@{Name='wrong-digest';Image=$exact;RepoDigests=@('golang@sha256:1111111111111111111111111111111111111111111111111111111111111111')},
    [ordered]@{Name='absent-evidence';Image=$exact;RepoDigests=@()},
    [ordered]@{Name='ambiguous-duplicate';Image=$exact;RepoDigests=@($engineDigest,$engineDigest)}
)
foreach ($case in $invalid) {
    $rejected = $false
    try { [void](Assert-ReleaseImageIdentityEvidence -Image $case.Image -RepoDigests $case.RepoDigests) } catch { $rejected = $true }
    if (!$rejected) { throw "$($case.Name) image evidence unexpectedly passed" }
}

Write-Output 'Image admission unit PASS canonical=PASS short-name=REJECT mutable-tag=REJECT wrong-repository=REJECT wrong-digest=REJECT absent=REJECT ambiguous=REJECT'
