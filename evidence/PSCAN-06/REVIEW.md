# PSCAN-06 Skeptical Review

## Review boundary

This was a separate read-only review pass after the definitive artifacts were
produced. It compared the activation commit, candidate diff, task path boundary,
workflow permissions/gates, verifier trust flow, generated manifests, archive
membership, SBOM/licence evidence, tests and reserved actions. It did not rely
on build success alone.

## Findings resolved before candidate freeze

1. Cosign v3.1.3 does not support the legacy `--offline` flag. The verifier now
   requires a separately acquired trusted root and its exact SHA-256, supplies
   it with `--trusted-root`, and relies on an enforced network-disabled caller
   boundary.
2. Current `actions/attest` v4.2.2 requires `artifact-metadata: write`; that
   permission is present only in the gated signing/attestation/draft job.
3. The trusted-root mutation test was initially placed outside the task's
   allowed test paths. The final candidate contains it only under
   `tests/integration/supply-chain/**`; the activation-to-candidate path audit
   passes.
4. The signing/draft gate initially required a PSCAN-07-valued variable. The
   final workflow instead requires exact value `PSCAN-06-SIGNING-APPROVED`.
   PSCAN-07 remains unselected and is required only for final publication.

## Final review result

- No unresolved local correctness, scope, pinning, redaction, licensing,
  reproducibility or fail-closed defect was found in candidate
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`.
- Build jobs are read-only and credential persistence is disabled. Only the
  still-unapproved gated job requests OIDC/attestation/draft permissions.
- There is no publish command, long-lived signing key, consuming-project data,
  mutable action reference, paid runner, TruffleHog material or Git remote.
- The complete plain-English offline verification and scanner input/output
  references include Mermaid diagrams and bounded parameters/payloads.
- Native Windows execution and every remote/signing fact remain explicitly
  unproven.

The candidate is locally ready for the exact owner gate, but PSCAN-06 is not
accepted or closed until remote identity/settings and a real keyless bundle,
attestation and draft are independently read back and verified.
