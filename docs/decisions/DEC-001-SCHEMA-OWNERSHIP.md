# DEC-001: Scanner-Owned Schemas and Project-Owned Policy Projection

- Status: accepted planning decision
- Date: 2026-08-30
- Authority: PASS-SPEC-001 sections 7.7, 9, 13, 16, 17 and 20
- Implementation: deferred to PSCAN-02 and PSCAN-05

## Decision

The scanner product owns and versions only these authoritative schema families:

- `scan-request`
- `scan-outcome`
- `scanner-release-manifest`
- `global-scanner-revocation`
- generic rule-pack schema

Each consuming project owns and versions its policy schema and instances,
allowlist schema and instances, receipt schema and instances, receipt-revocation
schema and instances, evidence-chain format, custody and retention rules.

The scanner will accept project-owned policy and allowlist inputs through a
strict, content-addressed projection boundary. The request binds their digests
and declared schema identifiers. A versioned projection adapter maps only the
minimum scanner-relevant concepts into an internal evaluation model. The
project artifact remains authoritative; the scanner repository does not copy,
publish or reinterpret project instances.

## Precedence

```text
non-overridable scanner safety invariants
    -> global scanner-release revocations and mandatory protections
    -> validated project-policy projection
    -> narrow validated project allowlist exceptions
```

Lower layers cannot weaken higher layers. Unknown major versions, missing
required fields, ambiguous mappings, digest/signature failures, unsupported
semantics, or conflicting projections fail closed. A genuine credential cannot
be allowlisted, and severity cannot downgrade a mandatory block.

## Compatibility

Every schema family uses an independent `major.minor` identifier. Compatible
minor additions are accepted only when unknown fields are safely ignorable and
meaning is unchanged. Unknown majors fail closed. Retirement preserves
readability for at least 90 days or two successful consumer release cycles,
whichever is longer, with a 180-day hard ceiling absent owner-approved
extension. Signed retirement notices name the replacement and last-supported
date.

## Consequences

- The same scanner release can serve structurally different projects without
  sharing policy or state.
- The product can publish reference examples for project-owned formats, but
  those examples are synthetic and explicitly non-authoritative.
- Receipt issuance and receipt-key handling remain outside the scanner.
- Scanner release trust cannot sign project receipts.
- PSCAN-02 defines scanner-owned schemas; PSCAN-05 implements and adversarially
  tests policy projection, precedence, compatibility, allowlists and isolation.
