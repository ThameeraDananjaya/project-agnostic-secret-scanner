package gitinput

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	coverageProbePrefix = "PSCAN_COVERAGE_MARKER_"
	projectionHeader    = "PSCAN_GITLEAKS_PROJECTION_V2"
	PreparationVersion  = "pscan.byte-overlap.v2"
	MaximumRuleSpan     = int64(4020)
	DetectorPayloadSize = int64(90_000)
	DetectorOverlap     = MaximumRuleSpan - 1
	MaximumDetectorFile = int64(100_000)
)

var (
	ErrIncompleteCoverage = errors.New("incomplete Git blob coverage")
	ErrResourceLimit      = errors.New("Git blob resource limit")
)

type ProjectionClass string

const (
	ProjectionHistoryBefore ProjectionClass = "history-before"
	ProjectionHistoryAfter  ProjectionClass = "history-after"
	ProjectionTrackedTree   ProjectionClass = "tracked-tree"
)

type ProjectionLimits struct {
	MaxBlobBytes  int64
	MaxBlobCount  int
	MaxTotalBytes int64
}

// BlobProjection binds one exact Git blob to one deterministic, path-preserving
// private projection. History records are namespaced by commit and parent edge;
// the original repository path remains the suffix presented to Gitleaks.
type BlobProjection struct {
	Class        ProjectionClass
	Commit       string
	Parent       string
	ParentIndex  int
	Path         string
	OID          string
	Mode         string
	Size         int64
	RelativePath string
}

type CoveragePlan struct {
	Range       RangeBinding
	ParentEdges []ParentEdge
	Entries     []BlobProjection
	PlanDigest  string
}

type ParentEdge struct {
	Commit string
	Parent string
	Index  int
}

// MaterializedProjection exposes only private roots, deterministic bindings,
// and the exact relative probe paths needed by the private engine decoder.
type MaterializedProjection struct {
	ScanRoot           string
	ProbeRoot          string
	PlanDigest         string
	ContentDigest      string
	ProbeDigest        string
	EntryCount         int
	ExpectedProbeFiles []string
	Files              []MaterializedFile
}

type MaterializedFile struct {
	RelativePath string
	OID          string
	Mode         string
	Size         int64
	Digest       string
	RawClass     RawClass
	Preparation  string
	Chunks       []DetectorChunk
}

type DetectorChunk struct {
	RelativePath string
	Start        int64
	End          int64
	PayloadSize  int64
	FileSize     int64
	Digest       string
	Marker       string
}

type rawChange struct {
	oldMode string
	newMode string
	oldOID  string
	newOID  string
	status  string
	path    string
}

