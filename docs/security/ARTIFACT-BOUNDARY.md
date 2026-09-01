# PSCAN-04 Artifact Security Boundary

## Denied input behavior

Traversal, absolute/drive/ADS paths, links, hard links, reparse points, devices,
FIFOs, duplicate and case-colliding paths, ZIP encryption, unsupported ZIP
methods, concatenated gzip streams, malformed terminal framing, unknown magic,
unsupported compression and unbound OCI/Docker material are non-pass before
detector invocation. Archive/member and workspace names are never copied to a
public surface.

Candidate executability and permission bits confer no authority. Candidate
files are opened only for bounded reads; their scripts, binaries, build files,
installers and package metadata are never invoked. OCI descriptor URLs and
inline data are rejected, and no provider or registry access exists.

## Disclosure boundary

Engine stdout/stderr are bounded private buffers. Decode panics are recovered
without formatting their payload. Artifact errors use one fixed message; the
typed diagnostic serializer admits only schema version, state, reason code and
permitted action. It rejects unknown fields, multiple documents and injected
frames. Acceptance sweeps raw, partial, base64, URL-base64, hexadecimal and
SHA-256 canary forms across success, finding, error, timeout, cancellation and
panic paths.

## Isolation boundary

The authoritative container has only loopback, an immutable root filesystem,
read-only source/module/detector mounts, no host environment, credentials,
provider configuration or Docker socket, no Linux capabilities, no-new-
privileges, default seccomp, and finite CPU, memory, PID and tmpfs resources.
DNS, public TCP and cloud-metadata endpoint attempts must all fail. Any missing
or conflicting isolation proof prevents acceptance.

The authoritative outcome serializer remains `internal/outcome`; the bounded
diagnostic envelope is not a second outcome schema. Because the public runner
cannot be changed in PSCAN-04, this task proves the internal stage and its
content-free diagnostic boundary, not end-to-end CLI exposure or consumer
integration.
