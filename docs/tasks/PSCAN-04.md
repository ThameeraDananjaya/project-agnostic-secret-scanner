# PSCAN-04 Task Specification

## Title

Artifact normalizer, OCI/archive coverage, redaction firewall, cleanup, and
network-disabled execution.

## State

Activated by the owner's exact `ACTIVATE PSCAN-04` command on 2026-09-01 and
claimed alone on 2026-09-02 from exact clean activation commit
`a21b030e1658f1f98ac4e4d001af12185d9ed311`. Implemented, corrected,
independently accepted and closed locally on 2026-09-02. No successor is
selected or activated.

## Objective

Safely normalize all initial release artifact classes without execution, make
redaction a mandatory boundary on every output path, destroy transient material
for every lifecycle outcome and prove authoritative engine execution has no
network or credentials.

## Preconditions

- PSCAN-10 outcome-proved primary coverage accepted, closed and clean at
  activation parent `14348b16d9ddb87f337391a4e47dd243f6e53bf2`.
- Exact PSCAN-04 owner selection and activation recorded in
  `evidence/PSCAN-04/ACTIVATION.md`, followed by a fresh implementation session
  from that activation commit.
- Supported archive/OCI parsing choices and licences reviewed from current
  primary sources before dependency intake, build, execution or admission.
- Existing pinned Gitleaks, rules, preparation/proof boundary and toolchain
  identity revalidated before scanner or test execution.

## Activation-only bundle

This session may change and commit only:

```text
docs/tasks/PSCAN-04.md
docs/tasks/PSCAN-04-READING-MAP.md
docs/tasks/TRACKER.md
evidence/PSCAN-04/ACTIVATION.md
```

No implementation, dependency intake, scanner/toolchain execution, test run or
successor work is authorized in the activation session.

## Allowed implementation paths

```text
go.mod
go.sum
internal/artifact/**
internal/redaction/**
internal/cleanup/**
internal/workspace/**
internal/engine/**
fixtures/archives/**
fixtures/containers/**
fixtures/adversarial/artifact/**
fixtures/adversarial/redaction/**
fixtures/adversarial/network/**
tests/unit/artifact/**
tests/unit/redaction/**
tests/unit/cleanup/**
tests/integration/artifact/**
tests/integration/offline/**
tests/acceptance/redaction/**
build/sandbox/**
licenses/**
THIRD_PARTY_NOTICES.md
docs/architecture/**
docs/security/**
docs/validation/**
docs/tasks/PSCAN-04.md
docs/tasks/TRACKER.md
evidence/PSCAN-04/**
```

## Forbidden scope

- No modification of PASS-OUTCOME-SPEC-001, PASS-SPEC-001, their source record,
  TRACEABILITY, DEC-001, DEC-002 or PSCAN-01 through PSCAN-03 and PSCAN-09/10
  historical evidence.
- No execution, build or dependency installation from candidate artifacts;
  candidate material remains data only.
- No project policy, allowlist, receipt, revocation, evidence-custody,
  deployment or retention authority.
- No TruffleHog assessment, legal work, download, material, integration,
  distribution or enablement.
- No GitHub workflow, remote creation, push, publication, release, signing,
  settings change, credential handling, spending, provider, deployment,
  production or go-live action.
- No real consuming-project artifact, source, finding, credential, customer
  data, repository identity, statistics or project-specific behavior.
- No selection, activation, claim or implementation of PSCAN-05, PSCAN-08 or
  any other successor.

## Deliverables

- Directory, ZIP, TAR, gzip/TGZ, OCI layout and Docker-save normalization with
  deterministic supported nesting.
- Rejection of traversal, absolute paths, unsafe links, reparse/device/hard-link
  hazards, duplicate/case collisions, encryption, malformed inputs and bombs.
- Exact depth/entry/expanded/file/ratio limits and explicit resource outcomes.
- Private engine channel and mandatory redaction of value, partial/encoded/hash/
  fingerprint, path/location, detector/provider, Git and allowlist details.
- Cleanup on pass, fail, indeterminate, timeout, cancellation and crash.
- Read-only mounts/root, network disabled, no secrets/tokens, dropped capability,
  no-new-privileges and bounded process/resource proof.

## Acceptance

Supported inputs are fully covered; unsupported and over-limit inputs never
skip/pass; canaries leak through no surface; output injection remains valid and
content-free; DNS/TCP/provider canaries cannot connect; every lifecycle leaves
no transient residue; CAP-3, CAP-4, CAP-7, CAP-8 and related acceptance rows pass
for this bounded implementation.

Stop after PSCAN-04 closeout with PSCAN-05 unselected.
