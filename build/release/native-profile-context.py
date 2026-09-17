"""Bounded same-user context comparison for the fixed diagnostic PowerShell host.

The runner PowerShell is an existing runtime assumption, not an immutable or
publisher-attested executable. Observed identity equality cannot exclude a
transient in-place write, and does not attest its dynamic runtime dependencies.
"""
import hashlib
import json
import os
from pathlib import Path
import re
import stat
import sys

TARGET = '/opt/microsoft/powershell/7/pwsh'
LABEL = 'pscan-native-diagnostic (unconfined)'
FIELDS = ('Uid', 'Gid', 'Groups', 'CapInh', 'CapPrm', 'CapEff', 'CapBnd', 'CapAmb',
          'NoNewPrivs', 'Seccomp', 'Seccomp_filters')
NAMESPACES = ('cgroup', 'ipc', 'mnt', 'net', 'pid', 'pid_for_children',
              'time', 'time_for_children', 'user', 'uts')
SOURCES = ('native-profile-context.py', 'native-profile-driver.ps1', 'native-host-facts.py',
           'test-docker-execution.ps1', 'docker-execution.ps1', 'execution-profile.ps1', 'native-fixture-diagnostics.ps1')


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


def read(path, maximum):
    with open(path, 'rb') as stream:
        raw = stream.read(maximum + 1)
    require(len(raw) <= maximum, 'Context read exceeds bound')
    return raw


def sha(raw):
    return hashlib.sha256(raw).hexdigest()


def metadata(info):
    return {name: getattr(info, 'st_' + name) for name in
            ('dev', 'ino', 'mode', 'uid', 'gid', 'size', 'mtime_ns', 'ctime_ns')}


def executable(pid):
    proc = '/proc/' + str(pid) + '/exe'
    require(os.readlink(proc) == TARGET, 'Caller is not the fixed runner PowerShell')
    fd = os.open(TARGET, os.O_RDONLY | os.O_NOFOLLOW | os.O_NONBLOCK)
    with os.fdopen(fd, 'rb') as stream:
        before = os.fstat(stream.fileno())
        require(stat.S_ISREG(before.st_mode) and before.st_uid == 0,
                'Unsupported runner PowerShell identity')
        require(not before.st_mode & 0o6000, 'Set-ID PowerShell unsupported')
        require('security.capability' not in os.listxattr(stream.fileno()), 'PowerShell file capabilities unsupported')
        raw = stream.read(16 * 1024 * 1024 + 1)
        require(len(raw) <= 16 * 1024 * 1024, 'PowerShell executable exceeds bound')
        after = os.fstat(stream.fileno())
        require(metadata(before) == metadata(after) == metadata(os.stat(TARGET, follow_symlinks=False)),
                'PowerShell changed during observation')
        running = os.stat(proc)
        require((running.st_dev, running.st_ino) == (before.st_dev, before.st_ino),
                'Caller executable differs from selected path')
    return {**metadata(before), 'sha256': sha(raw), 'path': TARGET,
            'meaning': 'observed-runner-runtime-not-immutable-or-publisher-attested'}


def parse_status(raw):
    entries = {}
    for line in raw.decode('ascii', errors='strict').splitlines():
        key, sep, value = line.partition(':')
        if sep and key in FIELDS:
            require(key not in entries, 'Duplicate context field')
            entries[key] = value.split()
    require(set(entries) == set(FIELDS), 'Missing context field')
    for key, values in entries.items():
        if key.startswith('Cap'):
            require(len(values) == 1 and re.fullmatch(r'[0-9a-f]{16}', values[0]), 'Invalid capability mask')
        else:
            count = 4 if key in ('Uid', 'Gid') else 1
            require((key == 'Groups' and len(values) <= 256) or len(values) == count,
                    'Invalid numeric context tuple')
            require(all(re.fullmatch(r'[0-9]{1,10}', v) for v in values), 'Invalid numeric context value')
    require(all(int(v) != 0 for v in entries['Uid']), 'Privileged caller unsupported')
    require(all(entries[k] == ['0000000000000000'] for k in ('CapInh', 'CapPrm', 'CapEff', 'CapAmb')),
            'Ambient or active privilege unsupported')
    return entries


def context(pid):
    root = '/proc/' + str(pid)
    require(read(root + '/comm', 128).strip() == b'pwsh', 'Expected calling PowerShell')
    fields = parse_status(read(root + '/status', 16384))
    namespaces = {}
    for name in NAMESPACES:
        value = os.readlink(root + '/ns/' + name)
        require(re.fullmatch(r'[a-z_]+:\[[0-9]{1,20}\]', value), 'Invalid namespace identity')
        namespaces[name] = value
    label = read(root + '/attr/current', 256).decode('ascii', errors='strict').removesuffix('\n')
    sources = {name: sha(read(Path(__file__).with_name(name), 1024 * 1024)) for name in SOURCES}
    return {'schema': 'pscan-native-profile-context-v1', 'status': fields,
            'namespaces': namespaces, 'label': label, 'executable': executable(pid), 'sources': sources}


def validate_baseline(value):
    require(isinstance(value, dict) and set(value) ==
            {'schema', 'status', 'namespaces', 'label', 'executable', 'sources'}, 'Invalid baseline keys')
    require(value['schema'] == 'pscan-native-profile-context-v1' and value['label'] == 'unconfined',
            'Baseline must be ordinary unconfined context')


def compare(baseline, actual):
    validate_baseline(baseline)
    require(actual['label'] == LABEL, 'Expected sole named profile label')
    for key in ('schema', 'status', 'namespaces', 'executable', 'sources'):
        require(actual[key] == baseline[key], 'Profile transition context mismatch: ' + key)


def main():
    require(sys.platform == 'linux' and os.geteuid() != 0, 'Ordinary Linux caller required')
    require(len(sys.argv) in (3, 5), 'Fixed context action required')
    action, path = sys.argv[1:3]
    require(Path(path).is_absolute(), 'Absolute baseline path required')
    if action == 'snapshot':
        require(len(sys.argv) == 3, 'Invalid snapshot arguments')
        value = context(os.getppid())
        validate_baseline(value)
        raw = json.dumps(value, sort_keys=True, separators=(',', ':')).encode()
        require(len(raw) <= 16384, 'Baseline exceeds bound')
        fd = os.open(path, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, 0o600)
        with os.fdopen(fd, 'wb') as stream:
            stream.write(raw)
        print(json.dumps({'schema': 'pscan-native-profile-check-v1', 'phase': 'baseline',
                          'sha256': sha(raw), 'context': value}, sort_keys=True))
    else:
        require(action == 'check' and len(sys.argv) == 5, 'Invalid check action')
        expected, phase = sys.argv[3:]
        require(re.fullmatch(r'[0-9a-f]{64}', expected) and phase in ('before', 'after'), 'Invalid check binding')
        raw = read(path, 16384)
        require(sha(raw) == expected, 'Baseline bytes changed')
        value = context(os.getppid())
        compare(json.loads(raw), value)
        print(json.dumps({'schema': 'pscan-native-profile-check-v1', 'phase': phase,
                          'baseline_sha256': expected, 'context': value}, sort_keys=True))


if __name__ == '__main__':
    main()