// PlanRange enumerates exact commits, every changed pre/post-image blob on each
// parent edge (including merge edges), and every blob in the exact head tree.
// It never depends on textual patches, attributes, diff drivers, or textconv.
func (g Git) PlanRange(ctx context.Context, repository, base, head string, firstRelease bool, privateHome string, limits ProjectionLimits) (CoveragePlan, error) {
	if limits.MaxBlobBytes <= 0 || limits.MaxBlobCount <= 0 || limits.MaxTotalBytes <= 0 {
		return CoveragePlan{}, fmt.Errorf("%w: invalid limits", ErrResourceLimit)
	}
	binding, err := g.BindRange(ctx, repository, base, head, firstRelease, privateHome)
	if err != nil {
		return CoveragePlan{}, err
	}
	plan := CoveragePlan{Range: binding, Entries: []BlobProjection{}}
	for commitIndex, commit := range binding.Commits {
		parentsRaw, runErr := g.run(ctx, "", privateHome, "-C", repository, "rev-list", "--parents", "-n", "1", commit, "--")
		if runErr != nil {
			return CoveragePlan{}, fmt.Errorf("%w: cannot enumerate commit parents", ErrIncompleteCoverage)
		}
		parentFields := strings.Fields(string(parentsRaw))
		if len(parentFields) == 0 || parentFields[0] != commit {
			return CoveragePlan{}, fmt.Errorf("%w: invalid commit parents", ErrIncompleteCoverage)
		}
		parents := parentFields[1:]
		for _, parent := range parents {
			if !oidPattern.MatchString(parent) {
				return CoveragePlan{}, fmt.Errorf("%w: invalid parent object", ErrIncompleteCoverage)
			}
		}
		if len(parents) == 0 {
			parents = []string{""}
		}
		for parentIndex, parent := range parents {
			plan.ParentEdges = append(plan.ParentEdges, ParentEdge{Commit: commit, Parent: parent, Index: parentIndex + 1})
			args := []string{"-C", repository, "diff-tree", "--raw", "-r", "-z", "--no-commit-id", "--no-renames", "--no-ext-diff", "--no-textconv", "--abbrev=64"}
			if parent == "" {
				args = append(args, "--root", commit, "--")
			} else {
				args = append(args, parent, commit, "--")
			}
			raw, diffErr := g.run(ctx, "", privateHome, args...)
			if diffErr != nil {
				return CoveragePlan{}, fmt.Errorf("%w: cannot enumerate commit blobs", ErrIncompleteCoverage)
			}
			changes, parseErr := parseRawChanges(raw)
			if parseErr != nil {
				return CoveragePlan{}, parseErr
			}
			for _, change := range changes {
				if err := validateChangedModes(change); err != nil {
					return CoveragePlan{}, err
				}
				if change.oldMode != "000000" {
					entry, entryErr := g.projectionEntry(ctx, repository, privateHome, ProjectionHistoryBefore, commit, parent, parentIndex+1, change.path, change.oldOID, change.oldMode, limits.MaxBlobBytes)
					if entryErr != nil {
						return CoveragePlan{}, entryErr
					}
					entry.RelativePath = historyProjectionPath(commitIndex+1, parentIndex+1, "before", change.path)
					plan.Entries = append(plan.Entries, entry)
				}
				if change.newMode != "000000" {
					entry, entryErr := g.projectionEntry(ctx, repository, privateHome, ProjectionHistoryAfter, commit, parent, parentIndex+1, change.path, change.newOID, change.newMode, limits.MaxBlobBytes)
					if entryErr != nil {
						return CoveragePlan{}, entryErr
					}
					entry.RelativePath = historyProjectionPath(commitIndex+1, parentIndex+1, "after", change.path)
					plan.Entries = append(plan.Entries, entry)
				}
			}
		}
	}
	treeRaw, err := g.run(ctx, "", privateHome, "-C", repository, "ls-tree", "-r", "-z", "--full-tree", head)
	if err != nil {
		return CoveragePlan{}, fmt.Errorf("%w: cannot enumerate tracked tree", ErrIncompleteCoverage)
	}
	treeEntries, err := parseTreeEntries(treeRaw)
	if err != nil || len(treeEntries) == 0 {
		return CoveragePlan{}, fmt.Errorf("%w: invalid or empty tracked tree", ErrIncompleteCoverage)
	}
	for _, tree := range treeEntries {
		entry, entryErr := g.projectionEntry(ctx, repository, privateHome, ProjectionTrackedTree, head, "", 0, tree.path, tree.oid, tree.mode, limits.MaxBlobBytes)
		if entryErr != nil {
			return CoveragePlan{}, entryErr
		}
		entry.RelativePath = filepath.ToSlash(filepath.Join("tree", filepath.FromSlash(tree.path)))
		plan.Entries = append(plan.Entries, entry)
	}
	if len(plan.Entries) == 0 {
		return CoveragePlan{}, fmt.Errorf("%w: no admitted blobs", ErrIncompleteCoverage)
	}
	if len(plan.Entries) > limits.MaxBlobCount {
		return CoveragePlan{}, fmt.Errorf("%w: admitted blob count exceeds maximum", ErrResourceLimit)
	}
	var totalBytes int64
	seen := make(map[string]struct{}, len(plan.Entries))
	seenFolded := make(map[string]string, len(plan.Entries))
	for _, entry := range plan.Entries {
		if entry.Size > limits.MaxTotalBytes-totalBytes {
			return CoveragePlan{}, fmt.Errorf("%w: admitted blob bytes exceed maximum", ErrResourceLimit)
		}
		totalBytes += entry.Size
		if _, exists := seen[entry.RelativePath]; exists {
			return CoveragePlan{}, fmt.Errorf("%w: conflicting projection", ErrIncompleteCoverage)
		}
		seen[entry.RelativePath] = struct{}{}
		folded := strings.ToLower(entry.RelativePath)
		if previous, exists := seenFolded[folded]; exists && previous != entry.RelativePath {
			return CoveragePlan{}, fmt.Errorf("%w: case-conflicting projection", ErrIncompleteCoverage)
		}
		seenFolded[folded] = entry.RelativePath
	}
	plan.PlanDigest = digestProjectionPlan(plan)
	return plan, nil
}

