# PSCAN-06 Correction C2 iteration 006 author validation 004

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`

## Exact third refined candidate

- Commit: `cf1679f9bca24887340fa4060f37d1ebff21f305`
- Tree: `6cb0198c71881baff9c601c31ad6aa84217f4a49`
- Direct third-review-evidence parent:
  `cf6b61d73910d7525a2213877654dfb900fe2cb2`
- Preserved rejected candidates:
  `8fbbf7b695aa7f0c995dd0a655d40ddd7fc16fb4`,
  `e6a767aab71db1d3f62063dded379b4701d2cb52` and
  `763730403ae538850e9d806ce9db82c815517b2d`
- Third review-record SHA-256:
  `F21FD594E757747EE5D3F9865FBE2C33D1FECE7AF09B28AA50778D355A183373`
- Production blob: `2fb43bf89df352e9153ba5d7ad23fd59cad77749`
- Production blob SHA-256:
  `49211CA52336BE859CD0C976719B5FB1F209F9E1235F1DA4926C6C21B52905C2`
- Third refined regression blob:
  `39c08202e32d7fc45ac57d90df0a11ef65ab048a`
- Third refined regression blob SHA-256:
  `E74D9EB349C4394D4AD8E8D6EB307FBD091D469791C84A690187FD1E869D8DAB`
- Third refinement-record blob:
  `c636aebd0a255dd9910e020df4d4a369846e331d`
- Third refinement-record SHA-256:
  `81B727EFCCC4BC270C2F40932205D7670412731D3075D1FBA56025DF22E9EFDB`

The candidate is the direct child of the evidence-bearing third rejection
commit and changes exactly two authority-permitted paths: the regression test
and append-only third refinement record. The production script remains
byte-identical to all earlier candidates. Every prior candidate, author record
and independent rejection remains immutable.

## Exact VariablePath correction

Reflection against the installed `System.Management.Automation.VariablePath`
type confirms `UserPath`, `IsVariable`, `IsDriveQualified`, `DriveName`,
`IsUnqualified`, `IsUnscopedVariable`, `IsGlobal`, `IsScript`, `IsLocal` and
`IsPrivate` are public. `UnqualifiedPath` exists but its getter is
assembly-internal and is not used by the regression.

The predicate now admits automatic PID only when a public scope flag and the
complete `UserPath` agree exactly, case-insensitively:

- unqualified plus `PID`;
- unscoped-variable plus `variable:PID`;
- global plus `global:PID`;
- script plus `script:PID`;
- local plus `local:PID`; or
- private plus `private:PID`.

Drive-qualified paths are excluded. No colon split or last-segment inference
remains. Consequently `${global:env:PID}` and `${local:foo:PID}` retain their
full distinct variable names and are not classified as automatic PID.

The established exact assignment-left walker is otherwise unchanged: it
unwraps attributed/typed and parenthesized pipeline/command-expression
wrappers, recursively visits multiple-target array elements, stops at member
and index lvalues, and fails closed on unsupported target shapes.

## Predicate, runtime and data-flow validation

The inseparable executable regression reported:

```text
Docker execution iteration-006 PID source regression PASS untyped=REJECT typed=REJECT parenthesized=REJECT multi-target=REJECT nested-wrapper=REJECT runtime-scopes=REJECT multi-colon-provider-controls=ALLOW member-index=ALLOW rhs-string-comment-ledger=ALLOW data-flow=PASS
```

All prior direct, typed, attributed, parenthesized, nested, background and
multiple-target positives remain. Unqualified, `variable:`, global, script,
local and private PID assignments, including braced forms, are parser/predicate
positives. Multi-colon true-variable names, environment/function providers,
unknown drives, member/index targets, RHS references, property names, strings,
comments and ledger fields remain negative controls.

Nineteen child processes run with `pwsh -NoProfile -NonInteractive` and a
10-second bound. All runtime-proven automatic PID scope forms reproduce the
read-only/constant collision; `${variable:env:PID}`,
`${global:env:PID}`, `${local:foo:PID}`, `${script:foo:PID}`,
`${private:foo:PID}`, `$env:PID`, `${env:PID}` and `${function:PID}` all exit
zero with `EXECUTED_WITHOUT_PID_COLLISION`.

The production data-flow proof still identifies exactly one renamed integer
binding from the parsed stat PID, exact identity
`"${linuxProcessIdentifier}:$startTime"`, exact two-field ledger shape, and
ledger bindings `PID=$linuxProcessIdentifier` and `StartTime=$startTime`.

## Exact source, parse and confinement checks

The exact-candidate source verifier reported:

```text
Exact source trust PASS commit=cf1679f9bca24887340fa4060f37d1ebff21f305 tree=6cb0198c71881baff9c601c31ad6aa84217f4a49 files=371 raw_equal=353 canonical_crlf=18 unexpected_mismatch=0
index=EXACT index_flags=EXACT unsupported_config=0 nonignored_untracked=0 ignored_only_graphify=66
```

Both implementation scripts parsed with zero errors. Exact range
`cf6b61d73910d7525a2213877654dfb900fe2cb2..cf1679f9bca24887340fa4060f37d1ebff21f305`
passed `git diff --check` and exact two-path confinement. Committed-byte hashing
reproduced the production, regression and refinement identities above.

## Local no-Docker native-boundary checks

The complete exact-candidate available-host matrix passed:

```text
Docker execution iteration-005 PASS isolated-clean=PASS exact-source=PASS containment=WINDOWS-JOB streams=PASS utf8=PASS process-tree=EMPTY replacement-race=REJECT second-invocation=REJECT
Docker execution iteration-005 hostile isolation PASS compatible=REJECT/NONZERO stale=REJECT/NONZERO fake-calls=0 docker-calls=0 trusted-results=0 fabricated-fields=0
```

This covers the clean Windows job-object boundary, exact native source,
131072-byte stream limits, UTF-8 behavior, fixed 15000 ms timeout, child and
grandchild cleanup, empty containment membership, executable replacement,
second invocation and compatible/stale ambient types. Final inventory found
zero surviving processes whose command line referenced `native-fixture.ps1`,
zero remaining `pscan-c2-iteration-005-*` directories and zero remaining
`pscan-docker-boundary-*` directories.

Tracked, staged and non-ignored untracked state was clean. All 66 ignored paths
remained confined to preserved `graphify-out/**`. Local tags remained:

- `v1.0.0` -> `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`
- `release-tooling-v1.0.0-c2` ->
  `faef8435322c9096df09b56969662411f17356ea`

## Limits and open proof

No Docker command, image or container action; network connection; remote read
or mutation; dependency acquisition; build; tag action; workflow dispatch or
rerun; credential action; signing; attestation; draft; release; deployment;
publication; or successor action occurred.

This is author validation only. It does not create acceptance or independently
approve the candidate. A fourth genuinely fresh skeptical task must inspect
exact candidate `cf1679f9bca24887340fa4060f37d1ebff21f305` and rerun
proportionate local no-Docker checks before any acceptance record.

Recovery R5 remains terminally failed and consumed. Actual-Linux,
genuine-Docker, dependency, build, reproducibility and artifact-integrity proof
remains open. PSCAN-06 remains open and unaccepted overall; PSCAN-07 remains
proposed and unselected; PSCAN-08 remains inactive and ineligible. No successor
is selected, activated, claimed or worked.
