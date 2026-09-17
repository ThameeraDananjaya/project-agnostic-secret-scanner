"""Fail closed before the single unsigned artifact upload. No provider access.

The repository administrator issues the revision-bound admission record only
after independently verifying accrued storage and zero-spend controls.
This program verifies that record and the actual complete distribution; it
does not promise provider upload capacity or treat the record as provider proof.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import stat
import sys
import time

PRODUCT = "a13c28fe7273bc8dc6545f97966a02889524eb4c"
PRODUCT_TREE = "217b711ddea51fd0ea7e808edd2e27fdecef8427"
REPOSITORY = "ThameeraDananjaya/project-agnostic-secret-scanner"
TAG = "release-tooling-v1.0.0-c2-sbom-v1"
WORKFLOW = ".github/workflows/release-build-unsigned.yml"
MAX_PAYLOAD = 240 * 1024 * 1024
ENVELOPE_RESERVE = 16 * 1024 * 1024
MAX_TRANSFER = MAX_PAYLOAD + ENVELOPE_RESERVE
TRANSFER_GIB_HOURS_MICROS = 6_000_000  # 0.25 GiB at configured one-day retention
MAX_INCLUDED_GIB_HOURS_MICROS = 500_000_000 * (28 * 24) * 1_000_000 // (2 ** 30)
MAX_RECORD_BYTES = 8192
FILES = frozenset("""
scanner-runner-linux-amd64 scanner-runner-windows-amd64.exe
scanner-release-verifier-linux-amd64 scanner-release-verifier-windows-amd64.exe
gitleaks-linux-amd64 gitleaks-windows-amd64.exe
project-agnostic-secret-scanner_v1.0.0_linux_amd64.tar.gz
project-agnostic-secret-scanner_v1.0.0_windows_amd64.zip
rules-gitleaks-v8.30.1.toml rules-gitleaks-ignore-empty-v1.txt
schema-scan-request-1.1.json schema-scan-outcome-1.0.json
schema-release-manifest-1.1.json schema-release-manifest-2.0.json
schema-release-manifest-2.1.json schema-release-manifest-2.2.json schema-release-manifest-2.3.json schema-release-manifest-2.4.json schema-release-manifest-2.5.json
schema-global-revocation-1.1.json schema-rule-pack-1.0.json
LICENSE.txt THIRD_PARTY_NOTICES.md GITLEAKS-LICENCE-MANIFEST.json
sbom.spdx.json TEST-SUMMARY.json LIMITATIONS.md COMPATIBILITY.json
global-revocations.json global-revocation-checkpoint.json BUILD-PROVENANCE.json
OFFLINE-VERIFICATION-RUNBOOK.md SCANNER-IO-REFERENCE.md CHECKSUMS.sha256
release-manifest.json
""".split())
RECORD_FIELDS = frozenset("""
schema repository owner_id tooling_revision issued_epoch expires_epoch
included_storage_gib_hours_micros accrued_storage_gib_hours_micros
net_cost_ceiling_usd net_storage_cost_usd actions_stop_usage packages_stop_usage
actions_budget_usd packages_budget_usd
storage_billing_verified artifact_retention_days evidence_sha256
""".split())


class Rejected(ValueError):
    pass


def require(value, reason):
    if not value:
        raise Rejected(reason)


def strict_object(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, "DUPLICATE_JSON_KEY")
        result[key] = value
    return result


def parse_json(raw):
    try:
        return json.loads(raw, object_pairs_hook=strict_object,
                          parse_constant=lambda _: (_ for _ in ()).throw(Rejected("NONFINITE_JSON")))
    except (UnicodeError, json.JSONDecodeError) as exc:
        raise Rejected("MALFORMED_JSON") from exc


def validate_admission(raw, revision, now):
    require(isinstance(raw, str) and 0 < len(raw.encode("utf-8")) <= MAX_RECORD_BYTES,
            "ADMISSION_RECORD_SIZE")
    require(re.fullmatch(r"[0-9a-f]{40}", revision or ""), "TOOLING_REVISION")
    record = parse_json(raw)
    require(isinstance(record, dict) and set(record) == RECORD_FIELDS, "ADMISSION_FIELDS")
    require(record["schema"] == "pscan-build-storage-admission-v2", "ADMISSION_SCHEMA")
    require(record["repository"] == REPOSITORY, "ADMISSION_REPOSITORY")
    for key in ("owner_id", "issued_epoch", "expires_epoch", "included_storage_gib_hours_micros",
                "accrued_storage_gib_hours_micros", "net_cost_ceiling_usd",
                "net_storage_cost_usd", "artifact_retention_days", "actions_budget_usd", "packages_budget_usd"):
        require(type(record[key]) is int, "ADMISSION_INTEGER:" + key)
    require(record["owner_id"] == 50274860 and record["tooling_revision"] == revision,
            "ADMISSION_IDENTITY")
    require(0 <= now - record["issued_epoch"] <= 3600 and
            now <= record["expires_epoch"] and
            0 < record["expires_epoch"] - record["issued_epoch"] <= 3600,
            "ADMISSION_FRESHNESS")
    require(record["net_cost_ceiling_usd"] == 0 and record["net_storage_cost_usd"] == 0 and
            record["actions_budget_usd"] == 0 and record["packages_budget_usd"] == 0 and
            record["actions_stop_usage"] is True and record["packages_stop_usage"] is True and
            record["storage_billing_verified"] is True and record["artifact_retention_days"] == 1,
            "ADMISSION_ZERO_SPEND_PROOF")
    require(isinstance(record["evidence_sha256"], str) and
            re.fullmatch(r"[0-9a-f]{64}", record["evidence_sha256"]) and
            record["evidence_sha256"] != "0" * 64, "ADMISSION_EVIDENCE_DIGEST")
    included = record["included_storage_gib_hours_micros"]
    accrued = record["accrued_storage_gib_hours_micros"]
    require(0 < included <= MAX_INCLUDED_GIB_HOURS_MICROS and 0 <= accrued <= included,
            "ADMISSION_BILLING_ALLOWANCE")
    require(included - accrued >= TRANSFER_GIB_HOURS_MICROS,
            "ADMISSION_INSUFFICIENT_INCLUDED_USAGE")
    return record


def no_links(path):
    for parent in (path, *path.parents):
        s = parent.lstat()
        require(not stat.S_ISLNK(s.st_mode) and not getattr(s, "st_file_attributes", 0) & 1024,
                "LINK_OR_REPARSE_PATH")


def identity(s):
    return (s.st_dev, s.st_ino, s.st_size, s.st_mtime_ns, s.st_ctime_ns, s.st_nlink)


def cross_api_identity(s):
    # Windows path stat and handle fstat may assign different creation/change
    # semantics to ctime. Compare ctime only within the same API before/after.
    return (s.st_dev, s.st_ino, s.st_size, s.st_mtime_ns, s.st_nlink)


def validate_distribution(directory, revision):
    root = Path(os.path.abspath(directory))
    no_links(root)
    require(root.is_dir(), "DISTRIBUTION_DIRECTORY")
    entries = list(root.iterdir())
    require(len(entries) == 35 and {p.name for p in entries} == FILES, "DISTRIBUTION_FILE_SET")
    sizes = {}
    for p in entries:
        no_links(p)
        s = p.lstat()
        require(stat.S_ISREG(s.st_mode) and s.st_nlink == 1 and s.st_size > 0,
                "IRREGULAR_OR_EMPTY_FILE")
        sizes[p.name] = s.st_size
    total = sum(sizes.values())
    require(total <= MAX_PAYLOAD, "PAYLOAD_LIMIT")
    require(sizes["release-manifest.json"] <= 2 * 1024 * 1024, "MANIFEST_SIZE")
    require(sizes["CHECKSUMS.sha256"] <= 16384, "CHECKSUM_SIZE")
    hashes = {}
    bodies = {}
    for p in sorted(entries):
        before = p.lstat()
        h = hashlib.sha256()
        with p.open("rb") as f:
            handle_before = os.fstat(f.fileno())
            require(cross_api_identity(before) == cross_api_identity(handle_before), "FILE_CHANGED")
            capture = bytearray() if p.name in ("release-manifest.json", "CHECKSUMS.sha256") else None
            count = 0
            while block := f.read(65536):
                count += len(block)
                require(count <= sizes[p.name], "FILE_GREW")
                h.update(block)
                if capture is not None:
                    capture.extend(block)
            require(count == sizes[p.name] and identity(handle_before) == identity(os.fstat(f.fileno())) and
                    identity(before) == identity(p.lstat()), "FILE_CHANGED")
        hashes[p.name] = h.hexdigest()
        if capture is not None:
            bodies[p.name] = bytes(capture)
    manifest = parse_json(bodies["release-manifest.json"])
    require(isinstance(manifest, dict) and manifest.get("manifestSchemaVersion") == "2.5" and
            manifest.get("releaseVersion") == "v1.0.0", "MANIFEST_VERSION")
    require(manifest.get("productSource") == {"tag": "v1.0.0", "commit": PRODUCT, "tree": PRODUCT_TREE},
            "MANIFEST_PRODUCT")
    tooling = manifest.get("releaseTooling")
    require(isinstance(tooling, dict) and tooling.get("commit") == revision and
            tooling.get("tag") == TAG and tooling.get("workflow") == WORKFLOW and
            tooling.get("workflowRef") == "refs/tags/" + TAG and
            tooling.get("workflowSha") == revision and tooling.get("trigger") == "workflow_dispatch",
            "MANIFEST_TOOLING")
    require(manifest.get("releaseState") == "unsigned-candidate" and
            "releaseIdentity" in manifest and manifest["releaseIdentity"] is None, "MANIFEST_UNSIGNED_STATE")
    require(manifest.get("buildIdentity") == {"repository": REPOSITORY, "repositoryOwnerId": 50274860,
            "workflow": WORKFLOW, "ref": "refs/tags/"+TAG, "workflowSha": revision,
            "trigger": "workflow_dispatch"}, "MANIFEST_BUILD_IDENTITY")
    assets = manifest.get("assets")
    require(isinstance(assets, list) and len(assets) == 34, "MANIFEST_ASSETS")
    seen = set()
    for asset in assets:
        require(isinstance(asset, dict) and isinstance(asset.get("path"), str), "MANIFEST_ASSET_FIELDS")
        name = asset["path"]
        require(name in FILES - {"release-manifest.json"} and name not in seen, "MANIFEST_ASSET_PATH")
        seen.add(name)
        require(type(asset.get("size")) is int and asset["size"] == sizes[name] and
                asset.get("sha256") == hashes[name], "MANIFEST_ASSET_BINDING")
    # The frozen PowerShell builder sorts by its culture. Validate the complete
    # unique binding set without substituting Python's different ordinal order.
    try:
        checksum_text = bodies["CHECKSUMS.sha256"].decode("ascii")
    except UnicodeError as exc:
        raise Rejected("CHECKSUM_BINDING") from exc
    require(checksum_text.endswith("\n") and "\r" not in checksum_text, "CHECKSUM_BINDING")
    checksum_names = set()
    for line in checksum_text[:-1].split("\n"):
        match = re.fullmatch(r"([0-9a-f]{64})  ([A-Za-z0-9_.-]+)", line)
        require(match is not None, "CHECKSUM_BINDING")
        digest, name = match.groups()
        require(name in FILES - {"release-manifest.json", "CHECKSUMS.sha256"} and
                name not in checksum_names and digest == hashes[name], "CHECKSUM_BINDING")
        checksum_names.add(name)
    require(checksum_names == FILES - {"release-manifest.json", "CHECKSUMS.sha256"}, "CHECKSUM_BINDING")
    require({p.name for p in root.iterdir()} == FILES, "DISTRIBUTION_CHANGED")
    return total, hashes["release-manifest.json"]


def check(directory, raw_admission, revision, now):
    record = validate_admission(raw_admission, revision, now)
    total, manifest_digest = validate_distribution(directory, revision)
    return {"schema": "pscan-build-storage-guard-result-v1", "status": "PASS",
            "tooling_revision": revision, "files": 35, "payload_bytes": total,
            "transfer_upper_bound_bytes": total + ENVELOPE_RESERVE,
            "maximum_transfer_bytes": MAX_TRANSFER,
            "manifest_sha256": manifest_digest,
            "admission_evidence_sha256": record["evidence_sha256"]}


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--distribution", required=True)
    args = parser.parse_args()
    try:
        result = check(args.distribution, os.environ.get("PSCAN_STORAGE_ADMISSION", ""),
                       os.environ.get("PSCAN_EXPECTED_TOOLING_REVISION", ""), int(time.time()))
        validate_admission(os.environ["PSCAN_STORAGE_ADMISSION"],
                           os.environ["PSCAN_EXPECTED_TOOLING_REVISION"], int(time.time()))
    except (Rejected, OSError, TypeError, KeyError) as exc:
        print("Storage admission rejected: " + str(exc), file=sys.stderr)
        return 1
    print(json.dumps(result, sort_keys=True))
    return 0


if __name__ == "__main__":
    sys.exit(main())
