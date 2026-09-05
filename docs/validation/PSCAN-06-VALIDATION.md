# PSCAN-06 Validation Contract

The exact implementation candidate must pass from a clean commit using the
official Go 1.27.1 archive and full-digest Docker image:

1. revalidate every acquisition digest and populate a ledger only during the
   explicit networked phase;
2. build Gitleaks with Go 1.27.0 and compare both platform binaries with the
   accepted PSCAN-10 digests;
3. build the runner, release verifier and deterministic packaging/SBOM tools
   for Linux amd64 and Windows amd64 with Go 1.27.1;
4. pass all Go tests and all-package vet serially (`-p=1`) so the inherited
   100-millisecond deadline adversarial case is not distorted by unrelated
   concurrent package builds; compile every package/test for Windows amd64 with
   networking disabled and the module cache read-only;
5. build the complete release candidate twice from the same source commit and
   prove every file byte-identical;
6. parse every JSON contract and generated JSON asset, verify SPDX 2.3
   structure, validate archive membership/path safety, reconcile all licence
   files, and reject unexpected/unpinned workflow actions;
7. pass synthetic adversarial cases for manifest, asset, licence, SBOM,
   revocation, identity, issuer, tag and bundle mutation before scanner use;
8. confirm protected authority hashes and allowed-path scope remain unchanged;
9. perform a separate skeptical diff/evidence review.

Native Windows execution, GitHub repository/settings read-back, workflow
execution, signing, attestation, draft upload and publication are not replaced
by local success. Each remains `UNPROVEN` until its separately approved gate.

The final local execution of this contract is recorded in
`evidence/PSCAN-06/VALIDATION.md`. It binds candidate
`a13c28fe7273bc8dc6545f97966a02889524eb4c` and does not promote any reserved
remote or signing state to passed evidence.

## Correction C1 validation addendum

The correction tooling must additionally prove from its exact clean candidate:

1. the cache canary completes write, same-filesystem atomic rename, exact read
   and delete before any dependency download or success ledger;
2. an actual Linux numeric UID/GID positive case passes while wrong-owner and
   read-only bind mounts reject without acquired bytes or a ledger;
3. the acquisition container retains a read-only root, `--cap-drop ALL` and
   no-new-privileges and uses no privileged mode, added capability, mode `0777`
   or host user namespace;
4. schema `2.0` exactly binds the immutable product-source and correction-
   tooling roles, while every required-field omission, identity mutation, role
   swap, old-schema masquerade and unknown major rejects;
5. Cosign verification constrains repository, workflow ref, workflow SHA,
   trigger, certificate identity and GitHub OIDC issuer;
6. two independent network-disabled builds use the completed module cache
   read-only and produce byte-identical assets; and
7. a forced-CRLF checkout proves that the canary, acquisition and build Docker
   shell payloads are converted to LF immediately before invocation, raw
   carriage-return payloads reject, the exact normalized payloads parse in the
   pinned network-disabled image, and the read-only cache still fails closed;
   both complete reproducibility builds must execute from separate actual CRLF
   checkouts while all build inputs are first materialized from and verified
   against the exact product and correction-tooling Git trees, including exact
   path sets, modes and raw blob SHA-256 values; archive, tree, byte, cleanliness
   or extraction ambiguity must fail closed; no integrity asset normalization
   is permitted;
8. every release entrypoint used by the recovery workflow and CRLF harness is
   first materialized from the exact correction-tooling commit and verified
   against its Git blob identity; before output creation, the trusted launcher
   must reject all untracked files (including ignored files), every assume-
   unchanged or skip-worktree entry, filesystem-monitor/untracked-cache or
   partial-clone shortcuts, index path/mode/object differences, missing or
   unsupported paths/modes, and tracked raw-byte differences other than a
   byte-proved canonical LF-to-CRLF checkout projection that is never executed
   or consumed as a build input; exact Git-object materialization remains the
   sole build input; isolated adversarial cases must prove every rejection
   emits zero build files and executes no tampered driver action; and
9. the shipped runner, Gitleaks engines and rules retain the previously
   accepted product-source digests.

Local author validation remains non-acceptance. Independent skeptical review
and every remote/signing/publication owner gate remain open.

## Correction C2 validation addendum

The C1 addendum remains historical and schema `2.0` remains unchanged. The C2
candidate must additionally prove:

1. host positive, wrong-owner and read-only cache cases execute before any
   Docker inspection, execution or pull and leave no canary, ledger or acquired
   bytes on rejection;
2. CRLF checkout, raw-CR rejection and exact LF normalization complete in a
   host-only phase that cannot require an image;
