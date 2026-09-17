# Scoped native-fixture host prerequisite

## Evidence and purpose

Diagnostic run 35239788149 at ef24f4ed593246cf35148e490a1afd2194a1b6bb
failed before the first native fixture executed. Its 66-byte stderr reports an
EPERM writing the UID map. Ubuntu 24.04.5, kernel 6.17.0-1022-azure, had AppArmor
enabled and unprivileged-user-namespace restriction set to 1. A time-correlated
kernel record denies `sys_admin` for `unshare` in `unprivileged_userns`; exact
fixture PID correlation is unavailable. No package was built.

Ubuntu's [24.04 release notes](https://documentation.ubuntu.com/release-notes/24.04/)
describe an application profile with `flags=(unconfined)` and `userns,` for
applications requiring namespaces. This proposal uses that supported mechanism
for the fixed trusted `/usr/bin/unshare` launcher on one disposable diagnostic VM.
It has not yet been executed or admitted to a release build.

## Exact permission scope

The helper contains a fixed ASCII policy, ABI 4.0, profile name
`pscan-native-diagnostic-unshare`, exact attachment `/usr/bin/unshare`,
`flags=(unconfined)`, and sole rule `userns,`. There are no local includes or
wildcard attachments. The installed root-owned ABI file is checked and hashed.

This permission applies to every matching invocation on that VM while loaded,
not only one fixture call. Descendants can inherit the named profile. It grants
the application's user-namespace permission while otherwise retaining its
previous unconfined access; it does not provide a new AppArmor sandbox. Existing
UID/PID namespaces, same-UID mapping, held-inode execution, session containment,
131072-byte stream caps, 15000-ms budget and 2000-ms cleanup remain unchanged.
The fixture still runs as the ordinary runner user. Only the fixed host setup
and removal helper runs with privilege. No shell, fixture, Docker operation or
build pipeline runs as root through this helper.

## Admission and conflict handling

The workflow must bind the helper to its exact Git blob, read those bytes once
under isolated system Python, verify the Git blob digest, then execute that same
in-memory byte sequence. Cleanup must execute the same admitted helper. The
helper also binds setup and cleanup to a SHA-256 of those admitted source bytes.

Before loading anything, require Ubuntu 24.04, global AppArmor enabled, both
namespace flags at their observed value 1, and an unconfined helper. Require the
launcher, parser, ABI file and their parents to be root-owned, non-symlink and
not writable by group or other users. Hash and retain their identities.

Recursively inspect kernel policy `name`, `attach`, `mode` and `sha256` using the
[kernel AppArmor introspection interface](https://github.com/torvalds/linux/blob/v6.17/security/apparmor/apparmorfs.c).
Reject nested policy namespaces, missing hashes, ambiguous names, the reserved
profile name, and any attachment that could match the launcher. The conservative
checker expands bounded brace alternatives and proves disjoint literal prefixes;
unknown or more complex potentially matching expressions reject. This may reject
a harmless existing profile. It never replaces or overrides one to make progress.

Compile the tiny fixed policy without loading it. Require an unchanged host and
inventory, exclusively create root-owned mode-0700 state under
`/run/pscan-native-userns-prerequisite`, then add exactly one policy. Parser
arguments disable configuration inheritance and all caches, stop on warnings and
errors, and use no parallel compiler jobs. See the
[parser commands](https://gitlab.com/apparmor/apparmor/-/blob/v4.0.1/parser/apparmor_parser.pod)
and [config-file implementation](https://gitlab.com/apparmor/apparmor/-/blob/v4.0.1/parser/parser_main.c).
Read back exact name, attachment, unconfined mode and kernel digest. Require every
other profile and every checked global/file identity to remain unchanged.

## Cleanup and failure behavior

An unconditional workflow cleanup step runs even after fixture failure. It checks
the source/policy bindings, root-owned state, host identities, unrelated profile
inventory, and exact installed kernel digest before removing only this policy.
Then it verifies the original inventory and global settings and deletes only its
three named state files and empty owned directory. No recursive deletion occurs.

Compile or preflight failure never loads policy. Failed or ambiguous addition
does not permit fixture execution. A policy present without a successful owned
digest record is not removed speculatively. Any altered host, unrelated policy,
ownership uncertainty or failed removal/readback leaves the job failed. The
disposable VM's destruction is the backstop, not reported as successful cleanup.
The helper never retries, replaces policy, changes a sysctl, stops AppArmor,
installs packages, writes persistent configuration, or executes arbitrary input.

Local inert tests cover conflicts/unknowns, strict readback, preflight/add ordering,
failure before ownership, source/hash drift and refusal to remove an unowned or
changed profile. They block native subprocesses. They prove control decisions,
not Linux parser compatibility or successful namespace execution. Both require
a separately reviewed fresh diagnostic invocation with live provider preflight.
