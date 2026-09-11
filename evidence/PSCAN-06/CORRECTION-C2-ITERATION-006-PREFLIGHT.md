# PSCAN-06 Correction C2 iteration 006 preflight and claim

Date: `2026-09-11`
Implementation session: `01a08fe9-ffdf-73e1-a8a9-c5c7d7424ba2`
Authority approval session: `01a08fae-2ee1-7553-a126-e635b70a34a6`

## Claim

This genuinely fresh saved-project session claims only PSCAN-06 Correction C2
iteration 006. It does not claim, select, activate or implement PSCAN-07,
PSCAN-08 or any other successor. Recovery R5 remains terminally failed and its
sole tag push and workflow dispatch remain consumed.

The controlling iteration authority is
`CORRECTION-C2-ITERATION-006-AUTHORITY.md`; its SHA-256 at the claim boundary
was `231A1E8C964CAB8FD73DA1667EFAAEE62AC86C509B0A425A29DC03C6756B3B36`.
The controlling `PASS-OUTCOME-SPEC-001` SHA-256 was
`8A034701867E366BF37ADEA6E4F47E3FFB7A53425463CEFF4FB280BCB4296E74`.
The immutable Recovery R5 failure SHA-256 was
`53C85B1E4ADCFDC3866B9D2194374B7098C296CA11ECDE3050963551566868F0`.

## Starting identity

- Saved-project root:
  `C:/OFFICE-DATA/Projects/Ongoing/project-agnostic-secret-scanner`.
- Direct saved checkout, not a worktree: Git directory and common directory
  both resolved to `.git`.
- Branch: `main`.
- Exact starting HEAD:
  `de162e8c347c725a0f041d3c8b0f51df1211c12d`.
- Exact starting tree:
  `6bbe945af6a12fe6227d1a28c56e254b00dfd82d`.
- Exact direct parent:
  `da45bb19480663aa21e9a67754ad509834dcaed7`.
- The implementation session differs from the authority approval session.
- Tracked and staged diffs were empty; non-ignored untracked status contained
  zero paths.
- All 66 ignored paths were confined to preserved `graphify-out/**`; none was
  deleted, moved or staged.

## Exact source trust

The complete authority-compliant every-byte check proved:

```text
Exact source trust PASS commit=de162e8c347c725a0f041d3c8b0f51df1211c12d tree=6bbe945af6a12fe6227d1a28c56e254b00dfd82d files=360 raw_equal=342 canonical_crlf=18 unexpected_mismatch=0
```

The check verified exact complete tree/index path, mode and object identity;
ordinary and fsmonitor index flags; SHA-1 object format; absence of unsupported
shortcut, sparse or promisor configuration; every tracked working-tree byte;
zero non-ignored untracked paths; and that every ignored path was confined to
the authority-excluded `graphify-out/**` surface.

The stock `invoke-exact-build.ps1 -PreflightOnly` probe was also attempted and
stopped at its stricter all-untracked prohibition because it does not implement
this iteration authority's explicit `graphify-out/**` exception. No ignored
path was removed or changed to make that probe pass. The complete check above
uses the same repository `source-trust.ps1` tree, index, flag, object and
canonical LF-to-CRLF byte algorithms with only the recorded ignored-path
exception.

## Bounded implementation surface

The production and regression sources were read completely before this claim.
The only planned production change is the required local parsed Linux process
identifier rename inside `Get-LinuxSessionMembers`, limited to its binding,
identity string and ledger `PID` value. The only planned regression change is
source-level proof that the function has no case-insensitive `PID` assignment
target and that the renamed binding feeds both dependent uses.

No Docker command, network action, remote read or mutation, dependency
acquisition, build, workflow action, tag change, signing, publication or
successor action occurred. Any need to change a forbidden path or broaden this
design remains terminal and returns to the owner.
