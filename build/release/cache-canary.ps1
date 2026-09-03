. (Join-Path $PSScriptRoot 'shell-payload.ps1')

function Invoke-ReleaseCacheCanary {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory = $true)][string]$Image,
        [Parameter(Mandatory = $true)][string]$ModuleCache,
        [Nullable[int]]$HostUID,
        [Nullable[int]]$HostGID,
        [switch]$ReadOnlyCache
    )

    $resolvedCache = (Resolve-Path -LiteralPath $ModuleCache).Path
    $tmpfs = '/work:rw,noexec,nosuid,nodev,size=16m,mode=0700'
    $arguments = @(
        'run', '--rm', '--pull=never', '--network', 'none', '--read-only',
        '--cap-drop', 'ALL', '--security-opt', 'no-new-privileges',
        '--pids-limit', '32', '--memory', '128m', '--memory-swap', '128m', '--cpus', '1'
    )
    if ($HostUID.HasValue -or $HostGID.HasValue) {
        if (!$HostUID.HasValue -or !$HostGID.HasValue -or $HostUID.Value -lt 0 -or $HostGID.Value -lt 0) {
            throw 'Cache canary requires a complete non-negative numeric host UID/GID pair'
        }
        $arguments += @('--user', "$($HostUID.Value):$($HostGID.Value)")
        $tmpfs += ",uid=$($HostUID.Value),gid=$($HostGID.Value)"
    }
    $arguments += @('--tmpfs', $tmpfs)
    $cacheMount = "type=bind,src=$resolvedCache,dst=/gomodcache"
    if ($ReadOnlyCache) { $cacheMount += ',readonly' }
    $arguments += @(
        '--mount', $cacheMount,
        '--workdir', '/work',
        $Image,
        '/usr/bin/env', '-i', 'PATH=/usr/bin:/bin', 'HOME=/work',
        '/bin/sh', '-ceu', @'
umask 077
canary=/gomodcache/.pscan-cache-canary-$$
renamed=/gomodcache/.pscan-cache-canary-ready-$$
cleanup() { rm -f "$canary" "$renamed"; }
trap cleanup EXIT HUP INT TERM
test ! -e "$canary"
test ! -e "$renamed"
printf '%s\n' 'PSCAN-06-C1-CACHE-CANARY' > "$canary"
mv "$canary" "$renamed"
test "$(cat "$renamed")" = 'PSCAN-06-C1-CACHE-CANARY'
rm "$renamed"
test ! -e "$canary"
test ! -e "$renamed"
trap - EXIT HUP INT TERM
'@
    )

    $arguments[$arguments.Count - 1] = ConvertTo-LFPosixShellPayload -Payload $arguments[$arguments.Count - 1]
    Assert-LFPosixShellPayload -Payload $arguments[$arguments.Count - 1]
    & docker @arguments
    if ($LASTEXITCODE -ne 0) {
        throw 'Module-cache write, atomic rename, read and delete canary failed before dependency acquisition'
    }
}
