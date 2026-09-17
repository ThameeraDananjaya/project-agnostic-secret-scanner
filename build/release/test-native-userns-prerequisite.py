"""Inert prerequisite admission/lifecycle tests; never run a host utility."""
import importlib.util
from pathlib import Path
import tempfile
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

    def test_conflicts_and_unknowns_reject(self):
        for attachment in ('/usr/bin/unshare', '/usr/bin/*', '/**', '<unknown>', '', '/usr/bin/unsh?re',
                           '/usr/bin/[ux]nshare', '/{usr/,}bin/unshare', '@{bin}/unshare', '/usr/bin/\\unshare',
                           '/**/not-unshare', '{broken', '/usr/bin/unshare\n'):
            with self.subTest(attachment=attachment):
                self.assertTrue(prereq.possible_attachment(attachment, 'existing'))
                with self.assertRaises(RuntimeError): prereq.admit_inventory([row(attachment=attachment)])

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
