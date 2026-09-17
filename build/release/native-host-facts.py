"""Read-only, allowlisted facts around the fixed Linux native fixture."""
import argparse
import json
import os
from pathlib import Path
import re
import selectors
import stat
import subprocess
import sys
import time

MAX_RECORD = 16384
STATUS_FIELDS = ('Uid', 'Gid', 'CapInh', 'CapPrm', 'CapEff', 'CapBnd', 'CapAmb',
                 'NoNewPrivs', 'Seccomp', 'Seccomp_filters')
AUDIT_FIELDS = ('apparmor', 'operation', 'profile', 'pid', 'comm', 'capability', 'capname',
                'error', 'requested_mask', 'denied_mask')


def unavailable(state):
    return {'state': state, 'value': None}


def file_metadata(path):
    """Availability facts only, not executable trust or selection admission."""
    try:
        info = os.lstat(path)
        return {'state': 'present', 'regular_file': stat.S_ISREG(info.st_mode),
                'symlink': stat.S_ISLNK(info.st_mode), 'uid': info.st_uid, 'gid': info.st_gid,
                'mode': stat.S_IMODE(info.st_mode), 'bytes': info.st_size,
                'device': info.st_dev, 'inode': info.st_ino}
    except FileNotFoundError: return unavailable('missing')
    except PermissionError: return unavailable('unreadable')
    except OSError: return unavailable('stat-error')


def bounded_read(path, maximum=4096):
    try:
        with open(path, 'rb') as f:
            raw = f.read(maximum + 1)
        if len(raw) > maximum:
            return unavailable('oversize')
        return {'state': 'present', 'value': raw.decode('utf-8', errors='strict').strip()}
    except FileNotFoundError:
        return unavailable('missing')
    except PermissionError:
        return unavailable('unreadable')
    except UnicodeError:
        return unavailable('invalid-encoding')
    except OSError:
        return unavailable('read-error')


def validated(field, pattern):
    if field['state'] != 'present':
        return field
    return field if re.fullmatch(pattern, field['value']) else unavailable('invalid-format')


def parse_status(field):
    if field['state'] != 'present':
        return {name: unavailable(field['state']) for name in STATUS_FIELDS}
    values = {}
    for line in field['value'].splitlines():
        key, sep, value = line.partition(':')
        if sep and key in STATUS_FIELDS:
            if key in values:
                return {name: unavailable('duplicate-fields') for name in STATUS_FIELDS}
            values[key] = value.strip()
    result = {}
    for name in STATUS_FIELDS:
        if name not in values:
            result[name] = unavailable('missing')
            continue
        value = values[name]
        if name in ('Uid', 'Gid'):
            parts = value.split()
            result[name] = {'state': 'present', 'effective': int(parts[1])} if len(parts) == 4 and all(re.fullmatch(r'[0-9]{1,10}', p) for p in parts) else unavailable('invalid-format')
        elif name.startswith('Cap'):
            result[name] = validated({'state': 'present', 'value': value}, r'[0-9a-fA-F]{1,16}')
        else:
            result[name] = validated({'state': 'present', 'value': value}, r'[0-9]{1,10}')
    return result


def parse_map(field):
    if field['state'] != 'present':
        return field
    lines = field['value'].splitlines()
    if len(lines) > 16:
        return unavailable('too-many-map-lines')
    rows = []
    for line in lines:
        parts = line.split()
        if len(parts) != 3 or not all(re.fullmatch(r'[0-9]{1,10}', p) for p in parts):
            return unavailable('invalid-map')
        rows.append([int(p) for p in parts])
    return {'state': 'present', 'value': rows}


def query(args, limit=8192, seconds=5):
    """Only main supplies fixed trusted version/journal argv; bounded raw capture."""
    try:
        p = subprocess.Popen(args, stdin=subprocess.DEVNULL, stdout=subprocess.PIPE, stderr=subprocess.PIPE,
                             env={'PATH': '/usr/bin:/bin', 'LC_ALL': 'C'}, start_new_session=False)
    except OSError:
        return {'state': 'unavailable', 'exit_code': None, 'stdout': b'', 'stderr_bytes': None}
    buffers = {'stdout': bytearray(), 'stderr': bytearray()}
    state = 'present'
    selector = selectors.DefaultSelector()
    try:
        for name, pipe in (('stdout', p.stdout), ('stderr', p.stderr)):
            os.set_blocking(pipe.fileno(), False)
            selector.register(pipe, selectors.EVENT_READ, name)
        deadline = time.monotonic() + seconds
        while selector.get_map():
            if time.monotonic() >= deadline:
                state = 'timeout'; break
            for key, _ in selector.select(min(0.05, max(0, deadline - time.monotonic()))):
                block = os.read(key.fileobj.fileno(), 4096)
                if not block:
                    selector.unregister(key.fileobj)
                    continue
                if len(buffers[key.data]) + len(block) > limit:
                    state = 'output-limit'; break
                buffers[key.data].extend(block)
            if state != 'present': break
        if state != 'present': p.kill()
        try: code = p.wait(timeout=2)
        except subprocess.TimeoutExpired:
            p.kill(); code = None; state = 'exit-unavailable'
        if state == 'present' and code != 0: state = 'unreadable'
        if state == 'present' and buffers['stderr']: state = 'stderr-present'
        return {'state': state, 'exit_code': code, 'stdout': bytes(buffers['stdout']),
                'stderr_bytes': len(buffers['stderr'])}
    finally:
        selector.close()
        p.stdout.close(); p.stderr.close()


