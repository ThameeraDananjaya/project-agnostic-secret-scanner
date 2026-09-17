# Disposable native diagnostic prerequisite

This diagnostic-only candidate creates the fixed root-level **unattached**
AppArmor profile `pscan-native-diagnostic`, ABI4, unconfined mode, with the sole
`userns,` permission. It does not configure the release build workflow or alter
production Docker execution. Runtime success remains unproved.

## Evidence and change of design

The original actual native fixture failed with UID-map EPERM on Ubuntu24.04.
Earlier scoped automatic-attachment candidates never loaded a policy: the
kernel exported 100 existing attachments as `<unknown>`, which cannot establish
disjointness. They remain preserved terminal evidence. The named design does
not participate in automatic executable attachment selection and preserves all
existing opaque metadata exactly; it does not declare those profiles disjoint.

Read-only run35246583859 at abe622b4333b906e9273b301eab83080c0e01086 observed
`apparmor_restrict_unprivileged_unconfined=0`, AppArmor enabled Y, user-namespace
restriction1, an ordinary unconfined caller, aa-exec mode0755, and the runner
PowerShell `/opt/microsoft/powershell/7/pwsh` mode0777. Each new VM must satisfy
live admission again. Never change any global flag to satisfy this prerequisite.

## Existing runtime assumption and exact authority

The release plan lines29-39 and docker-execution.ps1's bound-Docker function
require immutable Docker identity. Its Linux boundary separately requires
root-owned, non-writable setsid/unshare/shell. All those checks stay unchanged.
The release plan lines60-71 requires fresh non-profile PowerShell and newly
compiled exact native source. The Linux CleanNative fixture uses its current
PowerShell path and held descriptor; unlike its Windows branch, it does not
require Linux PowerShell publisher or non-writable-file admission.

An earlier proposal added a non-writable PowerShell condition. This candidate
replaces that diagnostic-only condition with explicit observation of the same
existing runner runtime. It requires the fixed canonical target to equal the
calling process executable, compares inode, metadata and SHA256 before/after,
and rejects set-ID and file-capability privilege. **Mode0777 is not immutable
trust. A held inode does not prevent in-place writes.** Observational equality
cannot exclude transient changes or attest dynamic runtime dependencies. The
runner/runtime remains an existing infrastructure assumption. No file mode,
runtime installation or production trust requirement is changed.

## Closed setup, selection and removal

The root helper's exact committed bytes are admitted before interpretation.
Root executes only fixed setup/removal. It requires unchanged AppArmor/global
flags, root-owned non-writable parents and files for setsid, unshare, aa-exec,
parser and ABI; set-ID/file-capabilities reject. All profile identities, ancestry,
opaque attachments, modes and kernel hashes are retained. Nested policy
namespaces, malformed inventory and any reserved-name collision reject.

The fixed parser compiles with no load, no cache/config inheritance, then adds
only the new profile. No replacement or persistent configuration is permitted.
Exclusive private state binds source, policy, complete prior inventory, host
facts and successful readback kernel digest. Bounded parser streams retain
terminal failures; uncertain add does not establish removal ownership.

The ordinary PowerShell caller snapshots UID/GID tuples, supplementary groups,
capability masks, NoNewPrivs/Seccomp, namespace identities, runtime observation
and exact diagnostic/native script hashes. A size-bounded, SHA256-bound data
baseline cannot supply executable code. It requires an unconfined non-root
caller with zero effective/permitted/inheritable/ambient capabilities.

The fixed `/usr/bin/aa-exec --profile pscan-native-diagnostic --` invocation uses
default change-on-exec, never `--immediate`, into the fixed fresh no-profile
PowerShell driver. Before native compilation, the driver requires identical
credentials, capabilities, security fields, namespaces, executable observation
and source hashes, with the sole exact label
`pscan-native-diagnostic (unconfined)`. Stacked, mixed, missing or unknown labels
reject. The driver runs the unchanged CleanNative fixture, retains its original
failure, and checks context plus bounded host facts afterward.

Profile permission may be inherited by descendants and may be selected by other
permitted local actors while loaded. It is not a new sandbox or per-user/per-call
authorization guarantee. Existing attachment transitions remain effective; no
transition is overridden to force success. The 131072-byte independent stream
caps, 15000ms full native budget, 2000ms cleanup, held-inode execution, stopped
PID-namespace init and complete lifecycle proof remain unchanged.

Always-run cleanup re-admits the same helper and exact owned state/kernel hash,
verifies all unrelated profiles and global/file bindings, removes only the owned
named profile and verifies original state before deleting exact state files.
Policy removal is not proof of fixture process containment. Any namespace,
process, pipe, ownership or cleanup uncertainty remains a terminal failure;
disposable VM destruction is only a backstop. No artifacts upload, Docker build,
signing, publication or successor activation occurs in this diagnostic.

## Validation and references

Inert tests cover named policy/readback, opaque preservation, ownership and
cleanup failures, unsupported flag rejection, exact versus stacked labels,
credentials/namespace/source changes, executable replacement/drift/privilege,
and original fixture failure preservation. No inert test invokes a host policy,
selector, namespace or actual native fixture.

Supported selection behavior:
- https://www.apparmor.net/profiles/profile-types-and-syntax/#unattached-profiles
- https://manpages.ubuntu.com/manpages/noble/man1/aa-exec.1.html
- https://gitlab.com/apparmor/apparmor/-/raw/v4.0.1/binutils/aa_exec.c
- https://gitlab.com/apparmor/apparmor/-/wikis/unprivileged_unconfined_restriction.md

Linuxv6.17 kernel sources explain introspection and label formats; they are not
an attestation of the exact Ubuntu Azure kernel binary:
- https://raw.githubusercontent.com/torvalds/linux/v6.17/security/apparmor/apparmorfs.c
- https://raw.githubusercontent.com/torvalds/linux/v6.17/security/apparmor/label.c
