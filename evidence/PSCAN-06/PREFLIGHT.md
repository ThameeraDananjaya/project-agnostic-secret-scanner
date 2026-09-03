# PSCAN-06 Current Preflight Evidence

## Event and baseline

- Preflight UTC: `2026-09-03T00:10:04Z`.
- Branch: `main`.
- Exact activation commit: `058ffcd446c6431b2e1afeed769d02c7b1f307f8`.
- Activation parent: `30d27856bdfd404fc190501be854b78bc5147f1f`.
- Activation tree: `1cfc26a1634afd010f2579b20754f4b648a13c45`.
- Git state before preflight: clean (`## main`).
- Git remotes: none configured.
- Activation diff: only the four PSCAN-06 activation-bundle paths recorded in
  `evidence/PSCAN-06/ACTIVATION.md`.
- `ACTIVATE PSCAN-06` was not run, replayed or recorded in this implementation
  session.
- The complete ordered reading map and all 60 listed authorities were read
  before this claim. Every tracked file already present under the allowed
  implementation paths was then inspected before editing: 96 files totaling
  377,544 bytes.

## Current official dependency and action facts

All references below were read-only checks against official project or provider
sources on 2026-09-03. Intake is limited to these immutable identities; mutable
major tags, `latest`, auto-update and pipe-to-shell installation are forbidden.

| Item | Selected immutable identity | Current primary-source result |
|---|---|---|
| Go | `go1.27.1`; Linux amd64 archive SHA-256 `63d339f0da5ab53635a56f2490a7984dfe12dfcff22ad749f63edaf590168445`; Windows amd64 archive SHA-256 `a3911b5e0e1b1053f25ed0675f4c1c6aad1e2bfcf253df2b9be4caabd2edd95d` | Current stable release in `https://go.dev/dl/?mode=json` |
| Gitleaks | `v8.30.1`; commit `83d9cd684c87d95d656c1458ef04895a7f1cbd8e`; MIT | Current upstream release; exact existing PSCAN-10 engine pin remains valid |
| Gitleaks Linux asset | `gitleaks_8.30.1_linux_x64.tar.gz`; SHA-256 `551f6fc83ea457d62a0d98237cbad105af8d557003051f41f3e7ca7b3f2470eb` | GitHub release API asset digest |
| Gitleaks Windows asset | `gitleaks_8.30.1_windows_x64.zip`; SHA-256 `d29144deff3a68aa93ced33dddf84b7fdc26070add4aa0f4513094c8332afc4e` | GitHub release API asset digest |
| Cosign | `v3.1.3`; commit `11926fa5bbbbde47e88fc006b625a17769b743b2`; Apache-2.0 | Current upstream release |
| Cosign Linux verifier | `cosign-linux-amd64`; SHA-256 `4629c757b7618056f8ddd7e2625ae9fdd94c0372a65049520bc7d9df9efc7f71` | GitHub release API asset digest |
| Cosign Windows verifier | `cosign-windows-amd64.exe`; SHA-256 `9fe59be0eca1271873ce019061335eb1ac419b7059202e797828467ddabe33be` | GitHub release API asset digest |
| `actions/checkout` | `3d3c42e5aac5ba805825da76410c181273ba90b1` (`v7.0.1`); MIT | Exact tag commit resolved through GitHub API |
| `actions/attest` | `1e69f48acb82d1966a394da916b4c1698aa569d6` (`v4.2.2`); MIT | Exact tag commit resolved through GitHub API |
| `actions/upload-artifact` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` (`v7.0.1`); MIT | Exact tag commit resolved through GitHub API |
| `actions/download-artifact` | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` (`v8.0.1`); MIT | Exact tag commit resolved through GitHub API |
| `sigstore/cosign-installer` | `6f9f17788090df1f26f669e9d70d6ae9567deba6` (`v4.1.2`); Apache-2.0 | Exact tag commit resolved through GitHub API |
| SBOM | SPDX 2.3 JSON | Official SPDX specification; deterministic document generation is local |

The existing Gitleaks licence corpus was independently reconciled before claim:
63 modules, 64 licence/notice files, zero extras, zero missing files, zero
digest mismatches, zero UTF-8 or carriage-return failures, and zero unclassified
licences. Family counts are Apache-2.0 11, BSD-style 20, CC0-1.0 1, MIT 27 and
MPL-2.0 5.

