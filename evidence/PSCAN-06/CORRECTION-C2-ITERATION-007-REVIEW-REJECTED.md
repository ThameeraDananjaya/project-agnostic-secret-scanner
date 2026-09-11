# PSCAN-06 Correction C2 iteration 007 independent review rejection

Date: `2026-09-11`
Independent review session: `01a0916e-95e6-7db0-9391-08232b2fc505`
Authority approval session: `01a09141-1ba5-7931-bedc-bee2402fbf95`
Implementation session: `01a09151-048a-7bf3-a3a6-2d0c26415f10`

## Outcome

- Result:
  `LOCAL_ACCEPTANCE_REJECTED_WORKFLOW_PRESERVATION_TEST_FALSE_FAILURE`
- Exact reviewed candidate:
  `230e1761f7c45d9629cecf360e7498b33d64ab6f`
- Exact reviewed tree:
  `5d74628d76f95e150c0ea1de51b4892af2350a7d`
- Direct authority parent:
  `b7a0fb211e4631f915e368174ae58022d1519182`
- Blocking findings: one
- Product, workflow, schema, verifier or test repair by reviewer: none

The append-only workflow and schema identities themselves match the selected
iteration-007 design. The required committed preservation regression does not.
Its workflow normalization globally replaces every `2.2` substring with
`2.1`. That changes an unchanged action-version comment from `# v4.2.2` to
`# v4.2.1`, so the exact committed Go test deterministically compares unequal
workflow text and fails. The candidate therefore cannot pass its required
offline Go suite or the authority's preservation-test criterion. Local
acceptance is rejected. The reviewer did not repair the candidate.

## Admission and authority

The review ran in the direct saved checkout at
`C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner` on branch
`main`. Git directory and common directory both resolve to `.git`. The review
session differs from both authority and implementation sessions. Exact HEAD,
tree and sole parent matched the identities above before review.

The controlling committed-byte hashes independently matched:

- `PASS-OUTCOME-SPEC-001`:
  `8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`;
- iteration-007 authority:
  `A17520815841415AD06357682423CF7C635B32AAFAB7EA2B0303E2540FE709E8`;
- Recovery R5 terminal failure:
  `53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`;
  and
- iteration-006 local acceptance:
  `9E045C7D22841EA89F98223848C8530396F910F729B2E25BAB20541C5F0CD6AC`.

The project goal and tracker were re-read at the start of the authority,
candidate, validation and verdict phases. PSCAN-06 remained the sole open
task. Recovery R5 remained terminal and consumed; Recovery R6 and all
successor work remained unauthorized.

The existing 108-node Graphify graph was used read-only for routing to the
product contract, PSCAN-06 task, release authority and schema ownership
surfaces. It predates iteration 007 and was not treated as source authority.
No Graphify refresh, result save, deletion, movement or staging occurred.

## Exact source and path confinement

The graphify-aware current-checkout verifier passed:

```text
commit=230e1761f7c45d9629cecf360e7498b33d64ab6f
tree=5d74628d76f95e150c0ea1de51b4892af2350a7d
files=379 raw_equal=361 canonical_crlf=18
index=EXACT index_flags=EXACT unsupported_config=0
nonignored_untracked=0 ignored_only_graphify=66
```

All ignored paths were confined to preserved `graphify-out/**`. A separate
local-only clone with `core.autocrlf=false` and `core.eol=lf` checked out the
exact candidate detached and passed the product source-trust verifier with
`379/379` tracked paths raw-equal, zero canonical projections, zero untracked
paths and zero ignored paths.

The candidate changes exactly 19 paths. All 19 are listed in the
iteration-007 authority; there are no extra or missing candidate paths.
`git diff --check` passes for the exact parent-to-candidate range.

## Preserved and rolled-forward identity evidence

Committed-byte verification reproduced:

```text
old workflow blob  3aa42627628b8b5298d854d29b3880cb19dc36ff
old workflow SHA-256 C5F40F1B32E87C005FE33EE607AF7E3D19EE4F0173C21619E21158C31AA0BDB4
old schema 2.1 blob f8b5232b8667636fec8f4f0c7f083e9a0297b0cb
old schema SHA-256 CAA9CD26665CC3A3550AFFEA7490A0F1F277A69B10B1F1537787616E8AB973CE
new workflow blob  ff5b55f32ced9f382b30e98746f9ed5e21f13a73
new workflow SHA-256 CFE3FB919CA26A7AE6CF2A8F56D22BA61B96242298E010940E7B863053C9FA1A
new schema 2.2 blob 8c20c4943072f80d57d3d3b9a3b094ad6c524fca
new schema SHA-256 347E424F23F48CF25A409328E3A2E6773F589104AC0129216977DE6D01ED872A
```

