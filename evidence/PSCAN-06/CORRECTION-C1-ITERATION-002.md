# PSCAN-06 Correction C1 iteration 002

Candidate `f4ad42664eabedba61226bccdea2bc64e1156b7e` passed every check callable
on the Windows author host. A subsequent skeptical workflow-path inspection
found that the future Ubuntu wrong-owner fixture used root ownership and mode
`0700`. The mapped non-root container would correctly reject its write, but the
non-root host PowerShell process would also be unable to enumerate the
directory afterward to prove that no canary or success ledger remained.

This iteration changes that negative fixture to root ownership and mode `0755`.
The mapped non-root container still cannot write, while the host can inspect the
empty directory. Mode `0777`, privilege, capability addition and host user
namespace remain forbidden and absent. No acquisition, product source, schema,
verifier, signing, remote or successor behavior changes.

The commit containing this record supersedes `f4ad426` as the Correction C1
tooling candidate. Because workflow bytes and the tooling tree change, both
network-disabled builds and all generated identity/digest evidence must be
regenerated from the new exact clean commit. No acceptance is claimed here.
