"""Inert prerequisite admission/lifecycle tests; never run a host utility."""
import importlib.util
import io
import json
from pathlib import Path
import tempfile
from types import SimpleNamespace
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('prerequisite', Path(__file__).with_name('native-userns-prerequisite.py'))
prereq = importlib.util.module_from_spec(spec)
spec.loader.exec_module(prereq)


def row(name='existing', attachment='/usr/bin/other', mode='enforce', digest='a' * 64):
    return {'name': name, 'attach': attachment, 'mode': mode, 'sha256': digest, 'lineage': [name]}


class Tests(unittest.TestCase):
    def setUp(self):
        self.block = patch.object(prereq.subprocess, 'run', side_effect=AssertionError('Host execution forbidden'))
        self.block.start()
        self.addCleanup(self.block.stop)
        self.popen_block = patch.object(prereq.subprocess, 'Popen', side_effect=AssertionError('Host execution forbidden'))
        self.popen_block.start()
        self.addCleanup(self.popen_block.stop)
        self.quiet = patch('builtins.print')
        self.quiet.start()
        self.addCleanup(self.quiet.stop)

    def test_conflicts_and_unknowns_reject(self):
        for attachment in ('/usr/bin/unshare', '/usr/bin/*', '/**', '<unknown>', '', '/usr/bin/unsh?re',
                           '/usr/bin/[ux]nshare', '/{usr/,}bin/unshare', '@{bin}/unshare', '/usr/bin/\\unshare',
                           '/**/not-unshare', '{broken', '/usr/bin/unshare\n'):
            with self.subTest(attachment=attachment):
                self.assertTrue(prereq.possible_attachment(attachment, 'existing'))
                with patch('builtins.print'), self.assertRaises(RuntimeError):
                    prereq.admit_inventory([row(attachment=attachment)])

    def test_conflict_diagnostics_are_bounded_explicit_and_do_not_admit(self):
        rows = [row(name='p' + str(n), attachment='/usr/bin/*' + '\u2603' * 2000) for n in range(30)]
        raw = prereq.conflict_record(rows)
        self.assertLessEqual(len(raw.encode()), 16384)
        record = json.loads(raw)
        self.assertEqual(record['total_conflicts'], 30)
        self.assertGreater(record['records_omitted'], 0)
        self.assertEqual(len(record['all_records_sha256']), 64)
        for item in record['records']:
            self.assertTrue(item['attachment']['truncated'])
            self.assertEqual(item['attachment']['utf8_bytes'], 6010)
            self.assertEqual(item['depth'], 0)
        with patch('builtins.print') as output, self.assertRaisesRegex(RuntimeError, 'Existing or ambiguous'):
            prereq.admit_inventory([row(attachment='/usr/bin/*')])
        self.assertEqual(json.loads(output.call_args.args[0])['total_conflicts'], 1)

    def test_attachment_reason_distinguishes_evidence_without_allowing_unknowns(self):
        for value, reason in (('/usr/bin/unshare', 'literal-overlap'), ('/usr/bin/*', 'possible-pattern-overlap'),
                              ('<unknown>', 'unknown-attachment'), ('@{bin}/unshare', 'unsupported-pattern')):
            self.assertEqual(prereq.attachment_reason(value, 'existing'), reason)
            self.assertTrue(prereq.possible_attachment(value, 'existing'))

    def test_parser_diagnostics_preserve_exit_and_capture_failures(self):
        excerpt = prereq.byte_excerpt(bytes([0, 27, 92, 255]) + b'x' * 1100)
        self.assertEqual(excerpt['captured_bytes'], 1104)
        self.assertEqual(excerpt['excerpt_bytes'], 1024)
        self.assertTrue(excerpt['truncated'])
        self.assertTrue(excerpt['value'].startswith('\\x00\\x1b\\x5c\\xff'))
        for state, code, success in (('complete', 0, True), ('complete', 1, False),
                                      ('timeout', 0, False), ('capture-limit', 0, False), ('exit-unavailable', None, False)):
            with patch.object(prereq, 'capture_parser', return_value={'state': state, 'input_state': 'complete', 'exit_code': code,
                              'stdout': b'', 'stderr': b'inert parser error\n'}), patch('builtins.print') as output:
                if success: prereq.parse('compile')
                else:
                    with self.assertRaises(RuntimeError): prereq.parse('compile')
                record = json.loads(output.call_args.args[0])
                self.assertEqual(record['state'], state)
                self.assertEqual(record['exit_code'], code)

    def capture_inert(self, chunks, clock, stdin=None, exit_code=0):
        class Pipe:
            def __init__(self, number): self.number = number
            def fileno(self): return self.number
            def close(self): pass
        class Selector:
            def __init__(self): self.keys = {}
            def register(self, pipe, events, name): self.keys[pipe.number] = SimpleNamespace(fileobj=pipe, data=name)
            def unregister(self, pipe): del self.keys[pipe.number]
            def get_map(self): return self.keys
            def select(self, timeout): return [(key, 1) for key in list(self.keys.values())]
            def close(self): pass
        kills = []
        process = SimpleNamespace(stdin=io.BytesIO() if stdin is None else stdin, stdout=Pipe(1), stderr=Pipe(2),
                                  wait=lambda timeout: exit_code, poll=lambda: exit_code, kill=lambda: kills.append(True))
        with patch.object(prereq.subprocess, 'Popen', return_value=process), \
             patch.object(prereq.selectors, 'DefaultSelector', Selector), patch.object(prereq.os, 'set_blocking'), \
             patch.object(prereq.os, 'read', side_effect=lambda fd, count: chunks[fd].pop(0)), \
             patch.object(prereq.time, 'monotonic', side_effect=clock):
            result = prereq.capture_parser('compile')
        return result, kills

    def test_bounded_parser_capture_with_inert_pipes(self):
        for chunks, clock, expected in (({1: [b'ok', b''], 2: [b'err', b'']}, [0] * 8, 'complete'),
                                        ({1: [b'x' * 4096] * 3, 2: [b'']}, [0] * 8, 'capture-limit'),
                                        ({1: [], 2: []}, [0, 11], 'timeout')):
            result, kills = self.capture_inert(chunks, clock)
            self.assertEqual(result['state'], expected)
            self.assertEqual(result['input_state'], 'complete')
            self.assertLessEqual(len(result['stdout']), 8192)
            self.assertEqual(bool(kills), expected != 'complete')

    def test_early_parser_input_closure_retains_stderr_exit_and_rejects_zero(self):
        class Input:
            def __init__(self, write_failure, close_failure):
                self.write_failure, self.close_failure = write_failure, close_failure
                self.closes = 0
            def write(self, value):
                if self.write_failure: raise BrokenPipeError('inert early exit')
                return len(value)
            def close(self):
                self.closes += 1
                if self.close_failure: raise BrokenPipeError('inert buffered close')
        for write_failure, close_failure in ((True, False), (False, True), (True, True)):
            for code in (0, 3):
                stream = Input(write_failure, close_failure)
                result, kills = self.capture_inert({1: [b'',], 2: [b'actual inert parser error\n', b'']},
                                                  [0] * 8, stream, code)
                self.assertEqual(result['state'], 'complete')
                self.assertEqual(result['input_state'], 'closed-before-input-complete')
                self.assertEqual(result['stderr'], b'actual inert parser error\n')
                self.assertEqual(result['exit_code'], code)
                self.assertEqual(stream.closes, 1)
                self.assertFalse(kills)
                with patch.object(prereq, 'capture_parser', return_value=result), patch('builtins.print') as output:
                    with self.assertRaises(RuntimeError): prereq.parse('compile')
                    record = json.loads(output.call_args.args[0])
                    self.assertEqual(record['input_state'], 'closed-before-input-complete')
                    self.assertEqual(record['exit_code'], code)
                    self.assertIn('actual inert parser error', record['stderr']['value'])

    def test_failed_policy_mutation_readback_never_establishes_ownership(self):
        own = row(prereq.NAME, prereq.TARGET, 'unconfined')
        for action in ('add', 'remove'):
            with patch.object(prereq, 'capture_parser', return_value={'state': 'complete', 'input_state': 'complete', 'exit_code': 1,
                              'stdout': b'', 'stderr': b'inert error'}), \
                 patch.object(prereq, 'inventory', return_value=[own]), patch.object(prereq, 'save') as save, \
                 patch('builtins.print') as output:
                with self.assertRaises(RuntimeError): prereq.parse(action)
                save.assert_not_called()
                self.assertEqual(json.loads(output.call_args.args[0])['purpose'], 'parser-failure-readback-not-ownership')

    def test_drift_reports_do_not_change_preservation_decisions(self):
        original = row()
        changed = dict(original, lineage=['other', 'existing'])
        with patch('builtins.print') as output:
            self.assertTrue(prereq.same_inventory([original], [original], 'inert'))
            output.assert_not_called()
            self.assertFalse(prereq.same_inventory([changed], [original], 'inert'))
            self.assertEqual(output.call_count, 2)
            for call in output.call_args_list:
                record = json.loads(call.args[0])
                self.assertEqual(record['changes']['total_records'], 1)
                self.assertNotIn('total_conflicts', record['changes'])
                self.assertLess(len(call.args[0].encode()), 32768)

    def test_provably_disjoint_literals_prefixes_and_plain_unattached(self):
        for attachment, name in (('/usr/bin/other', 'existing'), ('/opt/vendor/**', 'existing'),
                                 ('/{usr/,}bin/other', 'existing'), ('/usr/lib/{a,b}/**', 'existing'),
                                 ('unprivileged_userns', 'unprivileged_userns')):
            self.assertFalse(prereq.possible_attachment(attachment, name))
        prereq.admit_inventory([row()])
        with self.assertRaises(RuntimeError): prereq.admit_inventory([row(name=prereq.NAME)])

    def test_policy_and_commands_are_closed(self):
        self.assertEqual(prereq.POLICY, b'abi <abi/4.0>,\nprofile pscan-native-diagnostic-unshare /usr/bin/unshare flags=(unconfined) {\n  userns,\n}\n')
        for action, final in (('compile', '--add'), ('add', '--add'), ('remove', '--remove')):
            args = prereq.parser_command(action)
            self.assertEqual(args[0], '/usr/sbin/apparmor_parser')
            self.assertIn('--skip-cache', args)
            self.assertIn('--config-file=/dev/null', args)
            self.assertEqual(args[-1], final)
            self.assertNotIn('--replace', args)
        with self.assertRaises(RuntimeError): prereq.parser_command('replace')

    def test_readback_requires_exact_name_attachment_mode(self):
        own = row(prereq.NAME, prereq.TARGET, 'unconfined')
        self.assertEqual(prereq.own_row([row(), own]), own)
        for rows in ([], [own, own], [row(prereq.NAME)], [row(prereq.NAME, prereq.TARGET, 'complain')],
                     [dict(own, lineage=['other', prereq.NAME])]):
            with self.assertRaises(RuntimeError): prereq.own_row(rows)

    def test_kernel_short_names_preserve_parent_identity_and_all_metadata(self):
        # Synthetic files only: same child name under two distinct parents is valid.
        with tempfile.TemporaryDirectory(prefix='pscan-inert-policy-') as directory:
            root = Path(directory)
            (root / 'namespaces').mkdir()
            def write_entry(path, name, attachment, hashchar):
                path.mkdir(parents=True)
                for field, value in {'name': name, 'attach': attachment, 'mode': 'enforce',
                                     'sha256': hashchar * 64}.items():
                    (path / field).write_bytes((value + '\n').encode('utf-8'))
            write_entry(root / 'profiles/p1', 'parent-a', '/opt/a', 'a')
            write_entry(root / 'profiles/p2', 'parent-b', '/opt/b', 'b')
            write_entry(root / 'profiles/p1/profiles/c1', 'shared-child', 'shared-child', 'c')
            write_entry(root / 'profiles/p2/profiles/c2', 'shared-child', 'shared-child', 'd')
            with patch.object(prereq, 'POLICY_ROOT', root):
                rows = prereq.inventory()
                self.assertEqual([r['lineage'] for r in rows], [
                    ['parent-a'], ['parent-a', 'shared-child'], ['parent-b'], ['parent-b', 'shared-child']])
                self.assertEqual([r['sha256'] for r in rows], [c * 64 for c in 'acbd'])
                prereq.admit_inventory(rows)
                # Same name in the same ancestry really is ambiguous.
                write_entry(root / 'profiles/p1/profiles/c3', 'shared-child', 'shared-child', 'e')
                with self.assertRaisesRegex(RuntimeError, 'Ambiguous profile identity'):
                    prereq.inventory()

    def install_fixture(self, inventories, fail_action=None):
        calls, saved = [], {}
        def parser(action):
            calls.append(action)
            if action == fail_action: raise RuntimeError('inert parser failure')
        def save(name, value):
            self.assertNotIn(name, saved)
            saved[name] = value
        with patch.object(prereq.STATE.__class__, 'exists', return_value=False), \
             patch.object(prereq.STATE.__class__, 'is_symlink', return_value=False), \
             patch.object(prereq.STATE.__class__, 'mkdir'), patch.object(prereq, 'trusted'), \
             patch.object(prereq, 'host', return_value={'inert': 1}), \
             patch.object(prereq, 'inventory', side_effect=inventories), \
             patch.object(prereq, 'parse', side_effect=parser), patch.object(prereq, 'save', side_effect=save), \
             patch.object(prereq, 'emit'):
            try: prereq.install('b' * 64)
            except RuntimeError: return calls, saved, False
        return calls, saved, True

    def test_install_orders_compile_add_readback_and_ownership(self):
        before = [row()]
        own = row(prereq.NAME, prereq.TARGET, 'unconfined')
        calls, saved, success = self.install_fixture([before, before, before + [own]])
        self.assertTrue(success)
        self.assertEqual(calls, ['compile', 'add'])
        self.assertEqual(saved['loaded.json'], own)
        self.assertEqual(saved['baseline.json']['profiles'], before)

    def test_compile_failure_or_inventory_drift_never_adds(self):
        before = [row()]
        for inventories, failure in (([before], 'compile'), ([before, []], None)):
            calls, saved, success = self.install_fixture(inventories, failure)
            self.assertFalse(success)
            self.assertEqual(calls, ['compile'])
            self.assertFalse(saved)

    def test_add_failure_never_records_successful_ownership(self):
        calls, saved, success = self.install_fixture([[row()], [row()]], 'add')
        self.assertFalse(success)
        self.assertEqual(calls, ['compile', 'add'])
        self.assertIn('add-started.json', saved)
        self.assertNotIn('loaded.json', saved)

    def cleanup_fixture(self, current, stored_own, source='b' * 64, host_changed=False, final=None):
        baseline = [row()]
        state = {'source_sha256': 'b' * 64, 'policy_sha256': prereq.digest(prereq.POLICY),
                 'profiles': baseline, 'host': {'inert': 1}}
        def load(name):
            if name == 'baseline.json': return state
            if stored_own is None: raise RuntimeError('missing ownership record')
            return stored_own
        calls = []
        with patch.object(prereq.STATE.__class__, 'exists', return_value=True), \
             patch.object(prereq.STATE.__class__, 'iterdir', return_value=[]), \
             patch.object(prereq.STATE.__class__, 'unlink'), patch.object(prereq.STATE.__class__, 'rmdir'), \
             patch.object(prereq, 'load', side_effect=load), patch.object(prereq, 'trusted'), \
             patch.object(prereq, 'host', return_value={'inert': 2 if host_changed else 1}), \
             patch.object(prereq, 'inventory', side_effect=[current, baseline if final is None else final]), \
             patch.object(prereq, 'parse', side_effect=calls.append), patch.object(prereq, 'emit'):
            try: prereq.cleanup(source)
            except RuntimeError: return calls, False
        return calls, True

    def test_cleanup_removes_only_exact_owned_unchanged_profile(self):
        own = row(prereq.NAME, prereq.TARGET, 'unconfined')
        self.assertEqual(self.cleanup_fixture([row(), own], own), (['remove'], True))
        self.assertEqual(self.cleanup_fixture([row()], None), ([], True))
        for current, stored, kwargs in (([row(), own], None, {}),
                                        ([row(), own], dict(own, sha256='c' * 64), {}),
                                        ([row(), own], own, {'source': 'c' * 64}),
                                        ([row(), own], own, {'host_changed': True}),
                                        ([dict(row(), lineage=['changed-parent', 'existing']), own], own, {}),
                                        ([own], own, {})):
            self.assertEqual(self.cleanup_fixture(current, stored, **kwargs), ([], False))

    def test_cleanup_readback_drift_remains_failure(self):
        own = row(prereq.NAME, prereq.TARGET, 'unconfined')
        self.assertEqual(self.cleanup_fixture([row(), own], own, final=[row(), own]), (['remove'], False))


if __name__ == '__main__': unittest.main(verbosity=2)
