# PSCAN-06 Correction C2 iteration 007 preflight and claim

Date: `2026-09-11`
Implementation session: `01a09151-048a-7bf3-a3a6-2d0c26415f10`
Authority approval session: `01a09141-1ba5-7931-bedc-bee2402fbf95`

## Claim

This genuinely fresh saved-project session claims only PSCAN-06 Correction C2
iteration 007. It does not claim, select, activate or implement Recovery R6,
PSCAN-07, PSCAN-08 or any other successor. Recovery R5 remains terminally
failed and its sole tag push and workflow dispatch remain consumed.

The controlling iteration authority is
`CORRECTION-C2-ITERATION-007-AUTHORITY.md`; its SHA-256 at the claim boundary
was `A17520815841415AD06357682423CF7C635B32AAFAB7EA2B0303E2540FE709E8`.
The controlling `PASS-OUTCOME-SPEC-001` SHA-256 was
`8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.
The immutable Recovery R5 failure SHA-256 was
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.
The bounded iteration-006 acceptance SHA-256 was
`9E045C7D22841EA89F98223848C8530396F910F729B2E25BAB20541C5F0CD6AC`.

## Starting identity

- Saved-project root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
- Direct saved checkout, not a worktree: Git directory and common directory
  both resolved to `.git`.
- Branch: `main`.
- Exact starting HEAD:
  `b7a0fb211e4631f915e368174ae58022d1519182`.
- Exact starting tree:
  `bd2272041514675a8446bdde6733e359225c4e39`.
- Exact direct parent:
  `5ff14709520ba99b105ee515e5482455c8511efc`.
- The implementation session differs from the authority approval session.
- Tracked and staged diffs were empty; non-ignored untracked status contained
  zero paths.
- All 66 ignored paths were confined to preserved `graphify-out/**`; none was
  deleted, moved or staged.
- The existing ignored graph contained 108 nodes and routed the release
  identity inspection to the PSCAN-06 task, release authority and schema
  ownership surfaces. Graph output was treated only as a navigation aid.

## Exact source trust

The complete authority-compliant every-byte check proved:

```text
Exact graphify-aware source trust PASS commit=b7a0fb211e4631f915e368174ae58022d1519182 tree=bd2272041514675a8446bdde6733e359225c4e39 files=374 raw_equal=356 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

The check used the repository `source-trust.ps1` tree, index, index-flag,
object and canonical LF-to-CRLF algorithms. It verified SHA-1 object format,
the exact complete tree/index path, mode and object identity, ordinary and
fsmonitor index flags, absence of shortcut/sparse/promisor configuration,
every tracked working-tree byte, zero non-ignored untracked paths, and the
explicit authority exception only for ignored `graphify-out/**` material.

The locked local tag identities remained unchanged:

- `v1.0.0` -> `a13c28fe7273bc8dc6545f97966a02889524eb4c`;
- `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`; and
- `release-tooling-v1.0.0-c2` ->
  `faef8435322c9096df09b56969662411f17356ea`.

No local `release-tooling-v1.0.0-c2-r6` tag was created.

## Readings and bounded implementation surface

Before this claim the session read completely: `AGENTS.md`,
`PASS-OUTCOME-SPEC-001`, `TRACKER.md`, `PSCAN-06.md`, the PSCAN-06 reading
map, the iteration-007 authority, the Recovery R5 terminal failure and the
iteration-006 acceptance. The project objective and tracker were rechecked
immediately before this phase.

The implementation is limited to the authority's append-only workflow and
schema identities, the named builder/verifier/materialization and test routes,
the synchronized living documentation, and new iteration-007 evidence. The
old workflow, schema 2.1, PASS contracts, historical evidence, iteration-006
PID predicate/data flow, Docker operation table and all unrelated work remain
immutable.

No Docker command, network action, remote read or mutation, dependency
acquisition, tag action, workflow dispatch or rerun, signing, publication,
spending or successor action occurred. Any need to change an absent path or
broaden this design remains terminal and returns to the owner.
