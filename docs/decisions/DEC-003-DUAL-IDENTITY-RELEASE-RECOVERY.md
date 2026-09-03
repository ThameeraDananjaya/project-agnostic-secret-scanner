# DEC-003: Dual-identity release recovery

## Status

Accepted design for PSCAN-06 Correction C1; local implementation is not
PSCAN-06 acceptance and does not authorize a remote action.

## Context

The immutable product tag `v1.0.0` points to accepted product-source commit
`a13c28fe7273bc8dc6545f97966a02889524eb4c`. The first authorized release
workflow failed before building because its Linux container could not write the
host bind-mounted Go module cache. Correcting that tooling after the product
tag was locked creates two real identities: the bytes shipped as the scanner
product and the later tooling/workflow that reproducibly builds and signs those
bytes. Treating either identity as the other would misstate provenance.

## Decision

Release-manifest schema `2.0` makes both roles mandatory:

- `productSource` is exactly tag `v1.0.0`, commit
  `a13c28fe7273bc8dc6545f97966a02889524eb4c`, and tree
  `217b711ddea51fd0ea7e808edd2e27fdecef8427`;
- `releaseTooling` is tag `release-tooling-v1.0.0-c1`, its exact commit and
  tree, recovery workflow path, workflow ref, workflow SHA and
  `workflow_dispatch` trigger; and
- `releaseIdentity` binds the same recovery workflow ref/SHA/trigger plus the
  exact repository, numeric owner, GitHub OIDC issuer and certificate URI.

The build materializes product source from its exact Git object and builds the
runner, rules and product-owned contracts from those bytes. It builds the
release verifier and packaging controls from the separately identified clean
tooling commit. The verifier embeds the accepted tooling commit/tree at build
time and checks every manifest identity against that out-of-band policy.

The acquisition phase maps the invoking numeric UID/GID only on Linux, retains
a read-only root and dropped capabilities, and runs a synthetic cache write,
same-filesystem atomic rename, exact read and delete before dependency
downloads. Builds mount the completed cache read-only and disable networking.

Schema `1.0` and `1.1` remain immutable. The verifier continues to parse valid
`1.1` evidence under its legacy identity rules, but a `1.1` payload cannot add
or masquerade as the `2.0` roles. Unknown versions reject.

## Consequences and gates

The correction does not move or replace `v1.0.0` and introduces no replacement
product version. The correction-tooling tag may exist only after a later exact
owner gate at an independently accepted tooling commit. Push, tag creation,
workflow execution, signing, attestation, draft creation and publication remain
remote mutations outside local implementation authority.
