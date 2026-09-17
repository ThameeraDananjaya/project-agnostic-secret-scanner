"""One disposable diagnostic VM's scoped AppArmor prerequisite; not a build runner."""
import hashlib
import json
import os
from pathlib import Path
import re
import selectors
import stat
import subprocess
import sys
import time

NAME = 'pscan-native-diagnostic-unshare'
TARGET = '/usr/bin/unshare'
PARSER = '/usr/sbin/apparmor_parser'
STATE = Path('/run/pscan-native-userns-prerequisite')
POLICY_ROOT = Path('/sys/kernel/security/apparmor/policy')
POLICY = (f'abi <abi/4.0>,\nprofile {NAME} {TARGET} flags=(unconfined) {{\n  userns,\n}}\n').encode('ascii')
GLOBALS = {
    '/sys/module/apparmor/parameters/enabled': 'Y',
    '/proc/sys/kernel/apparmor_restrict_unprivileged_userns': '1',
    '/proc/sys/kernel/unprivileged_userns_clone': '1',
}


def require(condition, message):
    if not condition:
        raise RuntimeError(message)


def digest(raw):
    return hashlib.sha256(raw).hexdigest()


def read(path, maximum=4096):
    with open(path, 'rb') as stream:
        raw = stream.read(maximum + 1)
    require(len(raw) <= maximum, 'Read exceeds bound')
    return raw


def trusted(path, directory=False):
    path = Path(path)
    for parent in reversed(path.parents):
        info = parent.lstat()
        require(stat.S_ISDIR(info.st_mode) and info.st_uid == 0 and not info.st_mode & 0o022,
                'Untrusted parent directory: ' + str(parent))
    info = path.lstat()
    require((stat.S_ISDIR(info.st_mode) if directory else stat.S_ISREG(info.st_mode))
            and info.st_uid == 0 and not info.st_mode & 0o022, 'Untrusted path: ' + str(path))
    return info


def attachment_reason(attachment, name):
    """Conservative: admit only provably disjoint literal prefixes, never guess regexes."""
    if attachment == name and re.fullmatch(r'[A-Za-z0-9_.-]+', name):
        return 'unattached'  # kernel emits a plain unattached name when it has no xmatch
    if not attachment or len(attachment) > 4096 or any(c in attachment for c in ('\\', '@', '<', '>', '\n')):
        return 'unknown-attachment' if not attachment or attachment == '<unknown>' else 'unsupported-pattern'
    variants = [attachment]
    for _ in range(16):
        expanded = []
        changed = False
        for value in variants:
            match = re.search(r'\{([^{}]*)\}', value)
            if match:
                changed = True
                expanded.extend(value[:match.start()] + part + value[match.end():]
                                for part in match.group(1).split(','))
            else:
                expanded.append(value)
        if len(expanded) > 64:
            return 'unsupported-pattern'
        variants = expanded
        if not changed:
            break
    for value in variants:
        if not value.startswith('/') or any(c in value for c in '{}'):
            return 'unsupported-pattern'
        prefix = re.split(r'[*?\[]', value, maxsplit=1)[0]
        if prefix == value:
            if value == TARGET:
                return 'literal-overlap'
        elif TARGET.startswith(prefix):
            return 'possible-pattern-overlap'
    return 'disjoint'


def possible_attachment(attachment, name):
    return attachment_reason(attachment, name) not in ('unattached', 'disjoint')


