# Licensing Plan

## Product licence

The owner-selected licence for original scanner-runner source and repository
documentation is MIT. The root `LICENSE` records that choice. This decision
does not relicense third-party material.

## Third-party intake gate

Before source, binaries, rules, actions or generated assets enter a release,
the responsible task must record and independently verify:

1. exact upstream project, tag, commit and source location;
2. exact licence text, copyright notices and notice files;
3. source and binary redistribution, modification and attribution obligations;
4. whether build-time, linked, bundled or separate-process use changes duties;
5. release-asset digest and SBOM identity;
6. compatibility with distribution of the MIT-licensed runner;
7. required source offer, notice, disclosure or installation information;
8. action-specific licence boundaries, so a CLI licence never silently
   authorizes a separately licensed integration or action;
9. review date, reviewer, primary-source evidence and unresolved legal gates.

Every accepted release must contain the MIT licence, all exact third-party
licences/notices and an SBOM mutually bound by its signed release manifest.
Missing, stale or conflicting licence evidence blocks release.

## Planned component handling

| Component | PSCAN-01 status | Intake owner | Release rule |
|---|---|---|---|
| Original Go runner and repository content | MIT selected; root licence present | PSCAN-02 onward | Preserve MIT headers/notices where required |
| Gitleaks CLI/source | v8.30.1 historical source/config/licence intake retained; PSCAN-03 rejected | PSCAN-10 | Revalidate exact current primary sources and licence before new intake; no binary released; separately licensed GitHub Action remains absent |
| Go toolchain/modules | Go 1.27.0 historical build binding plus 63 compiled-module licence/notice sets recorded | PSCAN-02/10/06 as applicable | Historical intake is not a release SBOM or distribution approval; reproduce and re-review when rebuilt and at release |
| Cosign/Sigstore tooling and bundles | Design only; no signing performed | PSCAN-06 | Revalidate official licensing and redistribution; keep bundles as signed evidence |
| SBOM/build tooling and GitHub actions | None selected | PSCAN-06 | Pin exact versions/actions and record licence/notice evidence before workflow admission |
| TruffleHog | Absent; no assessment performed | PSCAN-08 only if eligible | No assessment, download, integration, distribution or enablement before accepted material-gap evidence and separate technical plus AGPL legal owner approval |

## Release evidence

PSCAN-06 must produce a machine-readable component inventory, full licence and
notice directory, SBOM, source/binary digest bindings and a human-readable
compatibility review. PSCAN-07 independently rejects any missing component,
notice, source binding or unresolved obligation.

No document in PSCAN-01 is legal advice or final legal approval. The final
retention schedule and every reserved legal question remain owner-controlled.

## PSCAN-03 intake review

On 2026-08-31, Gitleaks v8.30.1 was bound to exact commit
`83d9cd684c87d95d656c1458ef04895a7f1cbd8e`. Its MIT licence, config and every
root-level licence/notice file for the 63 modules in the compiled `go list
-deps` graph were preserved with SHA-256 identities. Two Go 1.27.0 builds per
platform were byte-identical. The upstream tag resolves to an unsigned commit,
so the binding is digest-based and does not claim signed-tag provenance.

No executable, release, SBOM, signature, publication or distribution decision
is made here. PSCAN-03 remains unaccepted because of the material history gap,
independently of the licence intake result.

The PSCAN-09 authority transition performs no intake, download, build,
execution, distribution or legal assessment. PASS-OUTCOME-SPEC-001 preserves
the MIT runner, exact third-party evidence, zero-spend and TruffleHog owner-gate
boundaries. A future activated PSCAN-10 must revalidate Gitleaks rather than
treating historical PSCAN-03 intake as current distribution approval.
