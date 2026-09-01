# PSCAN-04 Validation Contract

## Capability and acceptance mapping

| Contract | PSCAN-04 proof |
|---|---|
| CAP-3 / AT-07 / AT-08 | Private admission ledger, raw and expanded member projections, exact raw-range payload verification, class-isolated directory/ZIP/TAR/gzip/TGZ/OCI/Docker-save tests and actual pinned-Gitleaks clean/finding tests for every class |
| CAP-4 / AT-13 | Docker `--network none`, loopback-only assertion, failed DNS/public-TCP/cloud-provider metadata canaries, sealed environment and no credential-like keys |
| CAP-7 / AT-09 / AT-10 | Raw pre-transform magic/framing checks; unsafe, encrypted, malformed, unsupported, ambiguous and every over-limit condition is non-pass |
| CAP-8 / AT-11 / AT-12 | Typed diagnostic whitelist, authoritative outcome preservation, private engine capture, panic recovery, derived-canary oracle, injection rejection and one-document validation |
| T-05 / T-19 / AT-22 | Exact PR/release profile constants, isolated below/at/above arithmetic boundary tests, pre-extraction archive limits, exact Docker cgroup memory/no-swap checks, cleanup for every returned result and marker-plus-lease process-death recovery |

## Required commands

The authoritative command is `build/sandbox/validate.ps1` with the read-only
task module cache and exact Linux Gitleaks binary. It executes `go test
-count=1 ./...` and `go vet ./...` inside the pinned Go 1.27.0 container.
Separate Windows amd64 compile-only commands must build the artifact, cleanup
and existing workspace test packages with `CGO_ENABLED=0`; this does not claim
Windows runtime acceptance.

Acceptance additionally requires path-diff checks against the activation
commit, a clean `git diff --check`, no dependency graph change, exact detector
and container digests, and independent review of the completed patch.

The accepted claim is bounded to the internal PSCAN-04 stage. The public runner
remains outside the allowed paths and unavailable, so this validation cannot
be cited as public CLI, consumer, deployment, complete cross-platform or
production readiness.
