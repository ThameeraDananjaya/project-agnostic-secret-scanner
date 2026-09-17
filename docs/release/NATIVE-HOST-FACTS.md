# Read-only host facts for the fixed native fixture

The facts collector supports before/after snapshots around the unchanged
CleanNative fixture; earlier diagnostic revisions used a finally snapshot even
on failure. The current dedicated manual workflow performs only a read-only
named-profile compatibility snapshot, with no fixture or policy helper invocation.
The production boundary, namespace arguments, limits and cleanup are unchanged.

`native-host-facts.py` records the calling PowerShell process's effective UID/GID,
numeric UID/GID maps, capability masks, NoNewPrivs/Seccomp fields and current
AppArmor profile. It also records the two allowlisted image-version variables,
kernel release, util-linux version, AppArmor enabled/user-namespace restriction
flags and relevant namespace availability counters. It never dumps environment
variables, entire process status files, audit messages, source or credentials.

All file reads have byte limits; formats and keys are allowlisted. Missing,
unreadable, malformed and oversized facts remain explicit rather than becoming
zero, false or supported. Each JSON record is bounded to 16 KiB. Fixed trusted
version/journal queries have a five-second deadline, two-second exit allowance
and per-stream capture caps. No namespace is created by the facts collector.

After the fixture, an unprivileged read-only kernel-journal query requests at most
16 AppArmor denial entries naming `unshare` within the fixture time window.
Only allowlisted parsed denial fields are emitted; raw messages are discarded.
The record explicitly says that correlation is by time and command and that
exact fixture PID binding is unavailable. Lack of a visible matching record is
not proof of no denial. Journal warning/stderr, access failure, timeout or output
cap produces an explicit incomplete/unavailable result. There is no sudo fallback.

No profile, sysctl, permission, kernel setting or package is changed. This report
does not preselect AppArmor as the cause or authorize a correction. Its purpose
is to provide actual facts for a source-supported host prerequisite decision.

The compatibility snapshot additionally records the exact 0/1 value or explicit
unavailability of apparmor_restrict_unprivileged_unconfined, and lstat metadata
for the fixed /usr/bin/aa-exec selector and /opt/microsoft/powershell/7/pwsh
candidate. Symlink, regular-file, owner, mode and identity fields are availability
facts only. They do not prove parent-path trust, file capabilities, source hashes,
actual selection, unchanged credentials, namespace support or fixture success.
The selector is not invoked. Existing caller UID/GID, capability, namespace-map
and current-profile facts accompany this snapshot. No privileged fallback runs.

`test-native-host-facts.py` validates the bounded parsers and full before/after
records with inert providers. Subprocess creation is blocked during the record
tests. These tests do not query Linux, execute the fixture or prove host support.
