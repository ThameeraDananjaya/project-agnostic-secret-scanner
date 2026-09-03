# PSCAN-06 Correction C1 implementation record

## Claim

- Authority commit:
  `3fb1b0a55dc4f48dd35464c63c768f497efbc89b`
- Initial branch/status: `main`, clean
- Claimed work: PSCAN-06 Correction C1 only
- Candidate identity: the commit containing this implementation bundle
- Independent acceptance: not claimed
- Remote action: none

## Implemented boundary

- Pre-download Linux numeric UID/GID cache ownership and semantic canary,
  preserving a read-only root, dropped capabilities and no-new-privileges.
- Explicit Linux positive, wrong-owner and read-only-cache regression harness;
  failed cases cannot create an acquisition ledger.
- Two network-disabled builds using the completed module cache read-only, with
  a complete byte comparison helper.
- Locked product-source materialization from tag `v1.0.0`, commit
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`, tree
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`.
- Separate correction-tooling tag/commit/tree/workflow/ref/SHA/trigger binding
  in release-manifest schema `2.0` and the offline trust policy.
- Exact Cosign certificate claim checks for repository, ref, workflow SHA,
  trigger and issuer in addition to certificate identity.
- A full-SHA-pinned, least-permission recovery workflow whose build half runs
  only at the proposed tooling tag and whose signing half remains behind a new
  exact owner-controlled gate.
- Adversarial schema and policy tests for omission, masquerade, unknown major,
  mutation and tag-role swapping.

The historical failed run, original implementation evidence, existing schemas
and product behavior were not rewritten. The correction introduces no scanner
detection, policy, allowlist, receipt, revocation, consuming-project or
successor work.

## Validation state at implementation commit

PowerShell AST parsing, schema JSON parsing and focused Go integration and
acceptance tests passed while authoring. Definitive clean-commit acquisition,
Linux cache regression, two-build reproduction and complete test evidence are
recorded separately after the candidate commit. Passing author checks do not
constitute independent skeptical acceptance.
