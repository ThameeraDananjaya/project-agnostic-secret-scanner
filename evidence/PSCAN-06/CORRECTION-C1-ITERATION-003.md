# PSCAN-06 Correction C1 iteration 003 contradiction and repair record

## Independent rejection carried forward

Candidate `f24b832ebe5f6749aa0ab910e1b9be279065ebb1` is rejected as a
portable Correction C1 candidate. Independent review reproduced a Windows
CRLF checkout transporting carriage returns inside the embedded POSIX shell
arguments passed to Docker. Observed failures included:

- `umask: Illegal number: 077` in the cache canary; and
- `sha256sum: invalid option -- '\r'` in acquisition and build payloads.

The passing author runs recorded in
`CORRECTION-C1-AUTHOR-VALIDATION-002.md` remain historical evidence for the LF
working tree that was executed, but they do not prove Windows-checkout
portability and cannot override the reproduced false-pass path. This new
record is append-only; no earlier evidence is rewritten.

## Repository-identity contradiction

The historical `CORRECTION-C1-PREFLIGHT.md` records repository
`ThameeraRA/project-agnostic-secret-scanner`. That statement is contradicted
by the verified authority and current configured remote, which identify
`ThameeraDananjaya/project-agnostic-secret-scanner`. The controlling product
contract, owner-gate record, release authority, recovery workflow and offline
verification policy all use `ThameeraDananjaya`. The historical preflight is
preserved unchanged and must not be used as repository-identity authority.

## Bounded repair

Iteration 003 converts CRLF and lone carriage returns to LF at runtime and
then asserts that no carriage return or NUL remains immediately before every
embedded POSIX shell payload is passed to Docker. Sandbox arguments and shell
commands are unchanged. A forced-CRLF regression covers the actual canary,
acquisition and build here-string payloads, parses each normalized payload in
the exact pinned network-disabled image, executes the cache canary from the
CRLF fixture, and proves a read-only cache still rejects without residue or a
success ledger. The recovery workflow runs this regression before acquisition.

The commit containing this record supersedes `f24b832` as the Correction C1
tooling candidate. Fresh acquisition and exact-candidate validation are still
required. Actual general-purpose Linux UID/GID execution, independent
skeptical acceptance and every remote/signing/release gate remain open. No
remote action, tag, workflow run, signature, release or successor action is
authorized or claimed.
