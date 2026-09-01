# PSCAN Task Tracker

| Task | State | Current authority |
|---|---|---|
| PSCAN-01 | Completed and accepted locally | Activation `a85a64e`; accepted implementation `5019489` |
| PSCAN-02 | Completed and independently accepted locally | Activation `48a8bd0`; accepted implementation `98768cc` |
| PSCAN-03 | Rejected and closed fail-closed; not accepted | `evidence/PSCAN-03/CLOSEOUT-REJECTED.md` |
| PSCAN-04 | Proposed; unselected | PASS-SPEC-001 only |
| PSCAN-05 | Proposed; unselected | PASS-SPEC-001 only |
| PSCAN-06 | Proposed; unselected; remote gates reserved | PASS-SPEC-001 only |
| PSCAN-07 | Proposed; unselected; signing/publication gates reserved | PASS-SPEC-001 only |
| PSCAN-08 | Inactive; technically and legally gated | Material-gap evidence plus separate owner approval required |
| PSCAN-09 | Activated; unclaimed; not implemented | Owner command plus PSCAN-09 activation bundle |

PSCAN-02 completed its activated lifecycle within bounded paths and is accepted
locally. Its exact limitations remain recorded in `evidence/PSCAN-02`.
PSCAN-03 established a material Gitleaks binary-history gap; correction C1 then
failed independent fragment, archive-classification and proof-isolation review.
It is rejected, closed fail-closed and unaccepted. Its closeout selected or
activated no successor.

The owner abandoned the contract-preserving option on 2026-09-01 and activated
PSCAN-09 to decommission PASS-SPEC-001 authority and establish an outcome-based
replacement contract while preserving its historical source and evidence.
PASS-SPEC-001 remains controlling until a fresh PSCAN-09 implementation is
independently accepted and closed. PSCAN-09 is the only activated task and
remains unclaimed and not implemented. PSCAN-04 through PSCAN-07 remain
proposed and unselected. PSCAN-08 remains inactive, unselected, and ineligible.
No successor is selected or activated.
