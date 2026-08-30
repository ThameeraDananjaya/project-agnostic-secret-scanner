# Gitleaks v8.30.1 Coverage Map

| Required class | Primary-engine observation | PSCAN-03 disposition |
|---|---|---|
| Ordinary tracked head files | Gitleaks `dir` consumes file bytes | Draft adapter plus exact complete manifest; not integrated |
| Text additions/deletions in exact Git range | Gitleaks `git` consumes `git log -p -U0` patches | Inspectable when Git emits patch text; exact range is bound |
| Binary-classified blobs present at head | Directory mode can receive current file bytes | Requires tracked-source manifest and size preflight |
| Binary-classified blob added/deleted only in history | Ordinary `git log -p` omits blob body | **Material gap; incomplete coverage; no PASS** |
| File above declared maximum | Gitleaks can skip it | Must be rejected before launch as resource-limit non-pass; no silent pass |
| Nested archives and container layers | Adapter forces archive depth zero | PSCAN-04 only; unavailable in PSCAN-03 |
| Git metadata | Native Git parser may observe commit/diff metadata | Raw output remains private; complete metadata claim not accepted |
| Symlink, special, unreadable or unmanifested source | Not safe to infer engine behavior | Manifest boundary rejects; explicit non-pass |

The binary-history row is accepted evidence of a primary-engine gap only. It
does not select a remediation, authorize a fallback engine or activate
PSCAN-08.
