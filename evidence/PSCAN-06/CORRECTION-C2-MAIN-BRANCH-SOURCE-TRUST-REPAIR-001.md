# PSCAN-06 Correction C2 main-branch source-trust repair 001

## Result and boundary

After Recovery R1 failed closed, the owner directed autonomous continuation of
the pending local work while excluding Docker start/stop and external service
purchases. This record covers only the local working-tree projection repair.
It does not reinterpret the failure, authorize a retry, or grant any remote,
Docker, signing, attestation, draft, publication or successor authority.

The two noncanonical tracked working-tree projections were restored from their
exact committed Git objects. No committed source, index entry, product logic or
workflow byte changed.

## Exact repair evidence

- Repository: `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`
- Branch: `main`
- HEAD: `ca996eed20294a1d36822dd453005a333218279a`
- Repair method: exact two-path materialization from the `HEAD` Git archive,
  followed by Git index metadata refresh.
- Temporary archive: removed after the bounded materialization.

### `build/gitleaks/build.ps1`

- Committed blob: `7d66c57c27767464f8a084211f37db50a85e4702`
- Pre-repair raw working-tree blob:
  `be1241c8a2ba77f4140270750264d1dd54100893`
- Post-repair raw working-tree blob:
  `314e9e30166e2c1a0ade5b0ee4928ddfa07980da`
- Post-repair bytes: `2937`
- Canonical LF-to-CRLF projection: `true`
- CRLF sequences: `59`
- Bare LF bytes: `0`
- Bare CR bytes: `0`

### `build/gitleaks/collect-licenses.ps1`

- Committed blob: `6b5fd1e520c0d8be694461397e162b0447164ac8`
- Pre-repair raw working-tree blob:
  `ef95129091f5c9c33900ca20423c42d1eaa4b645`
- Post-repair raw working-tree blob:
  `fcc203210a6d4c2ef0e341a423a2881a2df12dce`
- Post-repair bytes: `2482`
- Canonical LF-to-CRLF projection: `true`
- CRLF sequences: `48`
- Bare LF bytes: `0`
- Bare CR bytes: `0`

## Validation and remaining gate

- Immediately after the two-path repair and before this evidence update,
  `git diff --quiet`: exit `0`.
- Immediately after the two-path repair and before this evidence update,
  `git diff --cached --quiet`: exit `0`.
- Preserved Recovery R1 failure-record SHA-256:
  `1CC095F7FD16CF6CD8F004371D42F38BB24AFBEAF4F059B99992CFBB8B56BAF6`.
- The only intended new paths are this append-only repair record and the
  preserved Recovery R1 failure record.
- No Docker, network, GitHub, remote repository, workflow, artifact, signing,
  attestation, draft, publication, credential, purchase or successor action
  occurred.

The local source-trust obstruction is repaired, but Recovery R1 remains a
terminal failed attempt and authorizes no retry. PSCAN-06 remains open and
unaccepted overall. Any new proof-gate execution requires its own exact
recorded authority and a genuinely fresh execution session.