def inventory():
    # AppArmor's policy symlink is the documented kernel introspection interface.
    # Child namespaces are not part of this supported disposable-host profile.
    namespaces = POLICY_ROOT / 'namespaces'
    require(namespaces.is_dir() and not namespaces.is_symlink(), 'Unexpected namespace inventory path')
    require(not list(namespaces.iterdir()), 'Nested policy namespace unsupported')
    rows = []

    def visit(root, ancestry=()):
        require(len(ancestry) <= 8, 'Profile nesting bound exceeded')
        require(root.is_dir() and not root.is_symlink(), 'Unexpected profile inventory path')
        for entry in sorted(root.iterdir()):
            require(entry.is_dir() and not entry.is_symlink(), 'Unexpected profile entry')
            row = {field: read(entry / field).decode('utf-8', errors='strict').removesuffix('\n')
                   for field in ('name', 'attach', 'mode', 'sha256')}
            require(all(not any(ord(char) < 32 or ord(char) == 127 for char in value)
                        for value in row.values()), 'Unsupported profile metadata control byte')
            require(re.fullmatch(r'[0-9a-f]{64}', row['sha256']) is not None, 'Profile digest unavailable')
            require(row['mode'] in ('enforce', 'complain', 'kill', 'unconfined'), 'Unsupported profile mode')
            require(bool(row['name']), 'Empty profile identity')
            # Kernel name is base.name, not a globally unique hierarchical name.
            row['lineage'] = [*ancestry, row['name']]
            rows.append(row)
            require(len(rows) <= 1024, 'Profile count bound exceeded')
            children = entry / 'profiles'
            if children.exists():
                visit(children, tuple(row['lineage']))

    visit(POLICY_ROOT / 'profiles')
    require(len({tuple(row['lineage']) for row in rows}) == len(rows), 'Ambiguous profile identity')
    return sorted(rows, key=lambda row: tuple(row['lineage']))


def host():
    flags = {path: read(path, 32).decode('ascii').strip() for path in GLOBALS}
    if flags != GLOBALS:
        print(json.dumps({'schema': 'pscan-native-userns-global-flags-v1', 'expected': GLOBALS, 'actual': flags}, sort_keys=True))
    require(flags == GLOBALS, 'Required global protections unavailable')
    files = {}
    for path in (TARGET, PARSER, '/etc/apparmor.d/abi/4.0'):
        info = trusted(path)
        files[path] = {'sha256': digest(read(path, 16 * 1024 * 1024)),
                       'device': info.st_dev, 'inode': info.st_ino, 'mode': info.st_mode}
    require(read('/proc/self/attr/current', 256).strip() == b'unconfined', 'Privileged helper profile unsupported')
    require(read('/etc/os-release', 4096).find(b'VERSION_ID="24.04"') >= 0, 'Unsupported host version')
    return {'globals': flags, 'files': files}


def parser_command(action):
    require(action in ('compile', 'add', 'remove'), 'Unsupported parser action')
    common = [PARSER, '--config-file=/dev/null', '--base=/etc/apparmor.d',
              '--skip-cache', '--abort-on-error', '--Werror', '--jobs=0']
    return common + {'compile': ['--skip-kernel-load', '--add'],
                     'add': ['--add'], 'remove': ['--remove']}[action]


def parse(action):
    result = capture_parser(action)
    record = {'schema': 'pscan-native-userns-parser-v1', 'action': action,
              'state': result['state'], 'exit_code': result['exit_code'],
              'stdout': byte_excerpt(result['stdout']), 'stderr': byte_excerpt(result['stderr'])}
    print(json.dumps(record, ensure_ascii=True, sort_keys=True))
    if (result['state'] != 'complete' or result['exit_code'] != 0) and action in ('add', 'remove'):
        try:
            rows = [row for row in inventory() if row['name'] == NAME]
            print(conflict_record(rows, 'parser-failure-readback-not-ownership'))
        except (OSError, ValueError, RuntimeError):
            print(json.dumps({'schema': 'pscan-native-userns-readback-v1', 'state': 'unavailable-after-parser-failure'}))
    require(result['state'] == 'complete' and result['exit_code'] == 0, 'Fixed policy parser failed: ' + action)


def byte_excerpt(raw):
    part = raw[:1024]
    return {'captured_bytes': len(raw), 'excerpt_bytes': len(part), 'truncated': len(part) < len(raw),
            'encoding': 'escaped-bytes', 'sha256_of_captured': digest(raw),
            'value': ''.join(chr(b) if 32 <= b <= 126 and b != 92 else '\\x%02x' % b for b in part)}