// Materialize writes an exact private raw ledger and overlapping detector
// chunks. Each detector file stays below the pinned 100,000-byte fragment and
// preserves the original basename for path-bound rules.
func (g Git) Materialize(ctx context.Context, repository, privateHome, scanRoot, probeRoot string, plan CoveragePlan) (MaterializedProjection, error) {
	if err := g.verify(); err != nil {
		return MaterializedProjection{}, err
	}
	if !safeAbsolute(repository) || !safeAbsolute(privateHome) || !safeAbsolute(scanRoot) || !safeAbsolute(probeRoot) || scanRoot == probeRoot || plan.PlanDigest == "" || plan.PlanDigest != digestProjectionPlan(plan) {
		return MaterializedProjection{}, errors.New("invalid projection binding")
	}
	for _, root := range []string{scanRoot, probeRoot} {
		if _, err := os.Lstat(root); !errors.Is(err, os.ErrNotExist) {
			return MaterializedProjection{}, errors.New("projection root must not exist")
		}
		if err := os.Mkdir(root, 0o700); err != nil {
			return MaterializedProjection{}, errors.New("cannot create projection root")
		}
	}
	contentHash, probeHash := sha256.New(), sha256.New()
	expected := []string{}
	files := make([]MaterializedFile, 0, len(plan.Entries))
	for entryIndex, entry := range plan.Entries {
		if !safeProjectionRelative(entry.RelativePath) {
			return MaterializedProjection{}, fmt.Errorf("%w: unsafe projection path", ErrIncompleteCoverage)
		}
		target := filepath.Join(scanRoot, filepath.FromSlash(entry.RelativePath))
		if err := makePrivateParent(scanRoot, target); err != nil {
			return MaterializedProjection{}, err
		}
		if err := g.writeBlob(ctx, repository, privateHome, entry.OID, entry.Size, target); err != nil {
			return MaterializedProjection{}, err
		}
		oidRaw, oidErr := g.run(ctx, "", privateHome, "-C", repository, "hash-object", "--no-filters", "--", target)
		if oidErr != nil || strings.TrimSpace(string(oidRaw)) != entry.OID {
			return MaterializedProjection{}, fmt.Errorf("%w: projected blob OID mismatch", ErrIncompleteCoverage)
		}
		contentDigest, err := fileDigest(target)
		if err != nil {
			return MaterializedProjection{}, fmt.Errorf("%w: cannot bind projected blob", ErrIncompleteCoverage)
		}
		rawClass, classErr := ClassifyRawFile(target, entry.Path)
		if classErr != nil {
			return MaterializedProjection{}, fmt.Errorf("%w: %s", ErrUnsupportedRawClass, rawClass)
		}
		bound := MaterializedFile{RelativePath: entry.RelativePath, OID: entry.OID, Mode: entry.Mode, Size: entry.Size, Digest: contentDigest, RawClass: rawClass, Preparation: PreparationVersion}
		chunks, chunkErr := writeDetectorChunks(target, probeRoot, entry.Path, plan.PlanDigest, entryIndex, entry.OID, entry.Size)
		if chunkErr != nil {
			return MaterializedProjection{}, chunkErr
		}
		bound.Chunks = chunks
		writeDigestRecord(contentHash, entry.RelativePath, entry.OID, entry.Mode, strconv.FormatInt(entry.Size, 10), contentDigest, string(rawClass), PreparationVersion)
		for _, chunk := range chunks {
			writeDigestRecord(probeHash, chunk.RelativePath, strconv.FormatInt(chunk.Start, 10), strconv.FormatInt(chunk.End, 10), chunk.Marker, chunk.Digest)
			expected = append(expected, chunk.RelativePath)
		}
		files = append(files, bound)
	}
	return MaterializedProjection{ScanRoot: scanRoot, ProbeRoot: probeRoot, PlanDigest: plan.PlanDigest, ContentDigest: hex.EncodeToString(contentHash.Sum(nil)), ProbeDigest: hex.EncodeToString(probeHash.Sum(nil)), EntryCount: len(files), ExpectedProbeFiles: expected, Files: files}, nil
}

