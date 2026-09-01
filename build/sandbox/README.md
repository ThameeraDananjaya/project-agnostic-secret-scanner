# Authoritative PSCAN-04 sandbox

`validate.ps1` runs validation in the exact digest-pinned Go 1.27.0 image with
no network, a read-only root filesystem, read-only source/module/detector
mounts, all capabilities dropped, no-new-privileges, default seccomp and
bounded CPU, memory, PID and temporary filesystems. The environment is rebuilt
from an explicit non-secret allowlist.

`-Profile pr` binds the cgroup to 4 GiB with swap disabled; `-Profile release`
binds it to 8 GiB with swap disabled. Both use the contract's single scanner
subprocess declaration while retaining a finite 256-PID harness allowance for
the Go compiler and test process tree.

The module cache is populated and checksum-verified in a separate acquisition
step, then supplied read-only. Candidate artifacts are never mounted into this
validation container and are never built or executed.
