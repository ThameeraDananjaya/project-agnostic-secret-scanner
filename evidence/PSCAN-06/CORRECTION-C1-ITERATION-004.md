# PSCAN-06 Correction C1 iteration 004 contradiction and repair record

## Independent rejection carried forward

Candidate `5ca77226ed3996a8267caf02da366d0beb915c8d` is rejected as a
portable Correction C1 candidate. A complete independent build from an actual
Windows CRLF checkout passed the shell-payload regression but failed
`TestPinnedRuleAndCoverageIntegrityBindings`. The build consumed transformed
working-tree bytes through its `/src` bind mount instead of exact Git-object
bytes.

Local reproduction from exact candidate `5ca7722` proved that
`rules/generic/gitleaks-ignore-empty-v1.txt` is LF in Git with SHA-256
`3cdd737e34cb301ffc19861f903cf59441ff57453ddf70bd5b8c1bd63389aebe`,
but a `core.autocrlf=true` checkout materializes CRLF with SHA-256
`256b2dc08c5552cfce8f82b39ad50919017e249894391049d2f4ddf6fb2e9fa5`.
Plain `git archive` also inherited the host conversion setting, so an
unchecked archive is not accepted as exact-tree evidence.

The passing author runs recorded in
`CORRECTION-C1-AUTHOR-VALIDATION-003.md` remain historical evidence for the LF
checkout that was executed, but they do not prove checkout-independent byte
identity and cannot override the reproduced failure. This record is
append-only; no earlier evidence is rewritten.

## Bounded repair

Iteration 004 removes the correction-tooling working-tree bind from the build.
It materializes both the locked product tree and exact correction-tooling tree
with checkout conversion disabled, while separately reading every expected
blob through raw `git cat-file blob`. Before any compile, test, schema,
verifier or packaging action, the pinned network-disabled container requires:

- the expected commit-to-tree relationship;
- a clean exact source checkout;
- supported regular-file modes and conservative Git paths only;
- exact archive SHA-256 at the container boundary;
- successful archive extraction;
- an exact extracted path set and file count;
- exact executable/non-executable modes; and
- every extracted file SHA-256 equal to its raw Git blob bytes.

Schema 2.0 and release documentation copied into the bundle come only from the
verified materialization. The verifier, packager, SBOM tool and all Go tests
run only from the verified exact correction-tooling tree. Product runners,
rules, policies, schemas, licences and integrity assets come only from the
separately verified locked product tree. No integrity asset is normalized or
rewritten as a substitute for exact Git bytes.

The forced-CRLF regression now creates an actual local
`core.autocrlf=true` checkout of the exact candidate, proves the integrity
asset contains a carriage return there, and can execute the focused pinned-
asset integrity test plus the complete offline release build. The recovery
workflow performs both independent release builds through separate actual
CRLF checkouts.

The first exact-tree author attempt, commit `29e2b43`, failed closed before
compilation after all blob-byte and path-set checks passed but the extracted
mode check found Git's Windows tar default had emitted `0664` for a `100644`
blob. No build result from that attempt is accepted. The bounded correction
pins archive `tar.umask=0022`; the exact-mode verifier remains mandatory and
must prove `0644`/`0755` after extraction.

The next clean candidate `dcde01d` passed two complete forced-CRLF builds and
the 30-file distributable comparison. Before final evidence, inspection found
that archives created from bare tree objects carried run-time metadata and the
generated test summary omitted the actually executed `-p=1` flags. The final
bounded refinement archives the already verified commit object, making archive
metadata deterministic while the raw-blob/path/mode checks remain controlling,
and records each executed Go command exactly. Results from `dcde01d` are not
promoted to final candidate evidence; both builds must be repeated from the
new exact commit.

The commit containing this record supersedes `5ca7722` as the Correction C1
tooling candidate. Two clean forced-CRLF builds and exact-candidate validation
remain required after that commit exists. Actual general-purpose Linux UID/GID
execution, independent skeptical acceptance and every remote/signing/release
gate remain open. No remote action, tag, workflow run, signature, release or
successor action is authorized or claimed.
