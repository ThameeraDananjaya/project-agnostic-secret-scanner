param(
    [Parameter(Mandatory = $true)][string]$SourceRoot,
    [Parameter(Mandatory = $true)][string]$GoExecutable,
    [Parameter(Mandatory = $true)][string]$Destination
)

$ErrorActionPreference = 'Stop'
New-Item -ItemType Directory -Force -Path $Destination | Out-Null
Push-Location $SourceRoot
try {
    $lines = & $GoExecutable list -deps -f '{{with .Module}}{{if not .Main}}{{.Path}}|{{.Version}}|{{.Dir}}{{end}}{{end}}' .
    if ($LASTEXITCODE -ne 0) { throw 'Cannot enumerate compiled module dependency set' }
} finally {
    Pop-Location
}

$modules = @()
foreach ($line in ($lines | Sort-Object -Unique)) {
    $parts = $line -split '\|', 3
    if ($parts.Count -ne 3 -or [string]::IsNullOrWhiteSpace($parts[2])) { throw "Incomplete module identity: $line" }
    $licenseFiles = Get-ChildItem -LiteralPath $parts[2] -File | Where-Object Name -Match '^(LICENSE|LICENCE|COPYING|NOTICE)(\.|$|-)'
    if ($licenseFiles.Count -eq 0) { throw "No licence evidence for $($parts[0])@$($parts[1])" }
    $slug = ($parts[0] + '@' + $parts[1]) -replace '[^A-Za-z0-9._@+-]', '_'
    $moduleDestination = Join-Path $Destination $slug
    New-Item -ItemType Directory -Force -Path $moduleDestination | Out-Null
    $licenses = @()
    foreach ($license in ($licenseFiles | Sort-Object Name)) {
        $target = Join-Path $moduleDestination $license.Name
        if (Test-Path -LiteralPath $target) { Set-ItemProperty -LiteralPath $target -Name IsReadOnly -Value $false }
        $licenseText = [IO.File]::ReadAllText($license.FullName).Replace("`r`n", "`n")
        [IO.File]::WriteAllText($target, $licenseText, [Text.UTF8Encoding]::new($false))
        Set-ItemProperty -LiteralPath $target -Name IsReadOnly -Value $false
        $licenses += [ordered]@{
            file = $license.Name
            sha256 = (Get-FileHash -LiteralPath $target -Algorithm SHA256).Hash.ToLowerInvariant()
        }
    }
    $modules += [ordered]@{ module = $parts[0]; version = $parts[1]; licenses = $licenses }
}

$manifest = [ordered]@{
    schemaVersion = '1.0'
    scope = 'LF-normalized licence texts for modules compiled by github.com/zricethezav/gitleaks/v8 at 83d9cd684c87d95d656c1458ef04895a7f1cbd8e'
    moduleCount = $modules.Count
    modules = $modules
}
$json = ($manifest | ConvertTo-Json -Depth 8).Replace("`r`n", "`n")
[IO.File]::WriteAllText((Join-Path $Destination 'manifest.json'), $json + "`n", [Text.UTF8Encoding]::new($false))
