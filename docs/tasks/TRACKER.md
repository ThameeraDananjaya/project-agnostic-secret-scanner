# PSCAN Task Tracker

| Task | State | Current authority |
|---|---|---|
| PSCAN-01 | Completed and accepted locally | Activation `a85a64e`; accepted implementation `5019489` |
| PSCAN-02 | Completed and independently accepted locally | Activation `48a8bd0`; accepted implementation `98768cc` |
| PSCAN-03 | Rejected and closed fail-closed; not accepted | `evidence/PSCAN-03/CLOSEOUT-REJECTED.md` |
| PSCAN-04 | Completed and independently accepted locally | Exact activation `a21b030e1658f1f98ac4e4d001af12185d9ed311`; evidence-bearing closeout commit |
| PSCAN-05 | Completed and independently accepted locally | Exact activation `50b4186`; corrected evidence-bearing closeout commit |
| PSCAN-06 | Correction C2 pinned-image bootstrap recovery is owner-approved and authority-recorded; unclaimed, unimplemented, unaccepted and open | `evidence/PSCAN-06/CORRECTION-C2-AUTHORITY.md` |
| PSCAN-07 | Proposed; unselected; signing/publication gates reserved | PASS-OUTCOME-SPEC-001; follows PSCAN-06 |
| PSCAN-08 | Inactive; technically and legally gated | Material-gap evidence plus separate owner approval required |
| PSCAN-09 | Completed and independently accepted locally | Activation `f486989`; accepted closeout commit |
| PSCAN-10 | Completed and independently accepted locally | Activation `9053b37`; accepted correction `d10df99` |

PSCAN-02 completed its activated lifecycle within bounded paths and is accepted
locally. Its exact limitations remain recorded in `evidence/PSCAN-02`.
PSCAN-03 established a material Gitleaks binary-history gap; correction C1 then
failed independent fragment, archive-classification and proof-isolation review.
It is rejected, closed fail-closed and unaccepted. Its closeout selected or
activated no successor.

The owner abandoned the contract-preserving option on 2026-09-01 and activated
PSCAN-09. PSCAN-09 is now independently accepted and closed. It preserved
PASS-SPEC-001 and all PSCAN-01 through PSCAN-03 history, established
PASS-OUTCOME-SPEC-001 through DEC-002, and did not rehabilitate PSCAN-03.

The owner exactly activated PSCAN-10 on 2026-09-01 as the bounded
primary-coverage successor. It was claimed in a fresh implementation session
from exact clean activation commit
`9053b37d799b19e3d98aeca0ae853296971ab09d`, independently accepted after the
correction commit `d10df991d72e3fcc40b378830b50d1f258e16654`, and closed locally. PSCAN-04
was exactly activated by the owner on 2026-09-01;
it was claimed alone on 2026-09-02 in a fresh implementation session from the
exact activation commit. PSCAN-04 is implemented, independently accepted and
closed locally through its evidence-bearing closeout commit. PSCAN-05 was
exactly activated by the owner on 2026-09-02 and claimed alone in a fresh
implementation session from exact activation commit
`50b418609c1f9927c0ecd5d740aba6d7bff11f58`. Candidate
`ab3626ac11e4a62c915e209f1b9f17f098950291` was independently reviewed,
corrected within PSCAN-05 paths, revalidated and accepted through its local
evidence-bearing closeout commit. PSCAN-06 was exactly activated by the owner
on 2026-09-02 and was claimed alone on 2026-09-03 in a fresh implementation
session from
exact clean activation commit
`058ffcd446c6431b2e1afeed769d02c7b1f307f8` after the complete ordered
reading map and current official/read-only preflight. Candidate
`a13c28fe7273bc8dc6545f97966a02889524eb4c` passed two clean network-disabled
builds, byte-for-byte reproduction, complete local tests and skeptical review.
The owner approved the exact remote/signing gate on 2026-09-03. The public
repository and approved controls were created and read back, `main` was pushed,
and locked tag `v1.0.0` was created at the validated candidate. Authorized run
`33709197614` failed closed during pinned dependency acquisition because its
Linux Docker bind-mounted Go module cache was not writable. No build artifact,
signature, attestation, draft release or publication was created. PSCAN-06 is
unaccepted and open. On 2026-09-03 the owner approved the bounded
contract-preserving Correction C1 recovery. Its authority bundle records a
Linux cache-ownership correction and a dual-identity release design that keeps
the public locked `v1.0.0` tag fixed at the accepted product-source candidate
while separately binding correction tooling and workflow identity. The
correction was claimed alone in a fresh session from exact clean authority
commit `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`. Independent review rejected
candidate `f24b832ebe5f6749aa0ab910e1b9be279065ebb1` after reproducing CRLF
carriage-return corruption in its Docker POSIX shell payloads. Bounded
iteration 003 candidate `5ca77226ed3996a8267caf02da366d0beb915c8d` added runtime LF
normalization but was independently rejected when a complete CRLF-checkout
build consumed transformed working-tree integrity bytes and failed the pinned-
rule binding test. Bounded iteration 004 candidate
`a22579fd5af14473e1d49b5027f21ab591bbb589` removed checkout build inputs but
was independently rejected after assume-unchanged driver tampering bypassed
its status-only clean gate. Bounded iteration 005 candidate
`3fb7592889820fa2739a4a53588e073689621809` adds exact committed entrypoints and
complete hostile-state index/path/mode/raw-byte verification. Twelve isolated
adversarial cases reject before output or untrusted driver action, and two
independent complete builds from separate actual CRLF checkouts reproduce all
30 files byte-for-byte. The required actual-Linux host positive/wrong-owner
proof was attempted under the exact build-only gate, but run `33829598255`
failed closed before acquisition because the fresh Ubuntu runner lacked the
pinned Docker image required by the earlier CRLF regression. Both locked tags
remain unchanged; signing and all later steps were skipped. On 2026-09-04 the
owner approved bounded Correction C2 pinned-image bootstrap recovery. Its
authority preserves all historical evidence and existing schemas, permits only
an image-independent host canary and CRLF-normalization proof before the exact
canonical digest-pinned image bootstrap, and proposes a distinct immutable C2
tooling identity. Correction C2 is unclaimed and unimplemented and must begin
in a fresh session from its committed authority bundle. No correction remote
action is approved. PSCAN-07 remains proposed and
unselected. PSCAN-08 remains inactive,
unselected and ineligible because
accepted evidence did not establish a material required-class Gitleaks gap;
separate technical and AGPL owner approval would also be required. No
successor to PSCAN-06 is selected or activated.