func writeDetectorChunks(source, root, originalPath, planDigest string, entryIndex int, oid string, size int64) ([]DetectorChunk, error) {
	in, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot read raw ledger", ErrIncompleteCoverage)
	}
	defer in.Close()
	chunks := []DetectorChunk{}
	step := DetectorPayloadSize - DetectorOverlap
	for start, chunkIndex := int64(0), 0; start < size || size == 0 && chunkIndex == 0; start, chunkIndex = start+step, chunkIndex+1 {
		end := start + DetectorPayloadSize
		if end > size {
			end = size
		}
		path := filepath.ToSlash(filepath.Join("chunks", fmt.Sprintf("%08d", entryIndex), fmt.Sprintf("%08d", chunkIndex), filepath.FromSlash(originalPath)))
		if !safeProjectionRelative(path) {
			return nil, fmt.Errorf("%w: unsafe detector path", ErrIncompleteCoverage)
		}
		marker := coverageProbeMarker(planDigest, path, oid)
		header := fmt.Sprintf("%s\n%s\nentry=%08d chunk=%08d start=%d end=%d\n", projectionHeader, marker, entryIndex, chunkIndex, start, end)
		if int64(len(header))+(end-start) >= MaximumDetectorFile {
			return nil, fmt.Errorf("%w: detector chunk exceeds pinned fragment", ErrIncompleteCoverage)
		}
		target := filepath.Join(root, filepath.FromSlash(path))
		if err := makePrivateParent(root, target); err != nil {
			return nil, err
		}
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return nil, fmt.Errorf("%w: cannot create detector chunk", ErrIncompleteCoverage)
		}
		ok := false
		if _, err = io.WriteString(out, header); err == nil {
			_, err = in.Seek(start, io.SeekStart)
		}
		if err == nil {
			_, err = io.CopyN(out, in, end-start)
		}
		if closeErr := out.Close(); err == nil {
			err = closeErr
		}
		if err == nil {
			ok = true
		}
		if !ok {
			_ = os.Remove(target)
			return nil, fmt.Errorf("%w: cannot seal detector chunk", ErrIncompleteCoverage)
		}
		info, err := os.Lstat(target)
		digest, digestErr := fileDigest(target)
		if err != nil || digestErr != nil || !info.Mode().IsRegular() || info.Size() >= MaximumDetectorFile {
			return nil, fmt.Errorf("%w: cannot bind detector chunk", ErrIncompleteCoverage)
		}
		chunks = append(chunks, DetectorChunk{RelativePath: path, Start: start, End: end, PayloadSize: end - start, FileSize: info.Size(), Digest: digest, Marker: marker})
		if end == size {
			break
		}
	}
	return chunks, nil
}