def capture_parser(action):
    # Fixed trusted parser, fixed tiny policy; no shell, environment inheritance or retry.
    process = subprocess.Popen(parser_command(action), stdin=subprocess.PIPE, stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE, env={'PATH': '/usr/sbin:/usr/bin:/bin', 'LC_ALL': 'C'})
    buffers = {'stdout': bytearray(), 'stderr': bytearray()}
    state, code = 'complete', None
    selector = selectors.DefaultSelector()
    deadline = time.monotonic() + 10
    try:
        process.stdin.write(POLICY)
        process.stdin.close()
        for name, pipe in (('stdout', process.stdout), ('stderr', process.stderr)):
            os.set_blocking(pipe.fileno(), False)
            selector.register(pipe, selectors.EVENT_READ, name)
        while selector.get_map():
            if time.monotonic() >= deadline:
                state = 'timeout'; break
            for key, _ in selector.select(0.05):
                block = os.read(key.fileobj.fileno(), 4096)
                if not block:
                    selector.unregister(key.fileobj)
                elif len(buffers[key.data]) + len(block) > 8192:
                    state = 'capture-limit'; break
                else:
                    buffers[key.data].extend(block)
            if state != 'complete': break
        if state != 'complete': process.kill()
        try: code = process.wait(timeout=2)
        except subprocess.TimeoutExpired:
            state = 'exit-unavailable'
            process.kill()
        return {'state': state, 'exit_code': code, **{name: bytes(raw) for name, raw in buffers.items()}}
    finally:
        selector.close()
        if process.poll() is None:
            process.kill()
        for pipe in (process.stdin, process.stdout, process.stderr): pipe.close()


def save(name, value):
    raw = json.dumps(value, sort_keys=True, separators=(',', ':')).encode('utf-8')
    require(len(raw) <= 1024 * 1024, 'State record too large')
    fd = os.open(STATE / name, os.O_WRONLY | os.O_CREAT | os.O_EXCL | os.O_NOFOLLOW, 0o600)
    with os.fdopen(fd, 'wb') as stream:
        stream.write(raw)
        stream.flush()
        os.fsync(stream.fileno())


def load(name):
    trusted(STATE, directory=True)
    trusted(STATE / name)
    return json.loads(read(STATE / name, 1024 * 1024))


def admit_inventory(rows):
    require(not any(row['name'] == NAME for row in rows), 'Reserved profile name already exists')
    conflicts = [row for row in rows if possible_attachment(row['attach'], row['name'])]
    if conflicts:
        print(conflict_record(conflicts))
        raise RuntimeError('Existing or ambiguous executable attachment')


def conflict_record(conflicts, purpose='conflicts'):
    """Bounded evidence only; does not classify a conservative hit as actual overlap."""
    def excerpt(value, maximum):
        return {'value': value[:maximum], 'characters': len(value),
                'utf8_bytes': len(value.encode('utf-8')), 'truncated': len(value) > maximum,
                'sha256': digest(value.encode('utf-8'))}
    record = {'schema': 'pscan-native-userns-profile-metadata-v1', 'purpose': purpose,
              'total_records': len(conflicts), 'records': []}
    if purpose == 'conflicts':
        record.update(schema='pscan-native-userns-conflicts-v1', total_conflicts=len(conflicts),
                      meaning='existing-or-conservatively-ambiguous-attachment')
        reasons = [attachment_reason(row['attach'], row['name']) for row in conflicts]
        record['reason_counts'] = {reason: reasons.count(reason) for reason in sorted(set(reasons))}
    for row in conflicts[:16]:
        record['records'].append({'name': excerpt(row['name'], 128), 'attachment': excerpt(row['attach'], 512),
                                  'depth': len(row['lineage']) - 1,
                                  'identity_sha256': digest(json.dumps(row['lineage'], ensure_ascii=True).encode()),
                                  'mode': row['mode'], 'policy_sha256': row['sha256'],
                                  'reason': attachment_reason(row['attach'], row['name'])})
        if purpose != 'conflicts': record['records'][-1].pop('reason')
    while True:
        record['records_omitted'] = len(conflicts) - len(record['records'])
        raw = json.dumps(record, ensure_ascii=True, sort_keys=True, separators=(',', ':'))
        if len(raw.encode('ascii')) <= 16384:
            return raw
        record['records'].pop()


def own_row(rows):
    own = [row for row in rows if row['name'] == NAME]
    if own:
        print(json.dumps({'schema': 'pscan-native-userns-readback-v1',
                          'own_records': json.loads(conflict_record(own, 'own-policy-readback'))}, sort_keys=True))
    require(len(own) == 1 and own[0]['lineage'] == [NAME]
            and own[0]['attach'] == TARGET and own[0]['mode'] == 'unconfined',
            'Loaded profile readback mismatch')
    return own[0]


def others(rows):
    return [row for row in rows if row['name'] != NAME]


