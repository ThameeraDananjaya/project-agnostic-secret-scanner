"""Offline synthetic admission, integrity, size and filesystem rejection tests."""
import copy
import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

sys.path.insert(0, str(Path(__file__).resolve().parent))
import storage_guard as guard

REVISION = "1" * 40
NOW = 1_800_000_000


def admission():
    return dict(schema="pscan-build-storage-admission-v2", repository=guard.REPOSITORY,
                owner_id=50274860, tooling_revision=REVISION, issued_epoch=NOW,
                expires_epoch=NOW + 3600,
                included_storage_gib_hours_micros=guard.MAX_INCLUDED_GIB_HOURS_MICROS,
                accrued_storage_gib_hours_micros=125099, artifact_retention_days=1,
                net_cost_ceiling_usd=0, net_storage_cost_usd=0, actions_stop_usage=True,
                packages_stop_usage=True, storage_billing_verified=True,
                actions_budget_usd=0, packages_budget_usd=0,
                evidence_sha256="2" * 64)


def distribution(root):
    for name in guard.FILES - {"release-manifest.json", "CHECKSUMS.sha256"}:
        (root / name).write_bytes(("inert fixture: " + name + "\n").encode())
    hashes = {p.name: hashlib.sha256(p.read_bytes()).hexdigest() for p in root.iterdir()}
    (root / "CHECKSUMS.sha256").write_text("".join(f"{hashes[n]}  {n}\n" for n in sorted(hashes)), encoding="ascii", newline="\n")
    assets = [{"path": p.name, "size": p.stat().st_size, "sha256": hashlib.sha256(p.read_bytes()).hexdigest()}
              for p in sorted(root.iterdir())]
    manifest = dict(manifestSchemaVersion="2.3", releaseVersion="v1.0.0",
                    productSource=dict(tag="v1.0.0", commit=guard.PRODUCT, tree=guard.PRODUCT_TREE),
                    releaseTooling=dict(tag=guard.TAG, commit=REVISION, workflow=guard.WORKFLOW,
                                       workflowRef="refs/tags/" + guard.TAG, workflowSha=REVISION,
                                       trigger="workflow_dispatch"),
                    releaseState="unsigned-candidate", releaseIdentity=None,
                    buildIdentity=dict(repository=guard.REPOSITORY,repositoryOwnerId=50274860,workflow=guard.WORKFLOW,ref="refs/tags/"+guard.TAG,workflowSha=REVISION,trigger="workflow_dispatch"), assets=assets)
    (root / "release-manifest.json").write_text(json.dumps(manifest), encoding="utf-8")


class AdmissionTests(unittest.TestCase):
    def test_valid_and_exact_billing_allowance(self):
        record = admission()
        record["accrued_storage_gib_hours_micros"] = record["included_storage_gib_hours_micros"] - guard.TRANSFER_GIB_HOURS_MICROS
        self.assertEqual(guard.validate_admission(json.dumps(record), REVISION, NOW), record)
        record["accrued_storage_gib_hours_micros"] += 1
        with self.assertRaisesRegex(guard.Rejected, "INSUFFICIENT_INCLUDED_USAGE"):
            guard.validate_admission(json.dumps(record), REVISION, NOW)

    def test_every_missing_field_and_null(self):
        for key in admission():
            for mode in ("missing", "null"):
                with self.subTest(key=key, mode=mode):
                    record = admission()
                    if mode == "missing": del record[key]
                    else: record[key] = None
                    with self.assertRaises(guard.Rejected):
                        guard.validate_admission(json.dumps(record), REVISION, NOW)

    def test_malformed_duplicate_and_extra(self):
        values = ["", "{}\n{}", "[]", "null", '{"schema":1,"schema":2}',
                  json.dumps(dict(admission(), extra=True)), " " * 8193,
                  json.dumps(admission()).replace('"net_cost_ceiling_usd": 0', '"net_cost_ceiling_usd": NaN')]
        for raw in values:
            with self.subTest(raw=raw[:40]), self.assertRaises(guard.Rejected):
                guard.validate_admission(raw, REVISION, NOW)

    def test_admission_negative_matrix(self):
        changes = [("tooling_revision", "3" * 40), ("owner_id", 1), ("repository", "wrong"),
                   ("schema", "wrong"), ("issued_epoch", NOW + 1), ("issued_epoch", NOW - 3601),
                   ("expires_epoch", NOW - 1), ("expires_epoch", NOW + 3601),
                   ("included_storage_gib_hours_micros", guard.MAX_INCLUDED_GIB_HOURS_MICROS + 1),
                   ("included_storage_gib_hours_micros", 0), ("accrued_storage_gib_hours_micros", -1),
                   ("accrued_storage_gib_hours_micros", guard.MAX_INCLUDED_GIB_HOURS_MICROS + 1),
                   ("net_cost_ceiling_usd", 1), ("net_storage_cost_usd", 1),
                   ("actions_budget_usd", 1), ("packages_budget_usd", 1),
                   ("actions_stop_usage", False), ("actions_stop_usage", "true"),
                   ("packages_stop_usage", False), ("packages_stop_usage", 1),
                   ("artifact_retention_days", 0), ("artifact_retention_days", 2),
                   ("storage_billing_verified", False), ("storage_billing_verified", 1),
                   ("evidence_sha256", "0" * 64), ("evidence_sha256", "G" * 64)]
        changes += [(key, True) for key in ("owner_id", "issued_epoch", "expires_epoch",
                                           "included_storage_gib_hours_micros", "accrued_storage_gib_hours_micros",
                                           "net_cost_ceiling_usd", "net_storage_cost_usd", "artifact_retention_days",
                                           "actions_budget_usd", "packages_budget_usd")]
        for key, value in changes:
            with self.subTest(key=key, value=value):
                record = admission(); record[key] = value
                with self.assertRaises(guard.Rejected):
                    guard.validate_admission(json.dumps(record), REVISION, NOW)


class DistributionTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="pscan-storage-guard-test-")
        self.addCleanup(self.temp.cleanup)
        self.root = Path(self.temp.name)
        distribution(self.root)

    def check(self):
        return guard.check(self.root, json.dumps(admission()), REVISION, NOW)

    def test_valid_complete_distribution(self):
        before = {p.name: p.read_bytes() for p in self.root.iterdir()}
        result = self.check()
        self.assertEqual(result["payload_bytes"], sum(map(len, before.values())))
        self.assertEqual(result["transfer_upper_bound_bytes"], sum(map(len, before.values())) + 16 * 1024 * 1024)
        self.assertEqual(result["files"], 33)
        self.assertEqual(before, {p.name: p.read_bytes() for p in self.root.iterdir()})

    def test_extra_hidden_file(self):
        (self.root / ".unexpected").write_bytes(b"x")
        with self.assertRaisesRegex(guard.Rejected, "FILE_SET"): self.check()

    def test_missing_file(self):
        (self.root / "LICENSE.txt").unlink()
        with self.assertRaisesRegex(guard.Rejected, "FILE_SET"): self.check()

    def test_nested_directory(self):
        target = self.root / "LICENSE.txt"; target.unlink(); target.mkdir()
        with self.assertRaisesRegex(guard.Rejected, "IRREGULAR"): self.check()

    def test_empty_file(self):
        (self.root / "LICENSE.txt").write_bytes(b"")
        with self.assertRaisesRegex(guard.Rejected, "IRREGULAR"): self.check()

    def test_payload_over_cap_rejected_before_hashing(self):
        with (self.root / "LICENSE.txt").open("wb") as f: f.truncate(guard.MAX_PAYLOAD + 1)
        with self.assertRaisesRegex(guard.Rejected, "PAYLOAD_LIMIT"): self.check()

    def test_changed_asset(self):
        (self.root / "LICENSE.txt").write_bytes(b"tampered")
        with self.assertRaisesRegex(guard.Rejected, "ASSET_BINDING"): self.check()

    def test_manifest_identity_and_asset_matrix(self):
        path = self.root / "release-manifest.json"; original = json.loads(path.read_text())
        mutations = [lambda m: m.update(releaseVersion="v2.0.0"),
                     lambda m: m.update(releaseState="signing-pending"),
                     lambda m: m.update(releaseIdentity={}),
                     lambda m: m["buildIdentity"].update(workflowSha="3"*40),
                     lambda m: m["productSource"].update(commit="3" * 40),
                     lambda m: m["releaseTooling"].update(workflowSha="3" * 40),
                     lambda m: m["assets"].append(m["assets"][0]),
                     lambda m: m["assets"].__setitem__(1, m["assets"][0]),
                     lambda m: m["assets"][0].update(path="../escape"),
                     lambda m: m["assets"][0].update(size=True)]
        for index, mutate in enumerate(mutations):
            with self.subTest(index=index):
                changed = copy.deepcopy(original); mutate(changed); path.write_text(json.dumps(changed))
                with self.assertRaises(guard.Rejected): self.check()

    def test_checksum_rebound_in_manifest_still_rejects(self):
        path = self.root / "CHECKSUMS.sha256"; path.write_bytes(b"bad checksum content\n")
        mpath = self.root / "release-manifest.json"; manifest = json.loads(mpath.read_text())
        for asset in manifest["assets"]:
            if asset["path"] == path.name:
                asset.update(size=path.stat().st_size, sha256=hashlib.sha256(path.read_bytes()).hexdigest())
        mpath.write_text(json.dumps(manifest))
        with self.assertRaisesRegex(guard.Rejected, "CHECKSUM_BINDING"): self.check()

    def test_complete_checksum_set_accepts_builder_sort_order(self):
        path = self.root / "CHECKSUMS.sha256"
        lines = path.read_bytes().splitlines(keepends=True)
        path.write_bytes(b"".join(reversed(lines)))
        mpath = self.root / "release-manifest.json"; manifest = json.loads(mpath.read_text())
        for asset in manifest["assets"]:
            if asset["path"] == path.name:
                asset.update(size=path.stat().st_size, sha256=hashlib.sha256(path.read_bytes()).hexdigest())
        mpath.write_text(json.dumps(manifest))
        self.assertEqual(self.check()["status"], "PASS")

    def test_hardlink_rejects(self):
        target = self.root / "LICENSE.txt"; target.unlink()
        os.link(self.root / "LIMITATIONS.md", target)
        with self.assertRaisesRegex(guard.Rejected, "IRREGULAR"): self.check()

    def test_symlink_rejects(self):
        target = self.root / "LICENSE.txt"; target.unlink()
        try: target.symlink_to(self.root / "LIMITATIONS.md")
        except OSError as exc: self.skipTest("Host did not permit file symlink: " + str(exc))
        with self.assertRaisesRegex(guard.Rejected, "LINK_OR_REPARSE"): self.check()

    def test_cli_missing_admission_has_no_success_output(self):
        env = dict(os.environ); env.pop("PSCAN_STORAGE_ADMISSION", None)
        p = subprocess.run([sys.executable, "-B", str(Path(guard.__file__)), "--distribution", str(self.root)],
                           env=env, capture_output=True, timeout=10)
        self.assertEqual(p.returncode, 1)
        self.assertEqual(p.stdout, b"")
        self.assertIn(b"ADMISSION_RECORD_SIZE", p.stderr)


if __name__ == "__main__":
    unittest.main(verbosity=2)