// VerifyMaterializedProjection proves raw-ledger and prepared-chunk bijections,
// exact integrity, finite-overlap continuity, and absence of unbound files.
func VerifyMaterializedProjection(projection MaterializedProjection) error {
	if projection.EntryCount <= 0 || projection.EntryCount != len(projection.Files) || len(projection.ExpectedProbeFiles) == 0 || !lowerHexDigest(projection.PlanDigest) || !lowerHexDigest(projection.ContentDigest) || !lowerHexDigest(projection.ProbeDigest) {
		return errors.New("invalid materialized projection")
	}
	contentHash, probeHash := sha256.New(), sha256.New()
	rawExpected := make(map[string]fileBinding, len(projection.Files))
	chunkExpected := make(map[string]fileBinding, len(projection.ExpectedProbeFiles))
	expectedIndex := 0
	for _, file := range projection.Files {
		if !safeProjectionRelative(file.RelativePath) || !oidPattern.MatchString(file.OID) || !regularGitMode(file.Mode) || !lowerHexDigest(file.Digest) || file.Size < 0 || file.Preparation != PreparationVersion || file.RawClass != RawText && file.RawClass != RawBinary || len(file.Chunks) == 0 {
			return errors.New("invalid projection file binding")
		}
		if _, exists := rawExpected[file.RelativePath]; exists {
			return errors.New("duplicate projection file binding")
		}
		rawExpected[file.RelativePath] = fileBinding{size: file.Size, digest: file.Digest}
		writeDigestRecord(contentHash, file.RelativePath, file.OID, file.Mode, strconv.FormatInt(file.Size, 10), file.Digest, string(file.RawClass), file.Preparation)
		previousEnd := int64(0)
		for index, chunk := range file.Chunks {
			if expectedIndex >= len(projection.ExpectedProbeFiles) || projection.ExpectedProbeFiles[expectedIndex] != chunk.RelativePath || !safeProjectionRelative(chunk.RelativePath) || !lowerHexDigest(chunk.Digest) || !strings.HasPrefix(chunk.Marker, coverageProbePrefix) || chunk.PayloadSize != chunk.End-chunk.Start || chunk.FileSize >= MaximumDetectorFile || chunk.Start < 0 || chunk.End < chunk.Start || chunk.End > file.Size {
				return errors.New("invalid detector chunk binding")
			}
			if index == 0 && chunk.Start != 0 || index > 0 && chunk.Start != previousEnd-DetectorOverlap || index < len(file.Chunks)-1 && chunk.PayloadSize != DetectorPayloadSize {
				return errors.New("detector chunk mapping has a gap")
			}
			previousEnd = chunk.End
			if _, exists := chunkExpected[chunk.RelativePath]; exists {
				return errors.New("duplicate detector chunk binding")
			}
			chunkExpected[chunk.RelativePath] = fileBinding{size: chunk.FileSize, digest: chunk.Digest}
			writeDigestRecord(probeHash, chunk.RelativePath, strconv.FormatInt(chunk.Start, 10), strconv.FormatInt(chunk.End, 10), chunk.Marker, chunk.Digest)
			expectedIndex++
		}
		if previousEnd != file.Size {
			return errors.New("detector chunk mapping is incomplete")
		}
	}
	if expectedIndex != len(projection.ExpectedProbeFiles) || hex.EncodeToString(contentHash.Sum(nil)) != projection.ContentDigest || hex.EncodeToString(probeHash.Sum(nil)) != projection.ProbeDigest {
		return errors.New("projection aggregate digest mismatch")
	}
	if err := verifyBoundRoot(projection.ScanRoot, rawExpected); err != nil {
		return err
	}
	return verifyBoundRoot(projection.ProbeRoot, chunkExpected)
}

type fileBinding struct {
	size   int64
	digest string
}

func verifyBoundRoot(root string, expected map[string]fileBinding) error {
	if !safeAbsolute(root) {
		return errors.New("invalid projection root")
	}
	seen := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return errors.New("projection cannot be traversed")
		}
		if path == root || entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return errors.New("projection path cannot be normalized")
		}
		bound, ok := expected[filepath.ToSlash(rel)]
		info, infoErr := entry.Info()
		digest, digestErr := fileDigest(path)
		if !ok || infoErr != nil || digestErr != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() != bound.size || digest != bound.digest {
			return errors.New("projection integrity mismatch")
		}
		seen[filepath.ToSlash(rel)] = true
		return nil
	})
	if err != nil || len(seen) != len(expected) {
		return errors.New("projection is incomplete")
	}
	return nil
}

func (g Git) projectionEntry(ctx context.Context, repository, privateHome string, class ProjectionClass, commit, parent string, parentIndex int, path, oid, mode string, maxBlobBytes int64) (BlobProjection, error) {
	if !portableGitPath(path) || !oidPattern.MatchString(oid) || !regularGitMode(mode) {
		return BlobProjection{}, fmt.Errorf("%w: unsafe blob identity", ErrIncompleteCoverage)
	}
	typeRaw, err := g.run(ctx, "", privateHome, "-C", repository, "cat-file", "-t", oid)
	if err != nil || strings.TrimSpace(string(typeRaw)) != "blob" {
		return BlobProjection{}, fmt.Errorf("%w: object is not a blob", ErrIncompleteCoverage)
	}
	sizeRaw, err := g.run(ctx, "", privateHome, "-C", repository, "cat-file", "-s", oid)
	if err != nil {
		return BlobProjection{}, fmt.Errorf("%w: cannot read blob size", ErrIncompleteCoverage)
	}
	size, err := strconv.ParseInt(strings.TrimSpace(string(sizeRaw)), 10, 64)
	if err != nil || size < 0 {
		return BlobProjection{}, fmt.Errorf("%w: invalid blob size", ErrIncompleteCoverage)
	}
	if size > maxBlobBytes {
		return BlobProjection{}, fmt.Errorf("%w: admitted blob exceeds maximum", ErrResourceLimit)
	}
	return BlobProjection{Class: class, Commit: commit, Parent: parent, ParentIndex: parentIndex, Path: path, OID: oid, Mode: mode, Size: size}, nil
}

