# Linux build path validation

The tagless manual `release-validate-linux.yml` workflow binds its exact source
commit, workflow SHA/ref, repository/owner and run identity. It uses the real
source admission, cache/image admission, acquisition, offline build, packaging
and independent-byte-comparison path. It has no artifact upload or signing job.
`Validation` mode writes `BUILD-VALIDATION.json` and truthful validation
provenance, and returns before any candidate `release-manifest.json` emission.
The production mode retains its immutable tooling-tag equality gate. Existing
schemas and historical tags retain their original meaning.

This diagnostic does not deliver the final release packages. A successful run
would establish real temporary package bytes and their hashes; the separately
reviewed immutable candidate path still must deliver the unsigned packages.

## Source and host ownership

Acquisition admits the explicit entire repository/index/commit and proves exact
go.mod/go.sum blobs before cache mutation or network. It does not infer source
from the build/release-only script materialization. Linux cache, acquisition,
build and package operations use the admitted invoking UID/GID, including owned
tmpfs and private output directories. Missing, partial, negative or different
Linux identities fail closed. PowerShell nullable values are checked directly.

The host canary uses exclusive file creation, non-overwriting same-directory
rename, exact read and deletion under the invoking OS identity. Failures preserve
the operation, original cause and bounded diagnostics; only its owned path is
cleaned. This avoids filesystem-provider Hidden-attribute filtering without
changing permissions. The previous run discarded its underlying exception and
therefore does not prove which individual file operation failed.

## Separate client and daemon ownership

The native process boundary alone cannot establish daemon workload termination.
Each fixed container operation now creates an inert container with an internal
random name/label, binds its exact ID, verifies the created image and relevant
user/mount/network/privilege/resource configuration, starts only that ID once,
proves exit, removes that ID and proves its absence with an exact filtered query.
The public operation set remains closed. There is no caller-selected container
target, callback, prune, broad cleanup or unrelated-container deletion.

A failed or uncertain create is never started. Unknown ownership is never used
for deletion, and one empty inventory read cannot resolve a late create. The
failure remains terminal and reports uncertainty. Exact admitted IDs may be
removed forcibly during the original bounded terminal cleanup. An actual
nonzero workload exit remains nonzero, including expected negative controls;
OOM, identity mismatch or uncertainty cannot become successful containment.

One monotonic budget spans the whole protocol: 15000ms for all operation stages
and a single 2000ms grace shared by native and daemon cleanup. Nested calls do
not reset these clocks. Protocol and workload bytes share independent 131072-
byte stdout/stderr limits. The existing 35s fresh-child transport does not
extend those limits. The validated child result includes the exact removed
container identity, protocol-call count and aggregate stream counts; the caller
logs a bounded `PSCAN_DOCKER_LIFECYCLE` record.

Cold image pull/acquisition/full build may exceed 15 seconds. No budget extension,
retry or stage splitting is implied. Failure freezes the attempt and stops.

## Evidence limits

Inert tests drive the actual state machine with substituted I/O, actual fixed
argument builders and actual native shared budget/capture code. They prove
rejection and ordering properties, not a Docker daemon's behavior. The default
native orchestrator retains the complete isolated Windows/Linux fixture matrix
and now exercises the ancestor-stop gate instead of looking for a retired shell
self-stop string. A fresh actual Linux run remains necessary; none is claimed
by the existence of these files or passing local tests.
