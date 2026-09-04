# PSCAN-06 Correction C2 iteration 003 author validation

## Bound candidate

- Date: `2026-09-05`
- Authority commit:
  `5762d1ef6a7708629bad5f5bb33eba21ca1cdcd4`
- Implementation commit:
  `dfbe897e9e47632ee5ca9437650bd62eafbaf341`
- Implementation tree:
  `b15586e8ec5b8b799425998920c60348f01e385b`
- Session: `01a06e8c-6a86-74a2-819f-edcbc15fe6f8`
- Checkout:
  `C:/Users/ITDan/.codex/worktrees/pscan06-c2-i003-clean-01a06e8c`
- Checkout status before checks: clean tracked, untracked and ignored state
- Result: `AUTHOR_VALIDATION_PASS`; this is not acceptance

The earlier EOL-projected staging commit
`1f7edd38959e2eaf20aa8e9b20bb72dd8e948856` was rejected by source trust and is
not part of this candidate or any integration range.

## Exact-commit checks

### Source trust

The repository source-trust verifier bound exact commit `dfbe897` and tree
`b15586e`, enumerated 332 tracked files, proved all 332 raw working-tree files
equal to HEAD, recorded zero canonical-EOL projections and returned
`WorkingTreeInputsTrusted=true`.

`test-source-trust.ps1` then rejected all 12 isolated hostile states:
assume-unchanged, skip-worktree, staged change, unstaged change, untracked file,
ignored untracked file, fsmonitor config, untracked-cache config, driver
tampering, path mismatch, mode mismatch and blob mismatch. Every rejection
produced zero build outputs and no untrusted driver action. Its exact positive
case passed.

### Private admission and native process bounds

`test-image-admission.ps1` completed with:

```text
Image admission iteration-003 PASS private=BOUND state-matrix=PASS byte-caps=LIVE utf8=STRICT timeout=MONOTONIC process-tree=TERMINAL
```

The suite executed no Docker or network operation. It reran the iteration-002
present, conclusive-absence, pull, identity and untrusted matrix through the
test-only state boundary, including terminal timeout, overflow, read and
cleanup uncertainty with zero later action. Static and runtime probes rejected
the former callback names, callback/script-block parameters, dot-sourcing,
ambient functions/aliases and PATH/environment substitution.

The native fixtures used only:

- executable: `C:/Program Files/PowerShell/7/pwsh.exe`
- SHA-256:
  `362A356CE7F0940EC74F73A8FC2C990A2CC24A38A11C90BBD8ECA947110AD139`
- Authenticode: `Valid`
- signer: `Microsoft Corporation`

They proved below/exact/above stdout and stderr caps, exact simultaneous stream
capture, simultaneous overflow, split UTF-8, invalid and incomplete UTF-8,
immediate and nonzero exit, start failure, hang, flood, child and grandchild
tree termination, bounded return and no surviving recorded descendant.

### Host-only CRLF and containment

`test-crlf-shell-payloads.ps1 -Phase HostOnly` proved an actual CRLF checkout,
raw shell carriage-return rejection and LF normalization with
`docker=NOT_INVOKED`.

The authority-to-candidate path inventory contains only the 11 named files in
the authorized path groups. `git diff --check` passed. Contracts and all schema
paths are unchanged. The locked tags still resolve exactly to:

- `v1.0.0` -> `a13c28fe7273bc8dc6545f97966a02889524eb4c`
- `release-tooling-v1.0.0-c1` ->
  `3fb7592889820fa2739a4a53588e073689621809`

Committed script SHA-256 values:

- `admit-image.ps1`:
  `2262716A57960267A1D571ACC172D0C7A82FF3F30E4E8F08CBD0249143A12078`
- `image-admission.ps1`:
  `E411DC9EC2B8337011979037F1AED6BAA6D6D294C5AA763E10A51BA52EAD4129`
- `test-image-admission.ps1`:
  `D8D31508462896F0E8C9CBC3A865DF0C221DE2757ADBCF126F0B67A21EED4A79`

The two explicitly named temporary validation directories were verified under
the system temporary root and removed after the checks.

## Open proof and owner gates

Actual-Linux native process execution, Linux cache ownership, genuine Docker
identity/admission, any pull, fresh hosted-runner behavior, container proof,
dependency acquisition, reproducibility builds and independent skeptical
review remain open. No Docker, network, download, scanner, build, remote, push,
tag, setting, workflow, signing, attestation, draft, publication, spending,
credential, key or TruffleHog action occurred.

PSCAN-06 remains open and unaccepted. PSCAN-07 remains proposed and unselected;
PSCAN-08 remains inactive and separately gated. No successor is activated.