## Account, repository, capability and cost preflight

- Authenticated GitHub identity: `ThameeraDananjaya`, immutable user ID
  `50274860`, user account, not site administrator.
- Authentication is present for read-only preflight. The existing classic token
  reports `gist`, `read:org`, `repo` and `workflow` scopes. It is not copied,
  stored, printed, used by a build, or used for any mutation.
- Exact intended repository:
  `ThameeraDananjaya/project-agnostic-secret-scanner`.
- Read-only repository lookup result: HTTP 404. No repository ID, settings,
  ruleset, environment, workflow, release or immutable-release state exists to
  verify yet. Name availability is therefore provisional, not reserved.
- Intended visibility remains public under controlling authority. GitHub's
  current billing documentation states that standard GitHub-hosted runners are
  free for public repositories. Larger runners and storage overages are not
  permitted. The workflow design uses only standard Ubuntu runners and short
  retention; zero spend remains an action-time invariant.
- Public repositories on GitHub Free can use rulesets. Repository rules,
  environment protection, workflow permissions, immutable releases and any
  spend/budget controls must be read after repository creation and configured
  only under a separate exact owner approval.
- GitHub's current immutable-release control is repository configuration. When
  enabled, the release tag and assets become immutable after publication and a
  release attestation is generated. The required flow is draft first, attach
  all assets, independently verify, then obtain a separate exact publication
  approval.
- GitHub artifact attestations on Free/Pro/Team are available for public
  repositories. They require `id-token: write`, `contents: read` and
  `attestations: write`. The release workflow gives these permissions only to
  its gated provenance job; build/test jobs remain `contents: read`.
- GitHub Actions OIDC issuer is exactly
  `https://token.actions.githubusercontent.com`. The keyless certificate
  identity is restricted to the intended repository, exact release workflow
  path and exact `refs/tags/v1.0.0` invocation. Repository numeric identity and
  actual workflow claims must be captured and compared after the remote exists;
  missing or changed identity fails closed.

## Local capability preflight

- Native Go is not installed or callable. No native downloaded executable was
  attempted.
- Docker client `29.7.2` is installed with context `desktop-linux`, but the
  Linux engine pipe was not reachable at the preflight timestamp. The already
  pinned Go image could not be inspected. Docker Desktop may be started and used
  as the authorized isolated build/test path; no third-party binary runs before
  this preflight record.
- No scanner, toolchain, dependency, workflow or signing command ran during
  preflight.

After claim, Docker Desktop 4.87.0 started successfully with Engine 29.7.2.
The existing pinned image
`golang@sha256:ded31c68586d2e49e760acc2e65a884b23d032e9bbbed0ae0c55abd3fcaf4452`
resolved locally and reported Go 1.27.0. The explicit acquisition phase then
downloaded the official Go 1.27.1 Linux archive, official Go 1.27.0 Linux
archive and exact Gitleaks source archive, verified all three recorded SHA-256
values, and populated the module cache with network access. Its first host
command timed out after the downloads/cache were present but before the ledger;
the resumable script revalidated existing bytes, completed the cache operation
and wrote `acquisition-ledger.json`. No unverified partial file was admitted.

## Exact proposed remote changes still requiring approval

No item in this section is authorized by activation or this preflight:

1. Create the public repository at the exact intended owner/name.
2. Configure Actions to allow only required full-SHA-pinned actions, set the
   default workflow token to read-only, configure zero-spend/retention controls,
   and configure the exact protected `release` environment.
3. Configure main/tag rules and enable immutable releases, then read them back.
4. Add the remote and push the exact accepted local history.
5. Enable or invoke any workflow.
6. Run a keyless signing dry run or actual signing event.
7. Create a draft release, upload assets, attest, publish or create `v1.0.0`.

Each step needs a separate exact action-time owner approval and a fresh
read-back preflight. Final `v1.0.0` publication belongs to PSCAN-07, which
remains proposed and unselected.

## Fail-closed conclusion

Local implementation, synthetic tests, deterministic local packaging and
offline negative-path verification may proceed at zero cost. No remote fact is
promoted to passed evidence. Absent repository identity/settings, absent valid
keyless bundle and absent publication evidence remain explicit `UNPROVEN` or
`NOT_PERFORMED` states and cannot produce a release-verification pass.
