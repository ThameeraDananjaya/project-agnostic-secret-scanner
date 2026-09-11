# PSCAN-06 Correction C2 iteration 008 implementation attempt

Date: `2026-09-11`
Implementation session: `01a091b5-687a-7392-b366-4b71f8bbbe8d`
Starting authority commit: `4f86651de67a05db8ae0076c0c155bd184082af3`
Starting authority tree: `4ef9e344c1849c42509a61077e6a509fae8ef355`

## Result

`STOPPED_NO_CANDIDATE_UNRELATED_GOFMT_AND_COSIGN_BASELINE_FAILURES`.

The session applied the exact bounded workflow-normalization repair in the
working tree: it retained the four identity-specific workflow replacements,
added exact source occurrence assertions `1`, `1`, `8`, `3`, `1`, replaced
only the exact full manifest-schema assertion, required the exact unchanged
`actions/attest` line containing `# v4.2.2` in the old, new and normalized
workflows, forbade `# v4.2.1`, and retained byte-for-byte workflow equality.

Before any Go test, the required exact Go 1.27.1 `gofmt -d` check returned
native exit `1` and a 66-line formatter diff. The diff was outside the bounded
workflow-normalization repair. The session then exported the untouched exact
authority commit to a separate clean materialization and ran the same exact
formatter against its committed test file. It returned the same native exit
`1`, the same 66-line formatter diff and the same four pre-existing hunks:

```text
@@ -162,7 +162,9 @@
@@ -193,10 +195,14 @@
@@ -209,10 +215,16 @@
@@ -263,9 +275,9 @@
```

The untouched authority-commit file and the restored saved-checkout file both
have SHA-256
`29E715D5534B389BCB65E046E84E5F2B2828F3D400BB664AFB13B55F12F377D0`.
The exact clean materialization contains all 381 committed files at authority
commit `4f86651de67a05db8ae0076c0c155bd184082af3`, tree
`4ef9e344c1849c42509a61077e6a509fae8ef355`.

## Fail-closed boundary

The authority permits changes only to the workflow-normalization portion of
`TestIteration007RepositoryIdentityAgreementAndPreservation`. Applying the
formatter's other-test and non-normalization changes would broaden the repair
and violate the mandatory path-within-file boundary. The user instruction also
requires an immediate stop on an unrelated failure. The bounded working-tree
test change was therefore reverted exactly; no implementation candidate
survives. The claiming implementation session ran no Go test or `go vet`
command after the formatter gate failed.

A separate fresh verification continuation
`01a091b5-6872-7873-9647-ab48962b6356` independently used the exact supplied
offline Go route. On the untouched authority commit, the targeted command
failed at `release_test.go:262` with native exit `1`. The full supply-chain
package also exited `1`, failing both that preservation test and the
pre-existing Windows synthetic Cosign test at `release_test.go:66` with
`signature is invalid`.

While the task-owned bounded probe was present, the targeted command exited
`0`. The full package then exited `1` solely at the same Cosign test. This
proves the normalization repair itself, but it cannot satisfy the mandatory
clean formatter and full-package gates. The probe was restored byte-for-byte.
The continuation stopped without repository-wide tests, vet or no-Docker
harnesses and did not repair either out-of-scope failure.

No product, workflow, schema, builder, Docker materializer, verifier,
acceptance-test or historical-evidence byte changed. No Docker, dependency,
network, remote, tag, workflow, signing, publication, spending or successor
action occurred.

Iteration 008 is claimed but not implemented, author-validated, accepted or
complete. Recovery R5 remains terminal and consumed; Recovery R6 remains
unauthorized. A new exact owner decision is required to reconcile the frozen
Go 1.27.1 formatting and full-package Cosign gates with the iteration's
within-file scope before another implementation attempt.
