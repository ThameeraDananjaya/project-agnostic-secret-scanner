# Repository Working Agreement

## Controlling authority

PASS-SPEC-001 is the controlling product contract. Work may add stricter
controls but must not weaken, reinterpret, omit, or silently replace any
capability, constraint, non-goal, security boundary, schema boundary,
acceptance test, owner gate, or definition-of-done requirement.

The source document and verified digest are recorded in
`docs/spec/PASS-SPEC-001-SOURCE.md`. PSCAN-01 must place the canonical contract
inside this repository without changing its meaning.

## Product independence

This repository must contain no consuming-project source, policy, allowlist,
receipt, key, finding, credential, customer data, repository identity,
statistics, or project-specific behavior. Consuming projects are callers, not
product authority.

## Task lifecycle

- Work on exactly one activated PSCAN task per session.
- Exact activation syntax is a raw standalone line: `ACTIVATE PSCAN-XX`.
- An activation session creates only the activation bundle and stops.
- Implementation begins in a fresh session from the activation commit.
- Acceptance is evidence-based and path-bounded.
- No successor activates automatically.
- Proposed, selected, activated, claimed, implemented, accepted, and completed
  are distinct states.

The one-time empty-repository bootstrap and PSCAN-01 activation are authorized
by the owner's 2026-08-30 consolidated decision. This exception cannot activate
or implement any later task.

## Safety and owner gates

Do not create credentials, signing keys, paid subscriptions, remote resources,
GitHub repositories, releases, settings, or public artifacts without the
applicable recorded owner approval and live preflight. Do not download or enable
TruffleHog unless accepted evidence proves a material Gitleaks coverage gap and
the owner separately approves both the technical need and AGPL obligations.

Preserve user work. Inspect Git state first, stage only approved paths, never
reset or clean unrelated work, and never push or publish unless explicitly
authorized for that exact action.

Every missing, stale, incomplete, conflicting, unsupported, or untrustworthy
condition fails closed.