func (g Git) writeBlob(parent context.Context, repository, privateHome, oid string, expectedSize int64, target string) error {
	ctx, cancel := context.WithTimeout(parent, g.Timeout)
	defer cancel()
	temporary := target + ".partial"
	file, err := os.OpenFile(temporary, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("%w: cannot create blob projection", ErrIncompleteCoverage)
	}
	cmd := exec.CommandContext(ctx, g.Executable, "-C", repository, "cat-file", "blob", oid)
	cmd.Env = gitEnvironment(filepath.Dir(g.Executable), privateHome)
	cmd.Stdin = nil
	cmd.Stdout = file
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	closeErr := file.Close()
	if runErr != nil || closeErr != nil || ctx.Err() != nil || stderr.Len() != 0 {
		_ = os.Remove(temporary)
		return fmt.Errorf("%w: cannot materialize blob", ErrIncompleteCoverage)
	}
	info, err := os.Lstat(temporary)
	if err != nil || !info.Mode().IsRegular() || info.Size() != expectedSize {
		_ = os.Remove(temporary)
		return fmt.Errorf("%w: projected blob size mismatch", ErrIncompleteCoverage)
	}
	if err := os.Rename(temporary, target); err != nil {
		_ = os.Remove(temporary)
		return fmt.Errorf("%w: cannot seal blob projection", ErrIncompleteCoverage)
	}
	return nil
}

func parseRawChanges(raw []byte) ([]rawChange, error) {
	changes := []rawChange{}
	for len(raw) > 0 {
		metaEnd := bytes.IndexByte(raw, 0)
		if metaEnd <= 0 {
			return nil, fmt.Errorf("%w: malformed raw diff metadata", ErrIncompleteCoverage)
		}
		meta := string(raw[:metaEnd])
		raw = raw[metaEnd+1:]
		pathEnd := bytes.IndexByte(raw, 0)
		if pathEnd < 0 {
			return nil, fmt.Errorf("%w: malformed raw diff path", ErrIncompleteCoverage)
		}
		path := string(raw[:pathEnd])
		raw = raw[pathEnd+1:]
		fields := strings.Fields(meta)
		if len(fields) != 5 || !strings.HasPrefix(fields[0], ":") || len(fields[4]) != 1 || !strings.Contains("ADMT", fields[4]) {
			return nil, fmt.Errorf("%w: unsupported raw diff record", ErrIncompleteCoverage)
		}
		changes = append(changes, rawChange{oldMode: strings.TrimPrefix(fields[0], ":"), newMode: fields[1], oldOID: fields[2], newOID: fields[3], status: fields[4], path: path})
	}
	return changes, nil
}

type treeEntry struct{ mode, oid, path string }

func parseTreeEntries(raw []byte) ([]treeEntry, error) {
	entries := []treeEntry{}
	for len(raw) > 0 {
		end := bytes.IndexByte(raw, 0)
		if end < 0 {
			return nil, fmt.Errorf("%w: malformed tree record", ErrIncompleteCoverage)
		}
		record := raw[:end]
		raw = raw[end+1:]
		tab := bytes.IndexByte(record, '\t')
		if tab < 0 {
			return nil, fmt.Errorf("%w: malformed tree path", ErrIncompleteCoverage)
		}
		fields := strings.Fields(string(record[:tab]))
		path := string(record[tab+1:])
		if len(fields) != 3 || fields[1] != "blob" || !regularGitMode(fields[0]) || !oidPattern.MatchString(fields[2]) || !portableGitPath(path) {
			return nil, fmt.Errorf("%w: unsupported tracked-tree entry", ErrIncompleteCoverage)
		}
		entries = append(entries, treeEntry{mode: fields[0], oid: fields[2], path: path})
	}
	return entries, nil
}

