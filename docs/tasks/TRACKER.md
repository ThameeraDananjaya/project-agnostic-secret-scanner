# PSCAN Task Tracker

| Task | State | Current authority |
|---|---|---|
| PSCAN-01 | Completed and accepted locally | Activation `a85a64e`; accepted implementation `5019489` |
| PSCAN-02 | Completed and independently accepted locally | Activation `48a8bd0`; accepted implementation `98768cc` |
| PSCAN-03 | Rejected and closed fail-closed; not accepted | `evidence/PSCAN-03/CLOSEOUT-REJECTED.md` |
| PSCAN-04 | Proposed; unselected | PASS-OUTCOME-SPEC-001; requires accepted PSCAN-10 |
| PSCAN-05 | Proposed; unselected | PASS-OUTCOME-SPEC-001; follows PSCAN-04 |
| PSCAN-06 | Proposed; unselected; remote gates reserved | PASS-OUTCOME-SPEC-001; follows PSCAN-05 |
| PSCAN-07 | Proposed; unselected; signing/publication gates reserved | PASS-OUTCOME-SPEC-001; follows PSCAN-06 |
| PSCAN-08 | Inactive; technically and legally gated | Material-gap evidence plus separate owner approval required |
| PSCAN-09 | Completed and independently accepted locally | Activation `f486989`; accepted closeout commit |
| PSCAN-10 | Proposed; unselected | PASS-OUTCOME-SPEC-001 primary-coverage boundary |

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

PSCAN-10 is the single proposed bounded primary-coverage successor. It is
unselected, not activated, not claimed and not implemented. PSCAN-04 through
PSCAN-07 remain proposed and unselected. PSCAN-08 remains inactive, unselected
and ineligible pending accepted material-gap evidence plus separate technical
and AGPL owner approval. No successor is selected or activated.
