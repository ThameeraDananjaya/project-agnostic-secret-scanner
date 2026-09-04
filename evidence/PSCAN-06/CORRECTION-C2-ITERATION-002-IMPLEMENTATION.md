# PSCAN-06 Correction C2 iteration 002 implementation

Date: 2026-09-05

Authority base: `281bea031bb6bbaf1be3059977074df2a88ecdf4`

State: bounded local author candidate; not independently accepted

## Implemented boundary

- Replaced generic inspect-failure-as-absence with three explicit states:
  exact present, conclusive absence and untrusted failure. Conclusive absence
  requires a structured responsive-engine result followed by an empty
  exact-reference structured inventory.
- Every Docker invocation returns one boundary object containing that call's
  arguments, exit code, stdout, stderr and timeout state. Invalid contracts,
  nonzero exit, timeout, unexpected stderr, over-limit output and malformed
  structured data reject before state transition.
- The complete `RepoDigests` JSON value must be an array containing exactly one
  canonical engine identity. Null, scalar, empty, malformed, alias, duplicate,
  mixed, wrong-repository and wrong-digest evidence rejects.
- The only pull-capable command is inside `Invoke-ReleaseImageBootstrap`, after
  fixed host cache and host-only CRLF prerequisites. It permits one pull of the
  exact digest only for conclusive absence, then performs a fresh independent
  identity inspection without retry.
- The recovery workflow invokes this orchestrator with the exact source,
  revision and CRLF work directory. Legacy acquisition verification remains
  compatible only as a pull-free presence check; any historical pull switch
  request fails before Docker.
- The fake-engine harness records exact prerequisite/Docker order, command
  count and arguments and proves zero later action for all rejected states.

## Preserved boundary

No product source, locked tag, manifest schema, verifier, signature policy,
source-trust control, cache/container sandbox, `--pull=never`, `--network none`,
acquisition/reproducibility order, dependency, licence, SBOM, credential,
remote, publication or successor boundary changed.

This implementation record is an author claim only. Exact executed checks,
limitations and the implementation commit are recorded in the separate author-
validation record. Independent skeptical review and actual Linux execution
remain required before any acceptance or remote-gate proposal.