3. image admission accepts only the canonical `docker.io/library/golang`
   repository at digest
   `sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`,
   revalidates Docker `RepoDigests`, and rejects short-name ambiguity, mutable
   tags, wrong repositories, wrong digests, absence without explicit pull,
   pull failure and post-pull inspection mismatch before execution;
4. container cache and shell-parser proofs occur only after admission and use
   `--pull=never --network none`; dependency downloads occur only in the
   acquisition phase; both builds use `--network none` and a read-only cache;
5. schema `2.1`, the C2 tag/ref/workflow/SHA/trigger, and distinct product and
   tooling roles validate exactly, while valid historical schema `2.0` C1
   evidence continues to parse under its own policy; and
6. the locked product and C1 tags, existing schema bytes and all accepted
   product controls remain unchanged.

Author-local proof is not independent acceptance and cannot authorize a pull,
tag, push, workflow run, signing, attestation, draft or publication action.

## Correction C2 iteration 002 validation addendum

The iteration-002 candidate must additionally prove at the Docker command
boundary that:

1. a successful structured engine response and an exact filtered structured
   inventory are required before absence can be concluded; no inspect exit,
   stderr text, ambient exit state, missing output or malformed data is absence;
2. only an empty exact-reference inventory from the responsive engine reaches
   the one exact pull, and the pull is followed by a fresh independent
   repository-digest inspection with no retry;
3. every invocation captures its own exit code, stdout, stderr and timeout, with
   bounded output and rejection of any invalid boundary result;
4. the full `RepoDigests` value is an array containing exactly one canonical
   engine repository-digest string; null, scalar, empty, malformed, aliases,
   duplicates, mixed exact-plus-other, wrong-repository and wrong-digest values
   all reject;
5. a deterministic fake-engine harness proves exact command count, arguments
   and order for present, absent, daemon, permission, timeout, invalid-command,
   deceptive-text, malformed, pull-failure and post-pull failure cases, with
   zero later action for every rejection; and
6. the only pull-capable call is inside the mandatory orchestrator after its
   fixed host-cache and host-only CRLF prerequisites. Legacy acquisition
   verification cannot pull, even when its historical switch name is supplied.

These checks are network-free author validation only. Actual Linux execution
and independent skeptical acceptance remain mandatory.

## Correction C2 iteration 003 validation addendum

Iteration 003 additionally requires:

1. static and runtime rejection of callback, variable, function, alias,
   script-block, dot-source, PATH and environment substitution against the
   workflow-facing admission entrypoint;
2. the complete iteration-002 present, conclusive-absence, identity and
   untrusted state matrix through a structurally separate test-only model that
   shares only immutable production parsing rules;
3. native no-Docker fixtures at below, exact and above 131072 bytes on each
   stream, including simultaneous streams, with live byte rejection and no
   truncation;
4. strict UTF-8 proof for read-boundary splits, invalid sequences and incomplete
   final sequences;
5. immediate, nonzero, start-failure, hang, flood, simultaneous-flood and
   child/grandchild cases that return within the fixed 15000 ms command budget
   plus 2000 ms cleanup grace and leave no live fixture descendant; and
6. terminal representation of timeout, overflow, read, termination and pipe
   uncertainty, with zero pull and zero later action.

The native suite must run on Windows and actual Linux before independent
acceptance. A local author run on only one platform is expressly incomplete.

## Correction C2 iteration 004 validation addendum

Iteration 004 additionally requires:

1. exact static inventory showing one Docker process-creation boundary and no
   ambient Docker call or duplicate runner in workflow-reachable scripts;
2. a finite operation table whose complete arguments preserve the admission,
   cache, CRLF, acquisition, network-disabled build and packaging boundaries;
3. held-file executable identity, private empty configuration and work
   directories, a cleared minimum environment and rejection of replacement,
   PATH, alias, function, callback, arbitrary-argument and dot-source bypasses;
4. Windows suspended start followed by kill-on-close job assignment before
   resume, and actual-Linux private-session execution with identity-safe member
   tracking before Docker bytes execute;
5. success only after root exit, both bounded streams close and containment is
   empty; timeout, overflow, invalid UTF-8, read, exit, termination, membership
   or cleanup uncertainty remains terminal under the exact 131072-byte,
   15000 ms and 2000 ms limits; and
6. no-Docker native fixtures covering immediate/nonzero/start failure, exact
   stream boundaries, simultaneous streams, split/invalid/incomplete UTF-8,
   hang, child, detached child, grandchild, replacement race and terminal
   cleanup without a live or uncertain member.

Windows author proof alone remains incomplete. Actual Linux, genuine Docker
operation, reproducibility and independent skeptical proof remain mandatory
for acceptance and require their separately authorized execution boundaries.
