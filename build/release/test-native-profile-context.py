"""Inert context tests: no selector, namespace, policy or native process execution."""
import copy
import importlib.util
import io
from pathlib import Path
from types import SimpleNamespace
import unittest
from unittest.mock import patch

spec = importlib.util.spec_from_file_location('context', Path(__file__).with_name('native-profile-context.py'))
ctx = importlib.util.module_from_spec(spec)
spec.loader.exec_module(ctx)


def status():
    values = {k: '0000000000000000' if k.startswith('Cap') else '0' for k in ctx.FIELDS}
    values.update(Uid='1001 1001 1001 1001', Gid='1001 1001 1001 1001', Groups='1001 116')
    return values


def encoded(values):
    return ''.join(k + ':\t' + v + '\n' for k, v in values.items()).encode()


def baseline():
    return {'schema': 'pscan-native-profile-context-v1', 'label': 'unconfined',
            'status': ctx.parse_status(encoded(status())),
            'namespaces': {name: name + ':[1]' for name in ctx.NAMESPACES},
            'executable': {'mode': 0o100777, 'sha256': 'a' * 64, 'ino': 1},
            'sources': {name: 'b' * 64 for name in ctx.SOURCES}}


class Tests(unittest.TestCase):
    def test_executable_observation_rejects_replacement_privilege_and_drift(self):
        def info(**changes):
            values = dict(st_dev=1, st_ino=2, st_mode=0o100777, st_uid=0, st_gid=0,
                          st_size=5, st_mtime_ns=1, st_ctime_ns=1)
            values.update(changes)
            return SimpleNamespace(**values)
        class Stream(io.BytesIO):
            def fileno(self): return 10
        cases = [('stable', info(), info(), info(), info(), [], True),
                 ('setid', info(st_mode=0o104777), info(), info(), info(), [], False),
                 ('not-root-owned', info(st_uid=1001), info(), info(), info(), [], False),
                 ('capability', info(), info(), info(), info(), ['security.capability'], False),
                 ('write-drift', info(), info(st_mtime_ns=2), info(), info(), [], False),
                 ('path-replaced', info(), info(), info(st_ino=3), info(), [], False),
                 ('proc-differs', info(), info(), info(), info(st_ino=3), [], False)]
        for name, before, after, pathinfo, procinfo, xattrs, passed in cases:
            with self.subTest(name=name), patch.object(ctx.os, 'readlink', return_value=ctx.TARGET), \
                 patch.object(ctx.os, 'O_NOFOLLOW', 0x20000, create=True), \
                 patch.object(ctx.os, 'O_NONBLOCK', 0x800, create=True), \
                 patch.object(ctx.os, 'open', return_value=10), \
                 patch.object(ctx.os, 'fdopen', return_value=Stream(b'inert')), \
                 patch.object(ctx.os, 'fstat', side_effect=[before, after]), \
                 patch.object(ctx.os, 'stat', side_effect=[pathinfo, procinfo]), \
                 patch.object(ctx.os, 'listxattr', return_value=xattrs, create=True):
                if passed:
                    value = ctx.executable(123)
                    self.assertEqual(value['sha256'], ctx.sha(b'inert'))
                    self.assertEqual(value['mode'], 0o100777)
                else:
                    with self.assertRaises(RuntimeError): ctx.executable(123)

    def test_exact_context_allows_only_label_change(self):
        before = baseline()
        after = dict(before, label=ctx.LABEL)
        ctx.compare(before, after)
        # Equality is observation only; writable runtime is not advertised as immutable.
        self.assertEqual(before['executable']['mode'], 0o100777)

    def test_every_context_class_mismatch_rejects(self):
        for key in ('schema', 'status', 'namespaces', 'executable', 'sources'):
            with self.subTest(key=key):
                before = baseline()
                after = copy.deepcopy(before)
                after['label'] = ctx.LABEL
                after[key] = 'different'
                with self.assertRaisesRegex(RuntimeError, 'context mismatch'):
                    ctx.compare(before, after)

    def test_stacked_hidden_missing_or_wrong_mode_labels_reject(self):
        for label in ('unconfined', '', ctx.LABEL + '\n', 'pscan-native-diagnostic (enforce)',
                      'pscan-native-diagnostic//&unconfined (unconfined)',
                      ':other:pscan-native-diagnostic (unconfined)'):
            with self.subTest(label=label), self.assertRaises(RuntimeError):
                ctx.compare(baseline(), dict(baseline(), label=label))

    def test_privileged_or_malformed_status_rejects(self):
        for field, value in [('Uid', '0 1001 1001 1001'), ('Uid', '1001'), ('Groups', 'x'),
                             ('CapEff', '0000000000000001'), ('CapPrm', '0000000000000001'),
                             ('CapAmb', '0000000000000001'), ('CapInh', '0000000000000001'),
                             ('CapBnd', 'missing'), ('Seccomp', '-1'), ('NoNewPrivs', '')]:
            values = status()
            values[field] = value
            with self.subTest(field=field, value=value), self.assertRaises(RuntimeError):
                ctx.parse_status(encoded(values))
        values = status()
        del values['Seccomp_filters']
        with self.assertRaises(RuntimeError): ctx.parse_status(encoded(values))
        with self.assertRaises(RuntimeError): ctx.parse_status(encoded(status()) + b'Uid: 1 1 1 1\n')

    def test_changed_group_capability_and_namespace_ids_reject(self):
        changes = [('status', 'Groups', ['1001']), ('status', 'CapBnd', ['000001ffffffffff']),
                   ('status', 'NoNewPrivs', ['1']), ('status', 'Seccomp', ['2']),
                   ('namespaces', 'user', 'user:[2]'), ('namespaces', 'pid', 'pid:[2]'),
                   ('executable', 'sha256', 'c' * 64), ('sources', 'docker-execution.ps1', 'c' * 64)]
        for key, field, value in changes:
            after = copy.deepcopy(baseline())
            after['label'] = ctx.LABEL
            after[key][field] = value
            with self.subTest(key=key, field=field), self.assertRaises(RuntimeError):
                ctx.compare(baseline(), after)

    def test_unknown_baseline_or_selected_baseline_rejects(self):
        for before in (dict(baseline(), callback='anything'), dict(baseline(), label=ctx.LABEL), {}):
            with self.assertRaises(RuntimeError): ctx.validate_baseline(before)

    def test_context_collection_observes_parent_fields_and_sources(self):
        def read(path, maximum):
            path = str(path)
            if path.endswith('/comm'): return b'pwsh\n'
            if path.endswith('/status'): return encoded(status())
            if path.endswith('/attr/current'): return ctx.LABEL.encode() + b'\n'
            return b'inert committed bytes'
        with patch.object(ctx, 'read', side_effect=read), \
             patch.object(ctx.os, 'readlink', side_effect=lambda p: p.rsplit('/', 1)[-1] + ':[1]'), \
             patch.object(ctx, 'executable', return_value={'inert': True}) as exe:
            value = ctx.context(123)
        exe.assert_called_once_with(123)
        self.assertEqual(value['label'], ctx.LABEL)
        self.assertEqual(set(value['sources']), set(ctx.SOURCES))
        self.assertEqual(set(value['namespaces']), set(ctx.NAMESPACES))


if __name__ == '__main__': unittest.main(verbosity=2)
