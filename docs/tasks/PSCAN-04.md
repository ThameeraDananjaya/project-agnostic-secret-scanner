# PSCAN-04 Task Specification

## Title

Artifact normalizer, OCI/archive coverage, redaction firewall, cleanup, and
network-disabled execution.

## State

Proposed. Unselected. Not activated. Not claimed.

## Objective

Safely normalize all initial release artifact classes without execution, make
redaction a mandatory boundary on every output path, destroy transient material
for every lifecycle outcome and prove authoritative engine execution has no
network or credentials.

## Preconditions

- PSCAN-03 accepted and clean.
- Exact PSCAN-04 activation and fresh implementation session.
- Supported archive/OCI parsing choices and licences reviewed before admission.

## Allowed paths

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

No candidate execution/build/dependency install; no project policy or allowlist;
no TruffleHog; no receipt/signing; no GitHub workflow or remote action; no real
project artifact, finding, credential or customer data.

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
