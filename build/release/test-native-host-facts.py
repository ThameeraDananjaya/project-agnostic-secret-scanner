"""Inert parser/record tests only; subprocesses and real host probes are blocked."""
import contextlib
import importlib.util
import io
import json
from pathlib import Path
import tempfile
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('hostfacts', Path(__file__).with_name('native-host-facts.py'))
facts = importlib.util.module_from_spec(spec)
spec.loader.exec_module(facts)


def present(value): return {'state': 'present', 'value': value}


class ParserTests(unittest.TestCase):
    def test_status_only_allowlisted_fields(self):
        result = facts.parse_status(present('Name:\tDO-NOT-LOG\nUid:\t1001 1002 1003 1004\nGid:\t20 21 22 23\nCapEff:\t0000000000000000\nNoNewPrivs:\t1\nSeccomp:\t2\nSecret:\tDO-NOT-LOG'))
        self.assertEqual(set(result), set(facts.STATUS_FIELDS))
        self.assertEqual(result['Uid']['effective'], 1002)
        self.assertEqual(result['Gid']['effective'], 21)
        self.assertEqual(result['CapEff']['value'], '0000000000000000')
        self.assertEqual(result['CapAmb']['state'], 'missing')
        self.assertNotIn('DO-NOT-LOG', json.dumps(result))

    def test_status_malformed_and_duplicate(self):
        for raw in ('Uid: 1 2 x 4', 'Uid: 1 2 3', 'Uid: -1 2 3 4'):
            self.assertEqual(facts.parse_status(present(raw))['Uid']['state'], 'invalid-format')
        self.assertEqual(facts.parse_status(present('Uid: 1 2 3 4\nUid: 1 2 3 4'))['Uid']['state'], 'duplicate-fields')
        self.assertEqual(facts.parse_status(facts.unavailable('unreadable'))['Uid']['state'], 'unreadable')

    def test_maps_empty_valid_malformed_and_excess(self):
        self.assertEqual(facts.parse_map(present('0 0 4294967295'))['value'], [[0, 0, 4294967295]])
        self.assertEqual(facts.parse_map(present(''))['value'], [])
        self.assertEqual(facts.parse_map(present('1 2 bad'))['state'], 'invalid-map')
        self.assertEqual(facts.parse_map(present('\n'.join(['1 1 1'] * 17)))['state'], 'too-many-map-lines')
        self.assertEqual(facts.parse_map(facts.unavailable('missing'))['state'], 'missing')

    def test_bounded_reads(self):
        with tempfile.TemporaryDirectory() as directory:
            path = Path(directory) / 'inert'
            self.assertEqual(facts.bounded_read(path)['state'], 'missing')
            path.write_bytes(b'x' * 9)
            self.assertEqual(facts.bounded_read(path, 8)['state'], 'oversize')
            path.write_bytes(b'\xff')
            self.assertEqual(facts.bounded_read(path)['state'], 'invalid-encoding')
            path.write_bytes(b' 1\n')
            self.assertEqual(facts.bounded_read(path), present('1'))

    def test_validated_never_echoes_invalid_payload(self):
        self.assertEqual(facts.validated(present('secret\nline'), r'[01]'), facts.unavailable('invalid-format'))
        self.assertEqual(facts.validated(facts.unavailable('unreadable'), r'[01]')['state'], 'unreadable')

    def journal(self, rows):
        return {'state': 'present', 'exit_code': 0, 'stdout': ('\n'.join(json.dumps(r) for r in rows)).encode(), 'stderr_bytes': 0}

    def test_audit_allowlist_and_correlation_limits(self):
        result = facts.parse_denials(self.journal([
            {'__REALTIME_TIMESTAMP': '101000000', 'MESSAGE': 'audit: apparmor="DENIED" operation="capable" profile="unprivileged_userns" pid=123 comm="unshare" capability=7 capname="setuid" secret="DO-NOT-LOG"'},
            {'__REALTIME_TIMESTAMP': '101000000', 'MESSAGE': 'apparmor="DENIED" comm="unrelated" secret="DO-NOT-LOG"'},
            {'__REALTIME_TIMESTAMP': '90000000', 'MESSAGE': 'apparmor="DENIED" comm="unshare"'}]), 100, 102)
        self.assertEqual(len(result['records']), 1)
        self.assertEqual(result['records'][0]['capname'], 'setuid')
        self.assertEqual(result['fixture_pid_binding'], 'unavailable')
        self.assertNotIn('DO-NOT-LOG', json.dumps(result))
        self.assertNotIn('MESSAGE', json.dumps(result))

    def test_audit_missing_malformed_oversize_and_duplicate(self):
        for state in ('unavailable', 'unreadable', 'timeout', 'output-limit', 'stderr-present'):
            result = facts.parse_denials({'state': state, 'exit_code': None, 'stdout': b'', 'stderr_bytes': None}, 100, 102)
            self.assertEqual(result['state'], state)
        for raw in (b'not json', b'{}', b'\xff'):
            self.assertEqual(facts.parse_denials({'state': 'present', 'exit_code': 0, 'stdout': raw}, 100, 102)['state'], 'unparseable')
        duplicate = self.journal([{'__REALTIME_TIMESTAMP': '101000000', 'MESSAGE': 'apparmor="DENIED" comm="unshare" comm="unshare"'}])
        self.assertEqual(facts.parse_denials(duplicate, 100, 102)['state'], 'unparseable')
        self.assertEqual(facts.parse_denials(self.journal([]), 100, 102)['state'], 'no-correlated-records')

    def test_full_record_before_and_after_with_only_inert_providers(self):
        def fake_read(path, maximum=4096):
            if path.endswith('/comm'): return present('pwsh')
            if path.endswith('/status'): return present('Uid: 1001 1001 1001 1001\nGid: 1001 1001 1001 1001\nCapEff: 0000000000000000\nNoNewPrivs: 0\nSeccomp: 2')
            if path.endswith('_map'): return present('0 0 4294967295')
            if path.endswith('osrelease'): return present('6.8.0-test')
            if path.endswith('/enabled'): return present('Y')
            if path.endswith('/current'): return present('unconfined')
            return facts.unavailable('unreadable')
        calls = []
        def fake_query(args, **kwargs):
            calls.append(args)
            return {'state': 'present', 'exit_code': 0, 'stdout': b'unshare from util-linux 2.39.3\n' if args[0].endswith('unshare') else b'', 'stderr_bytes': 0}
        for phase in ('before', 'after'):
            output = io.StringIO()
            with patch.object(facts.sys, 'platform', 'linux'), patch.object(facts.sys, 'argv', ['probe', '--phase', phase, '--since', '299']), \
                 patch.object(facts.time, 'time', return_value=300), patch.object(facts.os, 'getppid', return_value=123), \
                 patch.object(facts, 'bounded_read', side_effect=fake_read), patch.object(facts, 'query', side_effect=fake_query), \
                 patch.dict(facts.os.environ, {'ImageOS': 'ubuntu24', 'ImageVersion': '20260907.300.1', 'SECRET': 'DO-NOT-LOG'}, clear=True), \
                 patch.object(facts.subprocess, 'Popen', side_effect=AssertionError('Native execution forbidden in inert tests')), \
                 contextlib.redirect_stdout(output):
                facts.main()
            record = json.loads(output.getvalue())
            self.assertLessEqual(len(output.getvalue().encode()), facts.MAX_RECORD + 1)
            self.assertNotIn('DO-NOT-LOG', output.getvalue())
            self.assertEqual(record['phase'], phase)
            self.assertEqual(record['apparmor']['userns_restriction']['state'], 'unreadable')
        self.assertEqual(sum(args[0].endswith('journalctl') for args in calls), 1)
        self.assertTrue(all(args[0] in ('/usr/bin/unshare', '/usr/bin/journalctl') for args in calls))


if __name__ == '__main__': unittest.main(verbosity=2)