The old workflow and schema 2.1 blobs are identical at accepted iteration-006
candidate `cf1679f9bca24887340fa4060f37d1ebff21f305` and at the reviewed
candidate. Line-by-line comparison found that the new 293-line workflow
differs from the old workflow only at the eleven authority-selected identity,
name, concurrency and schema-assertion lines. The new 151-line schema differs
from schema 2.1 only at the nine selected version/title/identifier and exact
tag/ref/path/certificate fields. Schema normalization is exact.

Static inspection confirms the active builder, embedded verifier, generic
verifier and schema agree on schema `2.2`, tooling tag
`release-tooling-v1.0.0-c2-r6`, workflow
`.github/workflows/release-recovery-v1.0.0-c2-r6.yml`, exact tag ref and
certificate identity. The generic verifier keeps separate 2.0/C1, 2.1/C2 and
2.2/R6 branches and exact per-version identities.

`build/release/docker-execution.ps1` at the authority parent has the accepted
iteration-006 blob `2fb43bf89df352e9153ba5d7ad23fd59cad77749`. Removing the one new exact
schema-2.2 copy line from the reviewed candidate reproduces that entire parent
file byte-for-byte. This proves preservation of `Get-LinuxSessionMembers`, its
accepted PID predicate and identity/ledger data flow, the native boundary and
the complete Docker operation table.

## Blocking committed test defect

The exact candidate contains this final normalization step in
`tests/integration/supply-chain/release_test.go`:

```go
normalizedWorkflow = strings.ReplaceAll(normalizedWorkflow, "2.2", "2.1")
```

Both old and new workflows contain the same line 277:

```text
uses: actions/attest@1e69f48acb82d1966a394da916b4c1698aa569d6 # v4.2.2
```

The test transforms the new-workflow copy of that line to `# v4.2.1` while
leaving the old workflow at `# v4.2.2`. The following equality check therefore
always reaches:

```go
t.Fatal("R6 workflow differs from the C2 workflow outside the selected identity/version fields")
```

An independent clean-materialization probe applied the exact ordered string
replacements from the committed test. It reported the sole residual at line
277 and exited natively with code `1`. The outer review assertion required
that nonzero exit and rejected any false pass.

This is not a defect in the selected workflow derivation; it is a defect in
the mandatory regression that is supposed to prove that derivation. It makes
the committed test suite fail for unchanged, accepted action-pin commentary
and prevents the candidate from satisfying mandatory implementation design
item 6 and required local-check item 4.

## Other independent checks

These checks passed from the exact clean materialization with explicit native
exit-code enforcement where a child process was used:

- PowerShell parsing for the changed builder and Docker materialization route,
  plus the CRLF payload harness;
- JSON parsing for release-manifest schemas 2.0, 2.1 and 2.2;
- complete iteration-006 PID source/data-flow regression;
- complete iteration-005 Windows job-object, stream, UTF-8, timeout,
  descendant cleanup, replacement-race, second-invocation and hostile
  ambient-type matrix;
- synthetic image-admission three-state/parser/closed-boundary matrix;
- host-only CRLF checkout, raw rejection and LF-normalization regression with
  `docker=NOT_INVOKED`; and
- all 12 hostile source-trust cases plus the exact positive case, with zero
  build output and no untrusted driver action.

Final inventory found zero surviving processes referencing
`native-fixture.ps1`, zero `pscan-c2-iteration-005-*` directories and zero
`pscan-docker-boundary-*` directories. The three existing local tags remained
at their recorded commits, and local tag
`release-tooling-v1.0.0-c2-r6` remained absent.

No provenance-bound Go toolchain/module cache path was supplied to this review
session. `go` and `gofmt` were absent from `PATH`. No offline YAML parser,
`actionlint`, `yq` or executable JSON Schema engine was available. Go
formatting/parsing/compilation/tests, executable YAML parsing and executable
JSON Schema validation therefore remain unproved; no dependency or tool was
downloaded. Those gaps do not weaken or defer the deterministic rejection
above.

## Preserved terminal boundary

No Docker command, network connection, remote read or mutation, dependency
acquisition, tag change, workflow dispatch or rerun, credential action,
signing, attestation, draft, release, deployment, publication, purchase,
subscription, nonzero spend or successor action occurred.

Recovery R5 remains terminally failed and consumed. Candidate
`230e1761f7c45d9629cecf360e7498b33d64ab6f` is not independently accepted.
Recovery R6 remains unauthorized. Actual-Linux, genuine-Docker, dependency,
build, reproducibility and artifact-integrity proof remains open. PSCAN-06
remains open and unaccepted overall; PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and ineligible. No successor is selected, activated,
claimed or worked.
