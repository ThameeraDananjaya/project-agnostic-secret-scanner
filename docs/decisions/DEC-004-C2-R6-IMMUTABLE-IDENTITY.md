# DEC-004: Correction C2 R6 immutable identity roll-forward

## Status

Accepted local design for PSCAN-06 Correction C2 iteration 007. Local
implementation and author validation are not independent acceptance and do not
authorize Recovery R6, a tag, a workflow run or another remote action.

## Context

Recovery R5 pushed and locked `release-tooling-v1.0.0-c2` at pre-fix candidate
`faef8435322c9096df09b56969662411f17356ea`, then its sole workflow dispatch
failed because of the PowerShell automatic-PID collision. Iteration 006 fixed
that defect and is independently accepted locally at candidate
`cf1679f9bca24887340fa4060f37d1ebff21f305`, but it intentionally retained the
old workflow and schema bytes.

The old workflow requires its tag ref, workflow ref, workflow SHA and tooling
commit to agree. Schema 2.1, the active builder and the verifier bind the same
tag and workflow path. Moving the protected tag would rewrite immutable
evidence, dispatching the tag would rerun the defect, and dispatching the fixed
commit under another ref would fail the identity gate.

## Decision

Iteration 007 adds, rather than replaces, one complete identity generation:

- tooling tag `release-tooling-v1.0.0-c2-r6`;
- workflow `.github/workflows/release-recovery-v1.0.0-c2-r6.yml`;
- workflow ref `refs/tags/release-tooling-v1.0.0-c2-r6`; and
- release-manifest schema `2.2`.

The new workflow is a byte-preserving derivation of the C2 workflow except for
its selected name, concurrency, tag/ref, workflow path, certificate identity
and schema-2.2 assertion. The new schema is a structural derivation of schema
2.1 except for its version/title/identifier and the same selected identity.

The active builder and embedded verifier emit and require schema 2.2 under the
new identity. The generic verifier keeps schema 2.0/C1 and schema 2.1/C2 under
their original paths and refs, adds schema 2.2/R6, and rejects every
cross-version tag, ref, workflow or certificate mixture. Release artifacts
continue to include schemas 2.0 and 2.1 and add schema 2.2.

`docker-execution.ps1` changes only its exact-tree materialization payload to
copy schema 2.2. The accepted `Get-LinuxSessionMembers` predicate/data flow,
native containment, nine-operation Docker table, stream/time bounds and image
identity remain unchanged.

## Consequences and gates

The existing workflow, schemas 2.0/2.1, all locked tags and all historical
evidence remain immutable. Schema version selection is an identity boundary;
it is not a compatibility alias.

Iteration 007 creates no tag and performs no workflow execution. Recovery R6
remains unauthorized until a fresh independent task accepts the exact local
candidate and a later separate R6 authority re-proves current remote identity,
zero-spend and stop-usage facts. Signing, attestation, draft creation,
publication and successor work remain separately gated.
