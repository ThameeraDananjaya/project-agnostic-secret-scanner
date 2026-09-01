# Product Constitution

## Mission

Deliver one independently versioned local secret-scanning product that uses a
thin MIT-licensed Go runner and Gitleaks CLI as its primary engine. Scanner
execution is network-disabled. Findings are transient and fully redacted.

## Non-negotiable invariants

1. The product remains independent of every consuming project.
2. GitHub native secret scanning and hosted scanning services are not product
   dependencies.
3. Gitleaks CLI is the primary engine.
4. TruffleHog remains absent and disabled unless PSCAN-08 becomes separately
   eligible through evidence and owner approval.
5. Candidate material is data only and is never executed.
6. Scanner subprocesses use argument arrays, never interpolated shell commands.
7. Authoritative scanner execution has no network access or credentials.
8. Raw or partially raw findings never cross the private process boundary.
9. Unsupported, skipped, incomplete, or untrusted coverage never produces pass.
10. Outcomes are content-free and cannot approve merge, deployment, production,
    legal compliance, or go-live.
11. Projects own their policies, allowlists, receipt trust, revocations,
    evidence custody, deployment gates, and retention decisions.
12. Public release trust and project receipt trust are separate.
13. Releases use exact versions and digests; `latest` and auto-update are
    prohibited.
14. All CAP-1 through CAP-20 criteria and the complete acceptance matrix must
    pass before first-consumer readiness.
15. No successor task activates automatically.
16. Exact Git-object or projected input is pass-capable only when every admitted
    byte is bound to actual detector inspection.
17. Raw classification precedes framing, extraction or transformation.
18. Pinned detector span/stream behavior, declared resource profiles and each
    claimed input class require isolated adversarial proof.

## Standing owner decisions

- Product name: `project-agnostic-secret-scanner`.
- Default branch: `main`.
- Product licence: MIT.
- Intended final release: public, signed GitHub `v1.0.0`.
- Intended GitHub owner: `ThameeraDananjaya`.
- Spending ceiling: zero.
- Release signing: Cosign keyless signing through the specifically approved
  GitHub release workflow.
- Public repository creation: deferred to PSCAN-06 after read-only readiness
  and cost preflight and separate approval.
- Release maintainer and final human approver: `ThameeraDananjaya`.
- PSCAN-08: inactive unless a material Gitleaks coverage gap is accepted and a
  separate technical and AGPL legal decision is approved.
- Final runbook: plain English Markdown with clear Mermaid diagrams and complete
  input/output parameter and payload documentation.
- Controlling contract after PSCAN-09 closeout: `PASS-OUTCOME-SPEC-001`;
  `PASS-SPEC-001` remains immutable historical evidence.
