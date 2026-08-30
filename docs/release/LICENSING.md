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
| Gitleaks CLI/source | Not downloaded or pinned | PSCAN-03 | Revalidate current official source and licence before intake; record exact source/binary evidence; do not use the separately licensed GitHub Action |
| Go toolchain/modules | Not downloaded or pinned | PSCAN-02/03/06 as applicable | Record exact toolchain/module licences and reproducible-build evidence |
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