func validateChangedModes(change rawChange) error {
	if !portableGitPath(change.path) || !oidOrZero(change.oldOID) || !oidOrZero(change.newOID) {
		return fmt.Errorf("%w: unsafe raw change", ErrIncompleteCoverage)
	}
	if change.oldMode != "000000" && !regularGitMode(change.oldMode) {
		return fmt.Errorf("%w: unsupported old Git mode", ErrIncompleteCoverage)
	}
	if change.newMode != "000000" && !regularGitMode(change.newMode) {
		return fmt.Errorf("%w: unsupported new Git mode", ErrIncompleteCoverage)
	}
	return nil
}

func regularGitMode(mode string) bool { return mode == "100644" || mode == "100755" }

func oidOrZero(oid string) bool {
	if oidPattern.MatchString(oid) {
		return true
	}
	return len(oid) > 0 && strings.Trim(oid, "0") == "" && (len(oid) == 40 || len(oid) == 64)
}

func portableGitPath(path string) bool {
	if !utf8.ValidString(path) || !safeRelative(path) || len(path) > 2048 {
		return false
	}
	for _, r := range path {
		if r < 0x20 || r == 0x7f || strings.ContainsRune(`:*?"<>|`, r) {
			return false
		}
	}
	for _, part := range strings.Split(path, "/") {
		if part == "" || strings.HasSuffix(part, ".") || strings.HasSuffix(part, " ") || windowsReservedName(part) {
			return false
		}
	}
	return true
}

func windowsReservedName(component string) bool {
	name := strings.ToUpper(strings.SplitN(component, ".", 2)[0])
	if name == "CON" || name == "PRN" || name == "AUX" || name == "NUL" {
		return true
	}
	if len(name) == 4 && (strings.HasPrefix(name, "COM") || strings.HasPrefix(name, "LPT")) && name[3] >= '1' && name[3] <= '9' {
		return true
	}
	return false
}

func historyProjectionPath(commitIndex, parentIndex int, side, path string) string {
	prefix := fmt.Sprintf("history/c%06d/p%03d/%s", commitIndex, parentIndex, side)
	return filepath.ToSlash(filepath.Join(prefix, filepath.FromSlash(path)))
}

func safeProjectionRelative(path string) bool {
	if !safeRelative(path) {
		return false
	}
	clean := filepath.Clean(filepath.FromSlash(path))
	return clean != "." && !filepath.IsAbs(clean) && !strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

func makePrivateParent(root, target string) error {
	rel, err := filepath.Rel(root, target)
	if err != nil || !safeProjectionRelative(filepath.ToSlash(rel)) {
		return fmt.Errorf("%w: projection escaped private root", ErrIncompleteCoverage)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return fmt.Errorf("%w: cannot create projection directory", ErrIncompleteCoverage)
	}
	return nil
}

func coverageProbeMarker(planDigest, relativePath, oid string) string {
	h := sha256.New()
	writeDigestRecord(h, "pscan.gitleaks-coverage-probe.v1", planDigest, relativePath, oid)
	return coverageProbePrefix + hex.EncodeToString(h.Sum(nil))
}

func copyWithMarker(source, destination, marker string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	ok := false
	defer func() {
		_ = out.Close()
		if !ok {
			_ = os.Remove(destination)
		}
	}()
	if _, err := io.WriteString(out, projectionHeader+"\n"+marker+"\n"); err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	ok = true
	return nil
}

func digestProjectionPlan(plan CoveragePlan) string {
	h := sha256.New()
	writeDigestRecord(h, "pscan.git-blob-plan.v2", plan.Range.Base, plan.Range.Head, plan.Range.MergeBase, strconv.FormatBool(plan.Range.FirstRelease), plan.Range.HistoryRangeDigest, plan.Range.TrackedTreeDigest, plan.Range.HeadTreeOID)
	for _, edge := range plan.ParentEdges {
		writeDigestRecord(h, edge.Commit, edge.Parent, strconv.Itoa(edge.Index))
	}
	for _, entry := range plan.Entries {
		writeDigestRecord(h, string(entry.Class), entry.Commit, entry.Parent, strconv.Itoa(entry.ParentIndex), entry.Path, entry.OID, entry.Mode, strconv.FormatInt(entry.Size, 10), entry.RelativePath)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func writeDigestRecord(writer io.Writer, values ...string) {
	var length [8]byte
	for _, value := range values {
		binary.BigEndian.PutUint64(length[:], uint64(len(value)))
		_, _ = writer.Write(length[:])
		_, _ = io.WriteString(writer, value)
	}
}