def same_inventory(actual, expected, phase):
    if actual == expected: return True
    old = {tuple(row['lineage']): row for row in expected}
    new = {tuple(row['lineage']): row for row in actual}
    changed_before = [old[key] for key in old if old[key] != new.get(key)]
    changed_after = [new[key] for key in new if new[key] != old.get(key)]
    for side, rows in (('before', changed_before), ('after', changed_after)):
        print(json.dumps({'schema': 'pscan-native-userns-inventory-drift-v1', 'phase': phase,
                          'side': side, 'changes': json.loads(conflict_record(rows, 'inventory-drift'))}, sort_keys=True))
    return False


def same_host(actual, expected, phase):
    if actual == expected: return True
    # These objects contain only the fixed globals and executable/ABI metadata.
    print(json.dumps({'schema': 'pscan-native-userns-host-drift-v1', 'phase': phase,
                      'before': expected, 'after': actual}, sort_keys=True))
    return False


def emit(phase, state, rows):
    print(json.dumps({'schema': 'pscan-native-userns-prerequisite-v1', 'phase': phase,
                      'source_sha256': state['source_sha256'], 'policy_sha256': digest(POLICY),
                      'profile_name': NAME, 'attachment': TARGET, 'profile_count': len(rows),
                      'inventory_sha256': digest(json.dumps(rows, sort_keys=True).encode()),
                      'global_restrictions_preserved': True}, sort_keys=True))


def install(source_hash):
    require(not STATE.exists() and not STATE.is_symlink(), 'Prerequisite state already exists')
    baseline_host = host()
    baseline = inventory()
    admit_inventory(baseline)
    parse('compile')
    require(same_host(host(), baseline_host, 'preflight') and same_inventory(inventory(), baseline, 'preflight'),
            'Host changed during preflight')
    trusted(STATE.parent, directory=True)
    STATE.mkdir(mode=0o700)
    state = {'source_sha256': source_hash, 'policy_sha256': digest(POLICY),
             'host': baseline_host, 'profiles': baseline}
    save('baseline.json', state)
    # A failure after this marker is ambiguous until kernel readback resolves it.
    save('add-started.json', {'started': True})
    parse('add')
    current = inventory()
    own = own_row(current)
    save('loaded.json', own)
    require(same_inventory(others(current), baseline, 'after-add') and same_host(host(), baseline_host, 'after-add'),
            'Host changed during policy addition')
    emit('installed', state, current)


def cleanup(source_hash):
    if not STATE.exists() and not STATE.is_symlink():
        print(json.dumps({'schema': 'pscan-native-userns-prerequisite-v1', 'phase': 'not-installed'}))
        return
    state = load('baseline.json')
    require(state['source_sha256'] == source_hash and state['policy_sha256'] == digest(POLICY),
            'Cleanup source does not match setup')
    require(same_host(host(), state['host'], 'before-cleanup'), 'Host changed before cleanup')
    current = inventory()
    require(same_inventory(others(current), state['profiles'], 'before-cleanup'), 'Unrelated policy changed before cleanup')
    if any(row['name'] == NAME for row in current):
        # Only a successful captured add owns this exact kernel policy hash.
        require(own_row(current) == load('loaded.json'), 'Owned policy changed or load ownership uncertain')
        parse('remove')
    require(same_inventory(inventory(), state['profiles'], 'after-remove')
            and same_host(host(), state['host'], 'after-remove'), 'Cleanup readback mismatch')
    emit('removed', state, state['profiles'])
    # Delete exact own files only, and only after verified removal/preservation.
    allowed = {'baseline.json', 'add-started.json', 'loaded.json'}
    require({p.name for p in STATE.iterdir()} <= allowed, 'Unexpected state file')
    for name in sorted(allowed):
        path = STATE / name
        if path.exists():
            trusted(path)
            path.unlink()
    STATE.rmdir()


def main():
    require(sys.platform == 'linux' and os.geteuid() == 0, 'Requires reviewed disposable Linux setup')
    require(len(sys.argv) == 2 and sys.argv[1] in ('install', 'cleanup'), 'Fixed action required')
    source_hash = globals().get('__source_sha256__')
    require(isinstance(source_hash, str) and re.fullmatch(r'[0-9a-f]{64}', source_hash),
            'Committed-byte bootstrap admission required')
    {'install': install, 'cleanup': cleanup}[sys.argv[1]](source_hash)


if __name__ == '__main__':
    main()