def parse_denials(result, since, until):
    metadata = {'state': result['state'], 'query_exit_code': result['exit_code'],
                'query_stderr_bytes': result.get('stderr_bytes'),
                'correlation': 'fixture-time-window-and-unshare-comm',
                'fixture_pid_binding': 'unavailable', 'records': []}
    if result['state'] != 'present': return metadata
    try:
        lines = result['stdout'].decode('utf-8', errors='strict').splitlines()
        if len(lines) > 16: raise ValueError('too many')
        for line in lines:
            row = json.loads(line)
            stamp = int(row['__REALTIME_TIMESTAMP'])
            if not since * 1_000_000 <= stamp <= until * 1_000_000: continue
            message = row.get('MESSAGE')
            if not isinstance(message, str) or len(message) > 4096: continue
            pairs = re.findall(r'\b([a-z_]+)=(?:"([^"\r\n]*)"|([^\s]+))', message)
            fields = {}
            for name, quoted, bare in pairs:
                if name in AUDIT_FIELDS:
                    if name in fields: raise ValueError('duplicate')
                    fields[name] = quoted or bare
            if fields.get('apparmor') != 'DENIED' or fields.get('comm') != 'unshare': continue
            safe = {'timestamp_us': stamp}
            for name in AUDIT_FIELDS:
                value = fields.get(name)
                safe[name] = value if isinstance(value, str) and re.fullmatch(r'[A-Za-z0-9_./(): -]{1,128}', value) else None
            metadata['records'].append(safe)
        if not metadata['records']: metadata['state'] = 'no-correlated-records'
    except (ValueError, KeyError, TypeError, UnicodeError):
        metadata['state'] = 'unparseable'; metadata['records'] = []
    return metadata


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--phase', choices=('before', 'after'), required=True)
    parser.add_argument('--since', type=int, required=True)
    args = parser.parse_args()
    now = int(time.time())
    if sys.platform != 'linux' or not 0 <= now - args.since <= 300:
        raise SystemExit('Unsupported diagnostic platform or time window')
    parent = os.getppid()
    root = '/proc/' + str(parent)
    caller = validated(bounded_read(root + '/comm', 128), r'pwsh')
    version = query(['/usr/bin/unshare', '--version'], limit=512)
    try:
        util = validated({'state': version['state'], 'value': version['stdout'].decode('ascii').strip()},
                         r'unshare from util-linux [0-9][A-Za-z0-9.+~-]{0,63}')
    except UnicodeError: util = unavailable('invalid-encoding')
    facts = {'schema': 'pscan-native-host-facts-v1', 'phase': args.phase, 'observed_epoch': now,
             'fixture_window_start_epoch': args.since, 'subject': 'calling-pwsh-process', 'caller': caller,
             'image': {name: validated({'state': 'present', 'value': os.environ[name]}, r'[A-Za-z0-9_.-]{1,64}')
                       if name in os.environ else unavailable('missing') for name in ('ImageOS', 'ImageVersion')},
             'kernel': validated(bounded_read('/proc/sys/kernel/osrelease', 256), r'[A-Za-z0-9_.+-]{1,128}'),
             'util_linux': util,
             'status': parse_status(bounded_read(root + '/status', 8192)),
             'uid_map': parse_map(bounded_read(root + '/uid_map')),
             'gid_map': parse_map(bounded_read(root + '/gid_map')),
             'apparmor': {
                 'enabled': validated(bounded_read('/sys/module/apparmor/parameters/enabled', 16), r'[YN]'),
                 'userns_restriction': validated(bounded_read('/proc/sys/kernel/apparmor_restrict_unprivileged_userns', 16), r'[01]'),
                 'current_profile': validated(bounded_read(root + '/attr/current', 512), r'[A-Za-z0-9_./(): -]{1,256}')},
             'named_profile_compatibility': {
                 'unprivileged_unconfined_restriction': validated(bounded_read('/proc/sys/kernel/apparmor_restrict_unprivileged_unconfined', 16), r'[01]'),
                 'selector': file_metadata('/usr/bin/aa-exec'),
                 'canonical_powershell_candidate': file_metadata('/opt/microsoft/powershell/7/pwsh'),
                 'meaning': 'availability-and-context-only-not-executable-admission'},
             'unprivileged_userns_clone': validated(bounded_read('/proc/sys/kernel/unprivileged_userns_clone', 32), r'[01]'),
             'max_user_namespaces': validated(bounded_read('/proc/sys/user/max_user_namespaces', 32), r'[0-9]{1,10}')}
    if args.phase == 'after':
        until = int(time.time()) + 1
        result = query(['/usr/bin/journalctl', '--dmesg', '--since=@' + str(args.since), '--until=@' + str(until),
                        '--no-pager', '--output=json', '--lines=16', '--grep=apparmor="DENIED".*comm="unshare"'], limit=16384)
        facts['launcher_denials'] = parse_denials(result, args.since, until)
    else:
        facts['launcher_denials'] = {'state': 'not-queried-before-fixture', 'records': []}
    raw = json.dumps(facts, sort_keys=True, separators=(',', ':'))
    if len(raw.encode()) > MAX_RECORD: raise SystemExit('Host facts record exceeds its fixed bound')
    print(raw)


if __name__ == '__main__': main()
