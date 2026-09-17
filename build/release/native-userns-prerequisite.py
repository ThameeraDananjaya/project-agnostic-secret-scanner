"""One disposable diagnostic VM's scoped AppArmor prerequisite; not a build runner."""
import hashlib
import json
import os
from pathlib import Path
import re
import stat
import subprocess
import sys

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
                'Untrusted parent directory')
    info = path.lstat()
    require((stat.S_ISDIR(info.st_mode) if directory else stat.S_ISREG(info.st_mode))
            and info.st_uid == 0 and not info.st_mode & 0o022, 'Untrusted path')
    return info


def possible_attachment(attachment, name):
    """Conservative: admit only provably disjoint literal prefixes, never guess regexes."""
    if attachment == name and re.fullmatch(r'[A-Za-z0-9_.-]+', name):
        return False  # kernel emits a plain unattached name when it has no xmatch
    if not attachment or len(attachment) > 4096 or any(c in attachment for c in ('\\', '@', '<', '>', '\n')):
        return True
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
            return True
        variants = expanded
        if not changed:
            break
    for value in variants:
        if not value.startswith('/') or any(c in value for c in '{}'):
            return True
        prefix = re.split(r'[*?\[]', value, maxsplit=1)[0]
        if prefix == value:
            if value == TARGET:
                return True
        elif TARGET.startswith(prefix):
            return True
    return False


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
    # Fixed trusted utility and tiny fixed source. No inherited stdin or shell.
    result = subprocess.run(parser_command(action), input=POLICY, stdout=subprocess.DEVNULL,
                            stderr=subprocess.DEVNULL, timeout=10, check=False,
                            env={'PATH': '/usr/sbin:/usr/bin:/bin', 'LC_ALL': 'C'})
    require(result.returncode == 0, 'Fixed policy parser failed: ' + action)


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
    require(not any(possible_attachment(row['attach'], row['name']) for row in rows),
            'Existing or ambiguous executable attachment')


def own_row(rows):
    own = [row for row in rows if row['name'] == NAME]
    require(len(own) == 1 and own[0]['lineage'] == [NAME]
            and own[0]['attach'] == TARGET and own[0]['mode'] == 'unconfined',
            'Loaded profile readback mismatch')
    return own[0]


def others(rows):
    return [row for row in rows if row['name'] != NAME]


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
    require(host() == baseline_host and inventory() == baseline, 'Host changed during preflight')
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
    require(others(current) == baseline and host() == baseline_host, 'Host changed during policy addition')
    emit('installed', state, current)


def cleanup(source_hash):
    if not STATE.exists() and not STATE.is_symlink():
        print(json.dumps({'schema': 'pscan-native-userns-prerequisite-v1', 'phase': 'not-installed'}))
        return
    state = load('baseline.json')
    require(state['source_sha256'] == source_hash and state['policy_sha256'] == digest(POLICY),
            'Cleanup source does not match setup')
    require(host() == state['host'], 'Host changed before cleanup')
    current = inventory()
    require(others(current) == state['profiles'], 'Unrelated policy changed before cleanup')
    if any(row['name'] == NAME for row in current):
        # Only a successful captured add owns this exact kernel policy hash.
        require(own_row(current) == load('loaded.json'), 'Owned policy changed or load ownership uncertain')
        parse('remove')
    require(inventory() == state['profiles'] and host() == state['host'], 'Cleanup readback mismatch')
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
